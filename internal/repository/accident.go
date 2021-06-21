package repository

import (
	"context"

	"github.com/wagaru/recodar-rest/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (repo *mongoRepo) StoreAccident(ctx context.Context, a *domain.Accident) (interface{}, error) {
	collection := repo.db.Collection(MONGO_ACCIDENT_COLLECTION)
	result, err := collection.InsertOne(ctx, a)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

func (repo *mongoRepo) StoreAccidents(ctx context.Context, as []*domain.Accident) (interface{}, error) {
	collection := repo.db.Collection(MONGO_ACCIDENT_COLLECTION)
	documents := []interface{}{}
	for _, a := range as {
		documents = append(documents, a)
	}
	result, err := collection.InsertMany(ctx, documents)
	if err != nil {
		return nil, err
	}
	return result.InsertedIDs, nil
}

func (repo *mongoRepo) GetAccidents(ctx context.Context, queryFilter *domain.QueryFilter) ([]*domain.Accident, error) {
	collection := repo.db.Collection(MONGO_ACCIDENT_COLLECTION)
	var cursor *mongo.Cursor
	var err error
	if queryFilter.InFullTextSearchMode {
		cursor, err = collection.Aggregate(ctx, repo.parseQueryFilterInFullText(queryFilter))
		if err != nil {
			return nil, err
		}

	} else {
		filters := repo.parseQueryFilter(queryFilter)
		cursor, err = collection.Find(ctx, filters["filters"], filters["options"].(*options.FindOptions))
		if err != nil {
			return nil, err
		}
	}
	var accidents []*domain.Accident
	if err := cursor.All(ctx, &accidents); err != nil {
		return nil, err
	}
	return accidents, nil
}

func (repo *mongoRepo) DeleteAccident(ctx context.Context, IDHex string) error {
	collection := repo.db.Collection(MONGO_ACCIDENT_COLLECTION)
	objectID, err := primitive.ObjectIDFromHex(IDHex)
	if err != nil {
		return err
	}
	_, err = collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}
	return nil
}

func (repo *mongoRepo) DeleteAccidents(ctx context.Context, IDs []string) error {
	collection := repo.db.Collection(MONGO_ACCIDENT_COLLECTION)
	IDsHex := make([]primitive.ObjectID, len(IDs))
	for _, id := range IDs {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			//TODO throw error
			continue
		}
		IDsHex = append(IDsHex, objectID)
	}
	_, err := collection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": IDsHex}})
	if err != nil {
		return err
	}
	return nil
}
