package http

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wagaru/recodar-rest/internal/domain"
)

func (delivery *httpDelivery) getAccidents(c *gin.Context) {
	queryFilter := domain.NewQueryFilter(c.Request.URL.Query())
	accidents, err := delivery.usecase.GetAccidents(context.Background(), queryFilter)
	if err != nil {
		WrapResponse(c, ErrorResponse{err: err})
		return

	}
	WrapResponse(c, SuccessResponse{data: accidents})
}

func (delivery *httpDelivery) postAccidents(c *gin.Context) {
	var a domain.Accident
	err := c.ShouldBindJSON(&a)
	if err != nil {
		WrapResponse(c, ErrorResponse{err: err, errMsg: "Invalid input"})
		return
	}
	userID, _ := c.Get("userId")
	err = delivery.usecase.StoreAccident(context.Background(), &a, userID.(string))
	if err != nil {
		WrapResponse(c, ErrorResponse{err: err, errMsg: "DB Error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
	WrapResponse(c, SuccessResponse{data: map[string]interface{}{
		"success": true,
	}})
}

func (delivery *httpDelivery) deleteAccident(c *gin.Context) {
	// TODO: check user can delete accident
	IDHex := c.Param("id")
	err := delivery.usecase.DeleteAccident(context.Background(), IDHex)
	if err != nil {
		WrapResponse(c, ErrorResponse{err: err})
		return
	}
	WrapResponse(c, SuccessResponse{status: http.StatusNoContent})
}

func (delivery *httpDelivery) deleteAccidents(c *gin.Context) {
	// TODO: check user can delete accidents

	type IDs struct {
		IDs []string `json:"ids"`
	}

	var ids IDs
	if err := c.ShouldBindJSON(&ids); err != nil {
		WrapResponse(c, ErrorResponse{err: err})
		return
	}
	err := delivery.usecase.DeleteAccidents(context.Background(), ids.IDs)
	if err != nil {
		WrapResponse(c, ErrorResponse{err: err})
	}
	WrapResponse(c, SuccessResponse{status: http.StatusNoContent})
}
