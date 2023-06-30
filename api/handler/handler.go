package handler

import (
	"log"
	"net/http"
	"time"

	"example.com/url-shortener/internal/model"
	"example.com/url-shortener/utils"
	"github.com/gin-gonic/gin"
	"github.com/mssola/useragent"
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
		// c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		utils.CjsonError(c, err)
		return
	}

	cookie := http.Cookie{
		Name:     "token",
		Value:    res.AccessToken,
		Expires:  time.Now().Add(10 * time.Hour),
		Path:     "/",
		Domain:   "https://reago.netlify.app",
		HttpOnly: true,
		Secure:   true,
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
		utils.CjsonError(c, err)
		return
	}

	cookie := http.Cookie{
		Name:     "token",
		Value:    res.AccessToken,
		Expires:  time.Now().Add(10 * time.Hour),
		Path:     "/",
		Domain:   "https://reago.netlify.app",
		HttpOnly: true,
		Secure:   true,
	}

	http.SetCookie(c.Writer, &cookie)

	c.JSON(http.StatusOK, res)
}

func (h *Handler) Logout(c *gin.Context) {

	token, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	err = h.service.Logout(c, token)
	if err != nil {
		utils.CjsonError(c, err)
		return
	}

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
		utils.CjsonError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"url_id": id})
}

func (h *Handler) GetAllURLs(c *gin.Context) {
	userID := c.GetString("user_id")

	res, err := h.service.GetAllURLs(c, userID)
	if err != nil {
		utils.CjsonError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) Refresh(c *gin.Context) {
	refreshToken := c.GetHeader("refresh-token")
	log.Println(refreshToken)
	access_token, err := h.service.RefreshAccessToken(c, refreshToken)
	if err != nil {
		utils.CjsonError(c, err)
		return
	}

	cookie := http.Cookie{
		Name:     "token",
		Value:    *access_token,
		Expires:  time.Now().Add(10 * time.Minute),
		Path:     "/",
		Domain:   "localhost",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(c.Writer, &cookie)

	c.JSON(http.StatusOK, gin.H{"success": "new access token set in cookie"})
}

func (h *Handler) RedirectURL(c *gin.Context) {
	key := c.Param("key")

	ua := c.Request.UserAgent()
	userAgent := useragent.New(ua)

	ip := c.ClientIP()
	if ip == "127.0.0.1" || ip == "::1" {
		ip = "157.51.198.201"
	}

	device := "Android"
	if userAgent.OSInfo().Name != "" {
		device = userAgent.OSInfo().Name
	}
	redirectURl, err := h.service.RedirectURL(c, key, ip, device)
	if err != nil {
		
			utils.CjsonError(c, err)
			return
		
	}

	
	c.JSON(http.StatusOK, redirectURl)
}
