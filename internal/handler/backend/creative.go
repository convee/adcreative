package backend

import (
	"github.com/convee/adcreative/internal/pkg/app"
	"github.com/convee/adcreative/internal/pkg/code"
	"github.com/convee/adcreative/internal/service"
	"github.com/spf13/cast"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Creative struct {
}

func (c *Creative) List(ctx *gin.Context) {
	var (
		appG = app.Gin{C: ctx}
	)
	page := cast.ToInt(ctx.Query("page"))
	perPage := cast.ToInt(ctx.Query("per_page"))
	creativeService := &service.Creative{
		Page:    page,
		PerPage: perPage,
	}
	data := creativeService.GetList()
	appG.Response(http.StatusOK, code.SUCCESS, data)
	return
}

//// Upload 手动送审
//func (c *Creative) Upload(ctx *gin.Context) {
//	var (
//		appG = app.Gin{C: ctx}
//	)
//	cids := ctx.Query("cids")
//	puid := cast.ToInt(ctx.Query("puid"))
//	creativeIds := strings.Split(cids, ",")
//	for _, creativeId := range creativeIds {
//		logger := log.GetLogger().
//			With(zap.String("method", "Upload")).
//			With(zap.String("request_id", middleware.GetRequestIDFromHeaders(ctx))).
//			With(zap.String("creative_id", creativeId))
//		// 添加到送审chan中
//		task.CrUploadProducer.Producer(crea, logger)
//	}
//	appG.Response(http.StatusOK, code.SUCCESS, nil)
//	return
//}
