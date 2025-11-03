package gzipper

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"slices"
	"strings"

	"go.uber.org/zap"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func Compressor(appLogger *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			appLogger.Info("Compressor started")

			acceptEncoding := r.Header.Values("Accept-Encoding")
			var isGzip bool
			for _, val := range acceptEncoding {
				if strings.Contains(val, "gzip") {
					isGzip = true
					break
				}
			}

			if !isGzip {
				appLogger.Info("Compressor stopped - GZIP not accepted")
				next.ServeHTTP(w, r)
				return
			}

			buf := &bytes.Buffer{}
			next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: buf}, r)

			contentType := w.Header().Get("Content-Type")
			isJSON := strings.Contains(contentType, "application/json")
			isHTML := strings.Contains(contentType, "text/html")

			if !isJSON && !isHTML {
				appLogger.Info("Compressor stopped - Content-Type not for compression")
				w.Write(buf.Bytes())
				return
			}

			w.Header().Set("Content-Encoding", "gzip")

			gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
			if err != nil {
				appLogger.Errorw("Compressor stopped - compression failed", "error", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			defer gz.Close()

			_, err = gz.Write(buf.Bytes())
			if err != nil {
				appLogger.Errorw("Compressor stopped - failed gz.Write", "error", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		})
	}
}

func Decompressor(appLogger *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			appLogger.Infow("Decompressor started")

			contentEncoding := r.Header.Values("Content-Encoding")
			isGzip := slices.Contains(contentEncoding, "gzip")
			if isGzip {

				appLogger.Info("Decompressor started GZIP decompression")

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
}
