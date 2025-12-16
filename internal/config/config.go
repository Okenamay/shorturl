package config

import (
	"encoding/json"
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
	// Порт gRPC по умолчанию
	grpcPort = ":3200"
)

// Cfg определяет структуру конфигурации всего приложения.
// Поля заполняются из командной строки (флаги) и переменных окружения
type Cfg struct {
	ShortIDLen       int    // Длина генерируемого короткого URL
	IdleTimeout      int    // Таймаут неактивности сервера
	ServerPort       string // Адрес и порт для запуска сервера (e.g., ":8080")
	ShortIDAddress   string // Базовый адрес для сокращенных URL (e.g., "http://localhost:8080")
	SaveFilePath     string // Путь к файлу для хранения URL (если используется)
	PostgreDSN       string // DSN для подключения к PostgreSQL
	MemMode          string // Режим хранения ("postgres", "savefile", "memstore")
	DBReinitialize   bool   // Флаг, указывающий на необходимость запуска миграций
	AuthorizationKey string // Секретный ключ для подписи JWT-токенов
	TokenExpiry      int    // Срок жизни JWT-токена в часах
	AuditFile        string // Путь к файлу для логов аудита
	AuditURL         string // URL для отправки логов аудита
	EnableHTTPS      bool   // Режим HTTPS
	TrustedSubnet    string // Доверенная подсеть (CIDR)
	ConfigPath       string // Путь к файлу конфигурации
	GRPCAddress      string // Адрес запуска gRPC сервера
}

// fileConfig описывает структуру JSON-файла конфигурации. Используем
// указатели, чтобы различать отсутствие значения и zero-value
type fileConfig struct {
	ServerAddress    *string `json:"server_address"`
	BaseURL          *string `json:"base_url"`
	FileStoragePath  *string `json:"file_storage_path"`
	DatabaseDSN      *string `json:"database_dsn"`
	EnableHTTPS      *bool   `json:"enable_https"`
	ShortIDLen       *int    `json:"short_id_len"`
	IdleTimeout      *int    `json:"idle_timeout"`
	DBReinitialize   *bool   `json:"db_reinitialize"`
	AuthorizationKey *string `json:"authorization_key"`
	TokenExpiry      *int    `json:"token_expiry"`
	AuditFile        *string `json:"audit_file"`
	AuditURL         *string `json:"audit_url"`
	TrustedSubnet    *string `json:"trusted_subnet"`
	GRPCAddress      *string `json:"grpc_address"`
}

func parseFlags() (*Cfg, error) {
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
	config.GRPCAddress = grpcPort

	// Преписываем всё флагами:
	flag.IntVar(&config.ShortIDLen, "l", config.ShortIDLen,
		"Длина короткого ID – целое число от 8 до 32")
	flag.IntVar(&config.IdleTimeout, "tio", config.IdleTimeout,
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
	flag.StringVar(&config.TrustedSubnet, "t", "",
		"Строковое представление бесклассовой адресации (CIDR)")
	flag.StringVar(&config.GRPCAddress, "g", config.GRPCAddress,
		"Адрес запуска gRPC сервера")

	// Флаги работы через файл конфигурации
	flag.StringVar(&config.ConfigPath, "c", "", "Путь к файлу конфигурации")
	flag.StringVar(&config.ConfigPath, "config", "", "Путь к файлу конфигурации")

	flag.Parse()

	// Собираем информацию о флагах, явно установленых пользователем, чтобы
	// конфиг из файла не перезатирал явно переданные флаги
	setFlags := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) {
		setFlags[f.Name] = true
	})

	// Определяем путь к конфигу - еслифлаг не задан, пробуем Env
	if config.ConfigPath == "" {
		if cfgPathEnv, ok := os.LookupEnv("CONFIG"); ok {
			config.ConfigPath = cfgPathEnv
		}
	}

	// Если путь к конфигу есть, загружаем и применяем
	if config.ConfigPath != "" {
		if err := loadConfigFromFile(config.ConfigPath, config, setFlags); err != nil {
			return nil, err
		}
	}

	// Переписываем значения - дефолтные и из файла конфигурации env'ами:
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
	if trustedSubnet, ok := os.LookupEnv("TRUSTED_SUBNET"); ok {
		config.TrustedSubnet = trustedSubnet
	}
	if grpcAddr, ok := os.LookupEnv("GRPC_PORT"); ok {
		config.GRPCAddress = grpcAddr
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

	return config, nil
}

var (
	once    sync.Once
	config  *Cfg
	initErr error
)

func InitConfig() (*Cfg, error) {
	once.Do(func() {
		config, initErr = parseFlags()
	})
	return config, initErr
}

func loadConfigFromFile(path string, cfg *Cfg, setFlags map[string]bool) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var fCfg fileConfig
	if err := json.NewDecoder(file).Decode(&fCfg); err != nil {
		return err
	}

	// Применяем значения из файла, если соответствующий флаг не был установлен
	// -a: ServerAddress
	if fCfg.ServerAddress != nil && !setFlags["a"] {
		cfg.ServerPort = *fCfg.ServerAddress
	}
	// -b: BaseURL
	if fCfg.BaseURL != nil && !setFlags["b"] {
		cfg.ShortIDAddress = *fCfg.BaseURL
	}
	// -f: FileStoragePath
	if fCfg.FileStoragePath != nil && !setFlags["f"] {
		cfg.SaveFilePath = *fCfg.FileStoragePath
	}
	// -d: DatabaseDSN
	if fCfg.DatabaseDSN != nil && !setFlags["d"] {
		cfg.PostgreDSN = *fCfg.DatabaseDSN
	}
	// -s: EnableHTTPS
	if fCfg.EnableHTTPS != nil && !setFlags["s"] {
		cfg.EnableHTTPS = *fCfg.EnableHTTPS
	}

	// Дополнительные поля
	if fCfg.ShortIDLen != nil && !setFlags["l"] {
		cfg.ShortIDLen = *fCfg.ShortIDLen
	}
	if fCfg.IdleTimeout != nil && !setFlags["t"] {
		cfg.IdleTimeout = *fCfg.IdleTimeout
	}
	if fCfg.DBReinitialize != nil && !setFlags["dbx"] {
		cfg.DBReinitialize = *fCfg.DBReinitialize
	}
	if fCfg.AuthorizationKey != nil && !setFlags["k"] {
		cfg.AuthorizationKey = *fCfg.AuthorizationKey
	}
	if fCfg.TokenExpiry != nil && !setFlags["txp"] {
		cfg.TokenExpiry = *fCfg.TokenExpiry
	}
	if fCfg.AuditFile != nil && !setFlags["audit-file"] {
		cfg.AuditFile = *fCfg.AuditFile
	}
	if fCfg.AuditURL != nil && !setFlags["audit-url"] {
		cfg.AuditURL = *fCfg.AuditURL
	}
	if fCfg.TrustedSubnet != nil && !setFlags["t"] {
		cfg.TrustedSubnet = *fCfg.TrustedSubnet
	}
	if fCfg.GRPCAddress != nil && !setFlags["g"] {
		cfg.GRPCAddress = *fCfg.GRPCAddress
	}

	return nil
}
