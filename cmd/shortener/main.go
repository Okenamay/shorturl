package main

import (
	"github.com/Okenamay/shorturl.git/internal/config"
	logger "github.com/Okenamay/shorturl.git/internal/logger/zap"
	"github.com/Okenamay/shorturl.git/internal/server/router"
	"github.com/Okenamay/shorturl.git/internal/storage/memselect"
	"github.com/Okenamay/shorturl.git/internal/worker"
)

func main() {
	appLogger, err := logger.InitLogger()
	if err != nil {
		// Если логгер не стартовал, мы не можем даже это залогировать. Паникуем.
		panic("failed to initialize logger: " + err.Error())
	}
	defer appLogger.Sync()

	conf := config.InitConfig()

	worker.Start(memselect.BatchDelete, appLogger)

	err = memselect.MemInit(conf, appLogger)
	if err != nil {
		appLogger.Errorw(err.Error(), "Main", "Initialize storage")
	}
	defer memselect.MemStop(conf, appLogger)

	appLogger.Infof("Starting server on port: %s", conf.ServerPort)

	err = router.Launch(conf, appLogger)
	if err != nil {
		appLogger.Fatalw(err.Error(), "Main", "Start server")
	}

	defer appLogger.Sync()
}
