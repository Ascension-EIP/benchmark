package main

import (
	"log"

	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/config"
	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/db"
	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/repository"
	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/service"
	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/transport/http/handler"
	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/transport/http/router"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := db.New(cfg.DB.DSN())
	if err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	r := router.New(userHandler)
	log.Fatal(r.Run())
}
