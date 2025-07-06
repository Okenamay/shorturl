package main

import (
	"github.com/Okenamay/shorturl.git/internal/config"
	logger "github.com/Okenamay/shorturl.git/internal/logger/zap"
	"github.com/Okenamay/shorturl.git/internal/server/router"
	"github.com/Okenamay/shorturl.git/internal/storage/memselect"
)

func main() {
	if err := logger.InitLogger(); err != nil {
		logger.Zap.Fatalw(err.Error(), "Main", "Start logger")
	}

	conf := config.InitConfig()

	err := memselect.MemInit(conf)
	if err != nil {
		logger.Zap.Errorw(err.Error(), "Main", "Initialize storage")
	}
	defer memselect.MemStop(conf)

	logger.Zap.Infow("Starting server on port: ", conf.ServerPort)

	err = router.Launch(conf)
	if err != nil {
		logger.Zap.Fatalw(err.Error(), "Main", "Start server")
	}

	defer logger.Zap.Sync()
}
