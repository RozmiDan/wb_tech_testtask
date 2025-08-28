package drophandler

import (
	"net/http"
	"os"

	"github.com/RozmiDan/wb_tech_testtask/internal/entity"
	"go.uber.org/zap"
)

func New(log *zap.Logger) http.HandlerFunc {
	baselog := log.With(zap.String("handler", "DropHandler"))

	return func(w http.ResponseWriter, r *http.Request) {
		// 1) забираем request_id
		ctx := r.Context()
		logger := baselog
		// 2) оборачиваем логгер
		if reqID, ok := ctx.Value(entity.RequestIDKey{}).(string); ok && reqID != "" {
			logger = logger.With(zap.String("request_id", reqID))
		}

		w.WriteHeader(http.StatusOK)
		logger.Error("Fatal error imitation")

		os.Exit(1)
	}
}
