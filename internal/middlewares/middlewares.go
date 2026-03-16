package middlewares

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type responseData struct {
	code int
	size int
}
type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (rw *loggingResponseWriter) WriteHeader(code int) {
	rw.responseData.code = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.responseData.size += size
	return size, err
}

func WriteHeader(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		logFn := func (w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			
			responseData := &responseData {
				code: 0,
				size: 0,
			}
			lw := loggingResponseWriter {
				ResponseWriter: w,
				responseData: responseData,
			}
			h.ServeHTTP(&lw, r)
			
			duration := time.Since(start)
			
			logger.Info("HTTP request",
				zap.String("uri", r.RequestURI),
				zap.String("method", r.Method),
				zap.Int("status", responseData.code),
				zap.Duration("duration", duration),
				zap.Int("size", responseData.size),
			)
		}
		return http.HandlerFunc(logFn)
	}
}