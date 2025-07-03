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
	sugar, err := logger.InitLogger()
	if err != nil {
		sugar.Errorw(err.Error(), "MemInit", "Start logger")
		return err
	}

	switch conf.MemMode {
	case "postgres":
		err = database.StartDB(conf)
		if err != nil {
			sugar.Errorw(err.Error(), "MemInit", "Init DB")
			return err
		}
		sugar.Info("MemInit", "Init DB OK")
	case "savefile":
		err := savefile.LoadFile(conf)
		if err != nil {
			sugar.Errorw(err.Error(), "MemInit", "Load savefile")
			return err
		}
		sugar.Info("MemInit", "Load savefile OK")
	case "memstore":
		sugar.Info("MemInit", "Memstore OK")
	default:
		sugar.Info("MemInit", "Wrong MemMode")
	}

	return nil
}

func MemStop(conf *config.Cfg) {
	sugar, err := logger.InitLogger()
	if err != nil {
		sugar.Errorw(err.Error(), "MemStop", "Start logger")
	}
	switch conf.MemMode {
	case "postgres":
		database.StopDB()
		sugar.Info("MemStop", "Stop DB OK")
	default:
		sugar.Info("MemStop", "Nothing to stop for this MemMode")
	}
}

var pingOK bool

func PingDB(conf *config.Cfg) (error, bool) {
	sugar, err := logger.InitLogger()
	if err != nil {
		sugar.Errorw(err.Error(), "PingDB", "Start logger")
		return err, pingOK
	}

	switch conf.MemMode {
	case "postgres":
		err = database.DBPing()
		if err != nil {
			sugar.Errorw(err.Error(), "PingDB", "Pinging DB")
			return err, pingOK
		}

		pingOK = true
	default:
		pingOK = false
	}

	return nil, pingOK
}

func StorePair(conf *config.Cfg, shortID, fullURL string) error {
	sugar, err := logger.InitLogger()
	if err != nil {
		sugar.Errorw(err.Error(), "StorePair", "Start logger")
		return err
	}

	memstorage.StoreURLIDPair(shortID, fullURL)

	switch conf.MemMode {
	case "postgres":
		err = database.AddOne(conf, shortID, fullURL)
		if err != nil {
			sugar.Errorw(err.Error(), "StorePair", "AddOne to DB")
			return err
		}
		sugar.Info("StorePair", "AddOne to DB OK")
	case "savefile":
		err := savefile.SaveFile(conf)
		if err != nil {
			sugar.Errorw(err.Error(), "StorePair", "Save savefile")
			return err
		}
		sugar.Info("StorePair", "Save savefile OK")
	case "memstore":
		sugar.Info("StorePair", "Memstore OK")
	default:
		sugar.Info("StorePair", "Wrong MemMode")
	}

	return nil
}

func CheckPair(conf *config.Cfg, queryID string) (string, error) {
	sugar, err := logger.InitLogger()
	if err != nil {
		sugar.Errorw(err.Error(), "CheckOne", "Start logger")
		return "", err
	}

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
	sugar, err := logger.InitLogger()
	if err != nil {
		sugar.Errorw(err.Error(), "ProcessBatch", "Start logger")
	}

	sugar.Info("BatchHandler. Start")

	var responseBatch []ResponseEntry

	for v := range requestBatch {
		tempCorrelationID := requestBatch[v].CorrelationID
		sugar.Infof("BatchHandler. tempCorrelationID run %v: %s", v, tempCorrelationID)
		tempOriginalURL := requestBatch[v].OriginalURL
		sugar.Infof("BatchHandler. tempOriginalURL run %v: %s", v, tempOriginalURL)

		tempShortURL, tempShortID := urlmaker.ProcessURL(conf, tempOriginalURL)
		sugar.Infof("BatchHandler. tempShortURL run %v: %s", v, tempShortURL)
		sugar.Infof("BatchHandler. tempShortID run %v: %s", v, tempShortID)

		responseBatch = append(responseBatch, ResponseEntry{
			CorrelationID: tempCorrelationID,
			ShortURL:      tempShortURL,
		})

		err := StorePair(conf, tempShortID, tempOriginalURL)
		if err != nil {
			sugar.Errorw(err.Error(), "BatchHandler", "StorePair from batch")
			return nil, err
		}
	}

	sugar.Info("BatchHandler. Processed batch")
	return responseBatch, nil
}
