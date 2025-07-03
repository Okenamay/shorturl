package memselect

import (
	"github.com/Okenamay/shorturl.git/internal/app/urlmaker"
	"github.com/Okenamay/shorturl.git/internal/config"
	logger "github.com/Okenamay/shorturl.git/internal/logger/zap"
	"github.com/Okenamay/shorturl.git/internal/storage/database"
	"github.com/Okenamay/shorturl.git/internal/storage/memstorage"
	"github.com/Okenamay/shorturl.git/internal/storage/savefile"
)

func MemInit(conf *config.Cfg) error {
	logger.Zap.Info("MemInit", "Assessing memory mode")

	switch conf.MemMode {
	case "postgres":
		err := database.StartDB(conf)
		if err != nil {
			logger.Zap.Errorw(err.Error(), "MemInit", "Init DB")
			return err
		}
		logger.Zap.Info("MemInit", "Init DB OK")
	case "savefile":
		err := savefile.LoadFile(conf)
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

var pingOK bool

func PingDB(conf *config.Cfg) (error, bool) {
	logger.Zap.Info("PingDB", "Pinging DB")

	switch conf.MemMode {
	case "postgres":
		err := database.DBPing()
		if err != nil {
			logger.Zap.Errorw(err.Error(), "PingDB", "Pinging DB")
			return err, pingOK
		}

		pingOK = true
	default:
		pingOK = false
	}

	return nil, pingOK
}

func StorePair(conf *config.Cfg, shortID, fullURL string) error {
	logger.Zap.Info("StorePair", "Running")

	memstorage.StoreURLIDPair(shortID, fullURL)

	switch conf.MemMode {
	case "postgres":
		err := database.AddOne(conf, shortID, fullURL)
		if err != nil {
			logger.Zap.Errorw(err.Error(), "StorePair", "AddOne to DB")
			return err
		}
		logger.Zap.Info("StorePair", "AddOne to DB OK")
	case "savefile":
		err := savefile.SaveFile(conf)
		if err != nil {
			logger.Zap.Errorw(err.Error(), "StorePair", "Save savefile")
			return err
		}
		logger.Zap.Info("StorePair", "Save savefile OK")
	case "memstore":
		logger.Zap.Info("StorePair", "Memstore OK")
	default:
		logger.Zap.Info("StorePair", "Wrong MemMode")
	}

	return nil
}

func CheckPair(conf *config.Cfg, queryID string) (string, error) {
	fullURL := memstorage.URLStore[queryID]

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
	logger.Zap.Info("BatchHandler. Start")

	var responseBatch []ResponseEntry

	for v := range requestBatch {
		tempCorrelationID := requestBatch[v].CorrelationID
		logger.Zap.Infof("BatchHandler. tempCorrelationID run %v: %s", v, tempCorrelationID)
		tempOriginalURL := requestBatch[v].OriginalURL
		logger.Zap.Infof("BatchHandler. tempOriginalURL run %v: %s", v, tempOriginalURL)

		tempShortURL, tempShortID := urlmaker.ProcessURL(conf, tempOriginalURL)
		logger.Zap.Infof("BatchHandler. tempShortURL run %v: %s", v, tempShortURL)
		logger.Zap.Infof("BatchHandler. tempShortID run %v: %s", v, tempShortID)

		responseBatch = append(responseBatch, ResponseEntry{
			CorrelationID: tempCorrelationID,
			ShortURL:      tempShortURL,
		})

		err := StorePair(conf, tempShortID, tempOriginalURL)
		if err != nil {
			logger.Zap.Errorw(err.Error(), "BatchHandler", "StorePair from batch")
			return nil, err
		}
	}

	logger.Zap.Info("BatchHandler. Processed batch")
	return responseBatch, nil
}
