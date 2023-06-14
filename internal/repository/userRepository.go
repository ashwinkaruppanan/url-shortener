package repository

import (
	"context"

	"example.com/url-shortener/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (u *userRepo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := u.db.Collection("user").FindOne(ctx, bson.M{"email": email}).Decode(&user)
	return &user, err
}

func (u *userRepo) CheckUniqueUrlKey(ctx context.Context, key string) (int64, error) {
	return u.db.Collection("url").CountDocuments(ctx, bson.M{"short_url_key": key})
}

func (u *userRepo) InsertUrl(ctx context.Context, url *model.Url) error {
	_, err := u.db.Collection("url").InsertOne(ctx, url)
	return err
}

func (u *userRepo) GetAllURLs(ctx context.Context, userID primitive.ObjectID) (*[]model.Url, error) {

	cursor, err := u.db.Collection("url").Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}

	var res []model.Url

	err = cursor.All(ctx, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
