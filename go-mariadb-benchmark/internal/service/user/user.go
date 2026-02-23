package user

import (
	"context"

	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/model"
	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/repo"
)

type Service struct {
	r repo.User
}

func New(r repo.User) *Service {
	return &Service{r: r}
}

func (s *Service) Create(c context.Context, user *model.User) error {
	return s.r.CreateUser(user)
}

func (s *Service) Get(c context.Context, id uint) (*model.User, error) {
	return s.r.GetUserByID(id)
}

func (s *Service) List(c context.Context) ([]model.User, error) {
	return s.r.ListAllUser()
}

func (s *Service) Update(c context.Context, user *model.User) error {
	return s.r.UpdateUser(user)
}

func (s *Service) Delete(c context.Context, id uint) error {
	return s.r.DeleteUser(id)
}
