package mainhandler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/RozmiDan/wb_tech_testtask/internal/entity"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// 1) GET order/<order_uid>

type OrderInfoGetter interface {
	GetOrderInfo(ctx context.Context, orderUID string) (*entity.OrderResponse, error)
}

// Get Order by UID
// @Summary      Get order by UID
// @Description  Возвращает информацию о заказе по order_uid.
// @Tags         orders
// @Param        order_uid   path      string  true  "Order UID"
// @Success      200  {object}  entity.OrderResponse
// @Failure      400  {object}  APIError  "invalid order_uid"
// @Failure      404  {object}  APIError  "order not found"
// @Failure      504  {object}  APIError  "timeout exceeded"
// @Failure      500  {object}  APIError  "unexpected internal error"
// @Router       /order/{order_uid} [get]
func New(log *zap.Logger, uc OrderInfoGetter) http.HandlerFunc {
	baselog := log.With(zap.String("handler", "MainHandler"))

	return func(w http.ResponseWriter, r *http.Request) {
		// 1) забираем request_id
		ctx := r.Context()
		logger := baselog

		// 2) оборачиваем логгер
		if reqID, ok := ctx.Value(entity.RequestIDKey{}).(string); ok && reqID != "" {
			logger = logger.With(zap.String("request_id", reqID))
		}

		// 3) достаем UID из URL
		orderUID := chi.URLParam(r, "order_uid")
		if err := uidParser(orderUID); err != nil {
			logger.Warn("invalid order_uid", zap.String("order_uid", orderUID))
			errDTO := APIError{
				Message: err.Error(),
			}
			http.Error(w, errDTO.Message, http.StatusBadRequest)

			return
		}

		// 4) вызываем usecase
		order, err := uc.GetOrderInfo(ctx, orderUID)
		if err != nil {
			switch {
			case errors.Is(ctx.Err(), context.DeadlineExceeded):
				logger.Error("timeout exceeded", zap.Error(err))
				errDTO := APIError{
					Message: "request took longer than the timelimit",
				}
				http.Error(w, errDTO.Message, http.StatusGatewayTimeout)

				return
			case errors.Is(err, entity.ErrorOrderNotFound):
				logger.Info("order not found", zap.String("order_uid", orderUID))
				errDTO := APIError{
					Message: "order not found",
				}
				http.Error(w, errDTO.Message, http.StatusNotFound)

				return
			}
			logger.Error("failed to get order", zap.Error(err))
			errDTO := APIError{
				Message: "unexpected internal error",
			}
			http.Error(w, errDTO.Message, http.StatusInternalServerError)

			return
		}

		// 5) формируем успешный ответ
		b, err := json.MarshalIndent(order, "", "	")
		if err != nil {
			logger.Error("error marshal response")
		}

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(b); err != nil {
			logger.Error("error sending the response")

			return
		}
	}
}

func uidParser(uid string) error {
	if len(uid) != 0 && uid == strings.ToLower(uid) {
		return nil
	}

	return errors.New("order_uid is not a valid UID")
}
