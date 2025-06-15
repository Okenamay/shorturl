package router

import (
	"net/http"
	"time"

	gzipper "github.com/Okenamay/shorturl.git/internal/app/middleware/gzipper"
	config "github.com/Okenamay/shorturl.git/internal/config"
	logger "github.com/Okenamay/shorturl.git/internal/logger/zap"
	handlers "github.com/Okenamay/shorturl.git/internal/server/handlers"
	"github.com/go-chi/chi/v5"
)

// Запуск HTTP-сервера и работа с запросами:
func Launch() error {
	router := chi.NewRouter()

	router.Use(logger.WithLogging)

	router.Post("/api/shorten", handlers.JSONHandler)
	router.With(gzipper.Decompressor, gzipper.Compressor).Post("/", handlers.ShortenHandler)
	router.With(gzipper.Decompressor, gzipper.Compressor).Get("/{id}", handlers.RedirectHandler)

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
