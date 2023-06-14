package handler

import (
	"net/http"
	"time"

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

	// c.SetCookie("token", res.AccessToken, 10*60, "/", "localhost", false, true)
	cookie := http.Cookie{
		Name:     "token",
		Value:    res.AccessToken,
		Expires:  time.Now().Add(10 * time.Minute),
		Path:     "/",
		Domain:   "localhost",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(c.Writer, &cookie)

	c.JSON(http.StatusCreated, *res)
}

func (h *Handler) Login(c *gin.Context) {
	var loginReq model.LoginUserReq

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.Login(c, &loginReq)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	cookie := http.Cookie{
		Name:     "token",
		Value:    res.AccessToken,
		Expires:  time.Now().Add(10 * time.Minute),
		Path:     "/",
		Domain:   "localhost",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(c.Writer, &cookie)

	c.JSON(http.StatusOK, res)
}

func (h *Handler) Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"success": "logged out successfully"})
}

func (h *Handler) CreatURL(c *gin.Context) {
	var urlReq model.CreateUrlReq

	err := c.ShouldBindJSON(&urlReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("user_id")

	id, err := h.service.CreateURL(c, userID, &urlReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"url_id": id})
}

func (h *Handler) GetAllURLs(c *gin.Context) {
	userID := c.GetString("user_id")

	res, err := h.service.GetAllURLs(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, res)
}
