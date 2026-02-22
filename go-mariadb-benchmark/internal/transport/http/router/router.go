package router

import (
	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/service"
	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/transport/http/handler"
	"github.com/gin-gonic/gin"
)

func New(
	user service.User,
	auth service.Auth,
	upload service.Upload,
) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/v1")
	{
		handler.NewUsersRoutes(v1, user)
		handler.NewAuthRoutes(v1, auth)
		handler.NewUploadRoutes(v1, upload)
	}

	return r
}
