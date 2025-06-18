package memselect

import (
	"github.com/Okenamay/shorturl.git/internal/config"
	logger "github.com/Okenamay/shorturl.git/internal/logger/zap"
	"github.com/Okenamay/shorturl.git/internal/storage/database"
	"github.com/Okenamay/shorturl.git/internal/storage/memstorage"
	"github.com/Okenamay/shorturl.git/internal/storage/savefile"
)

func MemInit(conf *config.Cfg) error {
	sugar, err := logger.InitLogger()
	if err != nil {
		sugar.Fatalw(err.Error(), "MemInit", "Start logger")
		return err
	}

	switch conf.MemMode {
	case "postgres":
		err = database.StartDB(conf)
		if err != nil {
			sugar.Fatalw(err.Error(), "MemInit", "Init DB")
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

var pingOK bool

func PingDB(conf *config.Cfg) (error, bool) {
	sugar, err := logger.InitLogger()
	if err != nil {
		sugar.Fatalw(err.Error(), "PingDB", "Start logger")
		return err, pingOK
	}

	switch conf.MemMode {
	case "postgres":
		err = database.DBPing(conf)
		if err != nil {
			sugar.Fatalw(err.Error(), "PingDB", "Pinging DB")
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
		sugar.Fatalw(err.Error(), "StorePair", "Start logger")
		return err
	}

	memstorage.StoreURLIDPair(shortID, fullURL)

	switch conf.MemMode {
	case "postgres":
		err = database.AddOne(conf, shortID, fullURL)
		if err != nil {
			sugar.Fatalw(err.Error(), "StorePair", "AddOne to DB")
			return err
		}
		sugar.Info("StorePair", "AddOne to DB OK")
	case "savefile":
		err := savefile.SaveFile(conf)
		if err != nil {
			sugar.Fatalw(err.Error(), "StorePair", "Save savefile")
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
		sugar.Fatalw(err.Error(), "CheckOne", "Start logger")
		return "", err
	}

	// Хотим читать из нужного источника, но зачем, если всё равно при
	// инициализации всё в мапу считали? Или обязательно?

	// switch conf.MemMode {
	// case "postgres":
	// 	fullURL, err = database.CheckOne(conf, queryID)
	// 	if err != nil {
	// 		sugar.Fatalw(err.Error(), "CheckOne", "Check DB")
	// 		return "", err
	// 	}
	// 	sugar.Info("CheckOne", "DB check OK")
	// default:
	// 	fullURL, exists = memstorage.URLStore[queryID]
	// 	sugar.Info("CheckOne", "Mem check OK")
	// }

	fullURL := memstorage.URLStore[queryID]

	return fullURL, nil
}

// func MemSaveAll(conf *config.Cfg, ) error {
// 	sugar, err := logger.InitLogger()
// 	if err != nil {
// 		sugar.Fatalw(err.Error(), "MemInit", "Start logger")
// 		return err
// 	}

// 	switch conf.MemMode {
// 	case "postgres":
// 		err = database.Init(conf)
// 		if err != nil {
// 			sugar.Fatalw(err.Error(), "MemInit", "Init DB")
// 			return err
// 		}

// 		err = database

// 		sugar.Info("MemInit", "Init DB OK")
// 	case "savefile":
// 		err := savefile.LoadFile(conf)
// 		if err != nil {
// 			sugar.Fatalw(err.Error(), "MemInit", "Load savefile")
// 			return err
// 		}
// 		sugar.Info("MemInit", "Load savefile OK")
// 	case "memstore":
// 		sugar.Info("MemInit", "Memstore OK")
// 	}

// 	return nil
// }
