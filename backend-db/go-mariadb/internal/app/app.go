package app

import (
	"log"

	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/config"
	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/db"
	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/repo"
	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/service/auth"
	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/service/upload"
	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/service/user"
	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/transport/http/router"
)

func Run(cfg *config.Config) {
	db, err := db.New(cfg.DB.DSN())
	if err != nil {
		log.Fatal(err)
	}
	repo := repo.New(db)

	userService := user.New(repo)
	authService := auth.New(repo, cfg.Auth)
	uploadService := upload.New(repo)

	r := router.New(cfg, userService, authService, uploadService)
	log.Fatal(r.Run())
}
