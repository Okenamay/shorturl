package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

var sugar *zap.SugaredLogger

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// Middleware для логирования запросов и ответов
func WithLogging(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}

		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		h.ServeHTTP(&lw, r)

		duration := time.Since(start)
		sugar.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"status", responseData.status,
			"duration", duration,
			"size", responseData.size,
		)
	}

	return http.HandlerFunc(logFn)
}

func InitLogger() (*zap.SugaredLogger, error) {
	if sugar != nil {
		return sugar, nil
	}

	logger, err := zap.NewProduction()
	if err != nil {
		// А можно ошибку передать наверх, в main, и там её обработать?
		return nil, err
		// panic(err)
	}
	defer logger.Sync()

	sugar = logger.Sugar()
	return sugar, nil
}
