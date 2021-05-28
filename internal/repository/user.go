package repository

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/wagaru/recodar-rest/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (repo *mongoRepo) FindUser(ctx context.Context, condition map[string]interface{}) (*domain.User, error) {
	collection := repo.db.Collection(MONGO_USER_COLLECTION)
	user := &domain.User{}
	err := collection.FindOne(context.Background(), condition).Decode(user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return &domain.User{}, nil
	}
	if err != nil {
		return &domain.User{}, err
	}
	return user, nil
}

func (repo *mongoRepo) FindUserById(ctx context.Context, IDHex string) (*domain.User, error) {
	collection := repo.db.Collection(MONGO_USER_COLLECTION)
	user := &domain.User{}
	objectID, err := primitive.ObjectIDFromHex(IDHex)
	if err != nil {
		return &domain.User{}, nil
	}
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return &domain.User{}, nil
	}
	if err != nil {
		return &domain.User{}, err
	}
	return user, nil
}

func (repo *mongoRepo) UpsertUser(ctx context.Context, filter map[string]interface{}, update map[string]interface{}) (*domain.User, error) {
	collection := repo.db.Collection(MONGO_USER_COLLECTION)
	f := bson.M{}
	for key, value := range filter {
		f[key] = value
	}
	u := bson.M{
		"updated_at": time.Now(),
	}
	for key, value := range update {
		u[key] = value
	}
	var user *domain.User
	err := collection.FindOneAndUpdate(ctx, f, bson.D{{"$set", u}}, options.FindOneAndUpdate().SetUpsert(true)).Decode(&user)
	log.Printf("upsert user %v, err %v", user, err)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// if it is an upsert, no new document will be returned, we need to find again
			if err = collection.FindOne(ctx, f).Decode(&user); err != nil {
				return &domain.User{}, err
			}
			return user, nil
		}
		return &domain.User{}, err
	}
	return user, nil
}

func (repo *mongoRepo) StoreUser(ctx context.Context, u *domain.User) (string, error) {
	collection := repo.db.Collection(MONGO_USER_COLLECTION)
	result, err := collection.InsertOne(ctx, u)
	if err != nil {
		return "", err
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
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
