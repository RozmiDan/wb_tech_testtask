package postgre

import (
	"context"

	"github.com/RozmiDan/wb_tech_testtask/internal/entity"
	"github.com/RozmiDan/wb_tech_testtask/pkg/postgres"
	"go.uber.org/zap"
)

type RatingRepository struct {
	pg     *postgres.Postgres
	log *zap.Logger
}

func New(pg *postgres.Postgres, logger *zap.Logger) *RatingRepository {
	return &RatingRepository{
		pg:     pg,
		log: logger.With(zap.String("layer", "Repository")),
	}
}

func (rr *RatingRepository) GetOrderByUID(ctx context.Context, orderUID string) (*entity.OrderInfo, error) {
	reqID, _ := ctx.Value(entity.RequestIDKey{}).(string)

	logger := rr.log.With(zap.String("func", "GetOrderByUID"))
	if reqID != "" {
		logger = logger.With(zap.String("request_id", reqID))
	}

	logger.Info("succsessful search in db")
	return &entity.OrderInfo{}, nil
}
