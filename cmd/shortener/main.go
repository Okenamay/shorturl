package main

import (
	"context"
	"time"

	"github.com/Okenamay/shorturl.git/internal/config"
	logger "github.com/Okenamay/shorturl.git/internal/logger/zap"
	"github.com/Okenamay/shorturl.git/internal/server/router"
	"github.com/Okenamay/shorturl.git/internal/storage/database"
	"github.com/Okenamay/shorturl.git/internal/storage/memselect"
)

func main() {
	if err := logger.InitLogger(); err != nil {
		logger.Zap.Fatalw(err.Error(), "Main", "Start logger")
	}

	conf := config.InitConfig()

	memselect.DeleteChan = make(chan memselect.DeleteTask, 1024)
	runDeletionWorker(memselect.DeleteChan)

	runHardDeleter()

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
	var deleteBuffer []memselect.DeleteTask

	flushBuffer := func() {
		if len(deleteBuffer) == 0 {
			return
		}

		logger.Zap.Infof("Flushing %d deletion tasks from buffer", len(deleteBuffer))

		tasksByUser := make(map[string][]string)
		for _, task := range deleteBuffer {
			tasksByUser[task.UserID] = append(tasksByUser[task.UserID], task.ShortIDs...)
		}

		for userID, shortIDs := range tasksByUser {
			err := memselect.BatchDelete(context.Background(), userID, shortIDs)
			if err != nil {
				logger.Zap.Errorw("Deletion worker failed to batch delete", "error", err, "userID", userID)
			}
		}

		deleteBuffer = nil
	}

	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case task := <-deleteChan:
				deleteBuffer = append(deleteBuffer, task)
				if len(deleteBuffer) >= 25 {
					logger.Zap.Info("Deletion buffer reached size threshold, flushing...")
					flushBuffer()
					// Reset the timer to prevent an immediate second flush.
					ticker.Reset(2 * time.Second)
				}
			case <-ticker.C:
				logger.Zap.Info("Deletion ticker fired, flushing buffer...")
				flushBuffer()
			}
		}
	}()
}

func runHardDeleter() {
	go func() {
		// Create a ticker that fires every 20 seconds.
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			ctx := context.Background()

			err := database.DeleteFlaggedURLs(ctx)
			if err != nil {

				logger.Zap.Errorw("Hard delete worker failed", "error", err)
			}
		}
	}()
}
