package usecase

import (
	"context"
	"errors"

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
	if err := order.ValidateOrder(); err != nil {
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
