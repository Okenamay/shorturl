package handlers

import (
	"context"
	"net/http"

	"github.com/Okenamay/shorturl.git/internal/config"
	logger "github.com/Okenamay/shorturl.git/internal/logger/zap"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PingHandler проверяет соединение с базой данных и отвечает на пинг:
func PingHandler(conf *config.Cfg) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sugar, _ := logger.InitLogger()
		sugar.Info("PingHandler. Start")

		DBPool, err := pgxpool.New(context.Background(), conf.PostgreDSN)
		if err != nil {
			sugar.Errorw("PingHandler. DB pool error", "error", err)
			http.Error(w, "Database connection error", http.StatusInternalServerError)
			return
		}

		err = DBPool.Ping(context.Background())
		if err != nil {
			sugar.Errorw("PingHandler. DB ping error", "error", err)
			http.Error(w, "Database connection error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("PONG"))
	}
}
