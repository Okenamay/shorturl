package savefile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/Okenamay/shorturl.git/internal/config"
	"github.com/Okenamay/shorturl.git/internal/storage/memstorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTest - вспомогательная функция для настройки тестового окружения,
// которая создает временный файл и сбрасывает глобальное хранилище
func setupTest(t *testing.T) *config.Cfg {
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test_save.json")

	conf := &config.Cfg{
		SaveFilePath: tempFile,
	}

	memstorage.Store = memstorage.NewURLMap()

	return conf
}

// TestSaveLoad_Roundtrip проверяет полный цикл: запись в файл и последующая
// загрузка
func TestSaveLoad_Roundtrip(t *testing.T) {
	conf := setupTest(t)

	testData := map[string]string{
		"id1": "https://url1.com",
		"id2": "https://url2.com",
	}
	for k, v := range testData {
		memstorage.Store.Set(k, v)
	}

	err := SaveFile(conf)
	require.NoError(t, err, "SaveFile не должен возвращать ошибку")

	memstorage.Store = memstorage.NewURLMap()
	require.Empty(t, memstorage.Store.GetAll(), "Хранилище должно быть пустым после сброса")

	err = LoadFile(conf)
	require.NoError(t, err, "LoadFile не должен возвращать ошибку")

	loadedData := memstorage.Store.GetAll()
	assert.Equal(t, testData, loadedData, "Загруженные данные должны соответствовать сохраненным")
}

// TestLoadFile_NotExist проверяет, что LoadFile не возвращает ошибку, если
// файл сохранения еще не существует
func TestLoadFile_NotExist(t *testing.T) {
	conf := setupTest(t)
	// Мы не создаем файл, conf.SaveFilePath указывает на несуществующий путь

	err := LoadFile(conf)
	require.NoError(t, err, "LoadFile должен вернуть nil, если файл не найден (IsNotExist)")

	assert.Empty(t, memstorage.Store.GetAll(), "Хранилище должно быть пустым, если файл не найден")
}

// TestLoadFile_CorruptedData проверяет, что LoadFile возвращает ошибку, если
// файл содержит невалидные (не JSON) данные
func TestLoadFile_CorruptedData(t *testing.T) {
	conf := setupTest(t)

	// Записываем мусор в файл
	junkData := []byte("this is not json\n{\"valid\": \"json\"}\nbut this line is bad")
	err := os.WriteFile(conf.SaveFilePath, junkData, 0644)
	require.NoError(t, err, "Не удалось записать тестовые данные в файл")

	// Пытаемся загрузить мусор
	err = LoadFile(conf)
	require.Error(t, err, "LoadFile должен вернуть ошибку при чтении поврежденного файла")

	assert.Empty(t, memstorage.Store.GetAll(), "Хранилище должно быть пустым, если файл поврежден")
}

// TestSaveFile_Truncate проверяет, что SaveFile перезаписывает (O_TRUNC)
// существующий файл, а не добавляет в него
func TestSaveFile_Truncate(t *testing.T) {
	conf := setupTest(t)

	junkData := []byte("old data that must be truncated")
	err := os.WriteFile(conf.SaveFilePath, junkData, 0644)
	require.NoError(t, err)

	memstorage.Store.Set("newID", "https://new.com")

	err = SaveFile(conf)
	require.NoError(t, err)

	fileData, err := os.ReadFile(conf.SaveFilePath)
	require.NoError(t, err)

	fileContent := string(fileData)

	assert.NotContains(t, fileContent, "old data")

	var rec record
	err = json.Unmarshal(fileData[:len(fileData)-1], &rec) // -1, чтобы убрать \n
	require.NoError(t, err)
	assert.Equal(t, "newID", rec.ShortURL)
	assert.Equal(t, "https://new.com", rec.OriginalURL)
}

// --- Бенчмарки ---

// prepareBenchFile создает "золотой" файл для бенчмарков загрузки
func prepareBenchFile(b *testing.B, numRecords int) *config.Cfg {
	b.Helper()
	tempDir := b.TempDir()
	tempFile := filepath.Join(tempDir, "bench_save.json")

	conf := &config.Cfg{
		SaveFilePath: tempFile,
	}

	// Наполняем хранилище
	store := memstorage.NewURLMap()
	for i := 0; i < numRecords; i++ {
		store.Set(fmt.Sprintf("short-%d", i), fmt.Sprintf("https://long-url-%d.com", i))
	}
	// Используем глобальное хранилище для SaveFile
	memstorage.Store = store

	// Сохраняем файл
	err := SaveFile(conf)
	if err != nil {
		b.Fatalf("Failed to prepare bench file: %v", err)
	}

	return conf
}

func BenchmarkSaveFile(b *testing.B) {
	tempDir := b.TempDir()
	tempFile := filepath.Join(tempDir, "bench_save.json")

	conf := &config.Cfg{
		SaveFilePath: tempFile,
	}

	// Наполняем хранилище данными один раз
	store := memstorage.NewURLMap()
	for i := 0; i < 100; i++ {
		store.Set(fmt.Sprintf("short-%d", i), fmt.Sprintf("https://long-url-%d.com", i))
	}
	memstorage.Store = store

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// SaveFile по своей природе перезаписывает (Truncate) файл,
		// поэтому нам не нужно его удалять в цикле
		if err := SaveFile(conf); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLoadFile(b *testing.B) {
	// Создаем "золотой" файл с 100 записями один раз
	conf := prepareBenchFile(b, 100)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Сбрасываем хранилище перед каждой загрузкой
		b.StopTimer()
		memstorage.Store = memstorage.NewURLMap()
		b.StartTimer()

		if err := LoadFile(conf); err != nil {
			b.Fatal(err)
		}
	}
}
