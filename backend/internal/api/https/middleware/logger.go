package middleware

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/AdityaTaggar05/Purgatorio/pkg/ctxkeys"
	"github.com/AdityaTaggar05/Purgatorio/pkg/response"
)

// Extending the default request reader interface to add custom fields
type customReadCloser struct {
	io.ReadCloser
	bytesRead int
}

func (r *customReadCloser) Read(p []byte) (int, error) {
	n, err := r.ReadCloser.Read(p)
	r.bytesRead += n
	return n, err
}

// Extending the default response reader interface to add custom fields
type customResponseWriter struct {
	http.ResponseWriter
	bytesWritten int
	statusCode   int
}

func (w *customResponseWriter) Write(p []byte) (int, error) {
	// If only w.Write is called, it by default means StatusOK
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
	n, err := w.ResponseWriter.Write(p)
	w.bytesWritten += n
	return n, err
}

func (w *customResponseWriter) WriteHeader(status int) {
	w.statusCode = status
	w.ResponseWriter.WriteHeader(status)
}

func RequestLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			logCtx := &response.LogContext{}

			// Intercepting Requests
			requestReader := &customReadCloser{ReadCloser: r.Body}
			r.Body = requestReader

			// Intercepting Responses
			responseReader := &customResponseWriter{ResponseWriter: w}

			next.ServeHTTP(responseReader, r.WithContext(context.WithValue(r.Context(), ctxkeys.Log, logCtx)))

			attrs := []any{
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Duration("duration", time.Since(start)),
				slog.Int("request_body_bytes", requestReader.bytesRead),
				slog.Int("response_body_bytes", responseReader.bytesWritten),
				slog.Int("response_status", responseReader.statusCode),
			}

			if logCtx.UserID != "" {
				attrs = append(attrs, slog.String("user_id", logCtx.UserID))
			}

			if logCtx.Error != nil {
				attrs = append(attrs, slog.Any("error", logCtx.Error))
			}

			logLevel := slog.LevelInfo
			if responseReader.statusCode >= 500 {
				logLevel = slog.LevelError
			} else if responseReader.statusCode >= 400 {
				logLevel = slog.LevelWarn
			}
			logger.Log(r.Context(), logLevel, "Served request", attrs...)
		})
	}
}
