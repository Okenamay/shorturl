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
	logger "github.com/Okenamay/shorturl.git/internal/logger/zap"
	"github.com/Okenamay/shorturl.git/internal/storage/memselect"
	"github.com/go-chi/chi/v5"
)

// Обработка запросов на сокращение URL:
func ShortenHandler(conf *config.Cfg) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(auth.UserIDContextKey).(string)

		queryBody, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, emsg.ErrorServer.Error(), http.StatusInternalServerError)
			return
		}

		CheckedURL, checkErr := checker.CheckURL(string(queryBody))

		if checkErr != nil {
			http.Error(w, checkErr.Error(), http.StatusBadRequest)
			return
		}

		fullURL := CheckedURL.String()
		newURL, shortID := urlmaker.ProcessURL(conf, fullURL)

		exists, err := memselect.StorePair(conf, userID, shortID, fullURL)
		if err != nil {
			http.Error(w, emsg.ErrorFileSave.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		if exists {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusCreated)
		}
		io.WriteString(w, newURL)
	}
}

// Обработка запроса на переход по полному URL:
func RedirectHandler(conf *config.Cfg) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Zap.Info("RedirectHandler. Start")

		queryID := chi.URLParam(r, "id")

		if len(queryID) != conf.ShortIDLen {
			http.Error(w, emsg.ErrorInvalidShortID.Error(), http.StatusNotFound)
			return
		}

		fullURL, err := memselect.CheckPair(conf, queryID)
		if err != nil {
			logger.Zap.Errorw("RedirectHandler. Failed to check URL/ShortID pair", "error", err)
			return
		}

		if fullURL == "" {
			logger.Zap.Errorw("RedirectHandler. Failed to find URL/ShortID pair", "error", err)
			http.Error(w, emsg.ErrorNotInDB.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Location", fullURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

// Выдадим авторизованным пользователям их URL:
func UserURLsHandler(conf *config.Cfg) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(auth.UserIDContextKey).(string)
		if !ok || userID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		urls, err := memselect.GetUserURLs(conf, userID)
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
