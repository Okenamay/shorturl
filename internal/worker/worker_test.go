package worker

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// mockCall - структура для захвата вызовов mock-функции
type mockCall struct {
	userID   string
	shortIDs []string
}

// TestWorker_SoftDelete покрывает всю логику softDeleter, включая сброс по
// буферу, сброс по тикеру и агрегацию задач
func TestWorker_SoftDelete(t *testing.T) {
	testLogger := zap.NewNop().Sugar()

	t.Run("it flushes when buffer is full", func(t *testing.T) {
		// Канал для получения вызовов из mock-функции
		callsCh := make(chan mockCall, 25)
		// Mock-функция, которая отправляет полученные данные в канал
		mockDeleter := func(ctx context.Context, appLogger *zap.SugaredLogger, userID string, shortIDs []string) error {
			callsCh <- mockCall{userID: userID, shortIDs: shortIDs}
			return nil
		}

		Start(mockDeleter, testLogger)

		expectedCalls := make(map[string]mockCall)
		// В коде размер буфера = 25
		for i := 0; i < 25; i++ {
			userID := fmt.Sprintf("user-buffer-%d", i)
			shortID := fmt.Sprintf("id-%d", i)
			task := mockCall{userID: userID, shortIDs: []string{shortID}}
			expectedCalls[userID] = task

			SendToDelete(task.userID, task.shortIDs)
		}

		// Мы ожидаем ровно 25 вызовов (по одному на каждого пользователя)
		receivedCalls := make(map[string]mockCall)
		timeout := time.After(2 * time.Second) // 2 сек таймаут на сброс

		for i := 0; i < 25; i++ {
			select {
			case call := <-callsCh:
				receivedCalls[call.userID] = call
			case <-timeout:
				t.Fatal("Тест упал по таймауту в ожидании сброса буфера")
			}
		}

		assert.Equal(t, expectedCalls, receivedCalls)
	})

	t.Run("it flushes on ticker", func(t *testing.T) {
		callsCh := make(chan mockCall, 1)
		mockDeleter := func(ctx context.Context, appLogger *zap.SugaredLogger, userID string, shortIDs []string) error {
			callsCh <- mockCall{userID: userID, shortIDs: shortIDs}
			return nil
		}

		Start(mockDeleter, testLogger)

		// Отправляем 1 задачу (недостаточно для заполнения буфера)
		task := mockCall{userID: "user-ticker", shortIDs: []string{"ticker-id"}}
		SendToDelete(task.userID, task.shortIDs)

		// Тикер в коде 2 секунды ждём 3
		timeout := time.After(3 * time.Second)

		select {
		case call := <-callsCh:
			// Проверяем, что пришла именно наша задача
			assert.Equal(t, task, call)
		case <-timeout:
			t.Fatal("Тест упал по таймауту в ожидании сброса по тикеру")
		}
	})

	t.Run("it aggregates tasks for the same user", func(t *testing.T) {
		callsCh := make(chan mockCall, 1)
		mockDeleter := func(ctx context.Context, appLogger *zap.SugaredLogger, userID string, shortIDs []string) error {
			callsCh <- mockCall{userID: userID, shortIDs: shortIDs}
			return nil
		}

		Start(mockDeleter, testLogger)

		// Отправляем 2 задачи для ОДНОГО пользователя
		SendToDelete("userOne", []string{"id-a"})
		SendToDelete("userOne", []string{"id-b"})

		// Ждем сброса по тикеру (3 сек)
		timeout := time.After(3 * time.Second)
		// Ожидаем, что задачи "склеятся"
		expectedCall := mockCall{userID: "userOne", shortIDs: []string{"id-a", "id-b"}}

		select {
		case call := <-callsCh:
			// Проверяем, что пришел 1 вызов с 2 ID
			assert.Equal(t, expectedCall, call)
		case <-timeout:
			t.Fatal("Тест упал по таймауту в ожидании сброса по тикеру")
		}
	})
}

// --- Бенчмарки ---

func BenchmarkSendToDelete(b *testing.B) {
	testLogger := zap.NewNop().Sugar()
	// Мок-функция, которая просто "поглощает" вызовы, чтобы воркер не блокировался
	mockDeleter := func(ctx context.Context, appLogger *zap.SugaredLogger, userID string, shortIDs []string) error {
		return nil
	}

	// Запускаем воркер один раз
	Start(mockDeleter, testLogger)

	b.ReportAllocs()
	b.ResetTimer()

	// Запускаем отправку в параллельных горутинах
	b.RunParallel(func(pb *testing.PB) {
		// Каждая горутина имеет свой уникальный userID
		userID := fmt.Sprintf("bench-user-%d", time.Now().UnixNano())
		for pb.Next() {
			SendToDelete(userID, []string{"bench-id-1", "bench-id-2"})
		}
	})
}

func BenchmarkSoftDeleter_FlushBuffer(b *testing.B) {
	testLogger := zap.NewNop().Sugar()
	var wg sync.WaitGroup

	// Мок-функция, которая сигнализирует о завершении
	mockDeleter := func(ctx context.Context, appLogger *zap.SugaredLogger, userID string, shortIDs []string) error {
		wg.Done()
		return nil
	}

	Start(mockDeleter, testLogger)

	// Создаем "пачку" из 25 задач (размер буфера)
	tasks := make([]DeleteTask, 25)
	for i := 0; i < 25; i++ {
		tasks[i] = DeleteTask{
			UserID:   fmt.Sprintf("bench-user-flush-%d", i),
			ShortIDs: []string{"id-1"},
		}
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Настраиваем WaitGroup на 25 вызовов (по 1 на пользователя)
		b.StopTimer()
		wg.Add(25)
		b.StartTimer()

		// Заполняем буфер
		for _, task := range tasks {
			SendToDelete(task.UserID, task.ShortIDs)
		}

		// Ждем, пока mockDeleter не будет вызван 25 раз
		wg.Wait()
	}
}
