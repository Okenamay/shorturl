package config

import (
	"flag"
	"os"
	"sync"

	logger "github.com/Okenamay/shorturl.git/internal/logger/zap"
)

// Дефолтные значения до применения флагов:
const (
	ShortIDLen  = 10                                             // Длина короткого идентификатора
	IdleTimeout = 600                                            // Таймаут сервера в секундах
	ServerPort  = ":8080"                                        // Адрес и порт сервера
	ShortIDAddr = "http://localhost:8080"                        // Адрес и порт для коротких ID
	SaveFile    = "/tmp/short-url-db.json"                       // Имя файла-хранилища
	PostgreDSN  = "postgresql://tester:1234@localhost:5432/pgdb" // DSN по умолчанию
	Verbose     = false                                          // Флаг детальности логов. !!! Временная заглушка
	MigrID      = ""                                             // "20250520160000"                               // Дефолтная миграция, заглушка
	MigrDir     = ""                                             // "up"                                           // Дефолтный роллбек, заглушка
	DBReinit    = true                                           // Флаг переинициализации БД при старте
	AuthKey     = "secret_key"                                   // Ключ авторизации.
)

type Cfg struct {
	ShortIDLen        int
	IdleTimeout       int
	ServerPort        string
	ShortIDServerPort string
	SaveFilePath      string
	PostgreDSN        string
	MemMode           string
	LogVerbose        bool
	MigrateID         string
	MigrateDirection  string
	DBReinitialize    bool
	AuthorizationKey  string
}

var useFile bool
var useDSN bool

func parseFlags() *Cfg {
	config := &Cfg{}

	flag.IntVar(&config.ShortIDLen, "l", ShortIDLen,
		"Длина короткого ID – целое число от 8 до 32")
	flag.IntVar(&config.IdleTimeout, "t", IdleTimeout,
		"Таймаут сервера – целое число, желательно от 10 до 600")
	flag.StringVar(&config.ServerPort, "a", ServerPort,
		"Адрес запуска сервера в формате host:port или :port")
	flag.StringVar(&config.ShortIDServerPort, "b", ShortIDAddr,
		"Адрес коротких ID в формате host:port/path")
	// Сделали дефолтным значением SaveFilePath "" – если флагом или переменной среды
	// не задали значение, то никуда не будем писать:
	flag.StringVar(&config.SaveFilePath, "f", "",
		"Адрес места хранения файла")
	// Аналогично с дефолтным значением PostgreDSN "" – если флагом или переменной среды
	// не задали значение, то DSN будет пустой:
	flag.StringVar(&config.PostgreDSN, "d", "",
		"DSN подключения к СУБД PostgreSQL")
	flag.BoolVar(&config.LogVerbose, "log", Verbose,
		"Вывод подробного лога (bool)")
	flag.StringVar(&config.MigrateID, "migid", MigrID,
		"ID миграции БД в формате YYYYMMDDHHMMSS")
	flag.StringVar(&config.MigrateDirection, "migdir", MigrDir,
		"Направление миграции БД (up = миграция, down = роллбек)")
	flag.BoolVar(&config.DBReinitialize, "dbx", DBReinit,
		"Реинициализация БД (bool)")
	flag.StringVar(&config.AuthorizationKey, "s", AuthKey,
		"Ключ для генерации JWT-токена")
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
		logger.Zap.Infof("EnvFilePath = %s", saveFilePath)
	}

	if postgreDSN, ok := os.LookupEnv("DATABASE_DSN"); ok && postgreDSN != "" {
		config.PostgreDSN = postgreDSN
		logger.Zap.Infof("EnvDSN = %s", postgreDSN)
	}

	if logVerbose, ok := os.LookupEnv("LOGGER_VERBOSE"); ok {
		switch logVerbose {
		case "true":
			config.LogVerbose = true
		default:
			config.LogVerbose = false
		}
		logger.Zap.Infof("EnvVerbose = %s", logVerbose)
	}

	if migrateID, ok := os.LookupEnv("MIGRATION_ID"); ok && migrateID != "" {
		config.MigrateID = migrateID
		logger.Zap.Infof("EnvMigrID = %s", migrateID)
	}

	if migrateDirection, ok := os.LookupEnv("MIGRATION_DIRECTION"); ok && migrateDirection != "" {
		config.MigrateDirection = migrateDirection
		logger.Zap.Infof("EnvMigrDir = %s", migrateDirection)
	}

	if dbReinitialize, ok := os.LookupEnv("LOGGER_VERBOSE"); ok {
		switch dbReinitialize {
		case "true":
			config.DBReinitialize = true
		default:
			config.DBReinitialize = false
		}
		logger.Zap.Infof("EnvVerbose = %s", dbReinitialize)
	}

	if authorizationKey, ok := os.LookupEnv("AUTH_SECRET_KEY"); ok && authorizationKey != "" {
		config.AuthorizationKey = authorizationKey
		logger.Zap.Infof("EnvKey = %s", authorizationKey)
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

	logger.Zap.Infof("config.SaveFilePath: %s. config.PostgreDSN: %s. "+
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
