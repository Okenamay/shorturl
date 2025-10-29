package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Okenamay/shorturl.git/internal/app/checker"
	emsg "github.com/Okenamay/shorturl.git/internal/app/errmsg"
	"github.com/Okenamay/shorturl.git/internal/app/middleware/auth"
	"github.com/Okenamay/shorturl.git/internal/app/urlmaker"
	"github.com/Okenamay/shorturl.git/internal/config"
	"github.com/Okenamay/shorturl.git/internal/storage/memselect"
	"github.com/Okenamay/shorturl.git/internal/worker"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// Обработка запросов на сокращение URL:
func ShortenHandler(conf *config.Cfg, appLogger *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		appLogger.Info("ShortenHandler started")
		userID, _ := r.Context().Value(auth.UserIDContextKey).(string)

		queryBody, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, emsg.ErrorServer.Error(), http.StatusInternalServerError)
			return
		}

		CheckedURL, checkErr := checker.CheckURL(string(queryBody), appLogger)

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

		w.Header().Set("Content-Type", "text/plain")
		if exists {
			appLogger.Warn("ShortenHandler stopped - already exists")
			w.WriteHeader(http.StatusConflict)
		} else {
			appLogger.Info("ShortenHandler finished")
			w.WriteHeader(http.StatusCreated)
		}
		io.WriteString(w, newURL)
	}
}

// Обработка запроса на переход по полному URL:
func RedirectHandler(conf *config.Cfg, appLogger *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		appLogger.Info("RedirectHandler started")

		queryID := chi.URLParam(r, "id")

		if len(queryID) != conf.ShortIDLen {
			http.Error(w, emsg.ErrorInvalidShortID.Error(), http.StatusNotFound)
			return
		}

		urlInfo, err := memselect.CheckPair(conf, queryID)
		if err != nil {
			appLogger.Errorw("RedirectHandler stopped - URL/ShortID pair check FAIL", "error", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		if urlInfo.OriginalURL == "" {
			appLogger.Errorw("RedirectHandler stopped - URL/ShortID pair find FAIL", "error", err)
			http.Error(w, emsg.ErrorNotInDB.Error(), http.StatusInternalServerError)
			return
		}

		if urlInfo.IsDeleted {
			appLogger.Warn("RedirectHandler stopped - URL deleted")
			w.WriteHeader(http.StatusGone)
			return
		}

		appLogger.Info("RedirectHandler finished")
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Location", urlInfo.OriginalURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

// Выдадим авторизованным пользователям их URL:
func UserURLsHandler(conf *config.Cfg, appLogger *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, authorized := auth.CheckAuth(r)
		if !authorized {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		urls, err := memselect.GetUserURLs(conf, appLogger, userID)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if len(urls) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(urls)
	}
}

// Удалим URLы юзера по списку из запроса:
func BatchDeleter(conf *config.Cfg) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, authorized := auth.CheckAuth(r)
		if !authorized {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var shortIDs []string
		if err := json.NewDecoder(r.Body).Decode(&shortIDs); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if len(shortIDs) == 0 {
			w.WriteHeader(http.StatusAccepted)
			return
		}

		worker.SendToDelete(userID, shortIDs)

		w.WriteHeader(http.StatusAccepted)
	}
}
