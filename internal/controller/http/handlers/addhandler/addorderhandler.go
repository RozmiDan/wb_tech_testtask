package addhandler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/RozmiDan/wb_tech_testtask/internal/entity"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type OrderInfoPoster interface {
	AddOrderInfo(ctx context.Context, order *entity.OrderInfo) error
}

func New(log *zap.Logger, uc OrderInfoPoster) http.HandlerFunc {
	baselog := log.With(zap.String("handler", "PingHandler"))

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
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1MB
		defer func() {
			if err := r.Body.Close(); err != nil {
				logger.Error("failed to close body:", zap.Error(err))
			}
		}()

		// достаем JSON
		var order *entity.OrderInfo
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()

		if err := dec.Decode(&order); err != nil {
			if errors.Is(err, io.EOF) {
				http.Error(w, "empty body", http.StatusBadRequest)

				return
			}
			http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)

			return
		}

		// запретим «лишние» данные после валидного JSON
		if dec.More() {
			http.Error(w, "unexpected data after JSON object", http.StatusBadRequest)

			return
		}

		// 5) order_uid в body должен совпасть с path
		if order.OrderUID != orderUID {
			http.Error(w, "order_uid in path and body must match", http.StatusBadRequest)

			return
		}

		// 4) вызываем usecase
		err := uc.AddOrderInfo(ctx, order)
		if err != nil {
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				logger.Error("timeout exceeded", zap.Error(err))
				http.Error(w, "request took longer than the timelimit", http.StatusGatewayTimeout)
				return
			}

			logger.Error("failed to add order", zap.Error(err))
			http.Error(w, "unexpected internal error", http.StatusInternalServerError)

			return
		}

		// 5) формируем успешный ответ
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(nil); err != nil {
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
