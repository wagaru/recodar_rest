package usecase

import (
	"context"
	"time"

	"github.com/wagaru/Recodar/server/internal/domain"
)

func (u *usecase) StoreAccident(ctx context.Context, accident *domain.Accident) error {
	now := time.Now()
	accident.CreatedAt = &now
	_, err := u.repo.StoreAccident(ctx, accident)
	return err
}

func (u *usecase) StoreAccidents(ctx context.Context, accidents []*domain.Accident) error {
	_, err := u.repo.StoreAccidents(ctx, accidents)
	return err
}

func (u *usecase) GetAccidents(ctx context.Context, queryFilter *domain.QueryFilter) ([]*domain.Accident, error) {
	accidents, err := u.repo.GetAccidents(ctx, queryFilter)
	if err != nil {
		return nil, err
	}
	return accidents, nil
}
