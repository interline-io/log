package log

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/interline-io/log/internal/middleware"
)

// Copy of chi request id middleware
func RequestIDMiddleware(next http.Handler) http.Handler {
	return middleware.RequestID(next)
}

// Re-export GetReqID for use outside the internal package
func GetReqID(ctx context.Context) string {
	return middleware.GetReqID(ctx)
}

// Glue between chi RequestID and Zerolog
func RequestIDLoggingMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		rlog := For(ctx).
			With().
			Str("request_id", middleware.GetReqID(ctx)).Logger()
		next.ServeHTTP(w, r.WithContext(WithLogger(ctx, rlog)))
	}
	return http.HandlerFunc(fn)
}

// Log request and duration
// Renamed from LoggingMiddleware
func DurationLoggingMiddleware(longQueryDuration int, getUserName func(context.Context) string) func(http.Handler) http.Handler {
	return LoggingMiddleware(longQueryDuration, getUserName)
}

// Log request and duration
func LoggingMiddleware(longQueryDuration int, getUserName func(context.Context) string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get context logger
			ctx := r.Context()
			rlog := For(ctx).With()

			// Add username
			if getUserName != nil {
				if u := getUserName(ctx); u != "" {
					rlog = rlog.Str("user", getUserName(ctx))
				}
			}

			// Set as context logger
			ctx = WithLogger(ctx, rlog.Logger())
			r = r.WithContext(ctx)

			// Get request body for logging if request is json and length under 20kb
			t1 := time.Now()
			var body []byte
			if r.Header.Get("content-type") == "application/json" && r.ContentLength < 1024*20 {
				body, _ = io.ReadAll(r.Body)
				r.Body = io.NopCloser(bytes.NewBuffer(body))
			}

			// Wrap context to get error code and errors
			wr := wrapResponseWriter(w)
			next.ServeHTTP(wr, r)

			// Extra logging of request body if duration > 1s
			durationMs := (time.Now().UnixNano() - t1.UnixNano()) / 1e6

			// Get logger from context again in case it was modified
			msg := For(ctx).Info().
				Int64("duration_ms", durationMs).
				Str("method", r.Method).
				Str("path", r.URL.EscapedPath()).
				Str("query", r.URL.Query().Encode()).
				Int("status", wr.status)

			// Add duration info
			if durationMs > int64(longQueryDuration) {
				// Verify it's valid json
				msg = msg.Bool("long_query", true)
				var x interface{}
				if err := json.Unmarshal(body, &x); err == nil {
					msg = msg.RawJSON("body", body)
				}
			}
			msg.Msg("request")
		})
	}
}

// https://blog.questionable.services/article/guide-logging-middleware-go/
// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured for logging.
type responseWriter struct {
	status      int
	wroteHeader bool
	http.ResponseWriter
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if !rw.wroteHeader {
		rw.status = code
		rw.wroteHeader = true
	}
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(response []byte) (int, error) {
	if !rw.wroteHeader {
		rw.status = http.StatusOK
		rw.wroteHeader = true
	}
	return rw.ResponseWriter.Write(response)
}
