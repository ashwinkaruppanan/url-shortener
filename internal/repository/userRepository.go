package repository

import (
	"context"

	"example.com/url-shortener/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepo struct {
	db *mongo.Database
}

func NewUserRepository(db *mongo.Database) model.UserRepositoryInterface {
	return &userRepo{
		db,
	}
}

func (u *userRepo) Signup(ctx context.Context, user *model.User) error {
	_, err := u.db.Collection("user").InsertOne(ctx, user)
	return err
}

func (u *userRepo) CheckUniqueEmail(ctx context.Context, email string) (int64, error) {
	return u.db.Collection("user").CountDocuments(ctx, bson.M{"email": email})
}
