package urlmaker

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Okenamay/shorturl.git/internal/app/hasher"
	"github.com/Okenamay/shorturl.git/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMakeFullURL(t *testing.T) {
	testConf := &config.Cfg{
		ShortIDAddress: "http://localhost:8080",
	}
	testShortID := "aBcDeF12"

	expectedURL := fmt.Sprintf("%s/%s", testConf.ShortIDAddress, testShortID)
	resultURL := MakeFullURL(testConf, testShortID)

	assert.Equal(t, expectedURL, resultURL)
}

func TestProcessURL(t *testing.T) {
	testConf := &config.Cfg{
		ShortIDAddress: "http://example.com",
		ShortIDLen:     8,
	}
	fullURL := "https://practicum.yandex.ru/"

	newURL, shortID := ProcessURL(testConf, fullURL)

	require.Len(t, shortID, testConf.ShortIDLen)

	expectedURL := fmt.Sprintf("%s/%s", testConf.ShortIDAddress, shortID)
	assert.Equal(t, expectedURL, newURL)

	assert.True(t, strings.HasSuffix(newURL, shortID))

	expectedShortID := hasher.ShortenURL(testConf, fullURL)
	assert.Equal(t, expectedShortID, shortID)
}

// --- Бенчмарки ---

func BenchmarkMakeFullURL(b *testing.B) {
	testConf := &config.Cfg{
		ShortIDAddress: "http://localhost:8080",
	}
	testShortID := "aBcDeF12"

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MakeFullURL(testConf, testShortID)
	}
}

func BenchmarkProcessURL(b *testing.B) {
	testConf := &config.Cfg{
		ShortIDAddress: "http://example.com",
		ShortIDLen:     8,
	}
	fullURL := "https://practicum.yandex.ru/some/very/long/path/to/benchmark"

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Результат присваиваем, чтобы компилятор не "выкинул" вызов
		_, _ = ProcessURL(testConf, fullURL)
	}
}
