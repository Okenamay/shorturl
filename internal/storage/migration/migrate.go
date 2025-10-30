package migration

import (
	"context"

	"github.com/Okenamay/shorturl.git/internal/config"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

type DBExecutor interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}

func MigrateLauncher(ctx context.Context, dbpool DBExecutor, conf *config.Cfg, appLogger *zap.SugaredLogger) error {
	appLogger.Info("MigrateLauncher started")

	if conf.MigrateID == "" {
		appLogger.Warn("MigrateLauncher stopped - migration disabled")
		return nil
	}

	appLogger.Infof("MigrateLauncher - attempting migration for ID: %s, direction: %s",
		conf.MigrateID, conf.MigrateDirection)

	migration := DeliverMigration(conf)
	if migration.ID == "" {
		appLogger.Warnf("MigrateLauncher stopped - unknown migration ID: %s", conf.MigrateID)
		return nil
	}

	switch conf.MigrateDirection {
	case "up":
		_, err := dbpool.Exec(ctx, migration.UpSQL)
		if err != nil {
			appLogger.Errorf("MigrateLauncher stopped - migration ID: %s FAIL", migration.ID, "error", err)
			return err
		}

		appLogger.Infof("MigrateLauncher finished - migration: %s OK", migration.ID)
		return nil
	case "down":
		_, err := dbpool.Exec(ctx, migration.DownSQL)
		if err != nil {
			appLogger.Errorf("MigrateLauncher stopped - rollback ID: %s FAIL", migration.ID, "error", err)
			return err
		}

		appLogger.Infof("MigrateLauncher finished - applied rollback: %s OK", migration.ID)
		return nil
	default:
		appLogger.Warnf("MigrateLauncher finished - incorrect migration direction: %s", conf.MigrateDirection)
		return nil
	}
}
