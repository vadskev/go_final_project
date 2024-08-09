package logger

import (
	"github.com/vadskev/go_final_project/internal/logger"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func New() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			timeStart := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			logger.Debug("got incoming HTTP request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", ww.Status()),
				zap.Int("length", ww.BytesWritten()),
				zap.Duration("time", time.Since(timeStart)),
			)

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
