package database

import (
	"context"
	"fmt"
	"time"

	"github.com/Okenamay/shorturl.git/internal/config"
	"github.com/Okenamay/shorturl.git/internal/storage/memstorage"
	"github.com/Okenamay/shorturl.git/internal/storage/migration"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var (
	DBPool *pgxpool.Pool
	err    error
)

// StartDB инициализирует пул соединений, загружает данные в память и запускает
// DBReinit по необходимости
func StartDB(conf *config.Cfg, appLogger *zap.SugaredLogger) error {
	appLogger.Info("StartDB started")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	DBPool, err = pgxpool.New(ctx, conf.PostgreDSN)
	if err != nil {
		appLogger.Errorw("StartDB - connect pgxpool FAIL", "error", err)
		return err
	}

	appLogger.Info("StartDB - pool init OK")

	if conf.DBReinitialize {
		if err := DBReinit(conf, appLogger); err != nil {
			return err
		}
	}

	rows, err := DBPool.Query(context.Background(), "SELECT url, short_id FROM urls")
	if err != nil {
		appLogger.Errorw("StartDB stopped - read table FAIL", "error", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var fullURL, shortID string
		if err := rows.Scan(&fullURL, &shortID); err != nil {
			appLogger.Errorw("StartDB stopped - row scan FAIL", "error", err)
			return err
		}

		memstorage.Store.Set(shortID, fullURL)
	}

	if err := rows.Err(); err != nil {
		appLogger.Errorw("StartDB stopped - iterate rows FAIL", "error", err)
		return err
	}

	appLogger.Infof("StartDB finished - loaded %d records into memory", len(memstorage.Store.GetAll()))

	return nil
}

// StopDB закрывает пул соединений с БД
func StopDB(appLogger *zap.SugaredLogger) {
	if DBPool != nil {
		DBPool.Close()
		appLogger.Info("StopDB finished - DB connection pool closed")
	}
}

// DBPing проверяет доступность БД
func DBPing(appLogger *zap.SugaredLogger) error {
	appLogger.Info("DBPing started")

	if DBPool == nil {
		appLogger.Error("DBPing stopped - DB pool not initialized")
		return fmt.Errorf("DBPing. DB pool not initialized")
	}

	err := DBPool.Ping(context.Background())
	if err != nil {
		appLogger.Errorw("DBPing stopped - DB ping FAIL", "error", err)
		return err
	}

	return nil
}

// DBReinit (пере)создает таблицу в БД
func DBReinit(conf *config.Cfg, appLogger *zap.SugaredLogger) error {
	appLogger.Info("DBReinit started - calling MigrateLauncher...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := migration.MigrateLauncher(ctx, conf, appLogger)
	if err != nil {
		appLogger.Errorw("DBReinit stopped - MigrateLauncher FAIL", "error", err)
		return err
	}

	appLogger.Info("DBReinit finished - migration complete")
	return nil
}
