package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wagaru/recodar-rest/internal/domain"
)

func (delivery *httpDelivery) getAccidents(c *gin.Context) {
	queryFilter := domain.NewQueryFilter(c.Request.URL.Query())
	accidents, err := delivery.usecase.GetAccidents(context.Background(), queryFilter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	//fmt.Println(accidents[0])
	c.JSON(http.StatusOK, accidents)
}

func (delivery *httpDelivery) postAccidents(c *gin.Context) {
	var a domain.Accident
	err := c.ShouldBindJSON(&a)
	if err != nil {
		//TODO
		fmt.Print(err)
		return
	}
	// fmt.Printf("%+v\n", a)
	err = delivery.usecase.StoreAccident(context.Background(), &a)
	if err != nil {
		fmt.Println(err)
	}
	c.String(http.StatusOK, "hello")
}
