package database

import (
	"context"

	"github.com/Okenamay/shorturl.git/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	DBPool *pgxpool.Pool
)

func Init(conf *config.Cfg) error {
	DBPool, err := pgxpool.New(context.Background(), conf.PostgreDSN)
	if err != nil {
		return err
	}

	defer DBPool.Close()
	return nil
}
