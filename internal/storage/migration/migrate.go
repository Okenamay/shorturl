package migration

import (
	"context"
	"database/sql"
	"embed"

	"github.com/Okenamay/shorturl.git/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

// MigrateLauncher запускает процесс миграции с помощью goose
func MigrateLauncher(ctx context.Context, conf *config.Cfg, appLogger *zap.SugaredLogger) error {
	appLogger.Info("MigrateLauncher started")

	if conf.PostgreDSN == "" {
		appLogger.Warn("MigrateLauncher skipped - Database URI not set")
		return nil
	}

	db, err := sql.Open("pgx", conf.PostgreDSN)
	if err != nil {
		appLogger.Errorw("MigrateLauncher stopped - DB connection open FAIL", "error", err)
		return err
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		appLogger.Errorw("MigrateLauncher stopped - DB ping FAIL", "error", err)
		return err
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		appLogger.Errorw("MigrateLauncher stopped - set goose dialect FAIL", "error", err)
		return err
	}

	appLogger.Info("MigrateLauncher - running DB migrations...")

	if err := goose.UpContext(ctx, db, "migrations"); err != nil {
		appLogger.Errorw("MigrateLauncher stopped - apply migrations FAIL", "error", err)
		return err
	}

	appLogger.Info("MigrateLauncher finished - DB migrations OK")
	return nil
}
