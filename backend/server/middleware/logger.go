package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/ashirt-ops/ashirt-server/backend/logging"
)

// InjectLogger is a thin middleware that injects a request-scoped logger into the context
// via logging.AddRequestLogger. Use this alongside an external HTTP logging middleware (e.g.
// weby's Logger) so that handlers and services can retrieve the logger via logging.ReqLogger(ctx).
func InjectLogger(baseLogger *slog.Logger) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, _ := logging.AddRequestLogger(r.Context(), baseLogger)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

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
			logger.Info("Incoming request", "method", r.Method, "url", r.URL, "from", r.RemoteAddr)
			ww := &responseWriterWrapper{w, 0, 200}

			next.ServeHTTP(ww, r.WithContext(ctx))
			logger.Info("Request Completed", "status", ww.status, "sizeInBytes", ww.size, "duration", time.Since(start))
		})
	}
}
