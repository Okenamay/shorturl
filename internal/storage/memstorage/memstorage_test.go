package memstorage

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewURLMap(t *testing.T) {
	u := NewURLMap()
	require.NotNil(t, u, "NewURLMap should not return nil")
	require.NotNil(t, u.m, "The internal map 'm' should be initialized")
}

func TestURLMap_SetAndGet(t *testing.T) {
	u := NewURLMap()
	shortID := "testID"
	fullURL := "https://example.com"

	u.Set(shortID, fullURL)

	val, ok := u.Get(shortID)
	require.True(t, ok, "Expected to find key '%s'", shortID)
	assert.Equal(t, fullURL, val, "Expected value for '%s' to be '%s'", shortID, fullURL)
}

func TestURLMap_GetMissing(t *testing.T) {
	u := NewURLMap()

	val, ok := u.Get("missingKey")
	assert.False(t, ok, "Expected not to find key 'missingKey'")
	assert.Empty(t, val, "Expected value for missing key to be empty string")
}

func TestURLMap_Overwrite(t *testing.T) {
	u := NewURLMap()
	shortID := "key1"
	fullURL1 := "https://example.com/v1"
	fullURL2 := "https://example.com/v2"

	u.Set(shortID, fullURL1)
	val1, ok1 := u.Get(shortID)
	require.True(t, ok1)
	require.Equal(t, fullURL1, val1)

	u.Set(shortID, fullURL2)
	val2, ok2 := u.Get(shortID)
	require.True(t, ok2)
	assert.Equal(t, fullURL2, val2, "Expected value to be overwritten")
}

func TestURLMap_GetAll(t *testing.T) {
	u := NewURLMap()
	data := map[string]string{
		"id1": "https://url1.com",
		"id2": "https://url2.com",
		"id3": "https://url3.com",
	}

	allEmpty := u.GetAll()
	require.NotNil(t, allEmpty)
	assert.Empty(t, allEmpty, "Expected GetAll on new map to be empty")

	for k, v := range data {
		u.Set(k, v)
	}
	allFull := u.GetAll()
	require.NotNil(t, allFull)
	assert.Equal(t, data, allFull, "Expected GetAll to return all set items")

	allFull["id1"] = "MODIFIED"
	originalVal, ok := u.Get("id1")
	require.True(t, ok)
	assert.Equal(t, "https://url1.com", originalVal, "Modifying the returned map should not affect the original map")
}

func TestURLMap_ConcurrentAccess(t *testing.T) {
	// Этот тест проверяет потокобезопасность URLMap - он запускает 100
	// горутин. 50 горутин пишут (Set), 50 горутин читают (Get). Если RWMutex
	// реализован правильно - тест пройдет без 'race condition'. Его желательно
	// запускать с флагом -race (go test -race ./...)

	u := NewURLMap()
	var wg sync.WaitGroup
	numGoroutines := 100

	u.Set("initialKey", "initialValue")

	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(n int) {
			defer wg.Done()

			if n%2 == 0 {
				// Писатели (Set)
				shortID := fmt.Sprintf("key-%d", n)
				fullURL := fmt.Sprintf("https://url-%d.com", n)
				u.Set(shortID, fullURL)
			} else {
				// Читатели (Get)
				// Они просто читают одно и то же значение,
				// чтобы создать конкуренцию за RLock
				_, _ = u.Get("initialKey")
			}
		}(i)
	}

	wg.Wait()

	// Проверим, что все писатели отработали
	assert.Len(t, u.m, 51, "Expected 50 new keys + 1 initial key")
	val, ok := u.Get("key-50") // Проверяем одного из писателей
	require.True(t, ok)
	assert.Equal(t, "https://url-50.com", val)
}

// --- Бенчмарки ---

func BenchmarkSet(b *testing.B) {
	u := NewURLMap()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Преобразуем i в строку, чтобы ключ был разным
		key := fmt.Sprintf("key-%d", i)
		u.Set(key, "https://example.com")
	}
}

func BenchmarkGet(b *testing.B) {
	u := NewURLMap()
	u.Set("testKey", "https://example.com")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u.Get("testKey")
	}
}

func BenchmarkGetMissing(b *testing.B) {
	u := NewURLMap()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u.Get("missingKey")
	}
}

func BenchmarkGetAll(b *testing.B) {
	u := NewURLMap()
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key-%d", i)
		u.Set(key, "https://example.com")
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u.GetAll()
	}
}

func BenchmarkConcurrentSetGet(b *testing.B) {
	u := NewURLMap()
	u.Set("initialKey", "initialValue")
	b.ReportAllocs()
	b.ResetTimer()

	// Запускаем в N горутинах, b.N будет разделено между ними
	b.RunParallel(func(pb *testing.PB) {
		// Используем i как локальный счетчик для уникальности ключей
		var i int
		for pb.Next() {
			if i%2 == 0 {
				// Пишем
				key := fmt.Sprintf("key-%d", i)
				u.Set(key, "https://example.com")
			} else {
				// Читаем
				u.Get("initialKey")
			}
			i++
		}
	})
}
