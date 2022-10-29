package backend

import (
	"github.com/convee/adcreative/internal/pkg/app"
	"github.com/convee/adcreative/internal/pkg/code"
	"github.com/convee/adcreative/internal/service"
	logger "github.com/convee/adcreative/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"net/http"
)

type Position struct {
}

func (c *Position) List(ctx *gin.Context) {
	var (
		appG = app.Gin{C: ctx}
	)
	page := cast.ToInt(ctx.Query("page"))
	perPage := cast.ToInt(ctx.Query("per_page"))
	name := ctx.Query("name")
	position := ctx.Query("position")
	id := ctx.Query("id")
	publisherId := cast.ToInt(ctx.Query("publisher_id"))
	positionService := &service.Position{
		Page:        page,
		PerPage:     perPage,
		Name:        name,
		Position:    position,
		Id:          cast.ToInt(id),
		PublisherId: publisherId,
	}
	data := positionService.GetList()
	appG.Response(http.StatusOK, code.SUCCESS, data)
	return
}

type AddPositionForm struct {
	PublisherId                    int    `form:"publisher_id" validate:"required,numeric,exists_publisher"`
	Name                           string `form:"name" validate:"required"`
	Type                           string `form:"type" validate:"required"`
	Position                       string `form:"position" validate:"required"`
	MaterialInfo                   string `form:"material_info" validate:"required"`
	MediaType                      string `form:"media_type" validate:"required"`
	AdFormat                       *int   `form:"ad_format" validate:"required"`
	IsSupportDeeplink              *int   `form:"is_support_deeplink" validate:"required"`
	LandingChangeNeedRsync         *int   `form:"landing_change_need_rsync" validate:"required"`
	MonitorCodeChangeNeedRsync     *int   `form:"monitor_code_change_need_rsync" validate:"required"`
	MonitorPositionChangeNeedRsync *int   `form:"monitor_position_change_need_rsync" validate:"required"`
	IsCreativeBind                 *int   `form:"is_creative_bind" validate:"required"`
	PvLimit                        int    `form:"pv_limit"`
	ClLimit                        int    `form:"cl_limit"`
}

func (c *Position) Add(ctx *gin.Context) {

	var (
		addPositionForm AddPositionForm
		appG            = app.Gin{C: ctx}
	)
	validateErr := app.BindForm(ctx, &addPositionForm)
	if len(validateErr) > 0 {
		appG.BackendResponse(code.INVALID_PARAMS, validateErr[0], nil)
		return
	}
	positionService := &service.Position{
		PublisherId:                    addPositionForm.PublisherId,
		Name:                           addPositionForm.Name,
		Type:                           addPositionForm.Type,
		Position:                       addPositionForm.Position,
		MaterialInfo:                   addPositionForm.MaterialInfo,
		MediaType:                      addPositionForm.MediaType,
		AdFormat:                       addPositionForm.AdFormat,
		IsSupportDeeplink:              addPositionForm.IsSupportDeeplink,
		LandingChangeNeedRsync:         addPositionForm.LandingChangeNeedRsync,
		MonitorCodeChangeNeedRsync:     addPositionForm.MonitorCodeChangeNeedRsync,
		MonitorPositionChangeNeedRsync: addPositionForm.MonitorPositionChangeNeedRsync,
		IsCreativeBind:                 addPositionForm.IsCreativeBind,
		PvLimit:                        addPositionForm.PvLimit,
		ClLimit:                        addPositionForm.ClLimit,
	}
	result, err := positionService.Add()
	if err != nil {
		logger.Error("position_add_err", zap.Error(err))
		appG.Response(http.StatusOK, code.ERROR_CREATE_FAILED, nil)
		return
	}
	appG.Response(http.StatusOK, code.SUCCESS, result)
	return
}

type EditPositionForm struct {
	Id                             int    `form:"id" validate:"required"`
	PublisherId                    int    `form:"publisher_id" validate:"required,numeric,exists_publisher"`
	Name                           string `form:"name" validate:"required"`
	Type                           string `form:"type" validate:"required"`
	Position                       string `form:"position" validate:"required"`
	MaterialInfo                   string `form:"material_info" validate:"required"`
	MediaType                      string `form:"media_type" validate:"required"`
	AdFormat                       *int   `form:"ad_format" validate:"required"`
	IsSupportDeeplink              *int   `form:"is_support_deeplink" validate:"required"`
	LandingChangeNeedRsync         *int   `form:"landing_change_need_rsync" validate:"required"`
	MonitorCodeChangeNeedRsync     *int   `form:"monitor_code_change_need_rsync" validate:"required"`
	MonitorPositionChangeNeedRsync *int   `form:"monitor_position_change_need_rsync" validate:"required"`
	IsCreativeBind                 *int   `form:"is_creative_bind" validate:"required"`
	PvLimit                        int    `form:"pv_limit"`
	ClLimit                        int    `form:"cl_limit"`
}

func (c *Position) Edit(ctx *gin.Context) {
	var (
		editPositionForm EditPositionForm
		appG             = app.Gin{C: ctx}
	)
	validateErr := app.BindForm(ctx, &editPositionForm)
	if len(validateErr) > 0 {
		appG.BackendResponse(code.INVALID_PARAMS, validateErr[0], nil)
		return
	}
	positionService := &service.Position{
		Id:                             editPositionForm.Id,
		PublisherId:                    editPositionForm.PublisherId,
		Name:                           editPositionForm.Name,
		Type:                           editPositionForm.Type,
		Position:                       editPositionForm.Position,
		MaterialInfo:                   editPositionForm.MaterialInfo,
		MediaType:                      editPositionForm.MediaType,
		AdFormat:                       editPositionForm.AdFormat,
		IsSupportDeeplink:              editPositionForm.IsSupportDeeplink,
		LandingChangeNeedRsync:         editPositionForm.LandingChangeNeedRsync,
		MonitorCodeChangeNeedRsync:     editPositionForm.MonitorCodeChangeNeedRsync,
		MonitorPositionChangeNeedRsync: editPositionForm.MonitorPositionChangeNeedRsync,
		IsCreativeBind:                 editPositionForm.IsCreativeBind,
		PvLimit:                        editPositionForm.PvLimit,
		ClLimit:                        editPositionForm.ClLimit,
	}
	result, err := positionService.Edit()
	if err != nil {
		logger.Error("position_save_err", zap.Error(err))
		appG.Response(http.StatusOK, code.ERROR_UPDATE_FAILED, nil)
		return
	}
	appG.Response(http.StatusOK, code.SUCCESS, result)
	return
}

func (c *Position) Delete(ctx *gin.Context) {
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
	positionService := &service.Position{
		Id: id,
	}
	position, err := positionService.Delete()
	if err != nil {
		appG.Response(http.StatusOK, code.ERROR_DELETE_FAILED, nil)
		return
	}
	appG.Response(http.StatusOK, code.SUCCESS, position)
	return
}
