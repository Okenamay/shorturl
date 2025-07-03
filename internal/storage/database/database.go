package database

import (
	"context"
	"fmt"

	"github.com/Okenamay/shorturl.git/internal/config"
	logger "github.com/Okenamay/shorturl.git/internal/logger/zap"
	"github.com/Okenamay/shorturl.git/internal/storage/memstorage"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	DBPool      *pgxpool.Pool
	EntryExists bool
)

func StartDB(conf *config.Cfg) error {
	sugar, err := logger.InitLogger()
	if err != nil {
		sugar.Errorw(err.Error(), "StartDB", "Start logger")
	}
	sugar.Info("StartDB. Start")

	DBPool, err = pgxpool.New(context.Background(), conf.PostgreDSN)
	if err != nil {
		sugar.Errorw("StartDB. Failed to connect pgxpool", "error", err)
		return err
	}

	if conf.DBReinitialize {
		if err := DBReinit(conf); err != nil {
			return err
		}
	}

	rows, err := DBPool.Query(context.Background(), "SELECT url, short_id FROM urls")
	if err != nil {
		sugar.Errorw("StartDB. Failed to read table", "error", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var fullURL, shortID string
		if err := rows.Scan(&fullURL, &shortID); err != nil {
			sugar.Errorw("StartDB. Failed to scan row", "error", err)
			return err
		}
		memstorage.StoreURLIDPair(shortID, fullURL)
	}

	if err := rows.Err(); err != nil {
		sugar.Errorf("StartDB. Failed to iterate rows: %w", err)
		return err
	}

	sugar.Infof("StartDB. Loaded %d records into memory", len(memstorage.URLStore))

	return nil
}

func DBReinit(conf *config.Cfg) error {
	sugar, err := logger.InitLogger()
	if err != nil {
		sugar.Errorw(err.Error(), "DBReinit", "Start logger")
	}
	sugar.Info("DBReinit. Start")

	sql := fmt.Sprintf(`
    CREATE TABLE IF NOT EXISTS urls (
        id BIGSERIAL PRIMARY KEY,
        url VARCHAR(1024),
        short_id VARCHAR(%d)
    );
    CREATE UNIQUE INDEX IF NOT EXISTS idx_urls ON urls (url, short_id);
    `, conf.ShortIDLen)

	if DBPool == nil {
		sugar.Error("DBPing. DB pool not initialized")
		return nil
	}

	_, err = DBPool.Exec(context.Background(), sql)
	if err != nil {
		sugar.Errorf("DBReinit. Failed to create table: %w", err)
		return err
	}

	sugar.Info("DBReinit. Table created or already exists")
	return nil
}

func StopDB() {
	sugar, err := logger.InitLogger()
	if err != nil {
		sugar.Errorw(err.Error(), "StopDB", "Start logger")
	}
	if DBPool != nil {
		DBPool.Close()
		sugar.Info("StopDB. Database connection pool closed.")
	}
}

func DBPing() error {
	sugar, err := logger.InitLogger()
	if err != nil {
		sugar.Errorw(err.Error(), "DBPing", "Start logger")
	}
	sugar.Info("DBPing. Start")

	if DBPool == nil {
		sugar.Error("DBPing. DB pool not initialized")
		return nil
	}

	err = DBPool.Ping(context.Background())
	if err != nil {
		sugar.Errorw("DBPing. DB ping error", "error", err)
		return err
	}

	return nil
}

func AddOne(conf *config.Cfg, shortID, fullURL string) error {
	sugar, err := logger.InitLogger()
	if err != nil {
		sugar.Errorw(err.Error(), "AddOne", "Start logger")
	}
	sugar.Info("AddOne. Start")

	EntryExists = false

	if DBPool == nil {
		sugar.Error("AddOne. DB pool is not initialized")
		return nil
	}

	var exists bool
	err = DBPool.QueryRow(context.Background(), "SELECT EXISTS (SELECT 1 FROM urls WHERE url = $1)", fullURL).
		Scan(&exists)
	if err != nil {
		sugar.Errorw("AddOne. DB entry check error", "error", err)
		return err
	}
	if exists {
		sugar.Infof("AddOne. DB entry '%s' already exists", fullURL)
		EntryExists = true
		return nil
	}

	_, err = DBPool.Exec(context.Background(),
		"INSERT INTO urls (url, short_id) VALUES ($1, $2)",
		fullURL, shortID)
	if err != nil {
		sugar.Infof("AddOne. Failed to add entry: URL '%s', ShortID '%s'", fullURL, shortID)
		return err
	}

	sugar.Infof("AddOne. DB entry successful: URL '%s', ShortID '%s'", fullURL, shortID)

	return nil
}
