package postgre

import (
	"context"

	"github.com/RozmiDan/wb_tech_testtask/internal/entity"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

const (
	insertOrdersQuery = `
        INSERT INTO orders (
			order_uid, track_number, entry, locale, 
			internal_signature, customer_id, delivery_service, 
			shardkey, sm_id, date_created, oof_shard)
        VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		ON CONFLICT (order_uid) DO NOTHING
    `
	insertDeliveryQuery = `
		INSERT INTO deliveries (
    		order_uid, name, phone, zip, city, address, region, email) 
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		ON CONFLICT (order_uid) DO NOTHING
	`
	insertPaymentQuery = `
		INSERT INTO payments (
			order_uid, transaction, request_id, currency, provider,
			amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) 
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		ON CONFLICT (order_uid) DO NOTHING
	`
	insertItemQuery = `
		INSERT INTO items (
			order_uid, chrt_id, track_number, price, rid, name,
			sale, size, total_price, nm_id, brand, status) 
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
	`
)

func (rr *RatingRepository) SetOrder(ctx context.Context, order *entity.OrderInfo) error {
	reqID, _ := ctx.Value(entity.RequestIDKey{}).(string)
	logger := rr.log.With(zap.String("func", "SetOrder"))
	if reqID != "" {
		logger = logger.With(zap.String("request_id", reqID))
	}

	tx, err := rr.pg.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:       pgx.RepeatableRead,
		AccessMode:     pgx.ReadWrite,
		DeferrableMode: pgx.NotDeferrable,
	})
	if err != nil {
		logger.Error("begin tx failed", zap.Error(err))
		return entity.ErrorDBConnect
	}

	defer func() { _ = tx.Rollback(ctx) }()

	// 1) orders
	if cmdTg, err := tx.Exec(ctx, insertOrdersQuery, order.OrderUID, order.TrackNumber, order.Entry,
		order.Locale, order.InternalSignature, order.CustomerID,
		order.DeliveryService, order.ShardKey, order.SmID,
		order.DateCreated, order.OofShard,
	); err != nil {
		logger.Error("insert orders failed", zap.Error(err))
		return entity.ErrorInsertDB
	} else if cmdTg.RowsAffected() == 0 {
		logger.Info("orders upsert skipped (already exists)", zap.String("order_uid", order.OrderUID))
		return entity.ErrorOrderExists
	}

	// 2) deliveries
	if cmdTg, err := tx.Exec(ctx, insertDeliveryQuery, order.OrderUID,
		order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Address, order.Delivery.Region,
		order.Delivery.Email,
	); err != nil {
		logger.Error("insert deliveries failed", zap.Error(err))
		return entity.ErrorInsertDB
	} else if cmdTg.RowsAffected() == 0 {
		logger.Info("deliveries upsert skipped (already exists)", zap.String("order_uid", order.OrderUID))
	}

	// 3) payments
	if cmdTg, err := tx.Exec(ctx, insertPaymentQuery, order.OrderUID,
		order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency,
		order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDT,
		order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal,
		order.Payment.CustomFee,
	); err != nil {
		logger.Error("insert payments failed", zap.Error(err))
		return entity.ErrorInsertDB
	} else if cmdTg.RowsAffected() == 0 {
		logger.Info("payments upsert skipped (already exists)", zap.String("order_uid", order.OrderUID))
	}

	// 4) items
	for _, it := range order.Items {
		if cmdTg, err := tx.Exec(ctx, insertItemQuery,
			order.OrderUID, it.ChrtID, it.TrackNumber, it.Price, it.Rid, it.Name,
			it.Sale, it.Size, it.TotalPrice, it.NmID, it.Brand, it.Status,
		); err != nil {
			logger.Error("insert item failed", zap.Int64("chrt_id", it.ChrtID), zap.Error(err))
			return entity.ErrorInsertDB
		} else if cmdTg.RowsAffected() == 0 {
			logger.Info("item upsert skipped",
				zap.String("order_uid", order.OrderUID),
				zap.Int64("chrt_id", it.ChrtID),
			)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		logger.Error("commit failed", zap.Error(err))
		return entity.ErrorInsertDB
	}

	logger.Info("order successfully inserted", zap.String("order_uid", order.OrderUID))
	return nil
}
