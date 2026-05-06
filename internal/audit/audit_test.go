package audit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// newTestLogger создает логгер-пустышку для тестов
func newTestLogger() *zap.SugaredLogger {
	return zap.NewNop().Sugar()
}

func TestNewAuditor(t *testing.T) {
	logger := newTestLogger()

	tests := []struct {
		name      string
		filePath  string
		remoteURL string
		wantNil   bool
	}{
		{
			name:      "no_audit",
			filePath:  "",
			remoteURL: "",
			wantNil:   true,
		},
		{
			name:      "file_based_audit",
			filePath:  "test.log",
			remoteURL: "",
			wantNil:   false,
		},
		{
			name:      "url_based_audit",
			filePath:  "",
			remoteURL: "http://localhost",
			wantNil:   false,
		},
		{
			name:      "both_file_and_url",
			filePath:  "test.log",
			remoteURL: "http://localhost",
			wantNil:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auditor := NewAuditor(tt.filePath, tt.remoteURL, logger)
			if tt.wantNil {
				assert.Nil(t, auditor)
			} else {
				assert.NotNil(t, auditor)
				if auditor != nil {
					assert.Equal(t, tt.filePath, auditor.filePath)
					assert.Equal(t, tt.remoteURL, auditor.remoteURL)
				}
			}
		})
	}
}

func TestLogEventToFile(t *testing.T) {
	logger := newTestLogger()
	// Создаем временный файл
	tmpfile, err := os.CreateTemp("", "audit-*.log")
	require.NoError(t, err)
	filePath := tmpfile.Name()
	tmpfile.Close()
	defer os.Remove(filePath)

	auditor := NewAuditor(filePath, "", logger)
	require.NotNil(t, auditor)

	action, userID, url := "shorten", "user123", "https://google.com"
	auditor.LogEvent(context.Background(), action, userID, url)

	time.Sleep(100 * time.Millisecond)

	data, err := os.ReadFile(filePath)
	require.NoError(t, err)
	require.NotEmpty(t, data)

	var event AuditEvent
	err = json.Unmarshal(data, &event)
	require.NoError(t, err)

	assert.Equal(t, action, event.Action)
	assert.Equal(t, userID, event.UserID)
	assert.Equal(t, url, event.URL)
	assert.InDelta(t, time.Now().Unix(), event.Timestamp, 2)
}

func TestLogEventToURL(t *testing.T) {
	logger := newTestLogger()
	action, userID, url := "follow", "user456", "https://yandex.ru"

	eventChan := make(chan AuditEvent, 1)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var event AuditEvent
		err := json.NewDecoder(r.Body).Decode(&event)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		eventChan <- event
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	auditor := NewAuditor("", ts.URL, logger)
	require.NotNil(t, auditor)

	auditor.LogEvent(context.Background(), action, userID, url)

	select {
	case event := <-eventChan:
		assert.Equal(t, action, event.Action)
		assert.Equal(t, userID, event.UserID)
		assert.Equal(t, url, event.URL)
		assert.InDelta(t, time.Now().Unix(), event.Timestamp, 2)
	case <-time.After(2 * time.Second):
		t.Fatal("Таймаут: мок-сервер не получил событие аудита")
	}
}

func TestFileMutexConcurrent(t *testing.T) {
	logger := newTestLogger()
	tmpfile, err := os.CreateTemp("", "audit-concurrent-*.log")
	require.NoError(t, err)
	filePath := tmpfile.Name()
	tmpfile.Close()
	defer os.Remove(filePath)

	auditor := NewAuditor(filePath, "", logger)
	require.NotNil(t, auditor)

	numRoutines := 50
	var wg sync.WaitGroup
	wg.Add(numRoutines)

	// Запускаем 50 горутин, которые одновременно пишут
	for i := 0; i < numRoutines; i++ {
		go func(i int) {
			defer wg.Done()
			userID := fmt.Sprintf("user-%d", i)
			auditor.LogEvent(context.Background(), "shorten", userID, "https://concurrent.test")
		}(i)
	}

	wg.Wait() // Ждем, пока все LogEvent будут вызваны

	// Даем горутинам logToFile время на запись
	time.Sleep(200 * time.Millisecond)

	// Читаем файл
	data, err := os.ReadFile(filePath)
	require.NoError(t, err)

	// Проверяем, что в файле ровно 50 строк (событий)
	lines := bytes.Split(bytes.TrimSpace(data), []byte("\n"))
	assert.Len(t, lines, numRoutines)

	// Проверяем, что все ID пользователей на месте
	receivedUsers := make(map[string]bool)
	for _, line := range lines {
		var event AuditEvent
		err := json.Unmarshal(line, &event)
		require.NoError(t, err)
		receivedUsers[event.UserID] = true
	}

	assert.Len(t, receivedUsers, numRoutines)
	for i := 0; i < numRoutines; i++ {
		assert.True(t, receivedUsers[fmt.Sprintf("user-%d", i)])
	}
}

// --- Бенчмарки ---

func BenchmarkLogEventToFile(b *testing.B) {
	logger := newTestLogger()
	tmpfile, err := os.CreateTemp("", "bench-audit-*.log")
	if err != nil {
		b.Fatal(err)
	}
	filePath := tmpfile.Name()
	tmpfile.Close()
	defer os.Remove(filePath)

	auditor := NewAuditor(filePath, "", logger)
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// LogEvent сам по себе неблокирующий (запускает горутину), поэтому, по
		// сути, мы измеряем только скорость создания события и запуска горутины
		auditor.LogEvent(ctx, "benchmark", "bench-user", "https://bench.test")
	}
}

func BenchmarkLogEventToURL(b *testing.B) {
	logger := newTestLogger()
	// Создаем "пустой" сервер, который просто принимает запрос
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Читаем тело, чтобы r.Body был полностью обработан
		io.Copy(io.Discard, r.Body)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	auditor := NewAuditor("", ts.URL, logger)
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		auditor.LogEvent(ctx, "benchmark", "bench-user", "https://bench.test")
	}
}
