package service

import (
	"context"
	"errors"
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
		return nil, err
	}

	if count > 0 {
		return nil, errors.New("email already exist")
	}

	hashedPassword, err := utils.HashPassword(userReq.Password)
	if err != nil {
		return nil, err
	}

	s := model.User{
		UserID:     primitive.NewObjectID(),
		FullName:   userReq.FullName,
		Email:      userReq.Email,
		Password:   hashedPassword,
		Created_at: time.Now(),
	}

	err = u.repository.Signup(ctx, &s)
	if err != nil {
		return nil, err
	}

	accessToken, err := utils.GenerateAccessToken(&s, os.Getenv("ACCESS_TOKEN_SECRET"), 10)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken(&s, os.Getenv("REFRESH_TOKEN_SECRET"), 72)
	if err != nil {
		return nil, err
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
		return nil, errors.New("email not found")
	}

	err = utils.VerifyPassword(loginReq.Password, user.Password)
	if err != nil {
		return nil, errors.New("wrong password")
	}

	accessToken, err := utils.GenerateAccessToken(user, os.Getenv("ACCESS_TOKEN_SECRET"), 10)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken(user, os.Getenv("REFRESH_TOKEN_SECRET"), 72)
	if err != nil {
		return nil, err
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

	count, err := u.repository.CheckUniqueUrlKey(ctx, urlReq.ShortURLKey)

	if err != nil {
		return "", err
	}

	if count > 0 {
		return "", errors.New("key already used")
	}

	uID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return "", err
	}

	newUrl := &model.Url{
		UrlID:       primitive.NewObjectID(),
		UserID:      uID,
		Label:       urlReq.Label,
		LongURL:     urlReq.LongURL,
		ShortURLKey: urlReq.ShortURLKey,
		NoOfClicks:  0,
		Device:      nil,
		Location:    nil,
		CreatedAt:   time.Now(),
	}

	err = u.repository.InsertUrl(ctx, newUrl)
	if err != nil {
		return "", err
	}

	return newUrl.UserID.Hex(), nil
}

func (u *userServ) GetAllURLs(c context.Context, userID string) (*[]model.Url, error) {
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	return u.repository.GetAllURLs(ctx, uid)
}
