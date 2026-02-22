package upload

import (
	"context"

	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/dto/request"
	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/repo"
)

type Service struct {
	r repo.Upload
}

func New(r repo.Upload) *Service {
	return &Service{r: r}
}

func (s *Service) Upload(c context.Context, userId uint, req request.Upload) error {
	return nil
}
