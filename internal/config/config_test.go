package config

import (
	"flag"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func defaultConfig() *Cfg {
	return &Cfg{
		ShortIDLen:       shortIDLen,
		IdleTimeout:      idleTimeout,
		ServerPort:       serverPort,
		ShortIDAddress:   shortIDAddr,
		SaveFilePath:     saveFile,
		PostgreDSN:       postgreDSN,
		DBReinitialize:   dbReinit,
		AuthorizationKey: authKey,
		TokenExpiry:      tokenExp,
		AuditFile:        audFile,
		AuditURL:         audURL,
		EnableHTTPS:      enableHTTPS,
	}
}

func TestInitConfig(t *testing.T) {
	// Создаем временный файл конфига для тестов
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.json")
	configData := []byte(`{
		"server_address": "localhost:7000",
		"base_url": "http://config.url",
		"enable_https": true,
		"file_storage_path": "/tmp/config-db.json"
	}`)
	require.NoError(t, os.WriteFile(configFile, configData, 0644))

	// Определяем тестовые сценарии
	testCases := []struct {
		name     string
		args     []string
		env      map[string]string
		expected *Cfg
	}{
		{
			name: "default_values_(no_env,_no_flags)",
			args: []string{"cmd"},
			env:  map[string]string{},
			expected: func() *Cfg {
				cfg := defaultConfig()
				cfg.SaveFilePath = ""
				cfg.PostgreDSN = ""
				cfg.MemMode = "memstore"
				cfg.EnableHTTPS = false
				return cfg
			}(),
		},
		{
			name: "config_file_values",
			args: []string{"cmd", "-c", configFile},
			env:  map[string]string{},
			expected: func() *Cfg {
				cfg := defaultConfig()
				// Значения из файла
				cfg.ServerPort = "localhost:7000"
				cfg.ShortIDAddress = "http://config.url"
				cfg.EnableHTTPS = true
				cfg.SaveFilePath = "/tmp/config-db.json"
				cfg.PostgreDSN = ""
				cfg.ConfigPath = configFile
				cfg.MemMode = "savefile"
				return cfg
			}(),
		},
		{
			name: "values_from_flags",
			env:  map[string]string{},
			args: []string{
				"cmd",
				"-a", "localhost:9090",
				"-b", "http://test.url",
				"-f", "/tmp/flag-db.json",
				"-s",
			},
			expected: func() *Cfg {
				cfg := defaultConfig()
				// Значения из флагов
				cfg.ServerPort = "localhost:9090"
				cfg.ShortIDAddress = "http://test.url"
				cfg.SaveFilePath = "/tmp/flag-db.json"
				cfg.PostgreDSN = ""
				cfg.MemMode = "savefile"
				cfg.EnableHTTPS = true
				return cfg
			}(),
		},
		{
			name: "flag_overrides_config_file",
			args: []string{"cmd", "-config", configFile, "-a", ":9090"},
			env:  map[string]string{},
			expected: func() *Cfg {
				cfg := defaultConfig()
				// Флаг -a перебивает файл
				cfg.ServerPort = ":9090"
				// Остальное из файла
				cfg.ShortIDAddress = "http://config.url"
				cfg.EnableHTTPS = true
				cfg.SaveFilePath = "/tmp/config-db.json"
				cfg.PostgreDSN = ""
				cfg.ConfigPath = configFile
				cfg.MemMode = "savefile"
				return cfg
			}(),
		},
		{
			name: "DATABASE_DSN_from_env_overrides_-f_flag",
			env: map[string]string{
				"DATABASE_DSN":      "postgres://env-user:env-pass@host/env-db",
				"FILE_STORAGE_PATH": "/tmp/flag-db.json",
				"ENABLE_HTTPS":      "true",
			},
			args: []string{
				"cmd",
				"-f", "/tmp/ignored-path.json", // Этот флаг будет перезаписан env
			},
			expected: func() *Cfg {
				cfg := defaultConfig()
				// 1. Сначала применяется значение из флага
				cfg.SaveFilePath = "/tmp/ignored-path.json"
				// 2. Флаг -d не задан, получает default ""
				cfg.PostgreDSN = ""
				// 3. Затем env перезаписывает значения
				cfg.SaveFilePath = "/tmp/flag-db.json"
				cfg.PostgreDSN = "postgres://env-user:env-pass@host/env-db"
				// 4. MemMode вычисляется
				cfg.MemMode = "postgres"
				cfg.EnableHTTPS = true
				return cfg
			}(),
		},
		{
			name: "SERVER_ADDRESS_from_env_overrides_-a_flag",
			env: map[string]string{
				"SERVER_ADDRESS": "localhost:9999",
			},
			args: []string{
				"cmd",
				"-a", "localhost:1111", // Этот флаг будет перезаписан env
			},
			expected: func() *Cfg {
				cfg := defaultConfig()
				// 1. Сначала применяется значение из флага
				cfg.ServerPort = "localhost:1111"
				// 2. Флаги -f и -d не заданы, получают default ""
				cfg.SaveFilePath = ""
				cfg.PostgreDSN = ""
				// 3. Затем env перезаписывает значение
				cfg.ServerPort = "localhost:9999"
				// 4. MemMode вычисляется
				cfg.MemMode = "memstore"
				cfg.EnableHTTPS = false
				return cfg
			}(),
		},
		{
			name: "env_overrides_flag_and_config",
			args: []string{"cmd", "-c", configFile, "-a", ":9090"},
			env: map[string]string{
				"SERVER_ADDRESS": ":5000",
			},
			expected: func() *Cfg {
				cfg := defaultConfig()
				// Env перебивает Флаг и Файл
				cfg.ServerPort = ":5000"
				// Остальное из файла
				cfg.ShortIDAddress = "http://config.url"
				cfg.EnableHTTPS = true
				cfg.SaveFilePath = "/tmp/config-db.json"
				cfg.PostgreDSN = ""
				cfg.ConfigPath = configFile
				cfg.MemMode = "savefile"
				return cfg
			}(),
		},
		{
			name: "config_path_from_env",
			args: []string{"cmd"},
			env: map[string]string{
				"CONFIG": configFile,
			},
			expected: func() *Cfg {
				cfg := defaultConfig()
				// Значения из файла
				cfg.ServerPort = "localhost:7000"
				cfg.ShortIDAddress = "http://config.url"
				cfg.EnableHTTPS = true
				cfg.SaveFilePath = "/tmp/config-db.json"
				cfg.PostgreDSN = ""
				cfg.ConfigPath = configFile
				cfg.MemMode = "savefile"
				return cfg
			}(),
		},
	}

	// Сохраняем оригинальные os.Args, чтобы восстановить их после теста
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 1. Сбрасываем глобальные переменные
			config = nil
			initErr = nil
			once = sync.Once{}

			// 2. Сбрасываем состояние пакета flag
			flag.CommandLine = flag.NewFlagSet(tc.args[0], flag.ExitOnError)

			// 3. Устанавливаем переменные окружения
			for key, value := range tc.env {
				t.Setenv(key, value)
			}

			// 4. Устанавливаем os.Args
			os.Args = tc.args

			got, err := InitConfig()
			require.NoError(t, err)

			assert.Equal(t, tc.expected, got)
		})
	}
}
