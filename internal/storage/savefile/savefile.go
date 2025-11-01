package savefile

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/Okenamay/shorturl.git/internal/config"
	"github.com/Okenamay/shorturl.git/internal/storage/memstorage"
	"go.uber.org/zap"
)

type record struct {
	UUID        string `json:"uuid"`
	UserID      string `json:"user_id,omitempty"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

var fileMutex sync.Mutex

// SaveFile записывает всё содержимое memstorage.URLStore в файл
func SaveFile(conf *config.Cfg, appLogger *zap.SugaredLogger) error {
	appLogger.Info("SaveFile started")
	fileMutex.Lock()
	defer fileMutex.Unlock()

	dirPath := filepath.Dir(conf.SaveFilePath)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		appLogger.Errorw("SaveFile stopped - create directory FAIL", "error", err)
		return err
	}

	file, err := os.OpenFile(conf.SaveFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		appLogger.Errorw("SaveFile stopped - create file FAIL", "error", err)
		return err
	}
	defer file.Close()

	// Изменение: Используем потоковый энкодер
	encoder := json.NewEncoder(file)

	count := 1
	for shortID, fullURL := range memstorage.Store.GetAll() {
		rec := record{
			UUID:        strconv.Itoa(count),
			ShortURL:    shortID,
			OriginalURL: fullURL,
		}

		// Изменение: Убираем Marshal, string() и + "\n"
		if err := encoder.Encode(rec); err != nil {
			appLogger.Errorw("SaveFile stopped - write JSON to file FAIL", "error", err)
			return err
		}

		count++
	}

	appLogger.Info("SaveFile finished")
	return nil
}

// LoadFile загружает данные из файла в memstorage.URLStore
func LoadFile(conf *config.Cfg, appLogger *zap.SugaredLogger) error {
	appLogger.Info("LoadFile started")

	// Изменение: Используем os.Open для потокового чтения
	file, err := os.Open(conf.SaveFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		appLogger.Errorw("LoadFile stopped - read file FAIL", "error", err)
		return err
	}
	defer file.Close()

	// Изменение: Читаем файл построчно с помощью сканера
	scanner := bufio.NewScanner(file)

	// bufio.Scanner по умолчанию имеет ограничение на длину строки, увеличим
	// до 4 kiB
	const maxLineSize = 4096
	buf := make([]byte, maxLineSize)
	scanner.Buffer(buf, maxLineSize)

	for scanner.Scan() {
		line := scanner.Bytes() // Используем .Bytes() для экономии аллокаций
		if len(line) == 0 {
			continue
		}

		var rec record
		if err := json.Unmarshal(line, &rec); err != nil {
			// Логичнее пропустить плохую строку, чем прерывать всю загрузку
			appLogger.Errorf("LoadFile stopped - JSON demarshal FAIL, row ommited: %v\n", err)
			return err
		}

		memstorage.Store.Set(rec.ShortURL, rec.OriginalURL)
	}

	// Проверяем на ошибки самого сканера (например, если токен был слишком
	// длинный)
	if err := scanner.Err(); err != nil {
		appLogger.Errorw("LoadFile stopped - file scan FAIL", "error", err)
		return err
	}

	appLogger.Info("LoadFile finished")
	return nil
}
