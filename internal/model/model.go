package model

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateUserReq struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserRes struct {
	UserID   primitive.ObjectID `json:"user_id"`
	FullName string             `json:"full_name"`
}

type LoginUserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

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

type UserRepositoryInterface interface {
	Signup(ctx context.Context, user *User) error
	CheckUniqueEmail(ctx context.Context, email string) (int64, error)
}

type UserServiceInterface interface {
	Signup(c context.Context, userReq *CreateUserReq) (*CreateUserRes, error)
}
