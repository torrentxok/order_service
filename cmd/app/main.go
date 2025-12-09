package main

import (
	"log"

	"github.com/torrentxok/order_service/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("config: ", err)
	}

	log.Printf("Config loaded: %+v\n", cfg)
}
