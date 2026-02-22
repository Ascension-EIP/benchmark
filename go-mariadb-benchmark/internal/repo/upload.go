package repo

import "github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/model"

type Upload interface {
	AddFile(upload *model.Upload) error
}

func (r *Repo) AddFile(upload *model.Upload) error {
	return r.db.Create(upload).Error
}
