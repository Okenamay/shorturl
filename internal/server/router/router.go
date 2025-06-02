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

	noGzipR := router.PathPrefix("/api/shorten").Subrouter()
	noGzipR.HandleFunc("", handlers.JSONHandler).Methods("POST")
	noGzipR.Use(logger.WithLogging)

	gzipR := router.PathPrefix("/").Subrouter()
	gzipR.HandleFunc("", handlers.ShortenHandler).Methods("POST")
	gzipR.HandleFunc("{id}", handlers.RedirectHandler).Methods("GET")
	gzipR.Use(logger.WithLogging)
	gzipR.Use(gzipper.Decompressor)
	gzipR.Use(gzipper.Compressor)

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
