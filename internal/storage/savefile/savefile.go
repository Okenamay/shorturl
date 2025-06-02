package savefile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	config "github.com/Okenamay/shorturl.git/internal/config"
	memstorage "github.com/Okenamay/shorturl.git/internal/storage/memstorage"
)

type record struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// SaveFile записывает всё содержимое memstorage.URLStore в файл
func SaveFile() error {
	// Создаём файл, если его нет
	file, err := os.OpenFile(config.Cfg.SaveFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("ошибка создания файла: %w", err)
	}
	defer file.Close()

	// Записываем каждую пару из memstorage
	count := 1
	for shortID, fullURL := range memstorage.URLStore {
		rec := record{
			UUID:        strconv.Itoa(count),
			ShortURL:    shortID,
			OriginalURL: fullURL,
		}

		// Преобразуем запись в JSON
		jsonData, err := json.Marshal(rec)
		if err != nil {
			return fmt.Errorf("ошибка маршалинга JSON: %w", err)
		}

		// Записываем в файл
		_, err = file.WriteString(string(jsonData) + "\n")
		if err != nil {
			return fmt.Errorf("ошибка записи в файл: %w", err)
		}

		count++
	}
	return nil
}

// LoadFile загружает данные из файла в memstorage.URLStore
func LoadFile() error {
	// Читаем содержимое файла
	data, err := os.ReadFile(config.Cfg.SaveFilePath)
	if err != nil {
		return fmt.Errorf("ошибка чтения файла: %w", err)
	}

	// Разбиваем данные по строкам
	lines := bytes.Split(data, []byte("\n"))

	// Загружаем каждую запись
	for _, line := range lines {
		if len(line) == 0 {
			continue // пропускаем пустые строки
		}

		var rec record
		err := json.Unmarshal(line, &rec)
		if err != nil {
			return fmt.Errorf("ошибка демаршалинга JSON: %w", err)
		}

		// Добавляем в memstorage
		memstorage.URLStore[rec.ShortURL] = rec.OriginalURL
	}
	return nil
}
