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

		// Проверяем подключение к БД
		sugar.Info("PingHandler. тест 0")
		// err := database.DBPool.Ping(context.Background())
		DBPool, err := pgxpool.New(context.Background(), conf.PostgreDSN)
		if err != nil {
			sugar.Info("PingHandler. тест 5")
			sugar.Errorw("PingHandler. DB pool error", "error", err)
			http.Error(w, "Database connection error", http.StatusInternalServerError)
			return
		}

		err = DBPool.Ping(context.Background())
		sugar.Info("PingHandler. тест 1")
		if err != nil {
			sugar.Info("PingHandler. тест 2")
			sugar.Errorw("PingHandler. DB ping error", "error", err)
			http.Error(w, "Database connection error", http.StatusInternalServerError)
			return
		}

		sugar.Info("PingHandler. тест 3")
		w.WriteHeader(http.StatusOK)
		sugar.Info("PingHandler. тест 4")
		w.Write([]byte("PONG"))
	}
}
