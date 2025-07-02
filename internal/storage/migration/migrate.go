package migration

import (
	"context"

	"github.com/Okenamay/shorturl.git/internal/config"
	logger "github.com/Okenamay/shorturl.git/internal/logger/zap"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Run executes the migration logic based on the provided configuration and direction.
// direction should be "up" to apply or "down" to roll back.
func MigrateLauncher(ctx context.Context, dbpool *pgxpool.Pool, conf *config.Cfg) error {
	sugar, err := logger.InitLogger()
	if err != nil {
		sugar.Errorw(err.Error(), "MigrateLauncher", "Start logger")
	}
	sugar.Info("MigrateLauncher. Start")

	if conf.MigrateID == "" {
		sugar.Info("Migration disabled. Skipping.")
		return nil
	}

	sugar.Infof("Attempting migration for ID: %s, direction: %s",
		conf.MigrateID, conf.MigrateDirection)

	migration := DeliverMigration(conf)
	if migration.ID == "" {
		sugar.Infof("Unknown migration ID: %s", conf.MigrateID)
		return nil
	}

	switch conf.MigrateDirection {
	case "up":
		_, err := dbpool.Exec(ctx, migration.UpSQL)
		if err != nil {
			sugar.Errorf("Migration ID: %s failed: %v", migration.ID, err)
			return err
		}

		sugar.Infof("Successfully applied migration: %s", migration.ID)
		return nil
	case "down":
		_, err := dbpool.Exec(ctx, migration.DownSQL)
		if err != nil {
			sugar.Errorf("Rollback ID: %s failed: %v", migration.ID, err)
			return err
		}

		sugar.Infof("Successfully applied rollback: %s", migration.ID)
		return nil
	default:
		sugar.Infof("Incorrect migration direction: %s", conf.MigrateDirection)
		return nil
	}

}
