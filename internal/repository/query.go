package repository

import (
	"time"

	"github.com/wagaru/Recodar/server/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (repo *mongoRepo) parseQueryFilterInFullText(queryFilter *domain.QueryFilter) mongo.Pipeline {
	var sortStage, limitStage, skipStage, searchStage, nearStage, matchStage bson.D
	if queryFilter.Search != "" {
		clues := bson.D{
			{Key: "should", Value: bson.A{
				bson.M{
					"text": bson.M{
						"query": queryFilter.Search,
						"path":  "place.level_1",
					},
				},
				bson.M{
					"text": bson.M{
						"query": queryFilter.Search,
						"path":  "place.level_2",
					},
				},
				bson.M{
					"text": bson.M{
						"query": queryFilter.Search,
						"path":  "description",
					},
				},
			}},
			{Key: "minimumShouldMatch", Value: 1},
		}
		if queryFilter.Near != nil {
			clues = append(clues, bson.E{
				Key: "must", Value: bson.A{
					bson.M{
						"near": bson.M{
							"path": "place.geo_location",
							"origin": bson.M{
								"type":        "Point",
								"coordinates": bson.A{queryFilter.Near[0], queryFilter.Near[1]},
							},
							"pivot": 500,
						},
					},
				},
			})
		}
		searchStage = bson.D{{
			Key: "$search",
			Value: bson.M{
				"index":    "accidents_default",
				"compound": clues,
			},
		}}
	}

	if queryFilter.Near != nil && queryFilter.Search == "" {
		nearStage = bson.D{{
			Key: "$geoNear",
			Value: bson.M{
				"near": bson.M{
					"type":        "Point",
					"coordinates": bson.A{queryFilter.Near[0], queryFilter.Near[1]},
				},
				"distanceField": "dist.calculated",
				"spherical":     false,
				"maxDistance":   500,
				"key":           "place.geo_location",
			},
		}}
	}

	if queryFilter.AboutTime != nil {
		matchStage = bson.D{{
			Key: "$match",
			Value: bson.M{
				"approx_time": bson.M{
					"$gte": queryFilter.AboutTime.Add(time.Hour * -3),
					"$lte": queryFilter.AboutTime.Add(time.Hour * 3),
				},
			},
		}}
	}

	if queryFilter.Sort != nil {
		var ascending int
		switch queryFilter.Sort[1] {
		case "ascending":
			ascending = 1
		case "descending":
			ascending = -1
		default:
			ascending = 1
		}
		sortStage = bson.D{{Key: "$sort", Value: bson.D{{Key: queryFilter.Sort[0], Value: ascending}, {Key: "_id", Value: -1}}}}
	}

	if queryFilter.PerPage != 0 {
		limit := int64(queryFilter.PerPage)
		limitStage = bson.D{{Key: "$limit", Value: limit}}
	}

	if queryFilter.Page != 1 {
		skip := (int64(queryFilter.Page) - 1) * int64(queryFilter.PerPage)
		skipStage = bson.D{{Key: "$skip", Value: skip}}
	}

	pipeline := mongo.Pipeline{}
	if len(searchStage) > 0 {
		pipeline = append(pipeline, searchStage)
	}
	if len(nearStage) > 0 {
		pipeline = append(pipeline, nearStage)
	}
	if len(matchStage) > 0 {
		pipeline = append(pipeline, matchStage)
	}
	if len(sortStage) > 0 {
		pipeline = append(pipeline, sortStage)
	}
	if len(skipStage) > 0 {
		pipeline = append(pipeline, skipStage)
	}
	if len(limitStage) > 0 {
		pipeline = append(pipeline, limitStage)
	}
	return pipeline
}

func (repo *mongoRepo) parseQueryFilter(queryFilter *domain.QueryFilter) map[string]interface{} {

	opts := options.Find()
	filters := bson.D{}

	if queryFilter.Search != "" {
		// DO NOT USE
	}

	if queryFilter.Near != nil {
		filters = append(filters, bson.E{Key: "place.geo_location", Value: bson.M{
			"$near": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": bson.A{queryFilter.Near[0], queryFilter.Near[1]},
				},
				"$maxDistance": 500,
			},
		}})
	}

	if queryFilter.AboutTime != nil {
		filters = append(filters, bson.E{Key: "approx_time", Value: bson.M{
			"$gte": queryFilter.AboutTime.Add(time.Hour * -3),
			"$lte": queryFilter.AboutTime.Add(time.Hour * 3),
		}})
	}

	if queryFilter.Sort != nil {
		var ascending int
		switch queryFilter.Sort[1] {
		case "ascending":
			ascending = 1
		case "descending":
			ascending = -1
		default:
			ascending = 1
		}
		opts.SetSort(bson.D{{queryFilter.Sort[0], ascending}})
	}

	if queryFilter.PerPage != 0 {
		limit := int64(queryFilter.PerPage)
		opts.SetLimit(limit)
	}

	if queryFilter.Page != 1 {
		skip := (int64(queryFilter.Page) - 1) * int64(queryFilter.PerPage)
		opts.SetSkip(skip)
	}

	return map[string]interface{}{
		"options": opts,
		"filters": filters,
	}
}
