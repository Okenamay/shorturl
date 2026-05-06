package worker

import (
	"context"
	"sync"
	"time"

	"github.com/Okenamay/shorturl.git/internal/storage/database"
	"go.uber.org/zap"
)

type DeleteTask struct {
	UserID   string
	ShortIDs []string
}

var (
	DeleteChan chan DeleteTask
	wg         sync.WaitGroup
)

type BatchDeleteFunc func(ctx context.Context, appLogger *zap.SugaredLogger, userID string, shortIDs []string) error

func softDeleter(deleteChan <-chan DeleteTask, batchDeleter BatchDeleteFunc, appLogger *zap.SugaredLogger) {
	defer wg.Done() // Сообщаем, что воркер закончил работу
	var deleteBuffer []DeleteTask

	flushBuffer := func() {
		if len(deleteBuffer) == 0 {
			return
		}

		appLogger.Infof("softDeleter started - flushing %d deletion tasks from buffer", len(deleteBuffer))

		tasksByUser := make(map[string][]string)
		for _, task := range deleteBuffer {
			tasksByUser[task.UserID] = append(tasksByUser[task.UserID], task.ShortIDs...)
		}

		for userID, shortIDs := range tasksByUser {
			err := batchDeleter(context.Background(), appLogger, userID, shortIDs)
			if err != nil {
				appLogger.Errorw("softDeleter stopped - batch delete FAIL", "error", err, "userID", userID)
			}
		}

		appLogger.Info("softDeleter finished - deletion worker batch delete OK")
		deleteBuffer = nil
	}

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case task, ok := <-deleteChan:
			if !ok {
				// Канал закрыт, сбрасываем остатки и выходим
				appLogger.Info("softDeleter stopping - channel closed")
				flushBuffer()
				return
			}
			deleteBuffer = append(deleteBuffer, task)
			if len(deleteBuffer) >= 25 {
				appLogger.Info("softDeletion - buffer reached size threshold, flushing...")
				flushBuffer()
				ticker.Reset(2 * time.Second)
			}
		case <-ticker.C:
			appLogger.Info("softDeletion - deletion ticker fired, flushing...")
			flushBuffer()
		}
	}
}

func hardDeleter(appLogger *zap.SugaredLogger) {
	go func() {
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			ctx := context.Background()

			err := database.DeleteFlaggedURLs(ctx, appLogger)
			if err != nil {
				appLogger.Errorw("Hard delete worker failed", "error", err)
			}
		}
	}()
}

func Start(batchDeleter BatchDeleteFunc, appLogger *zap.SugaredLogger) {
	DeleteChan = make(chan DeleteTask, 1024)
	wg.Add(1)
	go softDeleter(DeleteChan, batchDeleter, appLogger)
	hardDeleter(appLogger)
}

// Stop корректно останавливает воркеры, дожидаясь обработки всех задач
func Stop(appLogger *zap.SugaredLogger) {
	appLogger.Info("Stopping workers...")
	close(DeleteChan)
	wg.Wait()
	appLogger.Info("Workers stopped")
}

func SendToDelete(userID string, shortIDs []string) {
	task := DeleteTask{
		UserID:   userID,
		ShortIDs: shortIDs,
	}
	DeleteChan <- task
}
