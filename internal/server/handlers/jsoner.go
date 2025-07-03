package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/Okenamay/shorturl.git/internal/app/checker"
	emsg "github.com/Okenamay/shorturl.git/internal/app/errmsg"
	"github.com/Okenamay/shorturl.git/internal/app/urlmaker"
	"github.com/Okenamay/shorturl.git/internal/config"
	logger "github.com/Okenamay/shorturl.git/internal/logger/zap"
	"github.com/Okenamay/shorturl.git/internal/storage/database"
	"github.com/Okenamay/shorturl.git/internal/storage/memselect"
)

type JSONRequest struct {
	URL string `json:"url"`
}

type JSONResponse struct {
	Result string `json:"result"`
}

// Обработка запроса на переход по JSON-запросу:
func JSONHandler(conf *config.Cfg) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Zap.Info("JSONHandler. Start")

		var request JSONRequest
		var buf bytes.Buffer

		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			logger.Zap.Error("JSONHandler. Body error")
			http.Error(w, emsg.ErrorServer.Error(), http.StatusInternalServerError)
			return
		}

		if err = json.Unmarshal(buf.Bytes(), &request); err != nil {
			logger.Zap.Errorw("JSONHandler. Unmarshal error", "Error", err)
			http.Error(w, emsg.ErrorServer.Error(), http.StatusInternalServerError)
			return
		}

		CheckedURL, checkErr := checker.CheckURL(request.URL)

		if checkErr != nil {
			http.Error(w, checkErr.Error(), http.StatusBadRequest)
			return
		}

		fullURL := CheckedURL.String()

		newURL, shortID := urlmaker.ProcessURL(conf, fullURL)

		err = memselect.StorePair(conf, shortID, fullURL)
		if err != nil {
			http.Error(w, emsg.ErrorFileSave.Error(), http.StatusInternalServerError)
			return
		}

		response := JSONResponse{
			Result: newURL,
		}

		data, err := json.Marshal(response)
		if err != nil {
			logger.Zap.Errorw("JSONHandler. Marshal error", "error", err)
			http.Error(w, emsg.ErrorServer.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("content-type", "application/json")
		if database.EntryExists {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusCreated)
		}
		w.Write(data)
		logger.Zap.Info("JSONHandler. Stop")
	}
}
