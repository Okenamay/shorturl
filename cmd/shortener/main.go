package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
	// _ "github.com/Okenamay/shorturl/internal/config"
)

// const (
// 	ShortIDLen  = 10      // Длина короткого идентификатора
// 	ServerPort  = ":8080" // Адрес и порт сервера
// 	IdleTimeout = 600     // Таймаут сервера
// )

var (
	URLStore = make(map[string]string) // Мапа для хранения пар ID – URL
)

// Набор сообщений об ошибках:
var (
	// ErrorMethodNowAllowed = errors.New("method not allowed")
	ErrorServer     = errors.New("server error")
	ErrorInvalidURL = errors.New("invalid URL")
	ErrorNoHost     = errors.New("no URL host found")
	ErrorHTTPS      = errors.New("invalid URL scheme")
	// ErrorBadRequest       = errors.New("bad request")
	// ErrorNotFound         = errors.New("URL not found")
	// ErrorURLTooLong       = errors.New("provided URL too long")
	// ErrorSaveFailed       = errors.New("failed to save URL")
	ErrorInvalidShortID = errors.New("invalid short ID")
	ErrorNotInDB        = errors.New("URL not found in database")
)

// Запуск HTTP-сервера и работа с запросами:
func Launch() {
	router := mux.NewRouter()

	router.HandleFunc("/", ShortenHandler).Methods("POST")
	router.HandleFunc("/{id}", RedirectHandler)

	server := http.Server{
		Addr:        cfg.ServerPort,
		Handler:     router,
		IdleTimeout: time.Duration(cfg.IdleTimeout) * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

// Обработка запросов на сокращение URL:
func ShortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	queryBody, readErr := io.ReadAll(r.Body)
	if readErr != nil {
		http.Error(w, ErrorServer.Error(), http.StatusInternalServerError)
		return
	}

	CheckedURL, checkErr := CheckURL(string(queryBody))

	if checkErr != nil {
		http.Error(w, checkErr.Error(), http.StatusUnprocessableEntity)
		return
	}

	fullURL := CheckedURL.String()

	shortID := AbbreviateURL(fullURL)

	newURL := MakeFullURL(r, cfg.ServerPort, shortID)

	StoreURLIDPair(shortID, fullURL)

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

	if len(queryID) != cfg.ShortIDLen {
		http.Error(w, ErrorInvalidShortID.Error(), http.StatusNotFound)
		return
	}

	fullURL, exists := URLStore[queryID]

	if !exists {
		http.Error(w, ErrorNotInDB.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Location", fullURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// Проверим URL на корректность:
func CheckURL(reqURL string) (*url.URL, error) {
	checkedURL, err := url.ParseRequestURI(reqURL)
	if err != nil {
		return nil, ErrorInvalidURL
	}

	if checkedURL.Scheme != "https" && checkedURL.Scheme != "http" {
		return nil, ErrorHTTPS
	}

	if checkedURL.Host == "" {
		return nil, ErrorNoHost
	}

	return checkedURL, nil
}

// Кодирование строки с URL в md5-сумму с обрезанием до ShortIDLen символов:
func AbbreviateURL(fullURL string) string {
	hash := md5.New()
	io.WriteString(hash, fullURL)

	shortID := hex.EncodeToString(hash.Sum(nil))

	return shortID[:cfg.ShortIDLen]
}

// Составление строки с сокращённым URL:
func MakeFullURL(r *http.Request, port string, shortID string) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	newURL := scheme + "://" + cfg.ShortIDServerPort + "/" + shortID

	return newURL
}

// Сохранение пары fullURL-shortID в URLStore:
func StoreURLIDPair(shortID, fullURL string) {
	URLStore[shortID] = fullURL
}

// Main:
func main() {
	parseFlags()

	log.Printf("Starting server on port %s", cfg.ServerPort)

	Launch()
}
