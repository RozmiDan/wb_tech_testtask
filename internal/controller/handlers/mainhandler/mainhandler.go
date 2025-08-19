package mainhandler

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/RozmiDan/wb_tech_testtask/internal/entity"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/v5/render"
	"go.uber.org/zap"
)

// 1) GET order/<order_uid>

type OrderInfoGetter interface {
	GetOrderInfo(ctx context.Context, orderUID string) (entity.OrderResponse, error)
}

func New(log *zap.Logger, uc OrderInfoGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1) забираем request_id
		ctx := r.Context()

		logger := log.With(zap.String("handler", "MainHandler"))

		// 2) оборачиваем логгер
		if reqID, ok := ctx.Value(entity.RequestIDKey{}).(string); ok && reqID != "" {
			logger = logger.With(zap.String("request_id", reqID))
		}

		// 3) достаем UID из URL
		orderUID := chi.URLParam(r, "order_uid")
		if ok := uidParser(orderUID); !ok {
			logger.Warn("invalid order_uid", zap.String("order_uid", orderUID))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, APIError{"invalid_order_uid", "order_id is not a valid UID"})
			return
		}

		// 4) вызываем usecase
		order, err := uc.GetOrderInfo(ctx, orderUID)
		if err != nil {
			switch {
			case errors.Is(ctx.Err(), context.DeadlineExceeded):
				logger.Error("timeout exceeded", zap.Error(err))
				render.Status(r, http.StatusGatewayTimeout)
				render.JSON(w, r, APIError{"timeout_exceeded", "request took longer than timelimit"})
				return

				// case errors.Is(err, entity.ErrOrderNotFound):
				// 	logger.Info("order not found", zap.String("order_uid", orderUID))
				// 	render.Status(r, http.StatusNotFound)
				// 	render.JSON(w, r, APIError{"not_found", "game not found"})
				// 	return
			}
			logger.Error("failed to get order", zap.Error(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, APIError{"internal error", "could not find order"})
			return
		}

		// 5) формируем успешный ответ
		resp := GetOrderResponse{
			Order: order,
		}
		render.Status(r, http.StatusOK)
		render.JSON(w, r, resp)
	}
}

func uidParser(uid string) bool {
	if len(uid) != 0 && uid == strings.ToLower(uid) {
		return true
	}
	return false
}

// if err := uc.PostRating(ctx, gameID, payload.UserID, payload.Rating); err != nil {
// 			switch {
// 			case errors.Is(err, entity.ErrBrokerUnavailable):
// 				logger.Error("broker unavailable", zap.Error(err))
// 				render.Status(r, http.StatusServiceUnavailable)
// 				render.JSON(w, r, ErrorResponse{
// 					Error: APIError{"service_unavailable", "unable to publish rating"},
// 				})
// 				return

// 			case errors.Is(err, entity.ErrGameNotFound):
// 				logger.Info("game not found", zap.String("game_id", gameID))
// 				render.Status(r, http.StatusNotFound)
// 				render.JSON(w, r, ErrorResponse{
// 					Error: APIError{"not_found", "game not found"},
// 				})
// 				return

// 			case errors.Is(err, context.DeadlineExceeded):
// 				logger.Error("timeout exceeded", zap.Error(err))
// 				render.Status(r, http.StatusGatewayTimeout)
// 				render.JSON(w, r, ErrorResponse{
// 					Error: APIError{"timeout_exceeded", "request took longer than 2 seconds"},
// 				})
// 				return
// 			}

// 			logger.Error("failed to post rating", zap.Error(err))
// 			render.Status(r, http.StatusInternalServerError)
// 			render.JSON(w, r, ErrorResponse{
// 				Error: APIError{"internal error", "could not submit rating"},
// 			})
// 			return
// 		}

// 		// 7) Успех — пустой ответ 200 OK
// 		render.Status(r, http.StatusOK)
// 		render.JSON(w, r, struct{}{})
// 	}
