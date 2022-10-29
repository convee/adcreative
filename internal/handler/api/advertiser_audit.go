package api

import (
	model2 "github.com/convee/adcreative/internal/model"
	"github.com/convee/adcreative/internal/pkg/app"
	"github.com/convee/adcreative/internal/pkg/code"
	"github.com/convee/adcreative/internal/routers/middleware"
	"github.com/convee/adcreative/internal/service"
	"github.com/convee/adcreative/internal/task"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AdvertiserAudit struct {
}

type UploadAdvertiserAuditJson struct {
	IsOnlyValid         int                  `json:"is_only_valid"`
	Publisher           string               `json:"publisher" validate:"required"`
	AdvertiserName      string               `json:"advertiser_name" validate:"required"`
	AdvertiserAuditInfo *AdvertiserAuditInfo `json:"advertiser_audit_info" validate:"required,min=1,dive,required"`
}

type SaveAdvertiserAuditJson struct {
	Publisher      string `json:"publisher" validate:"required"`
	AdvertiserName string `json:"advertiser_name"`
	AdvertiserId   int    `json:"advertiser_id" validate:"required"`
	Extra          string `json:"extra" validate:"required"`
}

type AdvertiserAuditInfo struct {
	CompanyName      string                 `json:"company_name" validate:""`
	CompanySummary   string                 `json:"company_summary" validate:""`
	WebsiteName      string                 `json:"website_name" validate:""`
	WebsiteAddress   string                 `json:"website_address" validate:""`
	WebsiteNumber    string                 `json:"website_number" validate:""`
	BusinessLicenser string                 `json:"business_licenser" validate:""`
	AuthorizeState   string                 `json:"authorize_state" validate:""`
	Industry         string                 `json:"industry" validate:""`
	Qualifications   []*Qualification       `json:"qualifications" validate:""`
	Extra            map[string]interface{} `json:"extra" validate:""`
}

type Qualification struct {
	FileName string `json:"file_name" validate:""`
	FileUrl  string `json:"file_url" validate:""`
}

func (c *AdvertiserAudit) Upload(ctx *gin.Context) {

	var (
		uploadAdvertiserAuditJson UploadAdvertiserAuditJson
		appG                      = app.Gin{C: ctx}
		results                   []map[string]interface{}
	)

	customer, _ := ctx.MustGet("customer").(*model2.Customer)
	validateErr := app.BindJson(ctx, &uploadAdvertiserAuditJson)
	if len(validateErr) > 0 {
		appG.Response(http.StatusOK, code.INVALID_PARAMS, validateErr)
		return
	}
	errMsg := ""
	var qualifications []model2.Qualification
	for _, q := range uploadAdvertiserAuditJson.AdvertiserAuditInfo.Qualifications {
		qualifications = append(qualifications, model2.Qualification{
			FileName: q.FileName,
			FileUrl:  q.FileUrl,
		})
	}
	advertiserAuditInfo := map[string]interface{}{
		"company_name":      uploadAdvertiserAuditJson.AdvertiserAuditInfo.CompanyName,
		"company_summary":   uploadAdvertiserAuditJson.AdvertiserAuditInfo.CompanySummary,
		"website_name":      uploadAdvertiserAuditJson.AdvertiserAuditInfo.WebsiteName,
		"website_address":   uploadAdvertiserAuditJson.AdvertiserAuditInfo.WebsiteAddress,
		"website_number":    uploadAdvertiserAuditJson.AdvertiserAuditInfo.WebsiteNumber,
		"business_licenser": uploadAdvertiserAuditJson.AdvertiserAuditInfo.BusinessLicenser,
		"authorize_state":   uploadAdvertiserAuditJson.AdvertiserAuditInfo.AuthorizeState,
		"industry":          uploadAdvertiserAuditJson.AdvertiserAuditInfo.Industry,
		"qualifications":    qualifications,
		"extra":             uploadAdvertiserAuditJson.AdvertiserAuditInfo.Extra,
	}
	advertiserAuditService := &service.AdvertiserAudit{
		IsOnlyValid:    uploadAdvertiserAuditJson.IsOnlyValid,
		Publisher:      uploadAdvertiserAuditJson.Publisher,
		AdvertiserName: uploadAdvertiserAuditJson.AdvertiserName,
		CustomerId:     customer.Id,
		Info:           advertiserAuditInfo,
		CheckRules:     true,
	}
	advertiserAudit, err := advertiserAuditService.CreateOrUpdate()
	if err != nil {
		errMsg = "送审失败"
		results = append(results, map[string]interface{}{
			"publisher":     uploadAdvertiserAuditJson.Publisher,
			"advertiser_id": nil,
			"err_code":      1,
			"err_msg":       err.Error(),
		})
		appG.Response(http.StatusOK, code.SUCCESS, results)
		return
	}
	// 添加到送审chan中
	if uploadAdvertiserAuditJson.IsOnlyValid != 1 {
		go task.AdvUploadProducer.Producer(advertiserAudit.Id, middleware.GetRequestIDFromHeaders(ctx))
	}
	results = append(results, map[string]interface{}{
		"publisher":     uploadAdvertiserAuditJson.Publisher,
		"advertiser_id": advertiserAudit.Id,
		"err_code":      0,
		"err_msg":       errMsg,
	})
	appG.Response(http.StatusOK, code.SUCCESS, results)
	return
}

func (c *AdvertiserAudit) Save(ctx *gin.Context) {

	var (
		saveAdvertiserAuditJson SaveAdvertiserAuditJson
		appG                    = app.Gin{C: ctx}
		results                 []map[string]interface{}
	)

	customer, _ := ctx.MustGet("customer").(*model2.Customer)
	validateErr := app.BindJson(ctx, &saveAdvertiserAuditJson)
	if len(validateErr) > 0 {
		appG.Response(http.StatusOK, code.INVALID_PARAMS, validateErr)
		return
	}
	errMsg := ""

	advertiserAuditService := &service.AdvertiserAudit{
		Publisher:      saveAdvertiserAuditJson.Publisher,
		AdvertiserName: saveAdvertiserAuditJson.AdvertiserName,
		CustomerId:     customer.Id,
		Extra:          saveAdvertiserAuditJson.Extra,
		CheckRules:     false,
		IsRsync:        1,
		ErrCode:        7,
		Status:         1,
		ErrMsg:         "CMS同步",
		AdvertiserId:   saveAdvertiserAuditJson.AdvertiserId,
	}
	_, err := advertiserAuditService.CreateOrUpdate()
	if err != nil {
		errMsg = "送审失败"
		results = append(results, map[string]interface{}{
			"publisher":     saveAdvertiserAuditJson.Publisher,
			"advertiser_id": saveAdvertiserAuditJson.AdvertiserId,
			"err_code":      1,
			"err_msg":       err.Error(),
		})
		appG.Response(http.StatusOK, code.SUCCESS, results)
		return
	}

	results = append(results, map[string]interface{}{
		"publisher":     saveAdvertiserAuditJson.Publisher,
		"advertiser_id": saveAdvertiserAuditJson.AdvertiserId,
		"err_code":      0,
		"err_msg":       errMsg,
	})
	appG.Response(http.StatusOK, code.SUCCESS, results)
	return
}

type QueryAdvertiserAuditJson struct {
	//Publisher    string `json:"publisher" validate:"required"`
	AdvertiserId []int `json:"advertiser_id" validate:"required,min=1,max=20,dive,min=1"`
}

func (c *AdvertiserAudit) Query(ctx *gin.Context) {

	var (
		queryAdvertiserAuditJson QueryAdvertiserAuditJson
		appG                     = app.Gin{C: ctx}
		results                  []map[string]interface{}
	)
	customer, _ := ctx.MustGet("customer").(*model2.Customer)

	validateErr := app.BindJson(ctx, &queryAdvertiserAuditJson)
	if len(validateErr) > 0 {
		appG.Response(http.StatusOK, code.INVALID_PARAMS, validateErr)
		return
	}

	advertiserAuditService := &service.AdvertiserAudit{
		//Publisher:     queryAdvertiserAuditJson.Publisher,
		AdvertiserIds: queryAdvertiserAuditJson.AdvertiserId,
		CustomerId:    customer.Id,
	}
	advertiserAudits := advertiserAuditService.GetAll()
	advertiserAuditMaps := make(map[int]model2.AdvertiserAudit)
	for _, advertiserAudit := range advertiserAudits {
		advertiserAuditMaps[advertiserAudit.Id] = *advertiserAudit
	}
	for _, advertiserId := range queryAdvertiserAuditJson.AdvertiserId {
		if advertiserAudit, ok := advertiserAuditMaps[advertiserId]; ok {
			results = append(results, map[string]interface{}{
				//"publisher":     queryAdvertiserAuditJson.Publisher,
				"advertiser_id": advertiserId,
				"status":        advertiserAudit.Status,
				"reason":        advertiserAudit.ErrMsg,
			})
		} else {
			results = append(results, map[string]interface{}{
				//"publisher":     queryAdvertiserAuditJson.Publisher,
				"advertiser_id": advertiserId,
				"status":        -1,
				"reason":        "广告主不存在",
			})
		}
	}
	appG.Response(http.StatusOK, code.SUCCESS, results)
	return
}
