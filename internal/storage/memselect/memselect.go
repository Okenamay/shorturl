package memselect

import (
	"context"

	"github.com/Okenamay/shorturl.git/internal/app/urlmaker"
	"github.com/Okenamay/shorturl.git/internal/config"
	logger "github.com/Okenamay/shorturl.git/internal/logger/zap"
	"github.com/Okenamay/shorturl.git/internal/storage/database"
	"github.com/Okenamay/shorturl.git/internal/storage/memstorage"
	"github.com/Okenamay/shorturl.git/internal/storage/savefile"
	pgx "github.com/jackc/pgx/v5"
)

func MemInit(conf *config.Cfg) error {
	logger.Zap.Info("MemInit", "Assessing memory mode")

	var err error

	switch conf.MemMode {
	case "postgres":
		err = database.StartDB(conf)
		if err != nil {
			logger.Zap.Errorw(err.Error(), "MemInit", "Init DB")
			return err
		}
		logger.Zap.Info("MemInit", "Init DB OK")
	case "savefile":
		err = savefile.LoadFile(conf)
		if err != nil {
			logger.Zap.Errorw(err.Error(), "MemInit", "Load savefile")
			return err
		}
		logger.Zap.Info("MemInit", "Load savefile OK")
	case "memstore":
		logger.Zap.Info("MemInit", "Memstore OK")
	default:
		logger.Zap.Info("MemInit", "Wrong MemMode")
	}

	return nil
}

func MemStop(conf *config.Cfg) {
	logger.Zap.Info("MemStop", "Stopping memory")

	switch conf.MemMode {
	case "postgres":
		database.StopDB()
		logger.Zap.Info("MemStop", "Stop DB OK")
	default:
		logger.Zap.Info("MemStop", "Nothing to stop for this MemMode")
	}
}

func PingDB(conf *config.Cfg) (error, bool) {
	pingOK := false

	logger.Zap.Info("PingDB", "Pinging DB")

	switch conf.MemMode {
	case "postgres":
		err := database.DBPing()
		if err != nil {
			logger.Zap.Errorw(err.Error(), "PingDB", "Pinging DB")
			return err, false
		}
		pingOK = true
	default:
		pingOK = false
	}

	return nil, pingOK
}

func StorePair(conf *config.Cfg, shortID, fullURL string) (bool, error) {
	logger.Zap.Info("StorePair", "Running")

	_, alreadyExists := memstorage.Store.Get(shortID)

	switch conf.MemMode {
	case "postgres":
		dbExists, err := database.AddOne(conf, shortID, fullURL)
		if err != nil {
			logger.Zap.Errorw(err.Error(), "StorePair", "AddOne to DB")
			return false, err
		}

		memstorage.Store.Set(shortID, fullURL)

		logger.Zap.Info("StorePair. Save to DB OK")
		return dbExists, nil
	case "savefile":
		memstorage.Store.Set(shortID, fullURL)

		if err := savefile.SaveFile(conf); err != nil {
			logger.Zap.Errorw(err.Error(), "StorePair", "Save savefile")
			return false, err
		}

		logger.Zap.Info("StorePair. Save savefile OK")
		return alreadyExists, nil
	case "memstore":
		memstorage.Store.Set(shortID, fullURL)

		logger.Zap.Info("StorePair. Memstore OK")
		return alreadyExists, nil
	default:
		logger.Zap.Info("StorePair. Wrong MemMode")
		return false, nil
	}
}

func CheckPair(conf *config.Cfg, queryID string) (string, error) {
	fullURL, _ := memstorage.Store.Get(queryID)

	return fullURL, nil
}

type RequestEntry struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ResponseEntry struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func ProcessBatch(conf *config.Cfg, requestBatch []RequestEntry) ([]ResponseEntry, error) {
	logger.Zap.Info("ProcessBatch. Start")

	var responseBatch []ResponseEntry

	for _, entry := range requestBatch {
		logger.Zap.Infof("ProcessBatch. Processing entry with CorrelationID: %s",
			entry.CorrelationID)

		shortURL, shortID := urlmaker.ProcessURL(conf, entry.OriginalURL)

		responseBatch = append(responseBatch, ResponseEntry{
			CorrelationID: entry.CorrelationID,
			ShortURL:      shortURL,
		})

		if _, err := StorePair(conf, shortID, entry.OriginalURL); err != nil {
			logger.Zap.Errorw(err.Error(), "ProcessBatch", "StorePair from batch")
			return nil, err
		}
	}

	logger.Zap.Info("ProcessBatch. Batch processed")
	return responseBatch, nil
}

func ProcessBatchTransaction(conf *config.Cfg, requestBatch []RequestEntry) ([]ResponseEntry, error) {
	logger.Zap.Info("ProcessBatchTransaction. Start")

	var responseBatch []ResponseEntry
	ctx := context.Background()

	var tx pgx.Tx
	var err error
	if conf.MemMode == "postgres" {
		tx, err = database.DBPool.Begin(ctx)
		if err != nil {
			logger.Zap.Errorw("ProcessBatchTransaction. Failed to begin transaction", "error", err)
			return nil, err
		}
		defer tx.Rollback(ctx)
	}

	for _, entry := range requestBatch {
		shortURL, shortID := urlmaker.ProcessURL(conf, entry.OriginalURL)

		responseBatch = append(responseBatch, ResponseEntry{
			CorrelationID: entry.CorrelationID,
			ShortURL:      shortURL,
		})

		if err := StorePairTransaction(ctx, tx, conf, shortID, entry.OriginalURL); err != nil {
			logger.Zap.Errorw(err.Error(), "ProcessBatchTransaction", "StorePairTransaction from batch")
			return nil, err
		}
	}

	if conf.MemMode == "savefile" {
		if err := savefile.SaveFile(conf); err != nil {
			logger.Zap.Errorw(err.Error(), "ProcessBatchTransaction", "Save savefile")
			return nil, err
		}
	}

	if tx != nil {
		if err := tx.Commit(ctx); err != nil {
			logger.Zap.Errorw("ProcessBatchTransaction. Failed to commit transaction", "error", err)
			return nil, err
		}
	}

	logger.Zap.Info("ProcessBatchTransaction. Batch processed successfully.")
	return responseBatch, nil
}

func StorePairTransaction(ctx context.Context, tx pgx.Tx, conf *config.Cfg, shortID, fullURL string) error {
	logger.Zap.Info("StorePairTransaction. Start")

	memstorage.Store.Set(shortID, fullURL)

	var err error
	switch conf.MemMode {
	case "postgres":
		err = database.AddOneTransaction(ctx, tx, shortID, fullURL)
		if err != nil {
			logger.Zap.Errorw(err.Error(), "StorePairTransaction", "AddOneTransaction to DB")
			return err
		}
	case "savefile":
	case "memstore":
	default:
		logger.Zap.Info("StorePairTransaction", "Wrong MemMode")
	}

	return nil
}

func GetUserURLs(conf *config.Cfg, userID string) ([]database.UserURL, error) {
	if conf.MemMode != "postgres" {
		return nil, nil
	}
	return database.GetUserURLs(userID)
}

// Перед отправкой хочу сделать следующее
// 1) DB через интерфейсы
// 2) unit-тесты
