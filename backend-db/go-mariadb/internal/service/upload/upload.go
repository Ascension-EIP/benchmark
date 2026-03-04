package upload

import (
	"context"
	"errors"
	"io"

	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/dto/request"
	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/model"
	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/repo"
)

type Service struct {
	r repo.Upload
}

func New(r repo.Upload) *Service {
	return &Service{r: r}
}

func (s *Service) Upload(c context.Context, userId uint, req request.Upload) error {
	if req.File == nil {
		return errors.New("missing file")
	}
	file, err := req.File.Open()
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	upload := model.Upload{
		UserID: userId,
		File:   data,
	}
	if err := s.r.AddFile(&upload); err != nil {
		return err
	}

	return nil
}
