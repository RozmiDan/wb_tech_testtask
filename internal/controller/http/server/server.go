package server

import (
	"context"
	"net/http"

	_ "github.com/RozmiDan/wb_tech_testtask/docs"
	"github.com/RozmiDan/wb_tech_testtask/internal/config"
	"github.com/RozmiDan/wb_tech_testtask/internal/controller/http/handlers/drophandler"
	"github.com/RozmiDan/wb_tech_testtask/internal/controller/http/handlers/mainhandler"
	"github.com/RozmiDan/wb_tech_testtask/internal/controller/http/handlers/pinghandler"
	custommiddleware "github.com/RozmiDan/wb_tech_testtask/internal/controller/http/middleware"
	"github.com/RozmiDan/wb_tech_testtask/internal/controller/http/webui"
	"github.com/RozmiDan/wb_tech_testtask/internal/entity"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

type UseCase interface {
	GetOrderInfo(ctx context.Context, orderUID string) (*entity.OrderResponse, error)
	AddOrderInfo(ctx context.Context, order *entity.OrderInfo) error
}

func InitServer(cfg *config.Config, logger *zap.Logger, uc UseCase) *http.Server {
	baseLog := logger.With(zap.String("layer", "Controller"))

	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.URLFormat)
	// router.Use(custommiddleware.PrometheusMiddleware)
	router.Use(custommiddleware.CustomLogger(baseLog, cfg.HTTPTimeout))

	router.Get("/swagger/*", httpSwagger.WrapHandler)
	// router.Handle("/metrics", promhttp.Handler())

	// static UI
	router.Get("/", webui.Index())
	router.Handle("/static/*", webui.Static())

	// GET http://localhost:8081/order/<order_uid>
	router.Get("/order/{order_uid}", mainhandler.New(baseLog, uc))
	router.Get("/service/drop", drophandler.New(baseLog))
	router.Post("/order/{order_uid}", pinghandler.New(baseLog, uc))

	server := &http.Server{
		Addr:         cfg.HTTPPort,
		Handler:      router,
		ReadTimeout:  cfg.HTTPTimeout,
		WriteTimeout: cfg.HTTPTimeout,
		IdleTimeout:  cfg.HTTPIdleTimeout,
	}

	return server
}
