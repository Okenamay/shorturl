package worker

import (
	"context"
	"time"

	logger "github.com/Okenamay/shorturl.git/internal/logger/zap"
	"github.com/Okenamay/shorturl.git/internal/storage/database"
)

type DeleteTask struct {
	UserID   string
	ShortIDs []string
}

var DeleteChan chan DeleteTask

type BatchDeleteFunc func(ctx context.Context, userID string, shortIDs []string) error

func softDeleter(deleteChan <-chan DeleteTask, batchDeleter BatchDeleteFunc) {
	var deleteBuffer []DeleteTask

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
			err := batchDeleter(context.Background(), userID, shortIDs)
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

					ticker.Reset(2 * time.Second)
				}
			case <-ticker.C:
				logger.Zap.Info("Deletion ticker fired, flushing buffer...")
				flushBuffer()
			}
		}
	}()
}

func hardDeleter() {
	go func() {
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

func Start(batchDeleter BatchDeleteFunc) {
	DeleteChan = make(chan DeleteTask, 1024)
	softDeleter(DeleteChan, batchDeleter)
	hardDeleter()
}

func SendToDelete(userID string, shortIDs []string) {
	task := DeleteTask{
		UserID:   userID,
		ShortIDs: shortIDs,
	}
	DeleteChan <- task
}
