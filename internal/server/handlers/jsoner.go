package handlers

import (
	"bytes"
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

type JSONRequest struct {
	URL string `json:"url"`
}

type JSONResponse struct {
	Result string `json:"result"`
}

// Обработка запроса на переход по JSON-запросу:
func JSONHandler(conf *config.Cfg, appLogger *zap.SugaredLogger, auditor *audit.Auditor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		appLogger.Info("JSONHandler started")

		userID, _ := r.Context().Value(auth.UserIDContextKey).(string)

		var request JSONRequest
		var buf bytes.Buffer

		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			appLogger.Errorw("JSONHandler stopped - read body FAIL", "error", err)
			http.Error(w, emsg.ErrorServer.Error(), http.StatusInternalServerError)
			return
		}

		if err = json.Unmarshal(buf.Bytes(), &request); err != nil {
			appLogger.Errorw("JSONHandler stopped - JSON unmarshal FAIL", "Error", err)
			http.Error(w, emsg.ErrorServer.Error(), http.StatusInternalServerError)
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

		data, err := json.Marshal(response)
		if err != nil {
			appLogger.Errorw("JSONHandler stopped - JSON marshal FAIL", "error", err)
			http.Error(w, emsg.ErrorServer.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("content-type", "application/json")
		if exists {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusCreated)
		}
		w.Write(data)
		appLogger.Info("JSONHandler finished")
	}
}
