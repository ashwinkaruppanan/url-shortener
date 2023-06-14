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

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	db := db.NewMongoDatabase()
	rep := repository.NewUserRepository(db)
	ser := service.NewUserService(rep)
	router.NewRouter(r, ser)

	if err := r.Run(os.Getenv("PORT")); err != nil {
		log.Fatal(err)
	}

}
