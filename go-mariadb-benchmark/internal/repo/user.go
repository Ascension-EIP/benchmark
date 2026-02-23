package repo

import (
	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/model"
)

type User interface {
	CreateUser(user *model.User) error
	GetUserByID(id uint) (*model.User, error)
	ListAllUser() ([]model.User, error)
	UpdateUser(user *model.User) error
	DeleteUser(id uint) error
}

func (r *Repo) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *Repo) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repo) ListAllUser() ([]model.User, error) {
	var users []model.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *Repo) UpdateUser(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *Repo) DeleteUser(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}
