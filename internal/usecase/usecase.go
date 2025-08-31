package usecase

import (
	"context"

	"github.com/RozmiDan/wb_tech_testtask/internal/entity"
	"go.uber.org/zap"
)

type RepoLayer interface {
	GetOrderByUID(ctx context.Context, orderUID string) (*entity.OrderInfo, error)
	SetOrder(ctx context.Context, order *entity.OrderInfo) error
	GetLatestOrders(ctx context.Context, limit int) ([]*entity.OrderInfo, error)
}

type OrderCache interface {
	Put(key string, val *entity.OrderResponse)
	Get(key string) *entity.OrderResponse
}

type UsecaseLayer struct {
	log   *zap.Logger
	db    RepoLayer
	cache OrderCache
}

func New(logger *zap.Logger, dbLayer RepoLayer, cache OrderCache) *UsecaseLayer {
	return &UsecaseLayer{
		log:   logger.With(zap.String("layer", "Usecase")),
		db:    dbLayer,
		cache: cache,
	}
}
