package urlmaker

import (
	"github.com/Okenamay/shorturl.git/internal/app/hasher"
	"github.com/Okenamay/shorturl.git/internal/config"
)

// Составление строки с сокращённым URL:
func MakeFullURL(conf *config.Cfg, shortID string) string {
	newURL := conf.ShortIDServerPort + "/" + shortID

	return newURL
}

// Делаем вывод с полным новым URL и shortID:
func ProcessURL(conf *config.Cfg, fullURL string) (newURL, shortID string) {
	shortID = hasher.ShortenURL(conf, fullURL)
	newURL = MakeFullURL(conf, shortID)

	return newURL, shortID
}
