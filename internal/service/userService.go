package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"example.com/url-shortener/internal/model"
	"example.com/url-shortener/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userServ struct {
	repository model.UserRepositoryInterface
}

func NewUserService(repository model.UserRepositoryInterface) model.UserServiceInterface {
	return &userServ{
		repository,
	}
}

func (u *userServ) Signup(c context.Context, userReq *model.CreateUserReq) (*model.SignupLoginUserRes, error) {
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	count, err := u.repository.CheckUniqueEmail(ctx, userReq.Email)
	if err != nil {
		return nil, &utils.AppError{Code: http.StatusInternalServerError, Message: "internal server error"}
	}

	if count > 0 {
		return nil, &utils.AppError{Code: http.StatusConflict, Message: "email already exist"}
	}

	hashedPassword, err := utils.HashPassword(userReq.Password)
	if err != nil {
		return nil, &utils.AppError{Code: http.StatusInternalServerError, Message: "internal server error"}
	}

	s := model.User{
		UserID:     primitive.NewObjectID(),
		FullName:   userReq.FullName,
		Email:      userReq.Email,
		Password:   hashedPassword,
		Created_at: time.Now(),
	}

	accessToken, err := utils.GenerateAccessToken(&s, os.Getenv("ACCESS_TOKEN_SECRET"), 10)
	if err != nil {
		return nil, &utils.AppError{Code: http.StatusInternalServerError, Message: "internal server error"}
	}

	refreshToken, err := utils.GenerateRefreshToken(&s, os.Getenv("REFRESH_TOKEN_SECRET"), 72)
	if err != nil {
		return nil, &utils.AppError{Code: http.StatusInternalServerError, Message: "internal server error"}
	}

	s.RefreshToken = refreshToken
	s.RefreshTokenIssuedAT = time.Now()

	err = u.repository.Signup(ctx, &s)
	if err != nil {
		return nil, &utils.AppError{Code: http.StatusInternalServerError, Message: "internal server error"}
	}

	res := &model.SignupLoginUserRes{
		UserID:       s.UserID,
		FullName:     s.FullName,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return res, nil
}

func (u *userServ) Login(c context.Context, loginReq *model.LoginUserReq) (*model.SignupLoginUserRes, error) {
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	user, err := u.repository.GetUserByEmail(ctx, loginReq.Email)
	if err != nil {
		return nil, &utils.AppError{Code: http.StatusUnauthorized, Message: "email not found"}
	}

	err = utils.VerifyPassword(loginReq.Password, user.Password)
	if err != nil {
		return nil, &utils.AppError{Code: http.StatusUnauthorized, Message: "wrong password"}
	}

	accessToken, err := utils.GenerateAccessToken(user, os.Getenv("ACCESS_TOKEN_SECRET"), 10)
	if err != nil {
		return nil, &utils.AppError{Code: http.StatusInternalServerError, Message: "internal server error"}
	}

	refreshToken, err := utils.GenerateRefreshToken(user, os.Getenv("REFRESH_TOKEN_SECRET"), 72)
	if err != nil {
		return nil, &utils.AppError{Code: http.StatusInternalServerError, Message: "internal server error"}
	}

	err = u.repository.UpdateRefreshTokenInDB(ctx, user.UserID, refreshToken, time.Now())

	if err != nil {
		return nil, &utils.AppError{Code: http.StatusInternalServerError, Message: "internal server error"}
	}

	res := &model.SignupLoginUserRes{
		UserID:       user.UserID,
		FullName:     user.FullName,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return res, nil
}

func (u *userServ) CreateURL(c context.Context, userID string, urlReq *model.CreateUrlReq) (string, error) {
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	wordSet := make(map[string]bool)
	words := []string{"signup", "login", "refresh", "logout", "create-url", "get-all-urls"}

	for _, word := range words {
		wordSet[word] = true
	}

	if wordSet[urlReq.ShortURLKey] {
		return "", &utils.AppError{Code: http.StatusBadRequest, Message: "endpoint is reserved, not allowed to use"}
	}

	count, err := u.repository.CheckUniqueUrlKey(ctx, urlReq.ShortURLKey)

	if err != nil {
		return "", &utils.AppError{Code: http.StatusInternalServerError, Message: "internal server error"}
	}

	if count > 0 {
		return "", &utils.AppError{Code: http.StatusBadRequest, Message: "endpoint already used"}
	}

	uID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return "", &utils.AppError{Code: http.StatusInternalServerError, Message: "internal server error"}
	}

	var temp = make(map[string]int)

	newUrl := &model.Url{
		UrlID:       primitive.NewObjectID(),
		UserID:      uID,
		Label:       urlReq.Label,
		LongURL:     urlReq.LongURL,
		ShortURLKey: urlReq.ShortURLKey,
		NoOfClicks:  0,
		Device:      temp,
		Location:    temp,
		CreatedAt:   time.Now(),
	}

	err = u.repository.InsertUrl(ctx, newUrl)
	if err != nil {
		return "", &utils.AppError{Code: http.StatusInternalServerError, Message: "internal server error"}
	}

	return newUrl.UserID.Hex(), nil
}

func (u *userServ) GetAllURLs(c context.Context, userID string) (*[]model.Url, error) {
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, &utils.AppError{Code: http.StatusInternalServerError, Message: "internal server error"}
	}

	res, err := u.repository.GetAllURLs(ctx, uid)
	if err != nil {
		return nil, &utils.AppError{Code: http.StatusInternalServerError, Message: "internal server error"}
	}

	return res, nil
}

