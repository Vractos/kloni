package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Vractos/kloni/pkg/metrics"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type CustomLogger struct {
	logger *metrics.Logger
}

func NewStructuredLogger(handler *metrics.Logger) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&CustomLogger{logger: handler})
}

func (cl *CustomLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &customLogEntry{logger: cl.logger}
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	requestURL := fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)
	msg := fmt.Sprintf("Request: %s %s", r.Method, r.URL.Path)
	logFields := []zap.Field{
		zap.String("url", requestURL),
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.String("ip", r.RemoteAddr),
		zap.String("proto", r.Proto),
		zap.String("agent", r.UserAgent()),
	}

	entry.logger.Info(msg, logFields...)
	return entry
}

type customLogEntry struct {
	logger *metrics.Logger
}

// Write implements middleware.LogEntry
func (cle *customLogEntry) Write(status int, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	msg := fmt.Sprintf("Response: %d %s", status, statusLabel(status))

	logFields := []zap.Field{
		zap.Int("status", status),
		zap.Int("bytes", bytes),
		zap.Float64("elapsed", float64(elapsed.Nanoseconds())/1000000.0),
	}

	switch {
	case status <= 0:
		cle.logger.Warn(msg, logFields...)
	case status < 400:
		cle.logger.Info(msg, logFields...)
	case status >= 400 && status < 500:
		cle.logger.Warn(msg, logFields...)
	case status >= 500:
		cle.logger.Error(msg, errors.New("error to handle request"), logFields...)
	default:
		cle.logger.Info(msg, logFields...)
	}
}

// Panic implements middleware.LogEntry
func (cle *customLogEntry) Panic(v interface{}, stack []byte) {
	logFields := []zap.Field{
		zap.String("stack", string(stack)),
		zap.String("panic", fmt.Sprintf("%+v", v)),
	}

	msg := fmt.Sprintf("%+v", v)

	cle.logger.Error(msg, errors.New("panic to handle request"), logFields...)
	middleware.PrintPrettyStack(v)
}

func statusLabel(status int) string {
	switch {
	case status >= 100 && status < 300:
		return "OK"
	case status >= 300 && status < 400:
		return "Redirect"
	case status >= 400 && status < 500:
		return "Client Error"
	case status >= 500:
		return "Server Error"
	default:
		return "Unknown"
	}
}
