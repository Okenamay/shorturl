package urlmaker

import (
	config "github.com/Okenamay/shorturl.git/internal/config"
)

// Составление строки с сокращённым URL:
func MakeFullURL(shortID string) string {
	newURL := config.Cfg.ShortIDServerPort + "/" + shortID

	return newURL
}
