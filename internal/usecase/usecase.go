package usecase

import (
	"context"

	"github.com/RozmiDan/wb_tech_testtask/internal/entity"
	"go.uber.org/zap"
)

type RepoLayer interface {
	GetOrderByUID(ctx context.Context, orderUID string) (*entity.OrderInfo, error)
	SetOrder(ctx context.Context, order *entity.OrderInfo) error
}

type UsecaseLayer struct {
	log *zap.Logger
	db  RepoLayer
}

func New(logger *zap.Logger, dbLayer RepoLayer) *UsecaseLayer {
	return &UsecaseLayer{
		log: logger.With(zap.String("layer", "Usecase")),
		db:  dbLayer,
	}
}