package main

import (
	"github.com/Okenamay/shorturl.git/internal/config"
	logger "github.com/Okenamay/shorturl.git/internal/logger/zap"
	router "github.com/Okenamay/shorturl.git/internal/server/router"

	_ "go.uber.org/zap"
)

// Main:
func main() {
	config.ParseFlags()

	err := logger.InitLogger()
	if err != nil {
		logger.Sugar.Fatalw(err.Error(), "event", "start logger")
	}

	logger.Sugar.Infow("Starting server on port: ", config.Cfg.ServerPort)

	err = router.Launch()
	if err != nil {
		logger.Sugar.Fatalw(err.Error(), "event", "start server")
	}
}
