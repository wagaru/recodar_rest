package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/wagaru/recodar-rest/internal/config"
	"github.com/wagaru/recodar-rest/internal/domain"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	MONGO_DATABASE            = "recodar"
	MONGO_USER_COLLECTION     = "users"
	MONGO_VIDEO_COLLECTION    = "videos"
	MONGO_ACCIDENT_COLLECTION = "accidents"
)

type Repository interface {
	// Disconnect
	Disconnect()

	// User
	FindUser(ctx context.Context, condition map[string]interface{}) (*domain.User, error)
	FindUserById(ctx context.Context, IDHex string) (*domain.User, error)
	StoreUser(ctx context.Context, u *domain.User) (string, error)
	UpdateUser(ctx context.Context, id string, u *domain.User) error
	UpsertUser(ctx context.Context, filter map[string]interface{}, update map[string]interface{}) (*domain.User, error)

	// Accident
	StoreAccident(ctx context.Context, a *domain.Accident) (interface{}, error)
	StoreAccidents(ctx context.Context, as []*domain.Accident) (interface{}, error)
	GetAccidents(ctx context.Context, queryFilter *domain.QueryFilter) ([]*domain.Accident, error)
	DeleteAccident(ctx context.Context, IDHex string) error
}
type mongoRepo struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewMongoRepo(config *config.Config) (Repository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return nil, fmt.Errorf("Connect mongo db failed: %w", err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, fmt.Errorf("Ping mongo db failed: %w", err)
	}

	return &mongoRepo{
		client: client,
		db:     client.Database(MONGO_DATABASE),
	}, nil
}

func (repo *mongoRepo) Disconnect() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	if err := repo.client.Disconnect(ctx); err != nil {
		log.Fatal(err)
	}
}
