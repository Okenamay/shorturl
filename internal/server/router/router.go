package router

import (
	"net/http"
	"time"

	"github.com/Okenamay/shorturl.git/internal/config"
	handlers "github.com/Okenamay/shorturl.git/internal/server/handlers"
	"github.com/gorilla/mux"
)

// Запуск HTTP-сервера и работа с запросами:
func Launch() error {
	router := mux.NewRouter()

	router.HandleFunc("/", handlers.ShortenHandler).Methods("POST")
	router.HandleFunc("/{id}", handlers.RedirectHandler)

	server := http.Server{
		Addr:        config.Cfg.ServerPort,
		Handler:     router,
		IdleTimeout: time.Duration(config.Cfg.IdleTimeout) * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		// panic(err)
		return err
	}

	return nil
}
