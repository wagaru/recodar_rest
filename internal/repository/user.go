package repository

import (
	"context"
	"errors"

	"github.com/wagaru/Recodar/server/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (repo *mongoRepo) GetUser(ctx context.Context, key string, value interface{}) (*domain.User, error) {
	collection := repo.db.Collection(MONGO_USER_COLLECTION)
	user := &domain.User{}
	err := collection.FindOne(context.Background(), bson.M{key: value}).Decode(user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return &domain.User{}, nil
	}
	if err != nil {
		return &domain.User{}, err
	}
	return user, nil
}

func (repo *mongoRepo) StoreUser(ctx context.Context, u *domain.User) (interface{}, error) {
	collection := repo.db.Collection(MONGO_USER_COLLECTION)
	result, err := collection.InsertOne(ctx, u)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

func (repo *mongoRepo) UpdateUser(ctx context.Context, idHex string, u *domain.User) error {
	collection := repo.db.Collection(MONGO_USER_COLLECTION)
	id, _ := primitive.ObjectIDFromHex(idHex)
	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.D{
			{"$set", bson.D{
				{"name", u.Name},
				{"email", u.Email},
				{"picture", u.Picture},
			}},
			{"$currentDate", bson.D{{"updated_at", true}}},
		},
	)
	if err != nil {
		return err
	}
	return nil
}
