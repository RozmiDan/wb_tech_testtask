package server

import (
	"context"
	"net/http"

	"github.com/RozmiDan/wb_tech_testtask/internal/config"
	"github.com/RozmiDan/wb_tech_testtask/internal/controller/handlers/mainhandler"
	custommiddleware "github.com/RozmiDan/wb_tech_testtask/internal/controller/middleware"
	"github.com/RozmiDan/wb_tech_testtask/internal/entity"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type UseCase interface {
	GetOrderInfo(ctx context.Context, orderUID string) (*entity.OrderResponse, error)
}

func InitServer(cfg *config.Config, logger *zap.Logger, uc UseCase) *http.Server {
	baseLog := logger.With(zap.String("layer", "Controller"))

	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.URLFormat)
	// router.Use(custommiddleware.PrometheusMiddleware)
	router.Use(custommiddleware.CustomLogger(baseLog, cfg.HTTPTimeout))

	// router.Get("/swagger/*", httpSwagger.WrapHandler)
	// router.Handle("/metrics", promhttp.Handler())

	// GET http://localhost:8081/order/<order_uid>
	router.Get("/order/{order_uid}", mainhandler.New(baseLog, uc))

	server := &http.Server{
		Addr:         cfg.HTTPPort,
		Handler:      router,
		ReadTimeout:  cfg.HTTPTimeout,
		WriteTimeout: cfg.HTTPTimeout,
		IdleTimeout:  cfg.HTTPIdleTimeout,
	}

	return server
}
