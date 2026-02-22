package handler

import (
	"net/http"

	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/service"
	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/transport/http/middleware"
	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	s service.Upload
}

func NewUploadRoutes(r *gin.RouterGroup, s service.Upload) {
	h := &UploadHandler{s: s}

	r.POST("/upload", middleware.Auth(), h.Upload)
}

func (h *UploadHandler) Upload(c *gin.Context) {
	c.Status(http.StatusOK)
}
