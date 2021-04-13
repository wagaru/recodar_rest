package usecase

import (
	"context"

	"github.com/wagaru/Recodar/server/internal/config"
	"github.com/wagaru/Recodar/server/internal/delivery/http/router"
	"github.com/wagaru/Recodar/server/internal/domain"
	"github.com/wagaru/Recodar/server/internal/repository"
)

type Usecase interface {
	// Login
	HandleUserLogin(session router.Session, info []byte, source string) error
	GetGoogleOauthURL(session router.Session) string
	GetGoogleOauthAccessToken(state, code string, session router.Session) (string, error)

	// Video
	StoreVideo(ctx context.Context, info map[string]interface{}) error

	// Accident
	StoreAccident(ctx context.Context, accident *domain.Accident) error
	StoreAccidents(ctx context.Context, accidents []*domain.Accident) error
	GetAccidents(ctx context.Context, queryFilter *domain.QueryFilter) ([]*domain.Accident, error)
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
