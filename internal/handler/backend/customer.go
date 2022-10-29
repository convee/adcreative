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

type Customer struct {
}

func (c *Customer) List(ctx *gin.Context) {
	var (
		appG = app.Gin{C: ctx}
	)
	page := cast.ToInt(ctx.Query("page"))
	perPage := cast.ToInt(ctx.Query("per_page"))
	name := ctx.Query("name")
	customerService := &service.Customer{
		Page:    page,
		PerPage: perPage,
		Name:    name,
	}
	data := customerService.GetList()
	appG.Response(http.StatusOK, code.SUCCESS, data)
	return
}

type AddCustomerForm struct {
	Name                  string `form:"name" validate:"required,unique_customer"`
	IsPrivate             *int   `form:"is_private" validate:"oneof=0 1"`
	CreativeCallbackUrl   string `form:"creative_callback_url" validate:"required"`
	AdvertiserCallbackUrl string `form:"advertiser_callback_url" validate:"required"`
}

func (c *Customer) Add(ctx *gin.Context) {

	var (
		addCustomerForm AddCustomerForm
		appG            = app.Gin{C: ctx}
	)
	validateErr := app.BindForm(ctx, &addCustomerForm)
	if len(validateErr) > 0 {
		appG.BackendResponse(code.INVALID_PARAMS, validateErr[0], nil)
		return
	}
	customerService := &service.Customer{
		IsPrivate:             addCustomerForm.IsPrivate,
		Name:                  addCustomerForm.Name,
		CreativeCallbackUrl:   addCustomerForm.CreativeCallbackUrl,
		AdvertiserCallbackUrl: addCustomerForm.AdvertiserCallbackUrl,
	}
	result, err := customerService.Add()
	if err != nil {
		logger.Error("customer_add_err", zap.Error(err))
		appG.Response(http.StatusOK, code.ERROR_CREATE_FAILED, nil)
		return
	}
	appG.Response(http.StatusOK, code.SUCCESS, result)
	return
}

type EditCustomerForm struct {
	Id                    int    `form:"id" validate:"required"`
	Name                  string `form:"name" validate:"required"`
	IsPrivate             *int   `form:"is_private" validate:"oneof=0 1"`
	CreativeCallbackUrl   string `form:"creative_callback_url" validate:"required"`
	AdvertiserCallbackUrl string `form:"advertiser_callback_url" validate:"required"`
}

func (c *Customer) Edit(ctx *gin.Context) {
	var (
		editCustomerForm EditCustomerForm
		appG             = app.Gin{C: ctx}
	)
	validateErr := app.BindForm(ctx, &editCustomerForm)
	if len(validateErr) > 0 {
		appG.BackendResponse(code.INVALID_PARAMS, validateErr[0], nil)
		return
	}
	customerService := &service.Customer{
		Id:                    editCustomerForm.Id,
		IsPrivate:             editCustomerForm.IsPrivate,
		Name:                  editCustomerForm.Name,
		CreativeCallbackUrl:   editCustomerForm.CreativeCallbackUrl,
		AdvertiserCallbackUrl: editCustomerForm.AdvertiserCallbackUrl,
	}
	result, err := customerService.Edit()
	if err != nil {
		logger.Error("customer_save_err", zap.Error(err))
		appG.Response(http.StatusOK, code.ERROR_UPDATE_FAILED, nil)
		return
	}
	appG.Response(http.StatusOK, code.SUCCESS, result)
	return
}

func (c *Customer) Delete(ctx *gin.Context) {
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
	customerService := &service.Customer{
		Id: id,
	}
	customer, err := customerService.Delete(id)
	if err != nil {
		logger.Error("customer_del_err", zap.Error(err))
		appG.Response(http.StatusOK, code.ERROR_DELETE_FAILED, nil)
		return
	}
	appG.Response(http.StatusOK, code.SUCCESS, customer)
	return
}
