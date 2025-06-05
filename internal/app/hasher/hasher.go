package hasher

import (
	"crypto/md5"
	"encoding/hex"
	"io"

	"github.com/Okenamay/shorturl.git/internal/config"
)

// Кодирование строки с URL в md5-сумму с обрезанием до ShortIDLen символов:
func ShortenURL(conf *config.Cfg, fullURL string) string {
	hash := md5.New()
	io.WriteString(hash, fullURL)

	shortID := hex.EncodeToString(hash.Sum(nil))

	return shortID[:conf.ShortIDLen]
}
