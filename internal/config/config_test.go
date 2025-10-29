package config

import (
	"flag"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
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
	}
}

func TestInitConfig(t *testing.T) {

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
				// Флаги -f и -d имеют default "", который перезаписывает
				// значения из констант (saveFile и postgreDSN)
				cfg.SaveFilePath = ""
				cfg.PostgreDSN = ""
				// MemMode вычисляется на основе финальных значений
				cfg.MemMode = "memstore"
				return cfg
			}(),
		},
		{
			name: "values from flags",
			env:  map[string]string{},
			args: []string{
				"cmd",
				"-a", "localhost:9090",
				"-b", "http://test.url",
				"-f", "/tmp/flag-db.json",
			},
			expected: func() *Cfg {
				cfg := defaultConfig()
				// Значения из флагов
				cfg.ServerPort = "localhost:9090"
				cfg.ShortIDAddress = "http://test.url"
				cfg.SaveFilePath = "/tmp/flag-db.json"
				// Флаг -d не был предоставлен, поэтому он получает свой default ""
				cfg.PostgreDSN = ""
				// MemMode вычисляется на основе финальных значений
				cfg.MemMode = "savefile"
				return cfg
			}(),
		},
		{
			name: "DATABASE_DSN_from_env_overrides_-f_flag",
			env: map[string]string{
				"DATABASE_DSN":      "postgres://env-user:env-pass@host/env-db",
				"FILE_STORAGE_PATH": "/tmp/flag-db.json",
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
				return cfg
			}(),
		},
	}

	// Сохраняем оригинальные os.Args, чтобы восстановить их после теста
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// --- Настройка ---

			// 1. Сбрасываем глобальные переменные
			// Это необходимо, чтобы sync.Once сработал снова
			config = nil
			once = sync.Once{}

			// 2. Сбрасываем состояние пакета flag
			// Это необходимо, т.к. flag.Parse() нельзя вызывать
			// несколько раз на одном и том же (глобальном) FlagSet.
			flag.CommandLine = flag.NewFlagSet(tc.args[0], flag.ExitOnError)

			// 3. Устанавливаем переменные окружения для этого теста
			// t.Setenv() автоматически очистит их после t.Run()
			for key, value := range tc.env {
				t.Setenv(key, value)
			}

			// 4. Устанавливаем os.Args для этого теста
			os.Args = tc.args

			// --- Выполнение ---
			got := InitConfig()

			assert.Equal(t, tc.expected, got)

			// --- Проверка ---
			// if !reflect.DeepEqual(tc.expected, got) {
			// 	// Выводим детальное сравнение в случае ошибки
			// 	t.Errorf("InitConfig() не совпадает:\nОжидалось: %+v\nПолучено:   %+v", tc.expected, got)
			// }
		})
	}
}
