package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var database = initDatabase()

func initDatabase() *mongo.Database {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal(err)
	}

	mongoURI := os.Getenv("MONGODB_URI")

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("DB connection successful")

	return client.Database(os.Getenv("DB_NAME"))
}

func CreateCollection(collectionName string) *mongo.Collection {
	return database.Collection(collectionName)
}
