package domain

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type QueryFilter struct {
	Page                 int
	PerPage              int
	Sort                 []string
	Search               string
	Near                 []float64
	AboutTime            *time.Time
	InFullTextSearchMode bool
}

func NewQueryFilter(query url.Values) *QueryFilter {
	queryFilter := &QueryFilter{
		Page:    1,
		PerPage: 10,
		Sort:    []string{"created_at", "desc"},
	}
	for key, value := range query {
		queryValue := value[len(value)-1]
		if queryValue != "" {
			switch key {
			case "page":
				page, _ := strconv.Atoi(queryValue)
				queryFilter.Page = page
			case "perPage":
				perPage, _ := strconv.Atoi(queryValue)
				queryFilter.PerPage = perPage
			case "sort":
				if fields := strings.Fields(queryValue); len(fields) == 2 {
					queryFilter.Sort = fields
				}
			case "search":
				queryFilter.Search = queryValue
				queryFilter.InFullTextSearchMode = true
			case "near":
				if fields := strings.Split(queryValue, ","); len(fields) == 2 {
					lat, _ := strconv.ParseFloat(fields[0], 64)
					lng, _ := strconv.ParseFloat(fields[1], 64)
					queryFilter.Near = []float64{lat, lng}
				}
			case "aboutTime":
				time, err := time.Parse(time.RFC3339, queryValue)
				if err != nil {
					//TODO
					fmt.Printf("parse error %v", err)
				} else {
					queryFilter.AboutTime = &time
				}
			}
		}
	}
	return queryFilter
}
