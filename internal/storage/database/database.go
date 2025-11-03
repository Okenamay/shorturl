package database

import (
	"context"
	"fmt"
	"time"

	"github.com/Okenamay/shorturl.git/internal/app/urlmaker"
	"github.com/Okenamay/shorturl.git/internal/config"
	"github.com/Okenamay/shorturl.git/internal/storage/memstorage"
	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func DeleteFlaggedURLs(ctx context.Context, appLogger *zap.SugaredLogger) error {
	appLogger.Info("DeleteFlaggedURLs starting")
	if DBPool == nil {
		return fmt.Errorf("database pool is not initialized")
	}

	res, err := DBPool.Exec(ctx, "DELETE FROM urls WHERE del_flag = true")
	if err != nil {
		appLogger.Errorw("DeleteFlaggedURLs stopped - delete flagged URLs FAIL", "error", err)
		return err
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected > 0 {
		appLogger.Infof("DeleteFlaggedURLs - deleted %d flagged URLs OK", rowsAffected)
	} else {
		appLogger.Warn("DeleteFlaggedURLs - no flagged URLs to delete")
	}

	return nil
}

var (
	DBPool *pgxpool.Pool
	err    error
)

type URLInfo struct {
	OriginalURL string
	IsDeleted   bool
}

func DBReinit(conf *config.Cfg, appLogger *zap.SugaredLogger) error {
	appLogger.Info("DBReinit started")

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
		appLogger.Errorw("DBReinit stopped - DB pool init FAIL", "error")
		return fmt.Errorf("DBReinit. DB pool is not initialized")
	}

	_, err := DBPool.Exec(context.Background(), sql)
	if err != nil {
		appLogger.Errorw("DBReinit stopped - create table FAIL", "error", err)
		return err
	}

	appLogger.Info("DBReinit finished - table created or already exists")
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

func BatchDelete(ctx context.Context, appLogger *zap.SugaredLogger, userID string, shortIDs []string) error {
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

	for range shortIDs {
		_, err := results.Exec()
		if err != nil {
			appLogger.Errorw("BatchDelete. Error in batch execution", "error", err)
		}
	}

	if closeErr := results.Close(); closeErr != nil {
		appLogger.Errorw("BatchDelete. Error closing batch results", "error", closeErr)
		return closeErr
	}

	return tx.Commit(ctx)
}

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

func StopDB(appLogger *zap.SugaredLogger) {
	if DBPool != nil {
		DBPool.Close()
		appLogger.Info("StopDB finished - DB connection pool closed")
	}
}

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

func AddOne(conf *config.Cfg, appLogger *zap.SugaredLogger, userID, shortID, fullURL string) (bool, error) {
	appLogger.Info("AddOne started")

	if DBPool == nil {
		appLogger.Error("AddOne stopped - DB pool not initialized")
		return false, fmt.Errorf("AddOne. DB pool is not initialized")
	}

	var exists bool
	err := DBPool.QueryRow(context.Background(),
		"SELECT EXISTS (SELECT 1 FROM urls WHERE url = $1)", fullURL).Scan(&exists)
	if err != nil {
		appLogger.Errorw("AddOne stopped - DB entry check FAIL", "error", err)
		return false, err
	}
	if exists {
		appLogger.Infof("AddOne - DB entry '%s' already exists", fullURL)
		return true, nil
	}

	_, errA := DBPool.Exec(context.Background(),
		"INSERT INTO urls (user_id, url, short_id) VALUES ($1, $2, $3)",
		userID, fullURL, shortID)
	if errA != nil {
		appLogger.Infof("AddOne - failed to add entry: user ID: '%s', URL '%s', ShortID '%s'",
			userID, fullURL, shortID)
		appLogger.Errorw("AddOne stopped - adding DB entry FAIL", "error", errA)
		return false, errA
	}

	appLogger.Infof("AddOne - DB entry successful: user ID: '%s', URL '%s', ShortID '%s'",
		userID, fullURL, shortID)
	return false, nil
}

func AddOneTransaction(ctx context.Context, tx pgx.Tx, appLogger *zap.SugaredLogger, userID, shortID, fullURL string) error {
	_, err := tx.Exec(ctx, "INSERT INTO urls (user_id, url, short_id) VALUES ($1, $2, $3) ON CONFLICT (url) DO NOTHING",
		userID, fullURL, shortID)
	if err != nil {
		appLogger.Errorw("AddOneTransaction stopped - adding DB entry FAIL", "error", err)
		return err
	}

	return nil
}

func GetUserURLs(conf *config.Cfg, appLogger *zap.SugaredLogger, userID string) ([]UserURL, error) {
	rows, err := DBPool.Query(context.Background(), "SELECT short_id, url FROM urls WHERE user_id = $1 AND del_flag = false", userID)
	if err != nil {
		appLogger.Errorw("GetUserURLs stopped - query FAIL", "error", err)
		return nil, err
	}
	defer rows.Close()

	var userURLs []UserURL
	var shortID string

	for rows.Next() {
		var u UserURL
		if err := rows.Scan(&shortID, &u.OriginalURL); err != nil {
			appLogger.Errorw("GetUserURLs - row scan FAIL", "error", err)
			return nil, err
		}

		u.ShortURL = urlmaker.MakeFullURL(conf, shortID)

		userURLs = append(userURLs, u)
	}

	if err := rows.Err(); err != nil {
		appLogger.Errorw("GetUserURLs - row iteration FAIL", "error", err)
		return nil, err
	}

	return userURLs, nil
}

type UserURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
