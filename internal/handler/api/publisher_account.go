package api

import (
	"github.com/convee/adcreative/internal/model"
	"github.com/convee/adcreative/internal/service"
	logger "github.com/convee/adcreative/pkg/log"
	"go.uber.org/zap"
	"net/http"

	"github.com/convee/adcreative/internal/pkg/app"
	"github.com/convee/adcreative/internal/pkg/code"
	"github.com/gin-gonic/gin"
)

type PublisherAccount struct {
}

type AddPublisherAccountJson struct {
	DspId       string `json:"dsp_id" validate:"required"`
	Token       string `json:"token" validate:"required"`
	PublisherId int    `json:"publisher_id" validate:"required,numeric,exists_publisher"`
	CallbackUrl string `json:"callback_url" validate:"url"`
}

func (pa *PublisherAccount) Add(ctx *gin.Context) {

	var (
		addPublisherAccountJson AddPublisherAccountJson
		appG                    = app.Gin{C: ctx}
	)
	customer, _ := ctx.MustGet("customer").(*model.Customer)
	validateErr := app.BindJson(ctx, &addPublisherAccountJson)
	if len(validateErr) > 0 {
		appG.Response(http.StatusOK, code.INVALID_PARAMS, validateErr)
		return
	}
	publisherAccountService := &service.PublisherAccount{
		CustomerId:  customer.Id,
		DspId:       addPublisherAccountJson.DspId,
		Token:       addPublisherAccountJson.Token,
		PublisherId: addPublisherAccountJson.PublisherId,
		CallbackUrl: addPublisherAccountJson.CallbackUrl,
	}
	result, err := publisherAccountService.CreateOrUpdate()
	if err != nil {
		logger.Error("api_publisher_account_add_err", zap.Error(err))
		appG.Response(http.StatusOK, code.ERROR_CREATE_FAILED, nil)
		return
	}
	appG.Response(http.StatusOK, code.SUCCESS, result)
	return
}
