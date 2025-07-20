package main

import (
	"context"
	"time"

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

	memselect.DeleteChan = make(chan memselect.DeleteTask, 1024)

	runDeletionWorker(memselect.DeleteChan)

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

func runDeletionWorker(deleteChan <-chan memselect.DeleteTask) {
	ticker := time.NewTicker(10 * time.Second)

	var deleteBuffer []memselect.DeleteTask

	go func() {
		for {
			select {
			case task := <-deleteChan:
				deleteBuffer = append(deleteBuffer, task)
			case <-ticker.C:
				if len(deleteBuffer) > 0 {
					tasksByUser := make(map[string][]string)
					for _, task := range deleteBuffer {
						tasksByUser[task.UserID] = append(tasksByUser[task.UserID], task.ShortIDs...)
					}

					for userID, shortIDs := range tasksByUser {
						err := memselect.BatchDelete(context.Background(), userID, shortIDs)
						if err != nil {
							logger.Zap.Errorw("Deletion worker failed to batch delete", "error", err)
						}
					}

					deleteBuffer = nil
				}
			}
		}
	}()
}
