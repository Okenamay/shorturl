package gzipper

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"strings"

	logger "github.com/Okenamay/shorturl.git/internal/logger/zap"
)

type gzipResponseWriter struct {
	w       http.ResponseWriter
	gzip    *gzip.Writer
	buffer  *bytes.Buffer
	written bool
}

func newGzipResponseWriter(w http.ResponseWriter) *gzipResponseWriter {
	return &gzipResponseWriter{
		w:       w,
		buffer:  &bytes.Buffer{},
		gzip:    nil,
		written: false,
	}
}

func (w *gzipResponseWriter) Header() http.Header {
	return w.w.Header()
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	if w.gzip == nil {
		w.buffer.Write(b)
	} else {
		return w.gzip.Write(b)
	}
	return len(b), nil
}

func (w *gzipResponseWriter) WriteHeader(statusCode int) {
	w.w.WriteHeader(statusCode)
	w.written = true

	if strings.Contains(w.w.Header().Get("Content-Type"), "text/html") ||
		strings.Contains(w.w.Header().Get("Content-Type"), "application/json") {

		w.gzip = gzip.NewWriter(w.w)
		w.w.Header().Set("Content-Encoding", "gzip")
	}
}

// Compressor middleware
func Compressor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sugar, _ := logger.InitLogger()

		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			sugar.Info("Compressor. GZIP accepted")

			gzw := newGzipResponseWriter(w)
			next.ServeHTTP(gzw, r)

			if gzw.gzip == nil {
				w.Header().Del("Content-Encoding")
				w.Write(gzw.buffer.Bytes())
			}
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

// Decompressor middleware
func Decompressor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sugar, _ := logger.InitLogger()

		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			sugar.Info("Decompressor. Starting GZIP decompression")

			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer gz.Close()

			r.Body = gz
		}

		next.ServeHTTP(w, r)
	})
}
