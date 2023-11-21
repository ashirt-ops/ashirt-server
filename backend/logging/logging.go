package logging

import (
	"context"
	"os"

	kitlog "github.com/go-kit/kit/log"
	"github.com/google/uuid"
)

// Logger is a generic logging interface. Currently wraps Go-kit's log.Logger
type Logger kitlog.Logger

var requestLoggerKey = &struct{ name string }{"requestLogger"}

var systemLogger kitlog.Logger

// SetupStdoutLogging creates a new logger that logs to standard out
func SetupStdoutLogging() Logger {
	w := kitlog.NewSyncWriter(os.Stdout)
	logger := kitlog.NewLogfmtLogger(w)
	logger = kitlog.With(logger, "timestamp", kitlog.DefaultTimestampUTC)
	SetSystemLogger(logger)
	return logger
}

func SetSystemLogger(logger Logger) {
	systemLogger = logger
}

func GetSystemLogger() Logger {
	return systemLogger
}

func With(logger Logger, keyvals ...interface{}) Logger {
	return kitlog.With(logger, keyvals...)
}

// NewNopLogger creates a logger that actually does not log. useful in situations where some logger
// is needed, but for whatever reason, a real logger is missing (or if you want to conditionally)
// disable logging.
func NewNopLogger() Logger {
	return kitlog.NewNopLogger()
}

// ReqLogger retrieves a stored initially stored with AddRequestLogger. This logger is tied to the
// request (or more specifically, the context). Assuming you only have a request handy, the code
// to retrieve the logger is: logging.ReqLogger(r.Context())
func ReqLogger(ctx context.Context) Logger {
	logger, ok := ctx.Value(requestLoggerKey).(Logger)
	if !ok {
		return NewNopLogger()
	}
	return logger
}

// AddRequestLogger adds a logger to this request. The logger will provide unique identification
// for any request in this stream (via the "ctx" field in the log). The logger can be retrieved
// via
func AddRequestLogger(ctx context.Context, baseLogger Logger) (context.Context, Logger) {
	requestUUID, _ := uuid.NewRandom()
	reqLogger := kitlog.With(baseLogger, "ctx", requestUUID.String())
	return context.WithValue(ctx, requestLoggerKey, reqLogger), reqLogger
}

// Fatal is an effective copy of go's log.Fatal, but using the logger provided, rather than
// using go's native logging. After writing the message, the code will exit with code 1
func Fatal(logger Logger, keyvals ...interface{}) {
	logger.Log(keyvals...)
	os.Exit(1)
}

// Log provides a shorthand for the following code:
// ReqLogger(myContext).Log(/*your values here*/)
func Log(ctx context.Context, keyvals ...interface{}) error {
	return ReqLogger(ctx).Log(keyvals...)
}

func LogWithoutAuth(keyvals ...interface{}) error {
	if systemLogger != nil {
		return systemLogger.Log(keyvals...)
	}
	return nil
}

// SystemLog provides a system-level logger, which is not tied to any request.
// this should be used in situations where either a context is not handy, but logging is important
// or for events that are not tied to a request (e.g. losing database connection)
func SystemLog(keyvals ...interface{}) error {
	if systemLogger != nil {
		return systemLogger.Log(keyvals...)
	}
	return nil
}
