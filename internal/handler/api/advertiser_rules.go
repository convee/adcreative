package api

import (
	"github.com/convee/adcreative/internal/pkg/app"
	"github.com/convee/adcreative/internal/pkg/code"
	"github.com/convee/adcreative/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AdvertiserRules struct {
}

func (c *AdvertiserRules) List(ctx *gin.Context) {
	var (
		appG = app.Gin{C: ctx}
	)
	publisher := ctx.Query("publisher")
	advertiserRulesService := &service.AdvertiserRules{
		Publisher: publisher,
	}
	data := advertiserRulesService.GetApiList()
	appG.Response(http.StatusOK, code.SUCCESS, data)
	return
}
