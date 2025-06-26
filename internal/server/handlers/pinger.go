package handlers

import (
	"net/http"

	"github.com/Okenamay/shorturl.git/internal/config"
	logger "github.com/Okenamay/shorturl.git/internal/logger/zap"
	"github.com/Okenamay/shorturl.git/internal/storage/memselect"
)

// PingHandler проверяет соединение с базой данных и отвечает на пинг:
func PingHandler(conf *config.Cfg) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sugar, err := logger.InitLogger()
		if err != nil {
			sugar.Errorw(err.Error(), "PingHandler", "Start logger")
		}
		sugar.Info("PingHandler. Start")

		err, pingOK := memselect.PingDB(conf)
		if err != nil {
			sugar.Errorw("PingHandler. DB ping error", "error", err)
			http.Error(w, "Database connection error", http.StatusInternalServerError)
			return
		} else {
			sugar.Info("PingHandler. DB ping success")
		}
		if !pingOK {
			sugar.Infof("Database DSN: %s. MemMode: %s.",
				conf.PostgreDSN, conf.MemMode)
			sugar.Errorw("PingHandler. DB ping failed", "error", err)
			http.Error(w, "Database not enabled", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("PONG"))
	}
}
