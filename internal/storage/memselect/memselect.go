package memselect

import (
	"context"

	"github.com/Okenamay/shorturl.git/internal/app/urlmaker"
	"github.com/Okenamay/shorturl.git/internal/config"
	"github.com/Okenamay/shorturl.git/internal/storage/database"
	"github.com/Okenamay/shorturl.git/internal/storage/memstorage"
	"github.com/Okenamay/shorturl.git/internal/storage/savefile"
	pgx "github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func MemInit(conf *config.Cfg, appLogger *zap.SugaredLogger) error {
	appLogger.Info("MemInit started - assessing memory mode")

	var err error

	switch conf.MemMode {
	case "postgres":
		err = database.StartDB(conf, appLogger)
		if err != nil {
			appLogger.Errorw("MemInit stopped - init DB FAIL", "error", err)
			return err
		}
		appLogger.Info("MemInit finished - init DB OK")
	case "savefile":
		err = savefile.LoadFile(conf, appLogger)
		if err != nil {
			appLogger.Errorw("MemInit stopped - load savefile failed", "error", err)
			return err
		}
		appLogger.Info("MemInit finished - load savefile OK")
	case "memstore":
		appLogger.Info("MemInit finished - memstore OK")
	default:
		appLogger.Info("MemInit finished - wrong MemMode")
	}

	return nil
}

func MemStop(conf *config.Cfg, appLogger *zap.SugaredLogger) {
	appLogger.Info("MemStop started")

	switch conf.MemMode {
	case "postgres":
		database.StopDB(appLogger)
		appLogger.Info("MemStop finished - stop DB OK")
	default:
		appLogger.Info("MemStop finished - nothing to stop for this MemMode")
	}
}

func PingDB(conf *config.Cfg, appLogger *zap.SugaredLogger) (error, bool) {
	pingOK := false

	appLogger.Info("PingDB started")

	switch conf.MemMode {
	case "postgres":
		err := database.DBPing(appLogger)
		if err != nil {
			appLogger.Errorw("PingDB stopped - ping DB FAIL", "error", err)
			return err, false
		}
		pingOK = true
	default:
		pingOK = false
	}

	appLogger.Info("PingDB finished - pinging DB OK")
	return nil, pingOK
}

func StorePair(conf *config.Cfg, appLogger *zap.SugaredLogger, userID, shortID, fullURL string) (bool, error) {
	appLogger.Info("StorePair started")

	_, alreadyExists := memstorage.Store.Get(shortID)

	switch conf.MemMode {
	case "postgres":
		dbExists, err := database.AddOne(conf, appLogger, userID, shortID, fullURL)
		if err != nil {
			appLogger.Errorw("StorePair stopped - AddOne to DB FAIL", "error", err)
			return false, err
		}

		memstorage.Store.Set(shortID, fullURL)

		appLogger.Info("StorePair finished - save to DB OK")
		return dbExists, nil
	case "savefile":
		memstorage.Store.Set(shortID, fullURL)

		if err := savefile.SaveFile(conf, appLogger); err != nil {
			appLogger.Errorw("StorePair stopped - save savefile FAIL", "error", err)
			return false, err
		}

		appLogger.Info("StorePair finished - save savefile OK")
		return alreadyExists, nil
	case "memstore":
		memstorage.Store.Set(shortID, fullURL)

		appLogger.Info("StorePair finished - memstore OK")
		return alreadyExists, nil
	default:
		appLogger.Info("StorePair finished - wrong MemMode")
		return false, nil
	}
}

func CheckPair(conf *config.Cfg, queryID string) (database.URLInfo, error) {
	if conf.MemMode == "postgres" {
		return database.GetURLInfo(queryID)
	}

	fullURL, _ := memstorage.Store.Get(queryID)
	return database.URLInfo{OriginalURL: fullURL, IsDeleted: false}, nil
}

func BatchDelete(ctx context.Context, appLogger *zap.SugaredLogger, userID string, shortIDs []string) error {
	return database.BatchDelete(ctx, appLogger, userID, shortIDs)
}

type RequestEntry struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ResponseEntry struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func ProcessBatchTransaction(conf *config.Cfg, appLogger *zap.SugaredLogger, requestBatch []RequestEntry, userID string) ([]ResponseEntry, error) {
	appLogger.Info("ProcessBatchTransaction started")

	var responseBatch []ResponseEntry
	ctx := context.Background()

	var tx pgx.Tx
	var err error
	if conf.MemMode == "postgres" {
		tx, err = database.DBPool.Begin(ctx)
		if err != nil {
			appLogger.Errorw("ProcessBatchTransaction stopped - begin transaction FAIL", "error", err)
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

		if err := StorePairTransaction(ctx, tx, conf, appLogger, userID, shortID, entry.OriginalURL); err != nil {
			appLogger.Errorw("ProcessBatchTransaction stopped - StorePairTransaction from batch FAIL", "error", err)
			return nil, err
		}
	}

	if conf.MemMode == "savefile" {
		if err := savefile.SaveFile(conf, appLogger); err != nil {
			appLogger.Errorw("ProcessBatchTransaction stopped - save savefile FAIL", "error", err)
			return nil, err
		}
	}

	if tx != nil {
		if err := tx.Commit(ctx); err != nil {
			appLogger.Errorw("ProcessBatchTransaction stopped - commit transaction FAIL", "error", err)
			return nil, err
		}
	}

	appLogger.Info("ProcessBatchTransaction finished - batch processed OK")
	return responseBatch, nil
}

func StorePairTransaction(ctx context.Context, tx pgx.Tx, conf *config.Cfg, appLogger *zap.SugaredLogger, userID, shortID, fullURL string) error {
	appLogger.Info("StorePairTransaction started")

	memstorage.Store.Set(shortID, fullURL)

	if conf.MemMode == "postgres" {
		if err := database.AddOneTransaction(ctx, tx, appLogger, userID, shortID, fullURL); err != nil {
			appLogger.Errorw("StorePairTransaction stopped - AddOneTransaction to DB FAIL", "error", err)
			return err
		}
	}
	return nil
}

func GetUserURLs(conf *config.Cfg, appLogger *zap.SugaredLogger, userID string) ([]database.UserURL, error) {
	if conf.MemMode != "postgres" {
		return nil, nil
	}
	return database.GetUserURLs(conf, appLogger, userID)
}
