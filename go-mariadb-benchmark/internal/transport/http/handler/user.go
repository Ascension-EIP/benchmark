package handler

import (
	"net/http"
	"strconv"

	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/model"
	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/service"
	"github.com/Ascension-EIP/benchmark/go-mariadb-benchmark/internal/transport/http/middleware"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	s service.User
}

func NewUsersRoutes(r *gin.RouterGroup, s service.User) {
	h := &UserHandler{s: s}

	users := r.Group("/users")
	{
		users.POST("/", middleware.Admin(), h.Create)
		users.GET("/", middleware.Admin(), h.List)
		users.GET("/:id", middleware.Admin(), h.GetByID)
		users.PUT("/:id", middleware.Admin(), h.Update)
		users.DELETE("/:id", middleware.Admin(), h.Delete)
	}
}

func (h *UserHandler) Create(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.s.Create(c.Request.Context(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := h.s.Get(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) List(c *gin.Context) {
	users, _ := h.s.List(c.Request.Context())
	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.ID = uint(id)
	if err := h.s.Update(c.Request.Context(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update user"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.s.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete user"})
		return
	}
	c.Status(http.StatusNoContent)
}
