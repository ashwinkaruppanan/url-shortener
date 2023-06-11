package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Url struct {
	LongURL    string         `json:"long_url"`
	ShortURL   string         `json:"short_url"`
	NoOfClicks int            `json:"no_of_clicks"`
	Device     map[string]int `json:"device"`
	Location   map[string]int `json:"location"`
}

type User struct {
	UserID     primitive.ObjectID `bson:"_id"`
	FullName   string             `json:"full_name"`
	Email      string             `json:"email"`
	Password   string             `json:"password"`
	Urls       []Url              `json:"urls"`
	Created_at time.Time          `json:"created_at"`
}
