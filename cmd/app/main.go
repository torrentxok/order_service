package main

import (
	"log"

	"github.com/torrentxok/order_service/internal/config"
	"github.com/torrentxok/order_service/internal/logger"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("config: ", err)
	}

	logger, err := logger.New("debug")
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	log.Printf("Config loaded: %+v\n", cfg)
}
