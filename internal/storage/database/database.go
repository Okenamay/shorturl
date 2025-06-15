package database

import (
	"context"
	"fmt"

	"github.com/Okenamay/shorturl.git/internal/config"
	logger "github.com/Okenamay/shorturl.git/internal/logger/zap"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	DBPool *pgxpool.Pool
)

func Init(conf *config.Cfg) error {
	sugar, _ := logger.InitLogger()
	sugar.Info("DB Init. Start")
	fmt.Println(conf.PostgreDSN)
	DBPool, err := pgxpool.New(context.Background(), conf.PostgreDSN)
	sugar.Info("DB Init. Make Pool")
	if err != nil {
		return err
	}
	fmt.Println(DBPool)
	defer DBPool.Close()
	return nil
}

func Close() {
	sugar, _ := logger.InitLogger()
	sugar.Info("DB Init. Defer Close")
	defer DBPool.Close()
}
