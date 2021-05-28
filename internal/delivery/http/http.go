package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wagaru/recodar-rest/internal/config"
	"github.com/wagaru/recodar-rest/internal/delivery/http/router"
	"github.com/wagaru/recodar-rest/internal/domain"
	"github.com/wagaru/recodar-rest/internal/usecase"
)

type httpDelivery struct {
	usecase usecase.Usecase
	router  *router.Router
	config  *config.Config
}

func NewHttpDelivery(usecase usecase.Usecase, config *config.Config) *httpDelivery {
	return &httpDelivery{
		usecase: usecase,
		router:  router.NewRouter(config),
		config:  config,
	}
}

func (delivery *httpDelivery) buildRoute() {
	api := delivery.router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.GET("/line", delivery.authLine)
			auth.GET("/line/callback", delivery.authLineCallback)
			auth.GET("/google", delivery.authGoogle)
			auth.GET("/google/callback", delivery.authGoogleCallback)
		}

		authRequired := api.Use(delivery.router.Middlewares["AuthRequired"])
		{
			authRequired.GET("/accidents", delivery.getAccidents)
			// authRequire.POST("/accidents", delivery.postAccidents)
		}
	}
}

func (delivery *httpDelivery) Run(port uint16) {
	delivery.buildRoute()
	delivery.router.Run(fmt.Sprintf(":%v", port))
}

func (delivery *httpDelivery) genAccidents(c *gin.Context) {
	c.Request.URL.Query()
	places := []map[string]interface{}{
		map[string]interface{}{
			"level_1":                  "新竹縣竹北市",
			"level_2":                  "光明六路",
			"geo_location.coordinates": []float64{121.00900982104032, 24.82745801255536},
			"approx_time":              "2021-03-21T13:00:00.000+00:00",
			"description":              "一台機車跟一台汽車在路口高速追撞",
			"accident_objects":         []string{"Motorcycle", "Automobile"},
			"videos.media_url":         "https://www.youtube.com/watch?v=eVIozKR9p50",
			"videos.media_id":          "eVIozKR9p50",
		},
		map[string]interface{}{
			"level_1":                  "新竹市東區",
			"level_2":                  "中正路",
			"geo_location.coordinates": []float64{120.9715940760574, 24.80227118111046},
			"approx_time":              "2021-03-20T16:00:00.000+00:00",
			"description":              "一台公車跟汽座在快車道發生擦撞，公車直行而汽車要右轉",
			"accident_objects":         []string{"Automobile"},
			"videos.media_url":         "https://www.youtube.com/watch?v=PALItzoZ5b0",
			"videos.media_id":          "PALItzoZ5b0",
		},
		map[string]interface{}{
			"level_1":                  "新竹縣竹東鎮",
			"level_2":                  "北興路二段",
			"geo_location.coordinates": []float64{121.09407708332822, 24.738113452682807},
			"approx_time":              "2021-03-21T00:00:00.000+00:00",
			"description":              "行人要穿越馬路時遭汽車撞擊",
			"accident_objects":         []string{"Pedestrian", "Automobile"},
			"videos.media_url":         "https://www.youtube.com/watch?v=4mH4-Ej5yjE",
			"videos.media_id":          "4mH4-Ej5yjE",
		},
		map[string]interface{}{
			"level_1":                  "新竹市東區",
			"level_2":                  "明湖路775巷",
			"geo_location.coordinates": []float64{120.96760743862308, 24.77449027085487},
			"approx_time":              "2021-03-20T03:00:00.000+00:00",
			"description":              "二台機車對撞，一台綠燈直行，一台闖紅燈",
			"accident_objects":         []string{"Motorcycle"},
			"videos.media_url":         "https://www.youtube.com/watch?v=zLB0vicG38I",
			"videos.media_id":          "zLB0vicG38I",
		},
		map[string]interface{}{
			"level_1":                  "新竹市東區",
			"level_2":                  "公道五路二段",
			"geo_location.coordinates": []float64{120.99279389814757, 24.805470026131193},
			"approx_time":              "2021-03-22T14:00:00.000+00:00",
			"description":              "二台汽車在路口發生撞擊，一台從公路五路左轉忠孝路，一台從忠孝路左轉公道五路",
			"accident_objects":         []string{"Automobile"},
			"videos.media_url":         "https://www.youtube.com/watch?v=uqyzNhSOqLw",
			"videos.media_id":          "uqyzNhSOqLw",
		},
	}
	now := time.Now()
	accidents := []*domain.Accident{}
	for _, place := range places {
		time, _ := time.Parse(time.RFC3339, place["approx_time"].(string))
		accident := &domain.Accident{
			Place: domain.Place{
				Level1: place["level_1"].(string),
				Level2: place["level_2"].(string),
				GeoLocation: domain.GeoJSON{
					"Point",
					[]float64{place["geo_location.coordinates"].([]float64)[0], place["geo_location.coordinates"].([]float64)[1]},
				},
			},
			ApproxTime:  time,
			Description: place["description"].(string),
			// AccidentObjects: []domain.Accident{
			// 	domain.Accident{place["accident_objects"].([]string)[0]}
			// },
			Videos: []domain.Video{
				domain.Video{
					MediaID:   place["videos.media_id"].(string),
					MediaURL:  place["videos.media_url"].(string),
					MediaType: "youtube",
				},
			},
			CreatedAt: &now,
		}
		for _, accidentObject := range place["accident_objects"].([]string) {
			accident.AccidentObjects = append(accident.AccidentObjects, domain.AccidentObject(accidentObject))
		}
		accidents = append(accidents, accident)
	}
	err := delivery.usecase.StoreAccidents(context.Background(), accidents)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, "success")
}
