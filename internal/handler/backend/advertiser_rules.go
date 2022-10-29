package backend

import (
	"github.com/convee/adcreative/internal/service"
	"net/http"

	"github.com/convee/adcreative/internal/pkg/app"
	"github.com/convee/adcreative/internal/pkg/code"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cast"

	"github.com/gin-gonic/gin"
)

type AdvertiserRules struct {
}

func (c *AdvertiserRules) List(ctx *gin.Context) {
	var (
		appG = app.Gin{C: ctx}
	)
	page := cast.ToInt(ctx.Query("page"))
	perPage := cast.ToInt(ctx.Query("per_page"))
	advertiserRulesService := &service.AdvertiserRules{
		Page:    page,
		PerPage: perPage,
	}
	data := advertiserRulesService.GetList()
	appG.Response(http.StatusOK, code.SUCCESS, data)
	return
}

type AddAdvertiserRulesForm struct {
	PublisherId *int   `form:"publisher_id" validate:"required"`
	Info        string `form:"info" validate:"required"`
}

func (c *AdvertiserRules) Add(ctx *gin.Context) {

	var (
		addAdvertiserRulesForm AddAdvertiserRulesForm
		appG                   = app.Gin{C: ctx}
	)
	validateErr := app.BindForm(ctx, &addAdvertiserRulesForm)
	if len(validateErr) > 0 {
		appG.BackendResponse(code.INVALID_PARAMS, validateErr[0], nil)
		return
	}
	advertiserRulesService := &service.AdvertiserRules{
		PublisherId: addAdvertiserRulesForm.PublisherId,
		Info:        addAdvertiserRulesForm.Info,
	}
	result, err := advertiserRulesService.Add()
	if err != nil {
		appG.Response(http.StatusOK, code.ERROR_CREATE_FAILED, nil)
		return
	}
	appG.Response(http.StatusOK, code.SUCCESS, result)
	return
}

type EditAdvertiserRulesForm struct {
	Id          int    `form:"id" validate:"required"`
	PublisherId *int   `form:"publisher_id" validate:"required"`
	Info        string `form:"info" validate:"required"`
}

func (c *AdvertiserRules) Edit(ctx *gin.Context) {
	var (
		editAdvertiserRulesForm EditAdvertiserRulesForm
		appG                    = app.Gin{C: ctx}
	)
	validateErr := app.BindForm(ctx, &editAdvertiserRulesForm)
	if len(validateErr) > 0 {
		appG.BackendResponse(code.INVALID_PARAMS, validateErr[0], nil)
		return
	}
	advertiserRulesService := &service.AdvertiserRules{
		Id:          editAdvertiserRulesForm.Id,
		PublisherId: editAdvertiserRulesForm.PublisherId,
		Info:        editAdvertiserRulesForm.Info,
	}
	result, err := advertiserRulesService.Edit()
	if err != nil {
		appG.Response(http.StatusOK, code.ERROR_UPDATE_FAILED, nil)
		return
	}
	appG.Response(http.StatusOK, code.SUCCESS, result)
	return
}

func (c *AdvertiserRules) Delete(ctx *gin.Context) {
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
	advertiserRulesService := &service.AdvertiserRules{
		Id: id,
	}
	advertiser, err := advertiserRulesService.Delete(id)
	if err != nil {
		appG.Response(http.StatusOK, code.ERROR_DELETE_FAILED, nil)
		return
	}
	appG.Response(http.StatusOK, code.SUCCESS, advertiser)
	return
}
