package http

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (delivery *httpDelivery) routeVideos(c *gin.Context) {
	// TODO: santize input
	test := map[string]interface{}{
		"test": "test",
	}
	delivery.usecase.StoreVideo(context.Background(), test)
	c.String(http.StatusOK, "hello")
}
