package main

import (
	"log"
	"os"

	"example.com/url-shortener/api/router"
	"example.com/url-shortener/db"
	"example.com/url-shortener/internal/repository"
	"example.com/url-shortener/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal(err)
	}
	r := gin.Default()
	r.Use(corsMiddleware())

	db := db.NewMongoDatabase()
	rep := repository.NewUserRepository(db)
	ser := service.NewUserService(rep)
	router.NewRouter(r, ser)

	if err := r.Run(os.Getenv("PORT")); err != nil {
		log.Fatal(err)
	}

}
