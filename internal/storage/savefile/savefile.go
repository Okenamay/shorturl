package savefile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/Okenamay/shorturl.git/internal/config"
	"github.com/Okenamay/shorturl.git/internal/storage/memstorage"
)

type record struct {
	UUID        string `json:"uuid"`
	UserID      string `json:"user_id,omitempty"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

var fileMutex sync.Mutex

// SaveFile записывает всё содержимое memstorage.URLStore в файл
func SaveFile(conf *config.Cfg) error {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	dirPath := filepath.Dir(conf.SaveFilePath)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.OpenFile(conf.SaveFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("ошибка создания файла: %w", err)
	}
	defer file.Close()

	count := 1
	for shortID, fullURL := range memstorage.Store.GetAll() {
		rec := record{
			UUID:        strconv.Itoa(count),
			ShortURL:    shortID,
			OriginalURL: fullURL,
		}

		jsonData, err := json.Marshal(rec)
		if err != nil {
			return fmt.Errorf("ошибка маршалинга JSON: %w", err)
		}

		if _, err := file.WriteString(string(jsonData) + "\n"); err != nil {
			return fmt.Errorf("ошибка записи в файл: %w", err)
		}
		count++
	}

	return nil
}

// LoadFile загружает данные из файла в memstorage.URLStore
func LoadFile(conf *config.Cfg) error {
	data, err := os.ReadFile(conf.SaveFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return fmt.Errorf("ошибка чтения файла: %w", err)
	}

	lines := bytes.Split(data, []byte("\n"))

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		var rec record
		if err := json.Unmarshal(line, &rec); err != nil {
			return fmt.Errorf("ошибка демаршалинга JSON: %w", err)
		}

		memstorage.Store.Set(rec.ShortURL, rec.OriginalURL)
	}

	return nil
}
