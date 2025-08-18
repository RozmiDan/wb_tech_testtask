package server

import (
	"net/http"

	"github.com/RozmiDan/wb_tech_testtask/internal/config"
	"github.com/RozmiDan/wb_tech_testtask/internal/controller/handlers/mainhandler"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

type UseCase interface {
}

func InitServer(cfg *config.Config, logger *zap.Logger, uc UseCase) *http.Server {
	logger = logger.With(zap.String("layer", "Controller"))

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middleware_metrics.PrometheusMiddleware)
	router.Use(middleware_logger.MyLogger(logger))

	router.Get("/swagger/*", httpSwagger.WrapHandler)
	router.Handle("/metrics", promhttp.Handler())

	// GET http://localhost:8081/order/<order_uid>
	router.Route("/order/{game_id}", mainhandler.)

	server := &http.Server{
		Addr:         cfg.HTTPPort,
		Handler:      router,
		ReadTimeout:  cfg.HTTPTimeout,
		WriteTimeout: cfg.HTTPTimeout,
		IdleTimeout:  cfg.HTTPIdleTimeout,
	}

	return server
}
