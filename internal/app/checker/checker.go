package checker

import (
	"net/url"

	emsg "github.com/Okenamay/shorturl.git/internal/app/errmsg"
	logger "github.com/Okenamay/shorturl.git/internal/logger/zap"
)

// Проверим URL на корректность:
func CheckURL(reqURL string) (*url.URL, error) {
	checkedURL, err := url.ParseRequestURI(reqURL)
	if err != nil {
		return nil, emsg.ErrorInvalidURL
	}

	if checkedURL.Scheme != "https" && checkedURL.Scheme != "http" {
		logger.Zap.Infof("Incorrect scheme: %s", checkedURL.Scheme)
		return nil, emsg.ErrorHTTPS
	}

	logger.Zap.Infof("Detected scheme: %s", checkedURL.Scheme)

	if checkedURL.Host == "" {
		return nil, emsg.ErrorNoHost
	}

	return checkedURL, nil
}
