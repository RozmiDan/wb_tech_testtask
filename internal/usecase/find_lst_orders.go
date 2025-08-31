package usecase

import (
	"context"
	"errors"

	"github.com/RozmiDan/wb_tech_testtask/internal/entity"
	"go.uber.org/zap"
)

func (u *UsecaseLayer) WarmCacheLatest(ctx context.Context, cacheCap int) error {
	// 1) забираем request_id
	reqID, _ := ctx.Value(entity.RequestIDKey{}).(string)
	// 2) оборачиваем логгер
	logger := u.log.With(zap.String("func", "WarmCacheLatest"))
	if reqID != "" {
		logger = logger.With(zap.String("request_id", reqID))
	}
	if cacheCap <= 0 {
		logger.Error("invalid cache capacity")

		return errors.New("invalid cache capacity value")
	}
	orders, err := u.db.GetLatestOrders(ctx, cacheCap)
	if err != nil {
		logger.Error("Cant find values", zap.Int("count", len(orders)))

		return err
	}
	for _, o := range orders {
		dto := mapOrderToResponse(o)
		u.cache.Put(dto.OrderUID, dto)
	}
	logger.Info("cache warmed", zap.Int("count", len(orders)))

	return nil
}
