package gzipper

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"slices"
	"strings"

	logger "github.com/Okenamay/shorturl.git/internal/logger/zap"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func Compressor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sugar, _ := logger.InitLogger()

		acceptEncoding := r.Header.Values("Accept-Encoding")
		var isGzip bool
		for _, val := range acceptEncoding {
			if strings.Contains(val, "gzip") {
				isGzip = true
				break
			}
		}

		if !isGzip {
			sugar.Info("Compressor. GZIP not accepted")
			next.ServeHTTP(w, r)
			return
		}

		buf := &bytes.Buffer{}
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: buf}, r)

		contentType := w.Header().Get("Content-Type")
		isJSON := strings.Contains(contentType, "application/json")
		isHTML := strings.Contains(contentType, "text/html")

		if !isJSON && !isHTML {
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

func Decompressor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sugar, _ := logger.InitLogger()

		contentEncoding := r.Header.Values("Content-Encoding")
		isGzip := slices.Contains(contentEncoding, "gzip")
		if isGzip {
			var reader io.ReadCloser
			sugar.Info("Decompressor. Starting GZIP decompression")

			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			defer gz.Close()

			reader = gz
			r.Body = reader
		}

		next.ServeHTTP(w, r)
	})
}
