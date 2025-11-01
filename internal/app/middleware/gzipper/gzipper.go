package gzipper

import (
	// "bytes" // Больше не нужен
	"compress/gzip"
	"net/http"
	"slices"
	"strings"
	"sync"

	"go.uber.org/zap"
)

// gzipWriter теперь потоковый
type gzipWriter struct {
	http.ResponseWriter
	gw *gzip.Writer
}

func (w *gzipWriter) Write(b []byte) (int, error) {
	// Пишем сразу в gzip.Writer, который пишет в http.ResponseWriter
	return w.gw.Write(b)
}

// Close закрывает gzip.Writer и возвращает его в пул
func (w *gzipWriter) Close() error {
	err := w.gw.Close()
	gzipWriterPool.Put(w.gw)
	return err
}

// Пул для gzip.Writer
var gzipWriterPool = sync.Pool{
	New: func() interface{} {
		// Создаем writer с уровнем сжатия по умолчанию.
		// Мы вызовем Reset() для него перед использованием.
		w, _ := gzip.NewWriterLevel(nil, gzip.BestSpeed)
		return w
	},
}

// --- ИСПРАВЛЕНИЕ 1: Пул для gzip.Reader ---
var gzipReaderPool = sync.Pool{
	New: func() interface{} {
		// ОШИБКА БЫЛА ЗДЕСЬ:
		// gzip.NewReader(nil) возвращает ридер, который паникует при Reset().
		// Правильный способ - вернуть пустую структуру.
		return new(gzip.Reader)
	},
}

func Compressor(appLogger *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			appLogger.Info("Compressor started")

			acceptEncoding := r.Header.Values("Accept-Encoding")
			isGzip := false
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

			// Потоковое сжатие
			// Получаем gzip.Writer из пула
			gz := gzipWriterPool.Get().(*gzip.Writer)
			// Устанавливаем ему реальный ResponseWriter
			gz.Reset(w)

			// Создаем наш кастомный writer
			gw := &gzipWriter{
				ResponseWriter: w,
				gw:             gz,
			}
			// Важно: Close() вернет writer в пул
			defer gw.Close()

			w.Header().Set("Content-Encoding", "gzip")
			next.ServeHTTP(gw, r)
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

				// --- ИСПРАВЛЕНИЕ 2: Логика Decompressor ---
				// 1. Получаем *gzip.Reader из пула
				gz := gzipReaderPool.Get().(*gzip.Reader)

				// 2. Вызываем Reset() с телом запроса
				if err := gz.Reset(r.Body); err != nil {
					appLogger.Errorw("Decompressor stopped - gzip.Reader.Reset FAIL", "error", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					gzipReaderPool.Put(gz) // Не забываем вернуть в пул при ошибке
					return
				}

				// 3. Устанавливаем кастомный ReadCloser
				// Он вернет ридер в пул, когда next.ServeHTTP вызовет r.Body.Close()
				r.Body = &gzipReadCloser{
					Reader: gz,
					pool:   &gzipReaderPool,
				}
				// --- КОНЕЦ ИСПРАВЛЕНИЯ ---
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Кастомный ReadCloser для возврата Reader'а в пул
type gzipReadCloser struct {
	*gzip.Reader
	pool *sync.Pool
}

func (grc *gzipReadCloser) Close() error {
	// Сначала закрываем gzip.Reader
	err := grc.Reader.Close()
	// Затем возвращаем его в пул
	grc.pool.Put(grc.Reader)
	return err
}
