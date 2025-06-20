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
	DBPool *pgxpool.Pool
)

// Создадим таблицу в базе данных:
func StartDB(conf *config.Cfg) error {
	sugar, _ := logger.InitLogger()
	sugar.Info("StartDB. Start")

	dbPool, err := pgxpool.New(context.Background(), conf.PostgreDSN)
	if err != nil {
		sugar.Errorw("Init. Failed to connect pgxpool", "error", err)
		return err
	}
	defer dbPool.Close()

	sql := fmt.Sprintf(`
    CREATE TABLE IF NOT EXISTS urls (
        id BIGSERIAL PRIMARY KEY,
        url VARCHAR(1024),
        short_id VARCHAR(%d)
    );
    CREATE UNIQUE INDEX IF NOT EXISTS idx_urls ON urls (url, short_id);
    `, conf.ShortIDLen)

	_, err = dbPool.Exec(context.Background(), sql)
	if err != nil {
		sugar.Errorw("Init. Failed to create table", "error", err)
		return fmt.Errorf("ошибка при создании таблицы: %w", err)
	}

	sugar.Info("StartDB. Table created or already exists")

	rows, err := dbPool.Query(context.Background(), "SELECT url, short_id FROM urls")
	if err != nil {
		sugar.Errorw("Init. Failed to read table", "error", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var fullURL, shortID string
		if err := rows.Scan(&fullURL, &shortID); err != nil {
			sugar.Errorw("Init. Failed to scan row", "error", err)
			return err
		}
		memstorage.StoreURLIDPair(shortID, fullURL)
	}

	if err := rows.Err(); err != nil {
		sugar.Errorw("Init. Failed to iterate rows", "error", err)
		return fmt.Errorf("ошибка при обработке строк: %w", err)
	}

	sugar.Infof("StartDB. Loaded %d records into memory", len(memstorage.URLStore))

	return nil
}

func DBPing(conf *config.Cfg) error {
	sugar, _ := logger.InitLogger()
	sugar.Info("DBPing. Start")

	dbPool, err := pgxpool.New(context.Background(), conf.PostgreDSN)
	if err != nil {
		sugar.Errorw("DBPing. DB pool error", "error", err)
		return err
	}
	defer dbPool.Close()

	err = dbPool.Ping(context.Background())
	if err != nil {
		sugar.Errorw("DBPing. DB ping error", "error", err)
		return err
	}

	return nil
}

var EntryExists bool

func AddOne(conf *config.Cfg, shortID, fullURL string) error {
	sugar, _ := logger.InitLogger()
	sugar.Info("AddOne. Start")

	dbPool, err := pgxpool.New(context.Background(), conf.PostgreDSN)
	if err != nil {
		sugar.Errorw("AddOne. DB pool error", "error", err)
		return err
	}
	defer dbPool.Close()

	var exists bool
	err = dbPool.QueryRow(context.Background(), "SELECT EXISTS (SELECT 1 FROM urls WHERE url = $1)", fullURL).
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

	_, err = dbPool.Exec(context.Background(),
		"INSERT INTO urls (url, short_id) VALUES ($1, $2)",
		fullURL, shortID)
	if err != nil {
		sugar.Errorw("AddOne. Error adding DB entry", "error", err)
		sugar.Infof("AddOne. Failed to add entry: URL '%s', ShortID '%s'", fullURL, shortID)
		return err
	}

	sugar.Infof("AddOne. DB entry successful: URL '%s', ShortID '%s'", fullURL, shortID)

	return nil
}
