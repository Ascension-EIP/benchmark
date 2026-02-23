package handler

import (
	"net/http"

	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/dto/request"
	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/dto/response"
	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	s service.Auth
}

func NewAuthRoutes(r *gin.RouterGroup, s service.Auth) {
	h := &AuthHandler{s: s}

	r.POST("/signup", h.Signup)
	r.POST("/login", h.Login)
}

func (h *AuthHandler) Signup(c *gin.Context) {
	var req request.Signup
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.s.Signup(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req request.Login
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := h.s.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.Login{JWTToken: token})
}
