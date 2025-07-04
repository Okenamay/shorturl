package migration

import (
	"context"

	"github.com/Okenamay/shorturl.git/internal/config"
	logger "github.com/Okenamay/shorturl.git/internal/logger/zap"
	"github.com/jackc/pgx/v5/pgxpool"
)

func MigrateLauncher(ctx context.Context, dbpool *pgxpool.Pool, conf *config.Cfg) error {
	logger.Zap.Info("MigrateLauncher. Start")

	if conf.MigrateID == "" {
		logger.Zap.Info("Migration disabled. Skipping.")
		return nil
	}

	logger.Zap.Infof("Attempting migration for ID: %s, direction: %s",
		conf.MigrateID, conf.MigrateDirection)

	migration := DeliverMigration(conf)
	if migration.ID == "" {
		logger.Zap.Infof("Unknown migration ID: %s", conf.MigrateID)
		return nil
	}

	switch conf.MigrateDirection {
	case "up":
		_, err := dbpool.Exec(ctx, migration.UpSQL)
		if err != nil {
			logger.Zap.Errorf("Migration ID: %s failed: %v", migration.ID, err)
			return err
		}

		logger.Zap.Infof("Successfully applied migration: %s", migration.ID)
		return nil
	case "down":
		_, err := dbpool.Exec(ctx, migration.DownSQL)
		if err != nil {
			logger.Zap.Errorf("Rollback ID: %s failed: %v", migration.ID, err)
			return err
		}

		logger.Zap.Infof("Successfully applied rollback: %s", migration.ID)
		return nil
	default:
		logger.Zap.Infof("Incorrect migration direction: %s", conf.MigrateDirection)
		return nil
	}

}
