package http

import (
	"context"

	"github.com/gin-gonic/gin"
)

func (delivery *httpDelivery) me(c *gin.Context) {
	userID, _ := c.Get("userId")
	user, err := delivery.usecase.FindUserById(context.Background(), userID.(string))
	if err != nil {
		WrapResponse(c, ErrorResponse{err: err})
		return
	}
	WrapResponse(c, SuccessResponse{data: user})
}
