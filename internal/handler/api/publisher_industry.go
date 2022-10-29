package api

import (
	"github.com/convee/adcreative/internal/pkg/app"
	"github.com/convee/adcreative/internal/pkg/code"
	"github.com/convee/adcreative/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PublisherIndustry struct {
}

func (c *PublisherIndustry) List(ctx *gin.Context) {
	var (
		appG = app.Gin{C: ctx}
	)
	publisher := ctx.Query("publisher")
	publisherIndustryService := &service.PublisherIndustry{
		Publisher: publisher,
	}
	data := publisherIndustryService.GetList()
	appG.Response(http.StatusOK, code.SUCCESS, data)
	return
}
