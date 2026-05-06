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

// gzipWriter - кастомный http.ResponseWriter, который сжимает данные "на лету"
// (потоково) и записывает их в нижележащий ResponseWriter.
type gzipWriter struct {
	http.ResponseWriter
	gw *gzip.Writer
}

// Write реализует интерфейс http.ResponseWriter, сжимая данные.
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

// Пул для gzip.Writer для переиспользования.
var gzipWriterPool = sync.Pool{
	New: func() any {
		// Создаем writer с уровнем сжатия по умолчанию.
		// Мы вызовем Reset() для него перед использованием.
		w, _ := gzip.NewWriterLevel(nil, gzip.BestSpeed)
		return w
	},
}

// Пул для gzip.Reader для переиспользования.
var gzipReaderPool = sync.Pool{
	New: func() any {
		// NewReader(nil) возвращает ридер, который паникует при Reset().
		// Правильный способ - вернуть пустую структуру.
		return new(gzip.Reader)
	},
}

// Compressor - это middleware, которое сжимает тело ответа (response body)
// методом gzip, если клиент отправил заголовок "Accept-Encoding: gzip".
// Сжатие происходит "на лету" (потоково) без буферизации всего ответа.
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

// Decompressor - это middleware, которое распаковывает тело запроса (request
// body), если оно было сжато методом gzip (т.е. содержит заголовок
// "Content-Encoding: gzip").
//
// Он заменяет r.Body на специальный io.ReadCloser, который читает
// распакованные данные и возвращает *gzip.Reader в пул при вызове Close().
func Decompressor(appLogger *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			appLogger.Infow("Decompressor started")

			contentEncoding := r.Header.Values("Content-Encoding")
			isGzip := slices.Contains(contentEncoding, "gzip")
			if isGzip {
				appLogger.Info("Decompressor started GZIP decompression")

				gz := gzipReaderPool.Get().(*gzip.Reader)

				if err := gz.Reset(r.Body); err != nil {
					appLogger.Errorw("Decompressor stopped - gzip.Reader.Reset FAIL", "error", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					gzipReaderPool.Put(gz)
					return
				}

				r.Body = &gzipReadCloser{
					Reader: gz,
					pool:   &gzipReaderPool,
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

// gzipReadCloser - это обертка над *gzip.Reader, которая
// реализует io.ReadCloser.
// Его задача - вернуть *gzip.Reader обратно в sync.Pool при вызове Close().
type gzipReadCloser struct {
	*gzip.Reader
	pool *sync.Pool
}

// Close возвращает gzip.Reader в пул, когда тело запроса закрывается.
func (grc *gzipReadCloser) Close() error {
	// Сначала закрываем gzip.Reader
	err := grc.Reader.Close()
	// Затем возвращаем его в пул
	grc.pool.Put(grc.Reader)
	return err
}
