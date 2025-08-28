package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/RozmiDan/wb_tech_testtask/internal/entity"
	"go.uber.org/zap"
)

func (u *UsecaseLayer) AddOrderInfo(ctx context.Context, order *entity.OrderInfo) error {
	// 1) забираем request_id
	reqID, _ := ctx.Value(entity.RequestIDKey{}).(string)
	// 2) оборачиваем логгер
	logger := u.log.With(zap.String("func", "AddOrderInfo"))
	if reqID != "" {
		logger = logger.With(zap.String("request_id", reqID))
	}
	// 3) валидируем order
	if err := validateOrder(order); err != nil {
		logger.Warn("invalid order payload", zap.Error(err))
		return entity.ErrInvalidInput
	}
	// 4) записываем в бд
	if err := u.db.SetOrder(ctx, order); err != nil {
		switch {
		case errors.Is(err, entity.ErrorOrderExists):
			logger.Info("order already exists (idempotent)", zap.String("order_uid", order.OrderUID))
			return entity.ErrAlreadyExists
		case errors.Is(err, entity.ErrorDBConnect):
			logger.Error("db connect failed", zap.Error(err))
			return entity.ErrInternal
		case errors.Is(err, entity.ErrorInsertDB):
			logger.Error("db insert failed", zap.Error(err))
			return entity.ErrInternal
		default:
			logger.Error("unexpected repo error", zap.Error(err))
			return entity.ErrInternal
		}
	}
	// 5) пишем в кэш 
	u.cache.Put(order.OrderUID, mapOrderToResponse(order))

	logger.Info("succsessfuly add order", zap.String("order_uid", order.OrderUID))
	return nil
}

func validateOrder(o *entity.OrderInfo) error {
	if o == nil {
		return errors.New("nil order")
	}

	if o.OrderUID == "" {
		return errors.New("empty order_uid")
	}
	if o.TrackNumber == "" {
		return errors.New("empty track_number")
	}
	if o.Entry == "" {
		return errors.New("empty entry")
	}
	if o.Locale == "" {
		return errors.New("empty locale")
	}
	if o.CustomerID == "" {
		return errors.New("empty customer_id")
	}
	if o.DeliveryService == "" {
		return errors.New("empty delivery_service")
	}
	if o.DateCreated.IsZero() {
		return errors.New("empty date_created")
	}

	// Delivery
	if o.Delivery.Name == "" {
		return errors.New("empty delivery.name")
	}
	if o.Delivery.Phone == "" {
		return errors.New("empty delivery.phone")
	}
	if o.Delivery.City == "" {
		return errors.New("empty delivery.city")
	}
	if o.Delivery.Address == "" {
		return errors.New("empty delivery.address")
	}
	if o.Delivery.Email == "" {
		return errors.New("empty delivery.email")
	}

	// Payment
	if o.Payment.Transaction == "" {
		return errors.New("empty payment.transaction")
	}
	if o.Payment.Currency == "" {
		return errors.New("empty payment.currency")
	}
	if o.Payment.Provider == "" {
		return errors.New("empty payment.provider")
	}
	if o.Payment.Amount <= 0 {
		return errors.New("invalid payment.amount")
	}
	if o.Payment.PaymentDT <= 0 {
		return errors.New("invalid payment.payment_dt")
	}
	if o.Payment.Bank == "" {
		return errors.New("empty payment.bank")
	}

	// Items
	if len(o.Items) == 0 {
		return errors.New("empty items")
	}
	for i, it := range o.Items {
		if it.ChrtID == 0 {
			return fmt.Errorf("item[%d]: empty chrt_id", i)
		}
		if it.TrackNumber == "" {
			return fmt.Errorf("item[%d]: empty track_number", i)
		}
		if it.Name == "" {
			return fmt.Errorf("item[%d]: empty name", i)
		}
		if it.Price <= 0 {
			return fmt.Errorf("item[%d]: invalid price", i)
		}
		if it.TotalPrice <= 0 {
			return fmt.Errorf("item[%d]: invalid total_price", i)
		}
	}

	return nil
}
