package postgre

import (
	"context"

	"github.com/RozmiDan/wb_tech_testtask/internal/entity"
	"go.uber.org/zap"
)

func (rr *RatingRepository) GetOrderByUID(ctx context.Context, orderUID string) (*entity.OrderInfo, error) {
	reqID, _ := ctx.Value(entity.RequestIDKey{}).(string)

	logger := rr.log.With(zap.String("func", "GetOrderByUID"))
	if reqID != "" {
		logger = logger.With(zap.String("request_id", reqID))
	}

	

	logger.Info("succsessful search in db")
	return &entity.OrderInfo{}, nil
}
