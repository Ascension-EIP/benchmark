package router

import (
	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/transport/http/handler"
	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/transport/http/middleware"
	"github.com/gin-gonic/gin"
)

func New(
	userHandler *handler.UserHandler,
) *gin.Engine {
	r := gin.New()

	users := r.Group("/users")
	{
		users.POST("/", middleware.Admin(), userHandler.Create)
		users.GET("/", middleware.Admin(), userHandler.List)
		users.GET("/:id", middleware.Admin(), userHandler.GetByID)
		users.PUT("/:id", middleware.Admin(), userHandler.Update)
		users.DELETE("/:id", middleware.Admin(), userHandler.Delete)
	}

	return r
}
