package http

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	status    int
	err       error
	errMsg    string
	errDetail string
}

type SuccessResponse struct {
	status int
	data   map[string]interface{}
}

func (delivery *httpDelivery) WrapResponse(c *gin.Context, responseType interface{}) {
	switch v := responseType.(type) {
	case ErrorResponse:
		status := v.status
		if status == 0 {
			status = http.StatusBadRequest
		}
		errMsg := v.errMsg
		if v.errMsg == "" {
			errMsg = "Unknown error"
		}
		log.Printf("[Response] Error: status %d, err %v, errMsg %s, errDetail: %s", status, v.err, errMsg, v.errDetail)
		c.AbortWithStatusJSON(status, gin.H{"error": errMsg})
		return
	case SuccessResponse:
		status := v.status
		if status == 0 {
			status = http.StatusOK
		}
		c.JSON(status, v.data)
		return
	}
}
