package repo

import "../../../../go-mariadb-benchmark/internal/repo/gorm.io/gorm"

type Repo struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repo {
	return &Repo{db: db}
}
