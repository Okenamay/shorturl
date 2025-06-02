package router

import (
	"net/http"
	"time"

	gzipper "github.com/Okenamay/shorturl.git/internal/app/middleware/gzipper"
	config "github.com/Okenamay/shorturl.git/internal/config"
	logger "github.com/Okenamay/shorturl.git/internal/logger/zap"
	handlers "github.com/Okenamay/shorturl.git/internal/server/handlers"
	"github.com/gorilla/mux"
)

// Запуск HTTP-сервера и работа с запросами:
func Launch() error {
	router := mux.NewRouter()

	router.Use(logger.WithLogging)
	router.Use(gzipper.Compressor)
	router.Use(gzipper.Decompressor)

	router.HandleFunc("/", handlers.ShortenHandler).Methods("POST")
	router.HandleFunc("/api/shorten", handlers.JSONHandler).Methods("POST")
	router.HandleFunc("/{id}", handlers.RedirectHandler).Methods("GET")

	server := http.Server{
		Addr:        config.Cfg.ServerPort,
		Handler:     router,
		IdleTimeout: time.Duration(config.Cfg.IdleTimeout) * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
