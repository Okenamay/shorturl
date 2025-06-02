package config

import (
	"flag"
	"os"
)

// Дефолтные значения до применения флагов:
const (
	ShortIDLen        = 10                      // Длина короткого идентификатора
	IdleTimeout       = 600                     // Таймаут сервера в секундах
	ServerPort        = ":8080"                 // Адрес и порт сервера
	ShortIDServerPort = "http://localhost:8080" // Адрес и порт для коротких ID
	SaveFile          = "savefile.txt"          // Имя файла-хранилища
)

var Cfg struct {
	ShortIDLen        int
	IdleTimeout       int
	ServerPort        string
	ShortIDServerPort string
	SaveFilePath      string
}

func ParseFlags() {
	flag.IntVar(&Cfg.ShortIDLen, "l", ShortIDLen,
		"Длина короткого ID – целое число от 8 до 32")
	flag.IntVar(&Cfg.IdleTimeout, "t", IdleTimeout,
		"Таймаут сервера – целое число, желательно от 10 до 600")
	flag.StringVar(&Cfg.ServerPort, "a", ServerPort,
		"Адрес запуска сервера в формате host:port или :port")
	flag.StringVar(&Cfg.ShortIDServerPort, "b", ShortIDServerPort,
		"Адрес коротких ID в формате host:port/path")
	flag.StringVar(&Cfg.SaveFilePath, "f", SaveFile,
		"Адрес коротких ID в формате host:port/path")
	flag.Parse()

	if servPort, ok := os.LookupEnv("SERVER_ADDRESS"); ok && servPort != "" {
		Cfg.ServerPort = servPort
	}

	if shortIDServPort, ok := os.LookupEnv("BASE_URL"); ok && shortIDServPort != "" {
		Cfg.ShortIDServerPort = shortIDServPort
	}

	if saveFilePath, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok && saveFilePath != "" {
		Cfg.SaveFilePath = saveFilePath
	}
}
