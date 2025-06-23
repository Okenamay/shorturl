package memselect

import (
	"github.com/Okenamay/shorturl.git/internal/app/urlmaker"
	"github.com/Okenamay/shorturl.git/internal/config"
	"github.com/Okenamay/shorturl.git/internal/storage/database"
	"github.com/Okenamay/shorturl.git/internal/storage/memstorage"
	"github.com/Okenamay/shorturl.git/internal/storage/savefile"
	"github.com/jackc/pgx/v5/pgxpool"

	logger "github.com/Okenamay/shorturl.git/internal/logger/zap"
)

type DB interface {
	Ping() error
	AddOne(shortID, fullURL string) error
}

type MemStorage interface {
	StoreURLIDPair(shortID, fullURL string)
}

type Memselect struct {
	db         DB
	memstorage MemStorage
}

func MemInit(conf *config.Cfg, pool *pgxpool.Pool) (*Memselect, error) {
	sugar, err := logger.InitLogger()
	if err != nil {
		sugar.Errorw(err.Error(), "MemInit", "Start logger")
		return err
	}

	switch conf.MemMode {
	case "postgres":
		db, err := database.StartDB(conf, pool, true)
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

	memselect := Memselect{
		db:         db,
		memstorage: memstorage,
	}

	return &memselect, nil
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
		err = database.DBPing(conf)
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

	// В теории, мы хотим читать из нужного источника, но зачем,
	// если всё равно при инициализации всё в мапу считали?
	// Или таки обязательно?

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
	sugar, _ := logger.InitLogger()
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
