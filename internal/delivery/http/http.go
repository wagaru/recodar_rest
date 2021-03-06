package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wagaru/recodar-rest/internal/config"
	"github.com/wagaru/recodar-rest/internal/domain"
	"github.com/wagaru/recodar-rest/internal/logger"
	"github.com/wagaru/recodar-rest/internal/usecase"
)

type httpDelivery struct {
	usecase              usecase.Usecase
	messageBrokerUsecase usecase.MessageBrokerUsecase
	router               *Router
	config               *config.Config
}

func NewHttpDelivery(usecase usecase.Usecase, messageBrokerUsecase usecase.MessageBrokerUsecase, config *config.Config) *httpDelivery {
	return &httpDelivery{
		usecase:              usecase,
		messageBrokerUsecase: messageBrokerUsecase,
		router:               NewRouter(config),
		config:               config,
	}
}

func (delivery *httpDelivery) buildRoute() {
	api := delivery.router.Group("/api/v1")
	api.Use(delivery.router.Middlewares["RateLimit"])
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
			authRequired.GET("/me", delivery.me)

			authRequired.GET("/accidents", delivery.getAccidents)
			authRequired.POST("/accidents", delivery.postAccidents)
			authRequired.POST("/accidents/delete", delivery.deleteAccidents)
			authRequired.DELETE("/accidents/:id", delivery.deleteAccident)
		}
	}

	// for test only
	// delivery.router.GET("genTest", delivery.genTestAccidents)
	delivery.router.GET("test", delivery.test)
}

func (delivery *httpDelivery) Run(port uint16) {
	delivery.buildRoute()
	logger.Logger.Printf("Listening and serving HTTP on %v", port)
	http.ListenAndServe(fmt.Sprintf(":%v", port), delivery.router)
}

func (delivery *httpDelivery) test(c *gin.Context) {
	message := domain.MessageUserLogin{
		ID: "123456",
	}
	messageEncoded, err := json.Marshal(message)
	if err != nil {
		logger.Logger.Printf("Encode message failed:%v", err)
		return
	}
	meta := &domain.RabbitMQMeta{
		ExchangeType: "fanout",
		ExchangeName: "fanout",
	}
	err = delivery.messageBrokerUsecase.SendMessages(meta, messageEncoded)
	if err != nil {
		logger.Logger.Printf("send message failed:%v", err)
		return
	}
	c.String(200, "OK")
}
