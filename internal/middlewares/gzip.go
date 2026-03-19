package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func GzipMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			defer gz.Close()
			r.Body = gz
		}
		acceptEncoding := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
		if acceptEncoding {
			gz, err := gzip.NewWriterLevel(w, gzip.HuffmanOnly)
			if err != nil {
				h.ServeHTTP(w, r)
				return
			}
			defer gz.Close()

			w = gzipWriter{
				ResponseWriter: w,
				Writer:         gz,
			}
			w.Header().Set("Content-Encoding", "gzip")
		}
		h.ServeHTTP(w, r)
	})
}
