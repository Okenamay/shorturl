package config

import (
	"flag"
	"os"
	"strconv"
	"sync"
)

// Дефолтные значения до применения флагов:
const (
	// Длина короткого идентификатора
	shortIDLen = 10
	// Таймаут сервера в секундах
	idleTimeout = 600
	// Адрес и порт сервера
	serverPort = ":8080"
	// Адрес и порт для коротких ID
	shortIDAddr = "http://localhost:8080"
	// Имя файла-хранилища
	saveFile = "/tmp/short-url-db.json"
	// DSN по умолчанию
	postgreDSN = "postgresql://tester:1234@localhost:5432/pgdb"
	// // Заглушка
	// migrID = ""
	// // Заглушка
	// migrDir = ""
	// Флаг переинициализации БД при старте
	dbReinit = true
	// Ключ авторизации
	authKey = "secret_key"
	// Срок истечения действия токена
	tokenExp = 24
	// Путь к файлу-приёмнику, в который сохраняются логи аудита
	audFile = ""
	// Полный URL удаленного сервера-приёмника, куда отправляются логи аудита.
	audURL = ""
)

// Cfg определяет структуру конфигурации
type Cfg struct {
	ShortIDLen       int
	IdleTimeout      int
	ServerPort       string
	ShortIDAddress   string
	SaveFilePath     string
	PostgreDSN       string
	MemMode          string
	MigrateID        string
	MigrateDirection string
	DBReinitialize   bool
	AuthorizationKey string
	TokenExpiry      int
	AuditFile        string
	AuditURL         string
}

func parseFlags() *Cfg {
	config := &Cfg{}

	// Инициализируцемся дефолтными значениями:
	config.ShortIDLen = shortIDLen
	config.IdleTimeout = idleTimeout
	config.ServerPort = serverPort
	config.ShortIDAddress = shortIDAddr
	config.SaveFilePath = saveFile
	config.PostgreDSN = postgreDSN
	// config.MigrateID = migrID
	// config.MigrateDirection = migrDir
	config.DBReinitialize = dbReinit
	config.AuthorizationKey = authKey
	config.TokenExpiry = tokenExp
	config.AuditFile = audFile
	config.AuditURL = audURL

	// Преписываем всё флагами:
	flag.IntVar(&config.ShortIDLen, "l", config.ShortIDLen,
		"Длина короткого ID – целое число от 8 до 32")
	flag.IntVar(&config.IdleTimeout, "t", config.IdleTimeout,
		"Таймаут сервера – целое число, желательно от 10 до 600")
	flag.StringVar(&config.ServerPort, "a", config.ServerPort,
		"Адрес запуска сервера в формате host:port или :port")
	flag.StringVar(&config.ShortIDAddress, "b", config.ShortIDAddress,
		"Адрес коротких ID в формате host:port/path")
	flag.StringVar(&config.SaveFilePath, "f", "",
		"Адрес места хранения файла")
	flag.StringVar(&config.PostgreDSN, "d", "",
		"DSN подключения к СУБД PostgreSQL")
	// flag.StringVar(&config.MigrateID, "migid", migrID,
	// 	"ID миграции БД в формате YYYYMMDDHHMMSS")
	// flag.StringVar(&config.MigrateDirection, "migdir", migrDir,
	// 	"Направление миграции БД (up = миграция, down = роллбек)")
	flag.BoolVar(&config.DBReinitialize, "dbx", config.DBReinitialize,
		"Реинициализация БД (bool)")
	flag.StringVar(&config.AuthorizationKey, "k", config.AuthorizationKey,
		"Ключ для генерации JWT-токена")
	flag.IntVar(&config.TokenExpiry, "txp", config.TokenExpiry,
		"Срок истечения токена, часов")
	flag.StringVar(&config.AuditFile, "audit-file", config.AuditFile,
		"Путь к файлу, в который сохраняются логи аудита")
	flag.StringVar(&config.AuditURL, "audit-url", config.AuditURL,
		"Полный URL удаленного сервера, куда отправляются логи аудита")

	flag.Parse()

	// Переписываем дефолтные env'ами:
	if servPort, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		config.ServerPort = servPort
	}
	if shortIDServPort, ok := os.LookupEnv("BASE_URL"); ok {
		config.ShortIDAddress = shortIDServPort
	}
	if saveFilePath, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		config.SaveFilePath = saveFilePath
	}
	if postgreDSN, ok := os.LookupEnv("DATABASE_DSN"); ok {
		config.PostgreDSN = postgreDSN
	}
	// if migrateID, ok := os.LookupEnv("MIGRATION_ID"); ok {
	// 	config.MigrateID = migrateID
	// }
	// if migrateDirection, ok := os.LookupEnv("MIGRATION_DIRECTION"); ok {
	// 	config.MigrateDirection = migrateDirection
	// }
	if dbReinitialize, ok := os.LookupEnv("DB_REINIT"); ok {
		config.DBReinitialize = (dbReinitialize == "true")
	}
	if authorizationKey, ok := os.LookupEnv("AUTH_SECRET_KEY"); ok {
		config.AuthorizationKey = authorizationKey
	}
	if tokenExpiryStr, ok := os.LookupEnv("TOKEN_EXPIRY"); ok {
		tokenExpiry, err := strconv.Atoi(tokenExpiryStr)
		if err == nil {
			config.TokenExpiry = tokenExpiry
		} else {
		}
	}
	if auditFile, ok := os.LookupEnv("AUDIT_FILE"); ok {
		config.AuditFile = auditFile
	}
	if auditURL, ok := os.LookupEnv("AUDIT_URL"); ok {
		config.AuditURL = auditURL
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

	return config
}

var (
	once   sync.Once
	config *Cfg
)

func InitConfig() *Cfg {
	once.Do(func() {
		config = parseFlags()
	})
	return config
}
