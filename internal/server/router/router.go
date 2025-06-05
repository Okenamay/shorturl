package router

import (
	"net/http"
	"time"

	"github.com/Okenamay/shorturl.git/internal/app/middleware/gzipper"
	"github.com/Okenamay/shorturl.git/internal/config"
	logger "github.com/Okenamay/shorturl.git/internal/logger/zap"
	"github.com/Okenamay/shorturl.git/internal/server/handlers"
	"github.com/go-chi/chi/v5"
)

// Запуск HTTP-сервера и работа с запросами:
func Launch(conf config.Cfg) error {
	router := chi.NewRouter()

	router.Use(logger.WithLogging)

	router.Post("/api/shorten", handlers.JSONHandler(conf))
	router.With(gzipper.Decompressor, gzipper.Compressor).Post("/", handlers.ShortenHandler(conf))
	router.With(gzipper.Decompressor, gzipper.Compressor).Get("/{id}", handlers.RedirectHandler(conf))

	server := http.Server{
		Addr:        conf.ServerPort,
		Handler:     router,
		IdleTimeout: time.Duration(conf.IdleTimeout) * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
