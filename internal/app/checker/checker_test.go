package checker

import (
	"testing"

	emsg "github.com/Okenamay/shorturl.git/internal/app/errmsg"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestCheckURL(t *testing.T) {
	appLogger := zap.NewNop().Sugar()

	tests := []struct {
		name      string
		rawURL    string
		wantError error
		wantHost  string
	}{
		{
			name:      "Valid HTTPS URL",
			rawURL:    "https://google.com/search?q=test",
			wantError: nil,
			wantHost:  "google.com",
		},
		{
			name:      "Valid HTTP URL",
			rawURL:    "http://example.com",
			wantError: nil,
			wantHost:  "example.com",
		},
		{
			name:      "Invalid Scheme (ftp)",
			rawURL:    "ftp://example.com",
			wantError: emsg.ErrorHTTPS,
		},
		{
			name:      "No Scheme",
			rawURL:    "example.com",
			wantError: emsg.ErrorInvalidURL,
		},
		{
			name:      "Invalid URL string",
			rawURL:    "not a valid url",
			wantError: emsg.ErrorInvalidURL,
		},
		{
			name:      "Email-like string",
			rawURL:    "test@example.com",
			wantError: emsg.ErrorInvalidURL,
		},
		{
			name:      "No Host",
			rawURL:    "http://",
			wantError: emsg.ErrorNoHost,
		},
		{
			name:      "Scheme with no host",
			rawURL:    "https://",
			wantError: emsg.ErrorNoHost,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkedURL, err := CheckURL(tt.rawURL, appLogger)

			if tt.wantError != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantError)
				assert.Nil(t, checkedURL)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, checkedURL)
				assert.Equal(t, tt.wantHost, checkedURL.Host)
			}
		})
	}
}

// --- Бенчмарки ---

var (
	benchLogger = zap.NewNop().Sugar()
	// Используем result, чтобы компилятор не оптимизировал вызов
	benchResult, _ = CheckURL("https://google.com", benchLogger)
)

func BenchmarkCheckURL_Valid(b *testing.B) {
	rawURL := "https://google.com/search?q=test&param=value"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult, _ = CheckURL(rawURL, benchLogger)
	}
}

func BenchmarkCheckURL_InvalidScheme(b *testing.B) {
	rawURL := "ftp://example.com"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult, _ = CheckURL(rawURL, benchLogger)
	}
}

func BenchmarkCheckURL_InvalidParse(b *testing.B) {
	rawURL := "not a valid url string at all"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult, _ = CheckURL(rawURL, benchLogger)
	}
}
