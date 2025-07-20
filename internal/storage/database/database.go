package database

import (
	"context"
	"fmt"
	"time"

	"github.com/Okenamay/shorturl.git/internal/app/urlmaker"
	"github.com/Okenamay/shorturl.git/internal/config"
	logger "github.com/Okenamay/shorturl.git/internal/logger/zap"
	"github.com/Okenamay/shorturl.git/internal/storage/memstorage"
	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	DBPool *pgxpool.Pool
	err    error
)

type URLInfo struct {
	OriginalURL string
	IsDeleted   bool
}

func DBReinit(conf *config.Cfg) error {
	logger.Zap.Info("DBReinit. Start")

	sql := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS urls (
		id BIGSERIAL PRIMARY KEY,
		user_id VARCHAR(36),
		url VARCHAR(1024) UNIQUE,
		short_id VARCHAR(%d) UNIQUE,
		del_flag BOOLEAN DEFAULT false
	);
	`, conf.ShortIDLen)

	if DBPool == nil {
		logger.Zap.Error("DBReinit. DB pool not initialized")
		return fmt.Errorf("DBReinit. DB pool is not initialized")
	}

	_, err := DBPool.Exec(context.Background(), sql)
	if err != nil {
		logger.Zap.Errorf("DBReinit. Failed to create table: %w", err)
		return err
	}

	logger.Zap.Info("DBReinit. Table created or already exists")
	return nil
}

func GetURLInfo(shortID string) (URLInfo, error) {
	var info URLInfo
	err := DBPool.QueryRow(context.Background(),
		"SELECT url, del_flag FROM urls WHERE short_id = $1", shortID).Scan(&info.OriginalURL, &info.IsDeleted)
	if err != nil {
		if err == pgx.ErrNoRows {
			return URLInfo{}, nil
		}
		return URLInfo{}, err
	}
	return info, nil
}

func BatchDelete(ctx context.Context, userID string, shortIDs []string) error {
	tx, err := DBPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Prepare(ctx, "batch_delete", "UPDATE urls SET del_flag = true WHERE user_id = $1 AND short_id = $2")
	if err != nil {
		return err
	}

	batch := &pgx.Batch{}
	for _, shortID := range shortIDs {
		batch.Queue("batch_delete", userID, shortID)
	}

	results := tx.SendBatch(ctx, batch)
	defer results.Close()

	for range shortIDs {
		_, err := results.Exec()
		if err != nil {
			logger.Zap.Errorw("BatchDelete. Error in batch execution", "error", err)
		}
	}

	return tx.Commit(ctx)
}

func StartDB(conf *config.Cfg) error {
	logger.Zap.Info("StartDB. Start")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	DBPool, err = pgxpool.New(ctx, conf.PostgreDSN)
	if err != nil {
		logger.Zap.Errorw("StartDB. Failed to connect pgxpool", "error", err)
		return err
	}

	logger.Zap.Info("StartDB. Pool initialized.")

	if conf.DBReinitialize {
		if err := DBReinit(conf); err != nil {
			return err
		}
	}

	rows, err := DBPool.Query(context.Background(), "SELECT url, short_id FROM urls")
	if err != nil {
		logger.Zap.Errorw("StartDB. Failed to read table", "error", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var fullURL, shortID string
		if err := rows.Scan(&fullURL, &shortID); err != nil {
			logger.Zap.Errorw("StartDB. Failed to scan row", "error", err)
			return err
		}

		memstorage.Store.Set(shortID, fullURL)
	}

	if err := rows.Err(); err != nil {
		logger.Zap.Errorf("StartDB. Failed to iterate rows: %w", err)
		return err
	}

	logger.Zap.Infof("StartDB. Loaded %d records into memory", len(memstorage.Store.GetAll()))

	return nil
}

func StopDB() {
	if DBPool != nil {
		DBPool.Close()
		logger.Zap.Info("StopDB. Database connection pool closed.")
	}
}

func DBPing() error {
	logger.Zap.Info("DBPing. Start")

	if DBPool == nil {
		logger.Zap.Error("DBPing. DB pool not initialized")
		return fmt.Errorf("DBPing. DB pool not initialized")
	}

	err := DBPool.Ping(context.Background())
	if err != nil {
		logger.Zap.Errorw("DBPing. DB ping error", "error", err)
		return err
	}

	return nil
}

func AddOne(conf *config.Cfg, userID, shortID, fullURL string) (bool, error) {
	logger.Zap.Info("AddOne. Start")

	if DBPool == nil {
		logger.Zap.Error("AddOne. DB pool is not initialized")
		return false, fmt.Errorf("AddOne. DB pool is not initialized")
	}

	var exists bool
	err := DBPool.QueryRow(context.Background(),
		"SELECT EXISTS (SELECT 1 FROM urls WHERE url = $1)", fullURL).Scan(&exists)
	if err != nil {
		logger.Zap.Errorw("AddOne. DB entry check error", "error", err)
		return false, err
	}
	if exists {
		logger.Zap.Infof("AddOne. DB entry '%s' already exists", fullURL)
		return true, nil
	}

	_, errA := DBPool.Exec(context.Background(),
		"INSERT INTO urls (user_id, url, short_id) VALUES ($1, $2, $3)",
		userID, fullURL, shortID)
	if errA != nil {
		logger.Zap.Infof("AddOne. Failed to add entry: user ID: '%s', URL '%s', ShortID '%s'",
			userID, fullURL, shortID)
		logger.Zap.Errorw("AddOne. Error adding DB entry", "error", errA)
		return false, errA
	}

	logger.Zap.Infof("AddOne. DB entry successful: user ID: '%s', URL '%s', ShortID '%s'",
		userID, fullURL, shortID)
	return false, nil
}

func AddOneTransaction(ctx context.Context, tx pgx.Tx, userID, shortID, fullURL string) error {
	_, err := tx.Exec(ctx, "INSERT INTO urls (user_id, url, short_id) VALUES ($1, $2, $3) ON CONFLICT (url) DO NOTHING",
		userID, fullURL, shortID)
	if err != nil {
		logger.Zap.Errorw("AddOneTransaction. Error adding DB entry", "error", err)
		return err
	}

	return nil
}

func GetUserURLs(conf *config.Cfg, userID string) ([]UserURL, error) {
	rows, err := DBPool.Query(context.Background(), "SELECT short_id, url FROM urls WHERE user_id = $1 AND del_flag = false", userID)
	if err != nil {
		logger.Zap.Errorw("GetUserURLs. Query failed", "error", err)
		return nil, err
	}
	defer rows.Close()

	var userURLs []UserURL
	var shortID string

	for rows.Next() {
		var u UserURL
		if err := rows.Scan(&shortID, &u.OriginalURL); err != nil {
			logger.Zap.Errorw("GetUserURLs. Row scan failed", "error", err)
			return nil, err
		}

		u.ShortURL = urlmaker.MakeFullURL(conf, shortID)

		userURLs = append(userURLs, u)
	}

	if err := rows.Err(); err != nil {
		logger.Zap.Errorw("GetUserURLs. Row iteration error", "error", err)
		return nil, err
	}

	return userURLs, nil
}

type UserURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
