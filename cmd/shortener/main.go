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

	sugar, err := logger.InitLogger()
	if err != nil {
		sugar.Fatalw(err.Error(), "Main", "Start logger")
	}

	sugar.Infow("Starting server on port: ", config.Cfg.ServerPort)

	err = router.Launch()
	if err != nil {
		sugar.Fatalw(err.Error(), "Main", "Start server")
	}
}
