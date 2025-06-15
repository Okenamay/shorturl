package savefile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Okenamay/shorturl.git/internal/config"
	"github.com/Okenamay/shorturl.git/internal/storage/memstorage"
)

type record struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// SaveFile записывает всё содержимое memstorage.URLStore в файл
func SaveFile(conf *config.Cfg) error {
	dirPath := filepath.Dir(conf.SaveFilePath)

	// Create the directory path if it doesn't exist
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.OpenFile(conf.SaveFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("ошибка создания файла: %w", err)
	}
	defer file.Close()

	count := 1
	for shortID, fullURL := range memstorage.URLStore {
		rec := record{
			UUID:        strconv.Itoa(count),
			ShortURL:    shortID,
			OriginalURL: fullURL,
		}

		jsonData, err := json.Marshal(rec)
		if err != nil {
			return fmt.Errorf("ошибка маршалинга JSON: %w", err)
		}

		_, err = file.WriteString(string(jsonData) + "\n")
		if err != nil {
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
		return fmt.Errorf("ошибка чтения файла: %w", err)
	}

	lines := bytes.Split(data, []byte("\n"))

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		var rec record
		err := json.Unmarshal(line, &rec)
		if err != nil {
			return fmt.Errorf("ошибка демаршалинга JSON: %w", err)
		}

		memstorage.URLStore[rec.ShortURL] = rec.OriginalURL
	}
	return nil
}
