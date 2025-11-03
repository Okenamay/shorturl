package handlers

import (
	"net/http"

	"github.com/Okenamay/shorturl.git/internal/config"
	"github.com/Okenamay/shorturl.git/internal/storage/memselect"
	"go.uber.org/zap"
)

// PingHandler обрабатывает GET /ping
// Используется для проверки доступности хранилища (PostgreSQL).
//
// Коды ответа:
// 200 OK: Если подключение к БД успешно.
// 500 Internal Server Error: Если БД недоступна или (для memstore/savefile)
// если режим работы не "postgres".
func PingHandler(conf *config.Cfg, appLogger *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		appLogger.Info("PingHandler started")

		err, pingOK := memselect.PingDB(conf, appLogger)
		if err != nil {
			appLogger.Errorw("PingHandler stopped - DB ping FAIL", "error", err)
			http.Error(w, "Database connection error", http.StatusInternalServerError)
			return
		}

		if !pingOK {
			appLogger.Warnf("PingHandler stopped - DB not configured for ping. MemMode: %s",
				conf.MemMode)
			http.Error(w, "Database not enabled", http.StatusInternalServerError)
			return
		}

		appLogger.Info("PingHandler finished - DB ping OK")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("PONG"))
	}
}
