package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Okenamay/shorturl.git/internal/config"
	logger "github.com/Okenamay/shorturl.git/internal/logger/zap"
	"github.com/Okenamay/shorturl.git/internal/server/router"
	"github.com/Okenamay/shorturl.git/internal/storage/database"
	"github.com/Okenamay/shorturl.git/internal/storage/memselect"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Main:
func main() {
	conf := config.InitConfig()

	sugar, err := logger.InitLogger()
	if err != nil {
		sugar.Fatalw(err.Error(), "Main", "Start logger")
	}

	dbPool, err := pgxpool.New(context.Background(), conf.PostgreDSN)
	if err != nil {
		sugar.Fatalw(err.Error(), "Main", "Create DB pool")
	}
	defer dbPool.Close()

	db, err := database.StartDB(conf, dbPool, true)
	if err != nil {
		sugar.Fatalw(err.Error(), "Main", "Initialize database")
	}

	memselect, err := memselect.MemInit(conf, db)
	if err != nil {
		fmt.Errorf("failed to initialize storage: %v", err)
		sugar.Errorw(err.Error(), "Main", "Initialize storage")
	}

	sugar.Infow("Starting server on port: ", conf.ServerPort)

	err = router.Launch(conf)
	if err != nil {
		sugar.Fatalw(err.Error(), "Main", "Start server")
	}
}

func init() {
	log.Fatal("Init")
}
