package config

import (
	"flag"
	"os"
	"sync"

	logger "github.com/Okenamay/shorturl.git/internal/logger/zap"
)

// Дефолтные значения до применения флагов:
const (
	ShortIDLen        = 10                                             // Длина короткого идентификатора
	IdleTimeout       = 600                                            // Таймаут сервера в секундах
	ServerPort        = ":8080"                                        // Адрес и порт сервера
	ShortIDServerPort = "http://localhost:8080"                        // Адрес и порт для коротких ID
	SaveFile          = "/tmp/short-url-db.json"                       // Имя файла-хранилища
	PostgreDSN        = "postgresql://tester:1234@localhost:5432/pgdb" // DSN по умолчанию
)

type Cfg struct {
	ShortIDLen        int
	IdleTimeout       int
	ServerPort        string
	ShortIDServerPort string
	SaveFilePath      string
	PostgreDSN        string
	MemMode           string
}

var useFile bool
var useDSN bool

func parseFlags() *Cfg {
	sugar, err := logger.InitLogger()
	if err != nil {
		sugar.Errorw(err.Error(), "Main", "Start logger")
	}

	config := &Cfg{}

	flag.IntVar(&config.ShortIDLen, "l", ShortIDLen,
		"Длина короткого ID – целое число от 8 до 32")
	flag.IntVar(&config.IdleTimeout, "t", IdleTimeout,
		"Таймаут сервера – целое число, желательно от 10 до 600")
	flag.StringVar(&config.ServerPort, "a", ServerPort,
		"Адрес запуска сервера в формате host:port или :port")
	flag.StringVar(&config.ShortIDServerPort, "b", ShortIDServerPort,
		"Адрес коротких ID в формате host:port/path")
	// Сделали дефолтным значением SaveFilePath "" – если флагом или переменной среды
	// не задали значение, то никуда не будем писать:
	flag.StringVar(&config.SaveFilePath, "f", "",
		"Адрес места хранения файла")
	// Аналогично с дефолтным значением PostgreDSN "" – если флагом или переменной среды
	// не задали значение, то DSN будет пустой:
	flag.StringVar(&config.PostgreDSN, "d", "",
		"DSN подключения к СУБД PostgreSQL")
	flag.Parse()

	var saveFilePath, postgreDSN string

	if servPort, ok := os.LookupEnv("SERVER_ADDRESS"); ok && servPort != "" {
		config.ServerPort = servPort
	}

	if shortIDServPort, ok := os.LookupEnv("BASE_URL"); ok && shortIDServPort != "" {
		config.ShortIDServerPort = shortIDServPort
	}

	if saveFilePath, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok && saveFilePath != "" {
		config.SaveFilePath = saveFilePath
		sugar.Infof("EnvFilePath = %s", saveFilePath)
	}

	if postgreDSN, ok := os.LookupEnv("DATABASE_DSN"); ok && postgreDSN != "" {
		config.PostgreDSN = postgreDSN
		sugar.Infof("EnvDSN = %s", postgreDSN)
	}

	// Проверим режим работы с данными и сформируем соотвествующий индикатор,
	// проверять будем по порядку:
	if config.PostgreDSN != "" {
		config.MemMode = "postgres"
	} else if config.SaveFilePath != "" {
		config.MemMode = "savefile"
	} else {
		config.MemMode = "memstore"
	}

	sugar.Infof("config.SaveFilePath: %s. config.PostgreDSN: %s. "+
		"useDSN: %t. useFile: %t. saveFilePath: %s. "+
		"postgreDSN: %s.",
		config.SaveFilePath, config.PostgreDSN, useDSN, useFile,
		saveFilePath, postgreDSN)

	return config
}

func InitConfig() *Cfg {
	var (
		once   sync.Once
		config *Cfg
	)

	once.Do(func() {
		config = parseFlags()
	})

	return config
}
