package usecase

import (
	"context"

	"github.com/wagaru/recodar-rest/internal/config"
	"github.com/wagaru/recodar-rest/internal/domain"
	"github.com/wagaru/recodar-rest/internal/repository"
)

type Usecase interface {
	// Accident
	StoreAccident(ctx context.Context, accident *domain.Accident, userID string) error
	StoreAccidents(ctx context.Context, accidents []*domain.Accident, userID string) error
	GetAccidents(ctx context.Context, queryFilter *domain.QueryFilter) ([]*domain.Accident, error)
	DeleteAccident(ctx context.Context, IDHex string) error
	DeleteAccidents(ctx context.Context, IDs []string) error

	// User
	FindUser(ctx context.Context, condition map[string]interface{}) (*domain.User, error)
	FindUserById(ctx context.Context, IDHex string) (*domain.User, error)
	StoreUser(ctx context.Context, user *domain.User) (string, error)
	UpsertUser(ctx context.Context, filter, update map[string]interface{}) (*domain.User, error)

	// Token
	GenerateJWTToken(ctx context.Context, user *domain.User) (string, error)
}

type usecase struct {
	repo   repository.Repository
	config *config.Config
}

func NewUsecase(repo repository.Repository, config *config.Config) Usecase {
	return &usecase{
		repo:   repo,
		config: config,
	}
}
