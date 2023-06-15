package router

import (
	"net/http"
	"os"

	"example.com/url-shortener/api/handler"
	"example.com/url-shortener/api/middleware"
	"example.com/url-shortener/internal/model"
	"github.com/gin-gonic/gin"
)

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func NewRouter(r *gin.Engine, ser model.UserServiceInterface) {
	h := handler.NewUserHandler(ser)

	r.Use(corsMiddleware())
	//Public routes
	public := r.Group("")
	public.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "home")
	})
	public.POST("/signup", h.Signup)
	public.POST("/login", h.Login)
	public.GET("/refresh", h.Refresh)
	public.GET("/:key", h.RedirectURL)

	// //Protected routes
	protected := r.Group("")
	protected.Use(middleware.AuthMiddleware(os.Getenv("ACCESS_TOKEN_SECRET")))
	protected.GET("/logout", h.Logout)
	protected.POST("/create-url", h.CreatURL)
	protected.GET("/get-all-urls", h.GetAllURLs)

}
