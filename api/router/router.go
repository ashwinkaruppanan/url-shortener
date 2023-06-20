package router

import (
	"net/http"
	"os"
	"time"

	"example.com/url-shortener/api/handler"
	"example.com/url-shortener/api/middleware"
	"example.com/url-shortener/internal/model"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(r *gin.Engine, ser model.UserServiceInterface) {
	h := handler.NewUserHandler(ser)

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:5173"
		},
		MaxAge: 12 * time.Hour,
	}))

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
