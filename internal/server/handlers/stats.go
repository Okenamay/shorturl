package handlers

import (
	"encoding/json"
	"net"
	"net/http"

	"github.com/Okenamay/shorturl.git/internal/config"
	"github.com/Okenamay/shorturl.git/internal/storage/memselect"
	"go.uber.org/zap"
)

// InternalStats возвращает статистику сервиса. Доступ разрешен только из
// доверенной подсети (TrustedSubnet)
func InternalStats(conf *config.Cfg, appLogger *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Если доверенная подсеть не задана, доступ запрещен всем.
		if conf.TrustedSubnet == "" {
			http.Error(w, "Access denied: trusted subnet not configured", http.StatusForbidden)
			return
		}

		// 2. Проверяем IP клиента (X-Real-IP)
		ipStr := r.Header.Get("X-Real-IP")
		if ipStr == "" {
			http.Error(w, "Access denied: missing X-Real-IP", http.StatusForbidden)
			return
		}

		ip := net.ParseIP(ipStr)
		if ip == nil {
			http.Error(w, "Access denied: invalid IP", http.StatusForbidden)
			return
		}

		// 3. Парсим CIDR и проверяем вхождение
		_, subnet, err := net.ParseCIDR(conf.TrustedSubnet)
		if err != nil {
			appLogger.Errorw("InternalStats - invalid trusted subnet config", "error", err)
			http.Error(w, "Internal configuration error", http.StatusInternalServerError)
			return
		}

		if !subnet.Contains(ip) {
			http.Error(w, "Access denied: IP not trusted", http.StatusForbidden)
			return
		}

		// 4. Получаем статистику через слой хранения
		urls, users, err := memselect.GetStats(r.Context(), conf, appLogger)
		if err != nil {
			appLogger.Errorw("InternalStats - get stats failed", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// 5. Отправляем JSON
		stats := struct {
			URLs  int `json:"urls"`
			Users int `json:"users"`
		}{
			URLs:  urls,
			Users: users,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(stats); err != nil {
			appLogger.Errorw("InternalStats - json encode failed", "error", err)
		}
	}
}
