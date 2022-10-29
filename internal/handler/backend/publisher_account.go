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

type PublisherAccount struct {
}

func (pa *PublisherAccount) List(ctx *gin.Context) {
	var (
		validate = validator.New()
		appG     = app.Gin{C: ctx}
		err      error
	)
	page := cast.ToInt(ctx.Query("page"))
	perPage := cast.ToInt(ctx.Query("per_page"))
	err = validate.Var(page, "required,gte=1")
	if err != nil {
		appG.Response(http.StatusOK, code.INVALID_PARAMS, nil)
		return
	}
	err = validate.Var(perPage, "required,gte=1")
	if err != nil {
		appG.Response(http.StatusOK, code.INVALID_PARAMS, nil)
		return
	}
	publisherAccountService := &service.PublisherAccount{
		Page:    page,
		PerPage: perPage,
	}
	data := publisherAccountService.GetList()
	appG.Response(http.StatusOK, code.SUCCESS, data)
	return
}

type AddPublisherAccountForm struct {
	DspId       string `form:"dsp_id" validate:"required"`
	Token       string `form:"token" validate:"required"`
	PublisherId int    `form:"publisher_id" validate:"required,numeric,exists_publisher"`
	CustomerId  int    `form:"customer_id" validate:"required,numeric"`
	CallbackUrl string `form:"callback_url" validate:"url"`
	Remark      string `form:"remark"`
}

func (pa *PublisherAccount) Add(ctx *gin.Context) {

	var (
		addPublisherAccountForm AddPublisherAccountForm
		appG                    = app.Gin{C: ctx}
	)
	validateErr := app.BindForm(ctx, &addPublisherAccountForm)
	if len(validateErr) > 0 {
		appG.BackendResponse(code.INVALID_PARAMS, validateErr[0], nil)
		return
	}
	publisherAccountService := &service.PublisherAccount{
		DspId:       addPublisherAccountForm.DspId,
		Token:       addPublisherAccountForm.Token,
		CustomerId:  addPublisherAccountForm.CustomerId,
		PublisherId: addPublisherAccountForm.PublisherId,
		CallbackUrl: addPublisherAccountForm.CallbackUrl,
		Remark:      addPublisherAccountForm.Remark,
	}
	result, err := publisherAccountService.Add()
	if err != nil {
		appG.Response(http.StatusOK, code.ERROR_CREATE_FAILED, nil)
		return
	}
	appG.Response(http.StatusOK, code.SUCCESS, result)
	return
}

type EditPublisherAccountForm struct {
	Id          int    `form:"id" validate:"required,numeric"`
	DspId       string `form:"dsp_id" validate:"required"`
	Token       string `form:"token" validate:"required"`
	PublisherId int    `form:"publisher_id" validate:"required,numeric,exists_publisher"`
	CustomerId  int    `form:"customer_id" validate:"required,numeric"`
	CallbackUrl string `form:"callback_url" validate:"url"`
	Remark      string `form:"remark"`
}

func (pa *PublisherAccount) Edit(ctx *gin.Context) {
	var (
		editPublisherAccountForm EditPublisherAccountForm
		appG                     = app.Gin{C: ctx}
	)
	validateErr := app.BindForm(ctx, &editPublisherAccountForm)
	if len(validateErr) > 0 {
		appG.BackendResponse(code.INVALID_PARAMS, validateErr[0], nil)
		return
	}
	publisherAccountService := &service.PublisherAccount{
		Id:          editPublisherAccountForm.Id,
		DspId:       editPublisherAccountForm.DspId,
		Token:       editPublisherAccountForm.Token,
		CustomerId:  editPublisherAccountForm.CustomerId,
		PublisherId: editPublisherAccountForm.PublisherId,
		CallbackUrl: editPublisherAccountForm.CallbackUrl,
		Remark:      editPublisherAccountForm.Remark,
	}
	result, err := publisherAccountService.Edit()
	if err != nil {
		appG.Response(http.StatusOK, code.ERROR_UPDATE_FAILED, nil)
		return
	}
	appG.Response(http.StatusOK, code.SUCCESS, result)
	return
}

func (pa *PublisherAccount) Delete(ctx *gin.Context) {
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
	publisherAccountService := &service.PublisherAccount{
		Id: id,
	}
	publisherAccount, err := publisherAccountService.Delete(id)
	if err != nil {
		appG.Response(http.StatusOK, code.ERROR_DELETE_FAILED, nil)
		return
	}
	appG.Response(http.StatusOK, code.SUCCESS, publisherAccount)
	return
}
