package utils

import (
	"errors"
	"time"

	"example.com/url-shortener/internal/model"
	"github.com/golang-jwt/jwt/v4"
)

func GenerateAccessToken(user *model.User, secret string, expiry int) (string, error) {
	exp := time.Now().Add(time.Hour * time.Duration(expiry))
	claims := &model.JwtCustomAccessClaims{
		UserID: user.UserID.Hex(),
		Name:   user.FullName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(secret))
}

func GenerateRefreshToken(user *model.User, secret string, expiry int) (string, error) {
	exp := time.Now().Add(time.Hour * time.Duration(expiry))
	claims := &model.JwtCustomRefreshClaims{
		UserID: user.UserID.Hex(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(secret))
}

func ValidateToken(token string, secret string) (string, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok && !t.Valid {
		return "", errors.New("invalid token")
	}

	return claims["user_id"].(string), nil
}
