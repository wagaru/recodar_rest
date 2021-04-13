package usecase

import (
	"context"
)

func (u *usecase) StoreVideo(ctx context.Context, info map[string]interface{}) error {
	// v := &domain.Video{
	// 	URL:        "https://youtu.be/o_alcd7iR7o",
	// 	City:       "新竹縣竹北市",
	// 	Roads:      []string{"光明六路", "自強南路", "自強北路"},
	// 	ApproxTime: time.Now().Add(time.Hour * -1),
	// 	GeoLocation: domain.GeoJSON{
	// 		Type:        "Point",
	// 		Coordinates: []float64{24.812734, 121.035558},
	// 	},
	// 	Description: "一台BMW 大7與Tesla擦撞",
	// 	CreatedAt:   time.Now(),
	// 	AccidentObjects: []domain.AccidentObject{
	// 		domain.AccidentObject(domain.Automobile),
	// 	},
	// }
	// _, err := u.repo.StoreVideo(ctx, v)
	// if err != nil {
	// 	return err
	// }
	return nil
}
