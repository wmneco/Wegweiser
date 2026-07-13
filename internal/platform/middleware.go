package platform

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"time"
)

// RequestIDHeader is the header used both to read an inbound correlation ID
// and to echo the request ID back on the response.
const RequestIDHeader = "X-Request-ID"

type contextKey int

const requestIDKey contextKey = iota

// RequestIDFromContext returns the request ID assigned by the RequestID
// middleware, or "" if none is present.
func RequestIDFromContext(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey).(string)
	return id
}

// requestID wraps next so every request carries a request ID: the inbound
// RequestIDHeader value if present, otherwise a freshly generated one. The ID
// is echoed on the response header and stored in the request context.
func requestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(RequestIDHeader)
		if id == "" {
			id = generateRequestID()
		}
		w.Header().Set(RequestIDHeader, id)
		ctx := context.WithValue(r.Context(), requestIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func generateRequestID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "unknown"
	}
	return hex.EncodeToString(b)
}

// recoverPanic wraps next so a panic in a handler is recovered, logged with
// the request ID, and turned into a 500 response in the standard error shape
// instead of crashing the connection.
func recoverPanic(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.Error("panic recovered",
					"error", rec,
					"request_id", RequestIDFromContext(r.Context()),
				)
				WriteError(w, http.StatusInternalServerError, "internal_error", "internal server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// statusRecorder captures the status code written to an http.ResponseWriter
// so requestLogger can log it after the handler completes.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

// requestLogger wraps next so every completed request emits a structured log
// entry with method, path, status, duration, and request ID.
func requestLogger(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		logger.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
			"duration", time.Since(start),
			"request_id", RequestIDFromContext(r.Context()),
		)
	})
}

// Middleware wraps next with the baseline stack applied to every request:
// request ID assignment, panic recovery, and request logging.
func Middleware(logger *slog.Logger, next http.Handler) http.Handler {
	return requestID(requestLogger(logger, recoverPanic(logger, next)))
}
