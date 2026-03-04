package main

import (
	"log"

	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/app"
	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config.Load: %s", err)
	}

	app.Run(cfg)
}
