package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/ashirt-ops/ashirt-server/backend/logging"
)

type responseWriterWrapper struct {
	http.ResponseWriter
	size   int
	status int
}

func (w *responseWriterWrapper) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWriterWrapper) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.size += n
	return n, err
}

func LogRequests(baseLogger *slog.Logger) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ctx, logger := logging.AddRequestLogger(r.Context(), baseLogger)
			logger.InfoContext(ctx, "Incoming request", "method", r.Method, "url", r.URL, "from", r.RemoteAddr)
			ww := &responseWriterWrapper{w, 0, 200}

			next.ServeHTTP(ww, r.WithContext(ctx))
			logger.InfoContext(ctx, "Request Completed", "status", ww.status, "sizeInBytes", ww.size, "duration", time.Since(start))
		})
	}
}
