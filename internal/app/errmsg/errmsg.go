package errmsg

import (
	"errors"
)

// Набор сообщений об ошибках:
var (
	ErrorServer         = errors.New("server error")
	ErrorInvalidURL     = errors.New("invalid URL")
	ErrorNoHost         = errors.New("no URL host found")
	ErrorHTTPS          = errors.New("invalid URL scheme")
	ErrorInvalidShortID = errors.New("invalid short ID")
	ErrorNotInDB        = errors.New("URL not found in database")
	ErrorFileSave       = errors.New("file save failed")
)
