package log

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// log request and duration
func LoggingMiddleware(longQueryDuration int, getUserName func(context.Context) string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Setup context logger
			// Add username
			rlog := Logger.With()
			if getUserName != nil {
				if u := getUserName(r.Context()); u != "" {
					rlog = rlog.Str("user", getUserName(r.Context()))
				}
			}
			// Add request ID
			// TODO
			rlogger := rlog.Logger()
			r = r.WithContext(WithLogger(r.Context(), rlogger))

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
			msg := rlogger.Info().
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
