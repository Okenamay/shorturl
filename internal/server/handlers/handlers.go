package handlers

import (
	"io"
	"net/http"

	checker "github.com/Okenamay/shorturl.git/internal/app/checker"
	emsg "github.com/Okenamay/shorturl.git/internal/app/errmsg"
	hasher "github.com/Okenamay/shorturl.git/internal/app/hasher"
	urlmaker "github.com/Okenamay/shorturl.git/internal/app/urlmaker"
	config "github.com/Okenamay/shorturl.git/internal/config"
	memstorage "github.com/Okenamay/shorturl.git/internal/storage/memstorage"
	"github.com/gorilla/mux"
)

// Обработка запросов на сокращение URL:
func ShortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	queryBody, readErr := io.ReadAll(r.Body)
	if readErr != nil {
		http.Error(w, emsg.ErrorServer.Error(), http.StatusInternalServerError)
		return
	}

	CheckedURL, checkErr := checker.CheckURL(string(queryBody))

	if checkErr != nil {
		http.Error(w, checkErr.Error(), http.StatusBadRequest)
		return
	}

	fullURL := CheckedURL.String()

	shortID := hasher.ShortenURL(fullURL)

	newURL := urlmaker.MakeFullURL(shortID)

	memstorage.StoreURLIDPair(shortID, fullURL)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, newURL)
}

// Обработка запроса на переход по полному URL:
func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	queryID := vars["id"]

	if len(queryID) != config.Cfg.ShortIDLen {
		http.Error(w, emsg.ErrorInvalidShortID.Error(), http.StatusNotFound)
		return
	}

	fullURL, exists := memstorage.URLStore[queryID]

	if !exists {
		http.Error(w, emsg.ErrorNotInDB.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Location", fullURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
