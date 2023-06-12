package service

import (
	"context"
	"errors"
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

func (u *userServ) Signup(c context.Context, userReq *model.CreateUserReq) (*model.CreateUserRes, error) {
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

	s := &model.User{
		UserID:     primitive.NewObjectID(),
		FullName:   userReq.FullName,
		Email:      userReq.Email,
		Password:   hashedPassword,
		Created_at: time.Now(),
	}

	err = u.repository.Signup(ctx, s)
	if err != nil {
		return nil, err
	}

	res := &model.CreateUserRes{
		UserID:   s.UserID,
		FullName: s.FullName,
	}

	return res, nil
}
