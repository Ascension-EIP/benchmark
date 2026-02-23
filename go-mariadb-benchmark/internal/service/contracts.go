package service

import (
	"context"

	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/dto/request"
	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/model"
)

type (
	User interface {
		Create(c context.Context, user *model.User) error
		Get(c context.Context, id uint) (*model.User, error)
		Update(c context.Context, user *model.User) error
		Delete(c context.Context, id uint) error
		List(c context.Context) ([]model.User, error)
	}

	Auth interface {
		Signup(c context.Context, req request.Signup) error
		Login(c context.Context, req request.Login) (string, error)
	}

	Upload interface {
		Upload(c context.Context, userID uint, req request.Upload) error
	}
)
