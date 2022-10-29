package backend

import (
	"github.com/convee/adcreative/internal/service"
	logger "github.com/convee/adcreative/pkg/log"
	"go.uber.org/zap"
	"net/http"

	"github.com/convee/adcreative/internal/pkg/app"
	"github.com/convee/adcreative/internal/pkg/code"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cast"

	"github.com/gin-gonic/gin"
)

type Publisher struct {
}

func (c *Publisher) List(ctx *gin.Context) {
	var (
		appG = app.Gin{C: ctx}
	)
	page := cast.ToInt(ctx.Query("page"))
	perPage := cast.ToInt(ctx.Query("per_page"))
	name := ctx.Query("name")
	id := ctx.Query("id")
	publisherService := &service.Publisher{
		Page:    page,
		PerPage: perPage,
		Name:    name,
		Id:      cast.ToInt(id),
	}
	data := publisherService.GetList()
	appG.Response(http.StatusOK, code.SUCCESS, data)
	return
}

type AddPublisherForm struct {
	Name                           string `form:"name" validate:"required"`
	IsRsyncAdvertiser              *int   `form:"is_rsync_advertiser" validate:"required"`
	IsRsyncCreative                *int   `form:"is_rsync_creative" validate:"required"`
	IsPublisherCdn                 *int   `form:"is_publisher_cdn" validate:"required"`
	IsCreativeBind                 *int   `form:"is_creative_bind" validate:"required"`
	MonitorCodeChangeNeedRsync     *int   `form:"monitor_code_change_need_rsync" validate:"required"`
	LandingChangeNeedRsync         *int   `form:"landing_change_need_rsync" validate:"required"`
	MonitorPositionChangeNeedRsync *int   `form:"monitor_position_change_need_rsync" validate:"required"`
	S2sStateInfo                   string `form:"s2s_state_info"`
	PubReturnInfo                  string `form:"pub_return_info"`
	PvLimit                        int    `form:"pv_limit"`
	ClLimit                        int    `form:"cl_limit"`
	Nickname                       string `form:"nickname"`
	AdvertiserUrls                 string `form:"advertiser_urls"`
	CreativeUrls                   string `form:"creative_urls"`
}

func (c *Publisher) Add(ctx *gin.Context) {

	var (
		addPublisherForm AddPublisherForm
		appG             = app.Gin{C: ctx}
	)
	validateErr := app.BindForm(ctx, &addPublisherForm)
	if len(validateErr) > 0 {
		appG.BackendResponse(code.INVALID_PARAMS, validateErr[0], nil)
		return
	}
	if addPublisherForm.AdvertiserUrls == "" {
		addPublisherForm.AdvertiserUrls = "{}"
	}
	if addPublisherForm.CreativeUrls == "" {
		addPublisherForm.CreativeUrls = "{}"
	}
	publisherService := &service.Publisher{
		Name:                           addPublisherForm.Name,
		IsRsyncAdvertiser:              addPublisherForm.IsRsyncAdvertiser,
		IsRsyncCreative:                addPublisherForm.IsRsyncCreative,
		IsPublisherCdn:                 addPublisherForm.IsPublisherCdn,
		IsCreativeBind:                 addPublisherForm.IsCreativeBind,
		MonitorCodeChangeNeedRsync:     addPublisherForm.MonitorCodeChangeNeedRsync,
		LandingChangeNeedRsync:         addPublisherForm.LandingChangeNeedRsync,
		MonitorPositionChangeNeedRsync: addPublisherForm.MonitorPositionChangeNeedRsync,
		S2sStateInfo:                   addPublisherForm.S2sStateInfo,
		PubReturnInfo:                  addPublisherForm.PubReturnInfo,
		PvLimit:                        addPublisherForm.PvLimit,
		ClLimit:                        addPublisherForm.ClLimit,
		Nickname:                       addPublisherForm.Nickname,
		AdvertiserUrls:                 addPublisherForm.AdvertiserUrls,
		CreativeUrls:                   addPublisherForm.CreativeUrls,
	}
	result, err := publisherService.Add()
	if err != nil {
		appG.Response(http.StatusOK, code.ERROR_CREATE_FAILED, nil)
		return
	}
	appG.Response(http.StatusOK, code.SUCCESS, result)
	return
}

