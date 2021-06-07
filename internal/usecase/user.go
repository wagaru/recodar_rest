package usecase

import (
	"context"

	"github.com/wagaru/recodar-rest/internal/domain"
)

func (usecase *usecase) FindUser(ctx context.Context, condition map[string]interface{}) (*domain.User, error) {
	return usecase.repo.FindUser(ctx, condition)
}

func (usecase *usecase) FindUserById(ctx context.Context, IDHex string) (*domain.User, error) {
	return usecase.repo.FindUserById(ctx, IDHex)
}

func (usecase *usecase) StoreUser(ctx context.Context, user *domain.User) (string, error) {
	return usecase.repo.StoreUser(ctx, user)
}

func (usecase *usecase) UpsertUser(ctx context.Context, filter, update map[string]interface{}) (*domain.User, error) {
	return usecase.repo.UpsertUser(ctx, filter, update)
}
