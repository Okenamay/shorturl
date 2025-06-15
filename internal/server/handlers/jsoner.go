package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	checker "github.com/Okenamay/shorturl.git/internal/app/checker"
	emsg "github.com/Okenamay/shorturl.git/internal/app/errmsg"
	hasher "github.com/Okenamay/shorturl.git/internal/app/hasher"
	urlmaker "github.com/Okenamay/shorturl.git/internal/app/urlmaker"
	logger "github.com/Okenamay/shorturl.git/internal/logger/zap"
	memstorage "github.com/Okenamay/shorturl.git/internal/storage/memstorage"
)

type JSONRequest struct {
	URL string `json:"url"`
}

type JSONResponse struct {
	Result string `json:"result"`
}

// Обработка запроса на переход по JSON-запросу:
func JSONHandler(w http.ResponseWriter, r *http.Request) {
	sugar, _ := logger.InitLogger()
	sugar.Info("JSONHandler. Start")
	if r.Method != http.MethodPost {
		sugar.Errorw("JSONHandler. Method error", "Method", r.Method)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var request JSONRequest
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		sugar.Error("JSONHandler. Body error")
		http.Error(w, emsg.ErrorServer.Error(), http.StatusInternalServerError)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &request); err != nil {
		sugar.Errorw("JSONHandler. Unmarshal error", "Error", err)
		http.Error(w, emsg.ErrorServer.Error(), http.StatusInternalServerError)
		return
	}

	CheckedURL, checkErr := checker.CheckURL(request.URL)

	if checkErr != nil {
		http.Error(w, checkErr.Error(), http.StatusBadRequest)
		return
	}

	fullURL := CheckedURL.String()

	shortID := hasher.ShortenURL(fullURL)

	newURL := urlmaker.MakeFullURL(shortID)

	memstorage.StoreURLIDPair(shortID, fullURL)

	response := JSONResponse{
		Result: newURL,
	}

	data, err := json.Marshal(response)
	if err != nil {
		sugar.Errorw("JSONHandler. Marshal error", "error", err)
		http.Error(w, emsg.ErrorServer.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
	sugar.Info("JSONHandler. Stop")
}
