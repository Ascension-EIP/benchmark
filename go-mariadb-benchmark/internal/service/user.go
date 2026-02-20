package service

import (
	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/model"
	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(user *model.User) error {
	return s.repo.Create(user)
}

func (s *UserService) Get(id uint) (*model.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) Update(user *model.User) error {
	return s.repo.Update(user)
}

func (s *UserService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *UserService) List() ([]model.User, error) {
	return s.repo.ListAll()
}
