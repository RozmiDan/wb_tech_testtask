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
	"github.com/RozmiDan/wb_tech_testtask/internal/controller/http/server"
	"github.com/RozmiDan/wb_tech_testtask/internal/controller/kafka"
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

	rootCtx, rootCancel := context.WithCancel(context.Background())
	defer rootCancel()

	// repo
	repo := postgre.New(pg, logger)

	// cache
	cache := lru_cache.NewLruCache[string, *entity.OrderResponse](cfg.CacheCap, nil)

	// usecase
	uc := usecase.New(logger, repo, cache)

	// Kafka
	kafkaConsumer := kafka.NewConsumer(cfg, uc, logger)

	go func() {
		kafkaConsumer.Start(rootCtx, cfg)
		logger.Info("kafka consumer stopped")
	}()

	// прогрев кэша
	warmCtx, cancel := context.WithTimeout(rootCtx, 2*time.Second)
	if err := uc.WarmCacheLatest(warmCtx, cfg.CacheCap); err != nil {
		logger.Warn("cache warm failed", zap.Error(err))
	} else {
		logger.Warn("cache warm success")
	}
	cancel()

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

	rootCancel()
	_ = kafkaConsumer.Close()

	ctx, cancel2 := context.WithTimeout(context.Background(), cfg.HTTPTimeout*time.Second)
	defer cancel2()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown error", zap.Error(err))
	} else {
		logger.Info("Server gracefully stopped")
	}

	logger.Info("Finishing programm")
}
