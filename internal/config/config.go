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
	// HTTPS по умолчанию выключен
	enableHTTPS = false
)

// Cfg определяет структуру конфигурации всего приложения.
// Поля заполняются из командной строки (флаги) и переменных окружения.
type Cfg struct {
	ShortIDLen       int    // Длина генерируемого короткого URL
	IdleTimeout      int    // Таймаут неактивности сервера
	ServerPort       string // Адрес и порт для запуска сервера (e.g., ":8080")
	ShortIDAddress   string // Базовый адрес для сокращенных URL (e.g., "http://localhost:8080")
	SaveFilePath     string // Путь к файлу для хранения URL (если используется)
	PostgreDSN       string // DSN для подключения к PostgreSQL
	MemMode          string // Режим хранения ("postgres", "savefile", "memstore")
	MigrateID        string // (Устарело, используется goose)
	MigrateDirection string // (Устарело, используется goose)
	DBReinitialize   bool   // Флаг, указывающий на необходимость запуска миграций
	AuthorizationKey string // Секретный ключ для подписи JWT-токенов
	TokenExpiry      int    // Срок жизни JWT-токена в часах
	AuditFile        string // Путь к файлу для логов аудита
	AuditURL         string // URL для отправки логов аудита
	EnableHTTPS      bool   // Режим HTTPS
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
	config.DBReinitialize = dbReinit
	config.AuthorizationKey = authKey
	config.TokenExpiry = tokenExp
	config.AuditFile = audFile
	config.AuditURL = audURL
	config.EnableHTTPS = enableHTTPS

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
	flag.BoolVar(&config.EnableHTTPS, "s", config.EnableHTTPS,
		"Включить HTTPS")

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
		}
	}
	if auditFile, ok := os.LookupEnv("AUDIT_FILE"); ok {
		config.AuditFile = auditFile
	}
	if auditURL, ok := os.LookupEnv("AUDIT_URL"); ok {
		config.AuditURL = auditURL
	}
	if enableHTTPS, ok := os.LookupEnv("ENABLE_HTTPS"); ok {
		config.EnableHTTPS = (enableHTTPS == "true" || enableHTTPS == "1")
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
