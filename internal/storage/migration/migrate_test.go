package migration

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/Okenamay/shorturl.git/internal/config"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestMigrateLauncher(t *testing.T) {
	appLogger := zap.NewNop().Sugar()
	ctx := context.Background()

	t.Run("Successful 'up' Migration", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("failed to create pgxmock pool: %v", err)
		}
		defer mock.Close()

		conf := &config.Cfg{MigrateID: "20250720152000", MigrateDirection: "up"}
		migration := DeliverMigration(conf)

		mock.ExpectExec(regexp.QuoteMeta(migration.UpSQL)).
			WillReturnResult(pgxmock.NewResult("EXEC", 1))

		err = MigrateLauncher(ctx, mock, conf, appLogger)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Successful 'down' Migration", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("failed to create pgxmock pool: %v", err)
		}
		defer mock.Close()

		conf := &config.Cfg{MigrateID: "20250718110100", MigrateDirection: "down"}
		migration := DeliverMigration(conf)

		mock.ExpectExec(regexp.QuoteMeta(migration.DownSQL)).
			WillReturnResult(pgxmock.NewResult("EXEC", 1))

		err = MigrateLauncher(ctx, mock, conf, appLogger)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failed 'up' Migration", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("failed to create pgxmock pool: %v", err)
		}
		defer mock.Close()

		conf := &config.Cfg{MigrateID: "20250520160000", MigrateDirection: "up"}
		migration := DeliverMigration(conf)
		dbErr := errors.New("DB 'up' error")

		mock.ExpectExec(regexp.QuoteMeta(migration.UpSQL)).
			WillReturnError(dbErr)

		err = MigrateLauncher(ctx, mock, conf, appLogger)
		assert.Error(t, err)
		assert.Equal(t, dbErr, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failed 'down' Migration", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("failed to create pgxmock pool: %v", err)
		}
		defer mock.Close()

		conf := &config.Cfg{MigrateID: "20250720152000", MigrateDirection: "down"}
		migration := DeliverMigration(conf)
		dbErr := errors.New("DB 'down' error")

		mock.ExpectExec(regexp.QuoteMeta(migration.DownSQL)).
			WillReturnError(dbErr)

		err = MigrateLauncher(ctx, mock, conf, appLogger)
		assert.Error(t, err)
		assert.Equal(t, dbErr, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Migration Disabled (No ID)", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("failed to create pgxmock pool: %v", err)
		}
		defer mock.Close()

		conf := &config.Cfg{MigrateID: ""}
		err = MigrateLauncher(ctx, mock, conf, appLogger)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Unknown Migration ID", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("failed to create pgxmock pool: %v", err)
		}
		defer mock.Close()

		conf := &config.Cfg{MigrateID: "unknown_id", MigrateDirection: "up"}
		err = MigrateLauncher(ctx, mock, conf, appLogger)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Invalid Direction", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("failed to create pgxmock pool: %v", err)
		}
		defer mock.Close()

		conf := &config.Cfg{MigrateID: "20250720152000", MigrateDirection: "sideways"}
		err = MigrateLauncher(ctx, mock, conf, appLogger)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// --- Бенчмарки ---

func BenchmarkMigrateLauncher(b *testing.B) {
	appLogger := zap.NewNop().Sugar()
	ctx := context.Background()

	// Настраиваем мок один раз
	mock, err := pgxmock.NewPool()
	if err != nil {
		b.Fatalf("failed to create pgxmock pool: %v", err)
	}
	defer mock.Close()

	// Настраиваем конфиги для разных сценариев
	confUp := &config.Cfg{MigrateID: "20250720152000", MigrateDirection: "up"}
	migrationUp := DeliverMigration(confUp)

	confDown := &config.Cfg{MigrateID: "20250718110100", MigrateDirection: "down"}
	migrationDown := DeliverMigration(confDown)

	confDisabled := &config.Cfg{MigrateID: ""}

	b.Run("Successful 'up' Migration", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Настраиваем ожидания мока внутри цикла
			mock.ExpectExec(regexp.QuoteMeta(migrationUp.UpSQL)).
				WillReturnResult(pgxmock.NewResult("EXEC", 1))

			_ = MigrateLauncher(ctx, mock, confUp, appLogger)

			// Проверяем ожидания, чтобы бенчмарк был честным
			if err := mock.ExpectationsWereMet(); err != nil {
				b.Fatalf("expectations not met: %v", err)
			}
		}
	})

	b.Run("Successful 'down' Migration", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			mock.ExpectExec(regexp.QuoteMeta(migrationDown.DownSQL)).
				WillReturnResult(pgxmock.NewResult("EXEC", 1))

			_ = MigrateLauncher(ctx, mock, confDown, appLogger)

			if err := mock.ExpectationsWereMet(); err != nil {
				b.Fatalf("expectations not met: %v", err)
			}
		}
	})

	b.Run("Migration Disabled (No ID)", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Нет ожиданий, т.к. функция должна выйти до Exec
			_ = MigrateLauncher(ctx, mock, confDisabled, appLogger)

			if err := mock.ExpectationsWereMet(); err != nil {
				b.Fatalf("expectations not met: %v", err)
			}
		}
	})
}
