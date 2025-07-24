package main

import (
	"github.com/Okenamay/shorturl.git/internal/config"
	logger "github.com/Okenamay/shorturl.git/internal/logger/zap"
	"github.com/Okenamay/shorturl.git/internal/server/router"
	"github.com/Okenamay/shorturl.git/internal/storage/memselect"
	"github.com/Okenamay/shorturl.git/internal/worker"
)

func main() {
	if err := logger.InitLogger(); err != nil {
		logger.Zap.Fatalw(err.Error(), "Main", "Start logger")
	}

	conf := config.InitConfig()

	worker.Start(memselect.BatchDelete)

	err := memselect.MemInit(conf)
	if err != nil {
		logger.Zap.Errorw(err.Error(), "Main", "Initialize storage")
	}
	defer memselect.MemStop(conf)

	logger.Zap.Infof("Starting server on port: %s", conf.ServerPort)

	err = router.Launch(conf)
	if err != nil {
		logger.Zap.Fatalw(err.Error(), "Main", "Start server")
	}

	defer logger.Zap.Sync()
}
