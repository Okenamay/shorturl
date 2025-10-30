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