func (u *userServ) RefreshAccessToken(c context.Context, refreshToken string) (*string, error) {
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	if refreshToken == "" {
		// c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return nil, &utils.AppError{Code: http.StatusUnauthorized, Message: "invalid refresh token"}
	}

	userID, err := utils.ValidateToken(refreshToken, os.Getenv("REFRESH_TOKEN_SECRET"))
	if err != nil {
		// c.JSON(http.StatusUnauthorized, gin.H{"error": "token error"})
		return nil, &utils.AppError{Code: http.StatusUnauthorized, Message: "invalid refresh token"}
	}

	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, &utils.AppError{Code: http.StatusInternalServerError, Message: "internal server error"}
	}

	res, err := u.repository.GetUserById(ctx, uid)

	if err != nil {
		return nil, &utils.AppError{Code: http.StatusInternalServerError, Message: "internal server error"}
	}

	if refreshToken != res.RefreshToken {
		return nil, &utils.AppError{Code: http.StatusUnauthorized, Message: "expired refresh token"}
	}

	accessToken, err := utils.GenerateAccessToken(res, os.Getenv("ACCESS_TOKEN_SECRET"), 10)
	if err != nil {
		return nil, &utils.AppError{Code: http.StatusInternalServerError, Message: "internal server error"}
	}

	return &accessToken, nil
}

func (u *userServ) Logout(c context.Context, token string) error {
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	uID, err := utils.ValidateToken(token, os.Getenv("ACCESS_TOKEN_SECRET"))
	if err != nil {
		return &utils.AppError{Code: http.StatusUnauthorized, Message: err.Error()}
	}

	userID, err := primitive.ObjectIDFromHex(uID)
	if err != nil {
		return &utils.AppError{Code: http.StatusInternalServerError, Message: "internal server error"}
	}
	err = u.repository.UpdateRefreshTokenInDB(ctx, userID, nil, time.Now())
	if err != nil {
		return &utils.AppError{Code: http.StatusInternalServerError, Message: "internal server error"}
	}
	return nil
}

type geolocation struct {
	Location string `json:"city"`
}

func (u *userServ) RedirectURL(c context.Context, key string, ip string, device string) (string, error) {
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	count, _ := u.repository.CheckUniqueUrlKey(ctx, key)
	if count < 1 {
		return "", &utils.AppError{Code: http.StatusBadRequest, Message: "invalid enpoint"}
	}

	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)

	res, err := http.Get(url)
	if err != nil {
		return "", &utils.AppError{Code: http.StatusInternalServerError, Message: "internal server error"}
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", &utils.AppError{Code: http.StatusInternalServerError, Message: "internal server error"}
	}

	var location geolocation

	err = json.Unmarshal(data, &location)
	if err != nil {
		return "", &utils.AppError{Code: http.StatusInternalServerError, Message: "internal server error"}
	}

	out, err := u.repository.UpdateUrlInfoInDB(ctx, key, device, location.Location)
	if err != nil {
		return "", &utils.AppError{Code: http.StatusInternalServerError, Message: "internal server error"}
	}

	return out.LongURL, nil
}
