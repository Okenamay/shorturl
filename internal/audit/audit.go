package audit

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
)

// AuditEvent - шаблон записи в логе аудита
type AuditEvent struct {
	Timestamp int64  `json:"ts"`
	Action    string `json:"action"`
	UserID    string `json:"user_id"`
	URL       string `json:"url"`
}

// Auditor - сервис логирования аудита
type Auditor struct {
	filePath   string
	remoteURL  string
	logger     *zap.SugaredLogger
	httpClient *http.Client
	fileMutex  sync.Mutex
}

// NewAuditor создает новый экземпляр сервиса аудита
func NewAuditor(filePath, remoteURL string, appLogger *zap.SugaredLogger) *Auditor {
	if filePath == "" && remoteURL == "" {
		appLogger.Info("NewAuditor finished - audit disabled (no file/URL provided)")
		return nil
	}

	appLogger.Infow("NewAuditor - audit service init",
		"file_path", filePath,
		"remote_url", remoteURL,
	)

	return &Auditor{
		filePath:  filePath,
		remoteURL: remoteURL,
		logger:    appLogger,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// LogEvent создает и отправляет событие аудита во все настроенные приемники,
// userID может быть пустым
func (a *Auditor) LogEvent(ctx context.Context, action, userID, originalURL string) {
	event := AuditEvent{
		Timestamp: time.Now().Unix(),
		Action:    action,
		UserID:    userID,
		URL:       originalURL,
	}

	if a.filePath != "" {
		go a.logToFile(event)
	}
	if a.remoteURL != "" {
		go a.logToURL(event)
	}
}

// logToFile безопасно записывает событие в файл
func (a *Auditor) logToFile(event AuditEvent) {
	data, err := json.Marshal(event)
	if err != nil {
		a.logger.Errorw("logToFile stopped - event marshal to file FAIL", "error", err)
		return
	}

	data = append(data, '\n')

	a.fileMutex.Lock()
	defer a.fileMutex.Unlock()

	f, err := os.OpenFile(a.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		a.logger.Errorw("logToFile stopped - open audit file FAIL", "error", err, "path", a.filePath)
		return
	}
	defer f.Close()

	if _, err := f.Write(data); err != nil {
		a.logger.Errorw("logToFile stopped - write to audit file FAIL", "error", err)
	}
}

// logToURL отправляет событие на удаленный сервер
func (a *Auditor) logToURL(event AuditEvent) {
	data, err := json.Marshal(event)
	if err != nil {
		a.logger.Errorw("logToURL - marshal event to URL FAIL", "error", err)
		return
	}

	// Используем context.Background(), т.к. родительский контекст r.Context()
	// может быть отменён до отправки аудита
	ctx, cancel := context.WithTimeout(context.Background(), a.httpClient.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.remoteURL, bytes.NewBuffer(data))
	if err != nil {
		a.logger.Errorw("logToURL stopped - create remote request FAIL", "error", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		a.logger.Errorw("logToURL stopped - send event to URL FAIL", "error", err, "url", a.remoteURL)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		a.logger.Warnw("logToURL stopped - remote server returned non-2xx status", "status", resp.StatusCode)
	}
}
