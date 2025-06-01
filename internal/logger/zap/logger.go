package logger

// import (
// 	"net/http"
// 	"time"

// 	"go.uber.org/zap"
// )

// var sugar *zap.SugaredLogger

// type (
// 	responseData struct {
// 		status int
// 		size   int
// 	}

// 	loggingResponseWriter struct {
// 		http.ResponseWriter
// 		responseData *responseData
// 	}
// )

// func (r *loggingResponseWriter) Write(b []byte) (int, error) {
// 	size, err := r.ResponseWriter.Write(b)
// 	r.responseData.size += size
// 	return size, err
// }

// func (r *loggingResponseWriter) WriteHeader(statusCode int) {
// 	r.ResponseWriter.WriteHeader(statusCode)
// 	r.responseData.status = statusCode
// }

// func WithLogging(h http.Handler) http.Handler {
// 	logFn := func(w http.ResponseWriter, r *http.Request) {
// 		start := time.Now()

// 		responseData := &responseData{
// 			status: 0,
// 			size:   0,
// 		}

// 		lw := loggingResponseWriter{
// 			ResponseWriter: w,
// 			responseData:   responseData,
// 		}

// 		h.ServeHTTP(&lw, r)

// 		duration := time.Since(start)
// 		sugar.Infoln(
// 			"uri", r.RequestURI,
// 			"method", r.Method,
// 			"status", responseData.status,
// 			"duration", duration,
// 			"size", responseData.size,
// 		)
// 	}

// 	return http.HandlerFunc(logFn)
// }

//         // эндпоинт /ping
//         uri := r.RequestURI
//         // метод запроса
//         method := r.Method

//         // точка, где выполняется хендлер pingHandler
//         h.ServeHTTP(w, r) // обслуживание оригинального запроса

//         // Since возвращает разницу во времени между start
//         // и моментом вызова Since. Таким образом можно посчитать
//         // время выполнения запроса.
//         duration := time.Since(start)

//         // отправляем сведения о запросе в zap
//         sugar.Infoln(
//             "uri", uri,
//             "method", method,
//             "duration", duration,

// func GetLogger() *zap.SugaredLogger {
// 	if sugar != nil {
// 		return sugar
// 	}

// 	logger, err := zap.NewDevelopment()
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer logger.Sync()

// 	sugar = logger.Sugar()
// 	return sugar
// }
