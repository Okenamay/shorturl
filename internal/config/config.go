package config

import (
	"flag"
	"os"
)

// Дефолтные значения до применения флагов:
const (
	ShortIDLen        = 10                                                           // Длина короткого идентификатора
	IdleTimeout       = 600                                                          // Таймаут сервера в секундах
	ServerPort        = ":8080"                                                      // Адрес и порт сервера
	ShortIDServerPort = "http://localhost:8080"                                      // Адрес и порт для коротких ID
	SaveFile          = "/tmp/short-url-db.json"                                     // Имя файла-хранилища
	PostgreDSN        = "postgres://tester:1234@localhost:5435/pgdb?sslmode=disable" // DSN по умолчанию
)

type Cfg struct {
	ShortIDLen        int
	IdleTimeout       int
	ServerPort        string
	ShortIDServerPort string
	SaveFilePath      string
	PostgreDSN        string
}

var config *Cfg

func parseFlags() {
	if config == nil {
		config = &Cfg{}
	}

	flag.IntVar(&config.ShortIDLen, "l", ShortIDLen,
		"Длина короткого ID – целое число от 8 до 32")
	flag.IntVar(&config.IdleTimeout, "t", IdleTimeout,
		"Таймаут сервера – целое число, желательно от 10 до 600")
	flag.StringVar(&config.ServerPort, "a", ServerPort,
		"Адрес запуска сервера в формате host:port или :port")
	flag.StringVar(&config.ShortIDServerPort, "b", ShortIDServerPort,
		"Адрес коротких ID в формате host:port/path")
	flag.StringVar(&config.SaveFilePath, "f", SaveFile,
		"Адрес места хранения файла")
	flag.StringVar(&config.PostgreDSN, "d", PostgreDSN,
		"DSN подключения к СУБД PostgreSQL")
	flag.Parse()

	if servPort, ok := os.LookupEnv("SERVER_ADDRESS"); ok && servPort != "" {
		config.ServerPort = servPort
	}

	if shortIDServPort, ok := os.LookupEnv("BASE_URL"); ok && shortIDServPort != "" {
		config.ShortIDServerPort = shortIDServPort
	}

	if saveFilePath, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok && saveFilePath != "" {
		config.SaveFilePath = saveFilePath
	}

	if postgreDSN, ok := os.LookupEnv("DATABASE_DSN"); ok && postgreDSN != "" {
		config.PostgreDSN = postgreDSN
	}
}

func InitConfig() *Cfg {
	if config != nil {
		return config
	}

	parseFlags()
	return config
}
