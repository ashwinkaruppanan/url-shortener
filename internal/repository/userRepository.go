package repository

import (
	"context"
	"time"

	"example.com/url-shortener/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (u *userRepo) GetUserById(ctx context.Context, userID primitive.ObjectID) (*model.User, error) {
	var user model.User
	err := u.db.Collection("user").FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	return &user, err
}

func (u *userRepo) UpdateRefreshTokenInDB(ctx context.Context, userID primitive.ObjectID, refreshToken any, at time.Time) error {
	_, err := u.db.Collection("user").UpdateOne(ctx, bson.M{"_id": userID}, bson.M{"$set": bson.M{"refresh_token": refreshToken, "refresh_token_issued_at": at}})
	return err
}

func (u *userRepo) UpdateUrlInfoInDB(ctx context.Context, key string, device string, location string) (*model.Url, error) {
	var url model.Url

	filter := bson.M{"short_url_key": key}
	update := bson.M{
		"$inc": bson.M{
			"no_of_clicks":         1,
			"device." + device:     1,
			"location." + location: 1,
		},
	}

	options := options.FindOneAndUpdate().SetReturnDocument(options.After)

	err := u.db.Collection("url").FindOneAndUpdate(ctx, filter, update, options).Decode(&url)
	if err != nil {
		return nil, err
	}

	return &url, nil
}
