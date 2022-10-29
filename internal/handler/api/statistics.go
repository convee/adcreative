package api

import (
	"github.com/convee/adcreative/internal/model"
	"github.com/convee/adcreative/internal/pkg/app"
	"github.com/convee/adcreative/internal/pkg/code"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type Statistics struct {
}

type StatisticsQuery struct {
	PublisherName string `json:"publisherName"`
	CreativeId    string `json:"creativeId"`
	AdvertiserId  string `json:"advertiserId"`
}

func (s *Statistics) Status(ctx *gin.Context) {

	var (
		appG    = app.Gin{C: ctx}
		results map[string]interface{}
		maps    map[string]interface{}
	)
	var params StatisticsQuery
	errMsg := app.BindJson(ctx, &params)
	if len(errMsg) > 0 {
		appG.Response(http.StatusInternalServerError, code.ERROR, results)
		return
	}

	maps = make(map[string]interface{})
	if len(params.PublisherName) > 0 {
		maps["publisher.name"] = strings.Split(params.PublisherName, ",")
	}
	if len(params.CreativeId) > 0 {
		maps["creative_id"] = strings.Split(params.CreativeId, ",")
	}
	if len(params.AdvertiserId) > 0 {
		maps["advertiser_id"] = strings.Split(params.AdvertiserId, ",")
	}

	lst, _ := model.GetStatisticsCreativeStats(maps)
	results = make(map[string]interface{})

	results["info"] = lst

	appG.Response(http.StatusOK, code.SUCCESS, results)
	return
}
