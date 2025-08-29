package db

import (
	"embed"
	"os"

	"github.com/RozmiDan/wb_tech_testtask/internal/config"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func SetupPostgres(cfg *config.Config, logger *zap.Logger) {

	pgxConf := pgx.ConnConfig{
		Host:     cfg.PostgresHost,
		Port:     cfg.PostgresPort,
		Database: cfg.PostgresDB,
		User:     cfg.PostgresUser,
		Password: cfg.PostgresPass,
	}

	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		logger.Error("can't set dialect in goose",
			zap.Error(err),
		)
		os.Exit(1)
	}

	db := stdlib.OpenDB(pgxConf)
	if err := goose.Up(db, "migrations"); err != nil {
		logger.Error("can't setup migrations",
			zap.Error(err),
		)
		os.Exit(1)
	}
}

