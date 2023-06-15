package model

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateUserReq struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupLoginUserRes struct {
	UserID       primitive.ObjectID `json:"user_id"`
	FullName     string             `json:"full_name"`
	AccessToken  string             `json:"access_token"`
	RefreshToken string             `json:"refresh_token"`
}

type Url struct {
	UrlID       primitive.ObjectID `json:"url_id" bson:"_id"`
	UserID      primitive.ObjectID `json:"user_id" bson:"user_id"`
	Label       string             `json:"label" bson:"label"`
	LongURL     string             `json:"long_url" bson:"long_url"`
	ShortURLKey string             `json:"short_url_key" bson:"short_url_key"`
	NoOfClicks  int                `json:"no_of_clicks" bson:"no_of_clicks"`
	Device      map[string]int     `json:"device"`
	Location    map[string]int     `json:"location"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
}

type CreateUrlReq struct {
	Label       string `json:"label"`
	LongURL     string `json:"long_url"`
	ShortURLKey string `json:"short_url_key"`
}

type User struct {
	UserID               primitive.ObjectID `bson:"_id"`
	FullName             string             `json:"full_name"`
	Email                string             `json:"email"`
	Password             string             `json:"password"`
	Created_at           time.Time          `json:"created_at"`
	RefreshToken         string             `json:"refresh_token" bson:"refresh_token"`
	RefreshTokenIssuedAT time.Time          `json:"refresh_token_issued_at" bson:"refresh_token_issued_at"`
}

type JwtCustomAccessClaims struct {
	Name   string `json:"name"`
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type JwtCustomRefreshClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type UserRepositoryInterface interface {
	Signup(ctx context.Context, user *User) error
	CheckUniqueEmail(ctx context.Context, email string) (int64, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserById(ctx context.Context, userID primitive.ObjectID) (*User, error)
	UpdateRefreshTokenInDB(ctx context.Context, userID primitive.ObjectID, refreshToken any, at time.Time) error

	CheckUniqueUrlKey(ctx context.Context, key string) (int64, error)
	InsertUrl(ctx context.Context, url *Url) error
	GetAllURLs(ctx context.Context, userID primitive.ObjectID) (*[]Url, error)
	UpdateUrlInfoInDB(ctx context.Context, key string, device string, location string) (*Url, error)
}

type UserServiceInterface interface {
	Signup(c context.Context, userReq *CreateUserReq) (*SignupLoginUserRes, error)
	Login(c context.Context, loginReq *LoginUserReq) (*SignupLoginUserRes, error)

	CreateURL(c context.Context, userID string, urlReq *CreateUrlReq) (string, error)
	GetAllURLs(c context.Context, userID string) (*[]Url, error)

	RefreshAccessToken(c context.Context, refreshToken string) (*string, error)
	Logout(c context.Context, userID string) error
	RedirectURL(c context.Context, key string, ip string, device string) (string, error)
}
