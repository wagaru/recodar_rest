package usecase

import (
	"context"
	"time"

	"github.com/wagaru/recodar-rest/internal/domain"
)

func (u *usecase) StoreAccident(ctx context.Context, accident *domain.Accident, userID string) error {
	now := time.Now()
	accident.CreatedAt = &now
	_, err := u.repo.StoreAccident(ctx, accident, userID)
	return err
}

func (u *usecase) StoreAccidents(ctx context.Context, accidents []*domain.Accident, userID string) error {
	_, err := u.repo.StoreAccidents(ctx, accidents, userID)
	return err
}

func (u *usecase) GetAccidents(ctx context.Context, queryFilter *domain.QueryFilter) ([]*domain.Accident, error) {
	accidents, err := u.repo.GetAccidents(ctx, queryFilter)
	if err != nil {
		return nil, err
	}
	return accidents, nil
}

func (u *usecase) DeleteAccident(ctx context.Context, IDHex string) error {
	return u.repo.DeleteAccident(ctx, IDHex)
}

func (u *usecase) DeleteAccidents(ctx context.Context, IDs []string) error {
	return u.repo.DeleteAccidents(ctx, IDs)
}
