package middleware

import (
	"net/http"

	"example.com/url-shortener/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "cookie not found"})
			c.Abort()
			return
		}

		user_id, err := utils.ValidateToken(token, secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid access token"})
			c.Abort()
			return
		}

		c.Set("user_id", user_id)
		c.Next()
	}
}
