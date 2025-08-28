package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RozmiDan/wb_tech_testtask/db"
	"github.com/RozmiDan/wb_tech_testtask/internal/config"
	"github.com/RozmiDan/wb_tech_testtask/internal/controller/server"
	"github.com/RozmiDan/wb_tech_testtask/internal/entity"
	"github.com/RozmiDan/wb_tech_testtask/internal/repo/postgre"
	"github.com/RozmiDan/wb_tech_testtask/internal/usecase"
	lru_cache "github.com/RozmiDan/wb_tech_testtask/pkg/cache"
	"github.com/RozmiDan/wb_tech_testtask/pkg/logger"
	"github.com/RozmiDan/wb_tech_testtask/pkg/postgres"
	"go.uber.org/zap"
)

func Run(cfg *config.Config) {
	// New logger
	logger := logger.NewLogger(cfg.Env, cfg.LogsPath)
	logger.Info("Starting programm")
	// Migrations
	db.SetupPostgres(cfg, logger)
	logger.Info("Migrations completed successfully\n")
	// Db connection
	pg, err := postgres.New(cfg.PostgresURL, postgres.MaxPoolSize(5))
	if err != nil {
		logger.Error("Cant open database", zap.Error(err))
		os.Exit(1)
	}
	defer pg.Close()
	// repo
	repo := postgre.New(pg, logger)
	// cache
	cache := lru_cache.NewLruCache[string, *entity.OrderResponse](cfg.CacheCap, nil)

	// Kafka

	// usecase
	uc := usecase.New(logger, repo, cache)

	// прогрев кэша
	warmCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := uc.WarmCacheLatest(warmCtx, cfg.CacheCap); err != nil {
		logger.Warn("cache warm failed", zap.Error(err))
	} else {
		logger.Warn("cache warm success")
	}

	// server
	server := server.InitServer(cfg, logger, uc)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Info("starting server", zap.String("port", cfg.HTTPPort))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server error", zap.Error(err))
			os.Exit(1)
		}
	}()

	<-stop
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.HTTPTimeout*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown error", zap.Error(err))
	} else {
		logger.Info("Server gracefully stopped")
	}

	logger.Info("Finishing programm")
}
