package main

import (
	"fmt"
	"log"

	"github.com/Okenamay/shorturl.git/internal/audit"
	"github.com/Okenamay/shorturl.git/internal/config"
	logger "github.com/Okenamay/shorturl.git/internal/logger/zap"
	"github.com/Okenamay/shorturl.git/internal/server/router"
	"github.com/Okenamay/shorturl.git/internal/storage/memselect"
	"github.com/Okenamay/shorturl.git/internal/worker"
)

// Переменные для информации о сборке (будут установлены через -ldflags)
var buildVersion string
var buildDate string
var buildCommit string

func main() {
	// Вспомогательная функция для вывода "N/A", если переменная пуста
	na := func(s string) string {
		if s == "" {
			return "N/A"
		}
		return s
	}

	// Выводим информацию о сборке в stdout
	fmt.Printf("Build version: %s\n", na(buildVersion))
	fmt.Printf("Build date: %s\n", na(buildDate))
	fmt.Printf("Build commit: %s\n", na(buildCommit))

	appLogger, err := logger.InitLogger()
	if err != nil {
		// Если логгер не стартовал, мы не можем даже это залогировать, паникуем.
		// Используем standard log.Fatalf, что разрешено в main.main
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer appLogger.Sync()

	conf := config.InitConfig()
	auditor := audit.NewAuditor(conf.AuditFile, conf.AuditURL, appLogger)

	worker.Start(memselect.BatchDelete, appLogger)

	err = memselect.MemInit(conf, appLogger)
	if err != nil {
		appLogger.Errorw(err.Error(), "Main", "Initialize storage")
	}
	defer memselect.MemStop(conf, appLogger)

	appLogger.Infof("Starting server on port: %s", conf.ServerPort)

	err = router.Launch(conf, appLogger, auditor)
	if err != nil {
		appLogger.Fatalw(err.Error(), "Main", "Start server")
	}

	defer appLogger.Sync()
}
