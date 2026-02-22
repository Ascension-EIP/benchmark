package repo

import "github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/model"

type Auth interface {
	CreateUser(user *model.User) error
	GetUserByUsername(username string) (*model.User, error)
}

func (r *Repo) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, "username = ?", username).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
