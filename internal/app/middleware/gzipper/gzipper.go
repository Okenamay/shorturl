package gzipper

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	logger "github.com/Okenamay/shorturl.git/internal/logger/zap"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

// Реализуем все методы интерфейса ResponseWriter
func (w *gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (w *gzipWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

// Compressor middleware
func Compressor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sugar, _ := logger.InitLogger()

		acceptEncoding := r.Header.Get("Accept-Encoding")
		if !strings.Contains(acceptEncoding, "gzip") {
			sugar.Info("Compressor. GZIP not accepted")
			next.ServeHTTP(w, r)
			return
		}

		buf := &bytes.Buffer{}
		next.ServeHTTP(&gzipWriter{ResponseWriter: w, Writer: buf}, r)

		contentType := w.Header().Get("Content-Type")
		if !strings.Contains(contentType, "application/json") &&
			!strings.Contains(contentType, "text/html") {
			sugar.Info("Compressor. Content-Type not for compression")
			w.Write(buf.Bytes())
			return
		}

		w.Header().Set("Content-Encoding", "gzip")

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			sugar.Infow("Compressor. Compression", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer gz.Close()

		_, err = gz.Write(buf.Bytes())
		if err != nil {
			sugar.Infow("Compressor. gz.Write", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

// Decompressor middleware
func Decompressor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sugar, _ := logger.InitLogger()

		contentEncoding := r.Header.Get("Content-Encoding")
		if contentEncoding == "gzip" {
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
