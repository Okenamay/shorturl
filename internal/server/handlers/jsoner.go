package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Okenamay/shorturl.git/internal/app/checker"
	emsg "github.com/Okenamay/shorturl.git/internal/app/errmsg"
	"github.com/Okenamay/shorturl.git/internal/app/middleware/auth"
	"github.com/Okenamay/shorturl.git/internal/app/urlmaker"
	"github.com/Okenamay/shorturl.git/internal/audit"
	"github.com/Okenamay/shorturl.git/internal/config"
	"github.com/Okenamay/shorturl.git/internal/storage/memselect"
	"go.uber.org/zap"
)

// JSONRequest - структура для входящего JSON-запроса на /api/shorten.
type JSONRequest struct {
	URL string `json:"url"`
}

// JSONResponse - структура для исходящего JSON-ответа от /api/shorten.
type JSONResponse struct {
	Result string `json:"result"`
}

// JSONHandler обрабатывает POST /api/shorten
// Принимает JSON-объект вида {"url": "..."} в теле запроса.
// Генерирует короткий URL, сохраняет его и возвращает в виде JSON-объекта
// {"result": "..."}.
//
// Эта реализация использует потоковые json.Decoder и json.Encoder для
// минимизации аллокаций памяти (не буферизует тело запроса/ответа).
//
// Коды ответа:
// 201 Created: Если URL успешно сокращен.
// 400 Bad Request: Если URL в теле невалидный или JSON некорректен.
// 409 Conflict: Если такой URL уже был сокращен ранее.
// 500 Internal Server Error: При ошибках сохранения.
func JSONHandler(conf *config.Cfg, appLogger *zap.SugaredLogger, auditor *audit.Auditor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		appLogger.Info("JSONHandler started")

		userID, _ := r.Context().Value(auth.UserIDContextKey).(string)

		var request JSONRequest

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			appLogger.Errorw("JSONHandler stopped - JSON decode FAIL", "Error", err)
			http.Error(w, emsg.ErrorServer.Error(), http.StatusBadRequest)
			return
		}

		CheckedURL, checkErr := checker.CheckURL(request.URL, appLogger)
		if checkErr != nil {
			http.Error(w, checkErr.Error(), http.StatusBadRequest)
			return
		}

		fullURL := CheckedURL.String()
		newURL, shortID := urlmaker.ProcessURL(conf, fullURL)

		exists, err := memselect.StorePair(conf, appLogger, userID, shortID, fullURL)
		if err != nil {
			http.Error(w, emsg.ErrorFileSave.Error(), http.StatusInternalServerError)
			return
		}

		if auditor != nil {
			auditor.LogEvent(r.Context(), "shorten", userID, fullURL)
		}

		response := JSONResponse{
			Result: newURL,
		}

		w.Header().Set("content-type", "application/json")
		if exists {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusCreated)
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			appLogger.Errorw("JSONHandler stopped - JSON encode FAIL", "error", err)
			return
		}

		appLogger.Info("JSONHandler finished")
	}
}
