package repository

import (
	"context"

	"github.com/wagaru/recodar-rest/internal/domain"
)

func (repo *mongoRepo) StoreVideo(ctx context.Context, v *domain.Video) (interface{}, error) {
	collection := repo.db.Collection(MONGO_VIDEO_COLLECTION)
	result, err := collection.InsertOne(ctx, v)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}
