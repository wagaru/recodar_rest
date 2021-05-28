package usecase

import (
	"context"
	"time"

	"github.com/wagaru/recodar-rest/internal/config"
	"github.com/wagaru/recodar-rest/internal/domain"
	"github.com/wagaru/recodar-rest/internal/repository"
)

type Usecase interface {
	// auth
	GetLineOAuthURL() string
	GetGoogleOAuthURL() string
	GetGoogleOAuthAccessToken(state, code string) (string, string, time.Time, error)

	// Video
	StoreVideo(ctx context.Context, info map[string]interface{}) error

	// Accident
	StoreAccident(ctx context.Context, accident *domain.Accident) error
	StoreAccidents(ctx context.Context, accidents []*domain.Accident) error
	GetAccidents(ctx context.Context, queryFilter *domain.QueryFilter) ([]*domain.Accident, error)

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
