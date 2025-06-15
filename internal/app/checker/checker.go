package checker

import (
	"net/url"

	emsg "github.com/Okenamay/shorturl.git/internal/app/errmsg"
)

// Проверим URL на корректность:
func CheckURL(reqURL string) (*url.URL, error) {
	checkedURL, err := url.ParseRequestURI(reqURL)
	if err != nil {
		return nil, emsg.ErrorInvalidURL
	}

	if checkedURL.Scheme != "https" && checkedURL.Scheme != "http" {
		return nil, emsg.ErrorHTTPS
	}

	if checkedURL.Host == "" {
		return nil, emsg.ErrorNoHost
	}

	return checkedURL, nil
}
