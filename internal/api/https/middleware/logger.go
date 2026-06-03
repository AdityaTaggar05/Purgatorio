package middleware

import (
	"io"
	"log/slog"
	"net/http"
	"time"
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
	statusCode int
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

			// Intercepting Requests
			requestReader := &customReadCloser{ReadCloser: r.Body}
			r.Body = requestReader

			// Intercepting Responses
			responseReader := &customResponseWriter{ResponseWriter: w}

			next.ServeHTTP(responseReader, r)

			attrs := []any{
				slog.Duration("duration", time.Since(start)),
				slog.Int("request_body_bytes", requestReader.bytesRead),
				slog.Int("response_body_bytes", responseReader.bytesWritten),
				slog.Int("response_status", responseReader.statusCode),
			}

			logger.Info("Served request", attrs...)
		})
	}
}