type EditPublisherForm struct {
	Id                             int    `form:"id" validate:"required"`
	Name                           string `form:"name" validate:"required"`
	IsRsyncAdvertiser              *int   `form:"is_rsync_advertiser" validate:"required"`
	IsRsyncCreative                *int   `form:"is_rsync_creative" validate:"required"`
	IsPublisherCdn                 *int   `form:"is_publisher_cdn" validate:"required"`
	IsCreativeBind                 *int   `form:"is_creative_bind" validate:"required"`
	MonitorCodeChangeNeedRsync     *int   `form:"monitor_code_change_need_rsync" validate:"required"`
	LandingChangeNeedRsync         *int   `form:"landing_change_need_rsync" validate:"required"`
	MonitorPositionChangeNeedRsync *int   `form:"monitor_position_change_need_rsync" validate:"required"`
	S2sStateInfo                   string `form:"s2s_state_info"`
	PubReturnInfo                  string `form:"pub_return_info"`
	PvLimit                        int    `form:"pv_limit"`
	ClLimit                        int    `form:"cl_limit"`
	Nickname                       string `form:"nickname"`
	AdvertiserUrls                 string `form:"advertiser_urls"`
	CreativeUrls                   string `form:"creative_urls"`
}

func (c *Publisher) Edit(ctx *gin.Context) {
	var (
		editPublisherForm EditPublisherForm
		appG              = app.Gin{C: ctx}
	)
	validateErr := app.BindForm(ctx, &editPublisherForm)
	if len(validateErr) > 0 {
		appG.BackendResponse(code.INVALID_PARAMS, validateErr[0], nil)
		return
	}
	if editPublisherForm.AdvertiserUrls == "" {
		editPublisherForm.AdvertiserUrls = "{}"
	}
	if editPublisherForm.CreativeUrls == "" {
		editPublisherForm.CreativeUrls = "{}"
	}
	publisherService := &service.Publisher{
		Id:                             editPublisherForm.Id,
		Name:                           editPublisherForm.Name,
		IsRsyncAdvertiser:              editPublisherForm.IsRsyncAdvertiser,
		IsRsyncCreative:                editPublisherForm.IsRsyncCreative,
		IsPublisherCdn:                 editPublisherForm.IsPublisherCdn,
		IsCreativeBind:                 editPublisherForm.IsCreativeBind,
		MonitorCodeChangeNeedRsync:     editPublisherForm.MonitorCodeChangeNeedRsync,
		LandingChangeNeedRsync:         editPublisherForm.LandingChangeNeedRsync,
		MonitorPositionChangeNeedRsync: editPublisherForm.MonitorPositionChangeNeedRsync,
		S2sStateInfo:                   editPublisherForm.S2sStateInfo,
		PubReturnInfo:                  editPublisherForm.PubReturnInfo,
		PvLimit:                        editPublisherForm.PvLimit,
		ClLimit:                        editPublisherForm.ClLimit,
		Nickname:                       editPublisherForm.Nickname,
		AdvertiserUrls:                 editPublisherForm.AdvertiserUrls,
		CreativeUrls:                   editPublisherForm.CreativeUrls,
	}
	result, err := publisherService.Edit()
	if err != nil {
		logger.Error("publisher_edit_err", zap.Error(err))
		appG.Response(http.StatusOK, code.ERROR_UPDATE_FAILED, nil)
		return
	}
	appG.Response(http.StatusOK, code.SUCCESS, result)
	return
}

func (c *Publisher) Delete(ctx *gin.Context) {
	var (
		appG     = app.Gin{C: ctx}
		validate = validator.New()
	)

	id := cast.ToInt(ctx.PostForm("id"))
	err := validate.Var(id, "required,gte=1")
	if err != nil {
		appG.Response(http.StatusOK, code.INVALID_PARAMS, nil)
		return
	}
	publisherService := &service.Publisher{
		Id: id,
	}
	publisher, err := publisherService.Delete()
	if err != nil {
		logger.Error("publisher_delete_err", zap.Error(err))
		appG.Response(http.StatusOK, code.ERROR_DELETE_FAILED, nil)
		return
	}
	appG.Response(http.StatusOK, code.SUCCESS, publisher)
	return
}
