package handler

import (
	"net/http"

	"example.com/url-shortener/internal/model"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service model.UserServiceInterface
}

func NewUserHandler(service model.UserServiceInterface) *Handler {
	return &Handler{
		service,
	}
}

func (h *Handler) Signup(c *gin.Context) {
	var user model.CreateUserReq
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.Signup(c, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, *res)
}
