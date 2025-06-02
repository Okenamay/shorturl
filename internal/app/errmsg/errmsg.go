package errmsg

import (
	"errors"
)

// Набор сообщений об ошибках:
var (
	// ErrorMethodNowAllowed = errors.New("method not allowed")
	ErrorServer     = errors.New("server error")
	ErrorInvalidURL = errors.New("invalid URL")
	ErrorNoHost     = errors.New("no URL host found")
	ErrorHTTPS      = errors.New("invalid URL scheme")
	// ErrorBadRequest       = errors.New("bad request")
	// ErrorNotFound         = errors.New("URL not found")
	// ErrorURLTooLong       = errors.New("provided URL too long")
	// ErrorSaveFailed       = errors.New("failed to save URL")
	ErrorInvalidShortID = errors.New("invalid short ID")
	ErrorNotInDB        = errors.New("URL not found in database")
)
