package logger

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestInitLogger(t *testing.T) {
	appLogger, err := InitLogger()
	assert.NoError(t, err)
	assert.NotNil(t, appLogger)
}

func TestWithLogging(t *testing.T) {
	var logBuffer bytes.Buffer
	writerSyncer := zapcore.AddSync(&logBuffer)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "" // Убираем время, чтобы не мешало
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		writerSyncer,
		zap.InfoLevel,
	)
	testLogger := zap.New(core).Sugar()

	// Создаем middleware с нашим тестовым логгером
	loggingMiddleware := WithLogging(testLogger)

	t.Run("With_WriteHeader", func(t *testing.T) {
		logBuffer.Reset()

		// Тестовый хендлер, который пишет статус и тело
		dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusAccepted) // 202
			w.Write([]byte("test_body"))       // 9 байт
		})

		wrappedHandler := loggingMiddleware(dummyHandler)

		// Выполняем запрос
		req := httptest.NewRequest(http.MethodGet, "/test-uri", nil)
		rec := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(rec, req)

		// Проверяем сам ответ
		assert.Equal(t, http.StatusAccepted, rec.Code)
		assert.Equal(t, "test_body", rec.Body.String())

		// Проверяем, что было залогировано
		var logEntry map[string]interface{}
		err := json.Unmarshal(logBuffer.Bytes(), &logEntry)
		require.NoError(t, err)

		// Проверяем поля в JSON-логе
		assert.Equal(t, "Request handled", logEntry["msg"])
		assert.Equal(t, "/test-uri", logEntry["uri"])
		assert.Equal(t, http.MethodGet, logEntry["method"])
		// JSON числа парсятся как float64
		assert.Equal(t, float64(http.StatusAccepted), logEntry["status"])
		assert.Equal(t, float64(9), logEntry["size"])
		assert.Contains(t, logEntry, "duration")
	})

	t.Run("No_WriteHeader", func(t *testing.T) {
		logBuffer.Reset() // Очищаем буфер

		// Хендлер, который *не* вызывает WriteHeader
		dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("default_ok")) // 10 байт
		})

		wrappedHandler := loggingMiddleware(dummyHandler)

		req := httptest.NewRequest(http.MethodPost, "/other-uri", nil)
		rec := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(rec, req)

		// Проверяем ответ
		// httptest.Recorder по умолчанию ставит 200, если Write вызван без
		// WriteHeader
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "default_ok", rec.Body.String())

		// Проверяем лог
		var logEntry map[string]interface{}
		err := json.Unmarshal(logBuffer.Bytes(), &logEntry)
		require.NoError(t, err)

		assert.Equal(t, "Request handled", logEntry["msg"])
		assert.Equal(t, "/other-uri", logEntry["uri"])
		assert.Equal(t, http.MethodPost, logEntry["method"])

		assert.Equal(t, float64(0), logEntry["status"])
		assert.Equal(t, float64(10), logEntry["size"])
		assert.Contains(t, logEntry, "duration")
	})
}
