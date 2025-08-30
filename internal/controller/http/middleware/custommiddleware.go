package custommiddleware

import (
	"context"
	"net/http"
	"time"

	"github.com/RozmiDan/wb_tech_testtask/internal/entity"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func CustomLogger(log *zap.Logger, httpTimeout time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		baselog := log.With(zap.String("component", "middleware/logger"))
		baselog.Info("logger middleware enabled")

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := middleware.GetReqID(r.Context())
			curLog := baselog.With(
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("request_id", reqID),
			)
			// ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			ctx := context.WithValue(r.Context(), entity.RequestIDKey{}, reqID)
			ctx, cancel := context.WithTimeout(ctx, httpTimeout)
			t1 := time.Now()
			defer func() {
				curLog.Info("request completed",
					// zap.Int("status", ww.Status()),
					zap.Duration("request time", time.Since(t1)),
				)
				cancel()
			}()
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// func PrometheusMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		//path := r.URL.Path
// 		method := r.Method

// 		rw := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

// 		prom_metrics.HTTPInFlight.WithLabelValues(method).Inc()
// 		defer prom_metrics.HTTPInFlight.WithLabelValues(method).Dec()

// 		histTimer := prometheus.NewTimer(prom_metrics.HTTPDuration.WithLabelValues(method))
// 		defer func() {
// 			histTimer.ObserveDuration()
// 		}()

// 		// timer := prometheus.NewTimer(prom_metrics.HTTPDuration.WithLabelValues(method, path))
// 		// defer timer.ObserveDuration()

// 		next.ServeHTTP(rw, r)

// 		status := strconv.Itoa(rw.Status())
// 		prom_metrics.HTTPRequests.WithLabelValues(method, status).Inc()
// 	})
// }
