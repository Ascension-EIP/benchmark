package handler

import (
	"net/http"

	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/config"
	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/dto/request"
	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/service"
	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/transport/http/middleware"
	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/transport/http/utils"
	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	s service.Upload
}

func NewUploadRoutes(r *gin.RouterGroup, cfg config.AuthConfig, s service.Upload) {
	h := &UploadHandler{s: s}

	r.POST("/upload", middleware.Auth(cfg), h.Upload)
}

func (h *UploadHandler) Upload(c *gin.Context) {
	userID, err := utils.GetFromContext[uint](c, "userID")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req request.Upload
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.s.Upload(c, userID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
