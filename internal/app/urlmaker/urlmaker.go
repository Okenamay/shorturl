package urlmaker

import (
	"net/http"

	"github.com/Okenamay/shorturl.git/internal/config"
)

// Составление строки с сокращённым URL:
func MakeFullURL(r *http.Request, port string, shortID string) string {
	newURL := config.Cfg.ShortIDServerPort + "/" + shortID

	return newURL
}
