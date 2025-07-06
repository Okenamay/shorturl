package config

import (
	"flag"
	"os"
	"testing"

	logger "github.com/Okenamay/shorturl.git/internal/logger/zap"
	"github.com/stretchr/testify/require"
)

func TestParseFlags(t *testing.T) {
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	}()

	if err := logger.InitLogger(); err != nil {
		logger.Zap.Fatalw(err.Error(), "Main", "Start logger")
	}
	defer logger.Zap.Sync()

	testCases := []struct {
		name    string
		envVars map[string]string
		args    []string
		want    Cfg
	}{
		{
			name:    "Default_values_(no_env,_no_flags)",
			envVars: map[string]string{},
			args:    []string{},
			want: Cfg{
				ShortIDLen:        ShortIDLen,
				IdleTimeout:       IdleTimeout,
				ServerPort:        ServerPort,
				ShortIDServerPort: ShortIDAddr,
				SaveFilePath:      "",
				PostgreDSN:        "",
				MemMode:           "memstore",
				LogVerbose:        Verbose,
				MigrateID:         MigrID,
				MigrateDirection:  MigrDir,
				DBReinitialize:    DBReinit,
				AuthorizationKey:  AuthKey,
			},
		},
		{
			name:    "Values from flags",
			envVars: map[string]string{},
			args: []string{
				"-a", "localhost:9090",
				"-b", "http://test.url",
				"-f", "/tmp/flag-db.json",
			},
			want: Cfg{
				ShortIDLen:        ShortIDLen,
				IdleTimeout:       IdleTimeout,
				ServerPort:        "localhost:9090",
				ShortIDServerPort: "http://test.url",
				SaveFilePath:      "/tmp/flag-db.json",
				PostgreDSN:        "",
				MemMode:           "savefile",
				LogVerbose:        Verbose,
				MigrateID:         MigrID,
				MigrateDirection:  MigrDir,
				DBReinitialize:    DBReinit,
				AuthorizationKey:  AuthKey,
			},
		},
		{
			name: "DATABASE_DSN_from_env_overrides_-f_flag",
			envVars: map[string]string{
				"DATABASE_DSN": "postgres://env-user:env-pass@host/env-db",
			},
			args: []string{
				"-f", "/tmp/flag-db.json",
			},
			want: Cfg{
				ShortIDLen:        ShortIDLen,
				IdleTimeout:       IdleTimeout,
				ServerPort:        ServerPort,
				ShortIDServerPort: ShortIDAddr,
				SaveFilePath:      "/tmp/flag-db.json",
				PostgreDSN:        "postgres://env-user:env-pass@host/env-db",
				MemMode:           "postgres",
				LogVerbose:        Verbose,
				MigrateID:         MigrID,
				MigrateDirection:  MigrDir,
				DBReinitialize:    DBReinit,
				AuthorizationKey:  AuthKey,
			},
		},
		{
			name: "SERVER_ADDRESS_from_env_overrides_-a_flag",
			envVars: map[string]string{
				"SERVER_ADDRESS": "localhost:9999",
			},
			args: []string{
				"-a", "localhost:1111",
			},
			want: Cfg{
				ShortIDLen:        ShortIDLen,
				IdleTimeout:       IdleTimeout,
				ServerPort:        "localhost:9999",
				ShortIDServerPort: ShortIDAddr,
				SaveFilePath:      "",
				PostgreDSN:        "",
				MemMode:           "memstore",
				LogVerbose:        Verbose,
				MigrateID:         MigrID,
				MigrateDirection:  MigrDir,
				DBReinitialize:    DBReinit,
				AuthorizationKey:  AuthKey,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(tc.name, flag.ExitOnError)

			for key, value := range tc.envVars {
				t.Setenv(key, value)
			}

			os.Args = append([]string{"test"}, tc.args...)

			got := parseFlags()

			require.NotNil(t, got)
			require.Equal(t, tc.want, *got)
		})
	}
}
