package api

import (
	"github.com/convee/adcreative/internal/pkg/app"
	"github.com/convee/adcreative/internal/pkg/code"
	"github.com/convee/adcreative/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Publisher struct {
}

func (c *Publisher) List(ctx *gin.Context) {
	var (
		appG = app.Gin{C: ctx}
	)
	publisher := ctx.Query("publisher")
	publisherService := &service.Publisher{
		Name: publisher,
	}
	data := publisherService.GetApiList()
	appG.Response(http.StatusOK, code.SUCCESS, data)
	return
}
