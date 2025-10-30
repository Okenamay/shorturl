package hasher

import (
	"testing"

	"github.com/Okenamay/shorturl.git/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestShortenURL(t *testing.T) {
	type args struct {
		conf    *config.Cfg
		fullURL string
	}
	tests := []struct {
		name         string
		args         args
		wantKnownVal string
		wantLen      int
	}{
		{
			name: "Test with ShortIDLen 10",
			args: args{
				conf:    &config.Cfg{ShortIDLen: 10},
				fullURL: "https://google.com",
			},
			wantKnownVal: "99999ebcfd", // md5("https://google.com")[:10]
			wantLen:      10,
		},
		{
			name: "Test with ShortIDLen 8",
			args: args{
				conf:    &config.Cfg{ShortIDLen: 8},
				fullURL: "https://yandex.ru",
			},
			wantKnownVal: "e9db20b2", // md5("https://yandex.ru")[:8]
			wantLen:      8,
		},
		{
			name: "Test with ShortIDLen 5",
			args: args{
				conf:    &config.Cfg{ShortIDLen: 5},
				fullURL: "https://example.com",
			},
			wantKnownVal: "c984d", // md5("https://example.com")[:5]
			wantLen:      5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortID := ShortenURL(tt.args.conf, tt.args.fullURL)

			assert.Equal(t, tt.wantLen, len(shortID))
			assert.Equal(t, tt.wantKnownVal, shortID)
			shortIDAgain := ShortenURL(tt.args.conf, tt.args.fullURL)
			assert.Equal(t, shortID, shortIDAgain)
		})
	}

	t.Run("Different URLs produce different hashes", func(t *testing.T) {
		conf := &config.Cfg{ShortIDLen: 10}
		url1 := "https://a.com"
		url2 := "https://b.com"

		hash1 := ShortenURL(conf, url1)
		hash2 := ShortenURL(conf, url2)

		assert.NotEqual(t, hash1, hash2)
	})
}
