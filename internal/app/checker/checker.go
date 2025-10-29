package checker

import (
	"net/url"

	emsg "github.com/Okenamay/shorturl.git/internal/app/errmsg"
	"go.uber.org/zap"
)

// Проверим URL на корректность:
func CheckURL(reqURL string, appLogger *zap.SugaredLogger) (*url.URL, error) {
	checkedURL, err := url.ParseRequestURI(reqURL)
	if err != nil {
		return nil, emsg.ErrorInvalidURL
	}

	if checkedURL.Scheme != "https" && checkedURL.Scheme != "http" {
		appLogger.Infof("Incorrect scheme: %s", checkedURL.Scheme)
		return nil, emsg.ErrorHTTPS
	}

	appLogger.Infof("Detected scheme: %s", checkedURL.Scheme)

	if checkedURL.Host == "" {
		return nil, emsg.ErrorNoHost
	}

	return checkedURL, nil
}
