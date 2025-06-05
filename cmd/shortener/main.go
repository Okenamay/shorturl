package main

import (
	"github.com/Okenamay/shorturl.git/internal/config"
	logger "github.com/Okenamay/shorturl.git/internal/logger/zap"
	"github.com/Okenamay/shorturl.git/internal/server/router"
	"github.com/Okenamay/shorturl.git/internal/storage/savefile"
)

// Main:
func main() {
	conf := config.InitConfig()

	sugar, err := logger.InitLogger()
	if err != nil {
		sugar.Fatalw(err.Error(), "Main", "Start logger")
	}

	err = savefile.LoadFile(conf)
	if err != nil {
		sugar.Errorw(err.Error(), "Main", "Load savefile")
	}

	sugar.Infow("Starting server on port: ", conf.ServerPort)

	err = router.Launch(conf)
	if err != nil {
		sugar.Fatalw(err.Error(), "Main", "Start server")
	}
}
