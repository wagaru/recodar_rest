package http

import (
	"fmt"

	"github.com/wagaru/recodar-rest/internal/config"
	"github.com/wagaru/recodar-rest/internal/delivery/http/router"
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
			authRequired.POST("/accidents", delivery.postAccidents)
		}
	}

	// for test only
	delivery.router.GET("genTest", delivery.genTestAccidents)
}

func (delivery *httpDelivery) Run(port uint16) {
	delivery.buildRoute()
	delivery.router.Run(fmt.Sprintf(":%v", port))
}
