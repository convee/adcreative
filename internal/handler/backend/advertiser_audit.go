package backend

import (
	"github.com/convee/adcreative/internal/service"
	"net/http"

	"github.com/convee/adcreative/internal/pkg/app"
	"github.com/convee/adcreative/internal/pkg/code"
	"github.com/spf13/cast"

	"github.com/gin-gonic/gin"
)

type AdvertiserAudit struct {
}

func (c *AdvertiserAudit) List(ctx *gin.Context) {
	var (
		appG = app.Gin{C: ctx}
	)
	page := cast.ToInt(ctx.Query("page"))
	perPage := cast.ToInt(ctx.Query("per_page"))
	advertiserAuditService := &service.AdvertiserAudit{
		Page:    page,
		PerPage: perPage,
	}
	data := advertiserAuditService.GetList()
	appG.Response(http.StatusOK, code.SUCCESS, data)
	return
}
