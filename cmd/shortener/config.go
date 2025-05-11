package main

import (
	"flag"
)

// Дефолтные значения до применения флагов:
const (
	ShortIDLen        = 10               // Длина короткого идентификатора
	IdleTimeout       = 600              // Таймаут сервера в секундах
	ServerPort        = ":8080"          // Адрес и порт сервера
	ShortIDServerPort = "localhost:8080" // Адрес и порт для коротких ID
)

var cfg struct {
	ShortIDLen        int
	IdleTimeout       int
	ServerPort        string
	ShortIDServerPort string
}

func parseFlags() {
	flag.IntVar(&cfg.ShortIDLen, "l", ShortIDLen,
		"Длина короткого ID – целое число от 8 до 32")
	flag.IntVar(&cfg.IdleTimeout, "t", IdleTimeout,
		"Таймаут сервера – целое число, желательно от 10 до 600")
	flag.StringVar(&cfg.ServerPort, "a", ServerPort,
		"Адрес запуска сервера в формате host:port или :port")
	flag.StringVar(&cfg.ShortIDServerPort, "b", ShortIDServerPort,
		"Адрес коротких ID в формате host:port/path")
	flag.Parse()
}
