package postgre

import (
	"github.com/RozmiDan/wb_tech_testtask/pkg/postgres"
	"go.uber.org/zap"
)

type RatingRepository struct {
	pg  *postgres.Postgres
	log *zap.Logger
}

func New(pg *postgres.Postgres, logger *zap.Logger) *RatingRepository {
	return &RatingRepository{
		pg:  pg,
		log: logger.With(zap.String("layer", "Repository")),
	}
}
