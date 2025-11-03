package router

import (
	"net/http"
	"time"

	"github.com/Okenamay/shorturl.git/internal/app/middleware/auth"
	"github.com/Okenamay/shorturl.git/internal/app/middleware/gzipper"
	"github.com/Okenamay/shorturl.git/internal/audit"
	"github.com/Okenamay/shorturl.git/internal/config"
	logger "github.com/Okenamay/shorturl.git/internal/logger/zap"
	"github.com/Okenamay/shorturl.git/internal/server/handlers"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// Launch запускает HTTP-сервер на основе предоставленной конфигурации.
// Эта функция блокирует выполнение до тех пор, пока сервер не будет остановлен.
func Launch(conf *config.Cfg, appLogger *zap.SugaredLogger, auditor *audit.Auditor) error {
	router := NewRouter(conf, appLogger, auditor)

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

// NewRouter создает новый экземпляр роутера (http.Handler) с настроенными
// middleware и всеми эндпоинтами приложения.
func NewRouter(conf *config.Cfg, appLogger *zap.SugaredLogger, auditor *audit.Auditor) http.Handler {
	router := chi.NewRouter()

	router.Use(logger.WithLogging(appLogger))
	router.Use(auth.Authenticator(conf))

	router.Get("/ping", handlers.PingHandler(conf, appLogger))
	router.Post("/api/shorten", handlers.JSONHandler(conf, appLogger, auditor))
	router.Post("/api/shorten/batch", handlers.BatchHandlerTransaction(conf, appLogger))
	router.Get("/api/user/urls", handlers.UserURLsHandler(conf, appLogger))
	router.Delete("/api/user/urls", handlers.BatchDeleter(conf))

	router.With(
		gzipper.Decompressor(appLogger),
		gzipper.Compressor(appLogger),
	).Post("/", handlers.ShortenHandler(conf, appLogger, auditor))
	router.With(
		gzipper.Decompressor(appLogger),
		gzipper.Compressor(appLogger),
	).Get("/{id}", handlers.RedirectHandler(conf, appLogger, auditor))

	return router
}
