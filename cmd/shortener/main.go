package main

import (
	"log"

	"github.com/Okenamay/shorturl.git/internal/config"
	router "github.com/Okenamay/shorturl.git/internal/server/router"

	_ "go.uber.org/zap"
)

// Main:
func main() {
	config.ParseFlags()

	log.Printf("Starting server on port %s", config.Cfg.ServerPort)

	err := router.Launch()
	if err != nil {
		log.Fatal(err)
	}
}
