package logging

import (
	"context"
	"io"
	"log/slog"
	"os"

	"github.com/google/uuid"
)

var requestLoggerKey = &struct{ name string }{"requestLogger"}

var systemLogger *slog.Logger

// SetupStdoutLogging creates a new logger that logs to standard out
func SetupStdoutLogging() *slog.Logger {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	SetSystemLogger(logger)
	return logger
}

func SetSystemLogger(logger *slog.Logger) {
	systemLogger = logger
}

func GetSystemLogger() *slog.Logger {
	return systemLogger
}

// NewNopLogger creates a logger that actually does not log. useful in situations where some logger
// is needed, but for whatever reason, a real logger is missing (or if you want to conditionally)
// disable logging.
func NewNopLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

// ReqLogger retrieves a stored initially stored with AddRequestLogger. This logger is tied to the
// request (or more specifically, the context). Assuming you only have a request handy, the code
// to retrieve the logger is: logging.ReqLogger(r.Context())
func ReqLogger(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(requestLoggerKey).(*slog.Logger)
	if !ok {
		return NewNopLogger()
	}
	return logger
}

// AddRequestLogger adds a logger to this request. The logger will provide unique identification
// for any request in this stream (via the "ctx" field in the log). The logger can be retrieved
// via
func AddRequestLogger(ctx context.Context, baseLogger *slog.Logger) (context.Context, *slog.Logger) {
	requestUUID, _ := uuid.NewRandom()
	reqLogger := baseLogger.With("ctx", requestUUID.String())
	return context.WithValue(ctx, requestLoggerKey, reqLogger), reqLogger
}

// Fatal is an effective copy of go's log.Fatal, but using the logger provided, rather than
// using go's native logging. After writing the message, the code will exit with code 1
func Fatal(ctx context.Context, logger *slog.Logger, msg string, keyvals ...interface{}) {
	logger.ErrorContext(ctx, msg, keyvals...)
	os.Exit(1)
}

func LogWithoutAuth(ctx context.Context, msg string, keyvals ...interface{}) {
	if systemLogger != nil {
		systemLogger.InfoContext(ctx, msg, keyvals...)
	}
}

// SystemLog provides a system-level logger, which is not tied to any request.
// this should be used in situations where either a context is not handy, but logging is important
// or for events that are not tied to a request (e.g. losing database connection)
func SystemLog(ctx context.Context, msg string, keyvals ...interface{}) {
	if systemLogger != nil {
		systemLogger.InfoContext(ctx, msg, keyvals...)
	}
}
