package postgre

import (
	"context"
	"database/sql"
	"time"

	"github.com/RozmiDan/wb_tech_testtask/internal/entity"
	"go.uber.org/zap"
)

const selectOrderQuery = ` SELECT
		o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature,
		o.customer_id, o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard,

		d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,

		p.transaction, p.request_id, p.currency, p.provider, p.amount,
		p.payment_dt, p.bank, p.delivery_cost, p.goods_total, p.custom_fee,

		i.chrt_id, i.track_number, i.price, i.rid, i.name AS item_name,
		i.sale, i.size, i.total_price, i.nm_id, i.brand, i.status
	FROM orders o
	LEFT JOIN deliveries d ON d.order_uid = o.order_uid
	LEFT JOIN payments   p ON p.order_uid = o.order_uid
	LEFT JOIN items      i ON i.order_uid = o.order_uid
	WHERE o.order_uid = $1
	ORDER BY i.chrt_id; 
`

func (rr *RatingRepository) GetOrderByUID(ctx context.Context, orderUID string) (*entity.OrderInfo, error) {
	reqID, _ := ctx.Value(entity.RequestIDKey{}).(string)

	logger := rr.log.With(zap.String("func", "GetOrderByUID"))
	if reqID != "" {
		logger = logger.With(zap.String("request_id", reqID))
	}

	rows, err := rr.pg.Pool.Query(ctx, selectOrderQuery, orderUID)
	if err != nil {
		logger.Error("query failed", zap.Error(err))
		return nil, entity.ErrorQueryFailed
	}
	defer rows.Close()

	var (
		order *entity.OrderInfo
		seen  bool
	)

	for rows.Next() {
		var (
			// order
			ordUID, trackNumber, entry, locale, internalSig, customerID, delivSvc, shardkey, oofShard string
			smID                                                                                      int
			dateCreated                                                                               time.Time
			// delivery
			dName, dPhone, dZip, dCity, dAddr, dRegion, dEmail sql.NullString
			// payment
			pTrans, pReqID, pCurr, pProv, pBank                 sql.NullString
			pAmount, pDelCost, pGoodsTotal, pCustomFee, pPaidDt sql.NullInt64

			// item
			chrtID                                sql.NullInt64
			itemTrack, rid, itemName, size, brand sql.NullString
			price, sale, totalPrice, nmID, status sql.NullInt64
		)

		if err := rows.Scan(
			// order
			&ordUID, &trackNumber, &entry, &locale, &internalSig,
			&customerID, &delivSvc, &shardkey, &smID, &dateCreated, &oofShard,
			// delivery
			&dName, &dPhone, &dZip, &dCity, &dAddr, &dRegion, &dEmail,
			// payment
			&pTrans, &pReqID, &pCurr, &pProv, &pAmount,
			&pPaidDt, &pBank, &pDelCost, &pGoodsTotal, &pCustomFee,
			// items
			&chrtID, &itemTrack, &price, &rid, &itemName,
			&sale, &size, &totalPrice, &nmID, &brand, &status,
		); err != nil {
			logger.Error("scan failed", zap.Error(err))
			return nil, entity.ErrorQueryFailed
		}

		if !seen {
			order = &entity.OrderInfo{
				OrderUID:          ordUID,
				TrackNumber:       trackNumber,
				Entry:             entry,
				Locale:            locale,
				InternalSignature: internalSig,
				CustomerID:        customerID,
				DeliveryService:   delivSvc,
				ShardKey:          shardkey,
				SmID:              smID,
				DateCreated:       dateCreated.UTC(),
				OofShard:          oofShard,
			}
			// delivery
			order.Delivery = entity.DeliveryInfo{
				Name:    dName.String,
				Phone:   dPhone.String,
				Zip:     dZip.String,
				City:    dCity.String,
				Address: dAddr.String,
				Region:  dRegion.String,
				Email:   dEmail.String,
			}
			// payment
			order.Payment = entity.PaymentInfo{
				Transaction:  pTrans.String,
				RequestID:    pReqID.String,
				Currency:     pCurr.String,
				Provider:     pProv.String,
				Amount:       pAmount.Int64,
				Bank:         pBank.String,
				DeliveryCost: pDelCost.Int64,
				GoodsTotal:   pGoodsTotal.Int64,
				CustomFee:    pCustomFee.Int64,
				PaymentDT:    pPaidDt.Int64,
			}
			order.Items = make([]entity.ItemInfo, 0, 1)
			seen = true
		}

		if !chrtID.Valid {
			continue
		}
		order.Items = append(order.Items, entity.ItemInfo{
			ChrtID:      chrtID.Int64,
			TrackNumber: itemTrack.String,
			Price:       price.Int64,
			Rid:         rid.String,
			Name:        itemName.String,
			Sale:        sale.Int64,
			Size:        size.String,
			TotalPrice:  totalPrice.Int64,
			NmID:        nmID.Int64,
			Brand:       brand.String,
			Status:      status.Int64,
		})
	}
	if err := rows.Err(); err != nil {
		logger.Error("rows error", zap.Error(err))
		return nil, entity.ErrorQueryFailed
	}
	if !seen {
		logger.Info("order not found", zap.String("order_uid", orderUID))
		return nil, entity.ErrorOrderNotFound
	}

	logger.Info("The request was completed successfully", zap.String("order_uid", orderUID))
	return order, nil
}
