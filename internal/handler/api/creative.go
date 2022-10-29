package api

import (
	model2 "github.com/convee/adcreative/internal/model"
	"github.com/convee/adcreative/internal/service"
	"net/http"
	"sync"
	"time"

	"github.com/convee/adcreative/internal/pkg/cache"
	"github.com/convee/adcreative/internal/routers/middleware"
	"github.com/convee/adcreative/pkg/log"
	"go.uber.org/zap"

	"github.com/convee/adcreative/internal/pkg/app"
	"github.com/convee/adcreative/internal/pkg/code"
	"github.com/convee/adcreative/internal/task"

	"github.com/gin-gonic/gin"
)

const (
	creativeUploadSuccess = 0
	creativeUploadFailed  = 1
)

type uploadResult struct {
	Position   string `json:"position"`
	CreativeId string `json:"creative_id"`
	ErrCode    int    `json:"err_code"`
	MediaCid   string `json:"media_cid"`
	ErrMsg     string `json:"err_msg"`
}
type Creative struct {
}

type UploadCreativeJson struct {
	Material []*Material `json:"material" validate:"required,min=1,dive,required"`
}

type BatchCheckJson struct {
	Material  []*CheckCreativeJson `json:"material" validate:"required,min=1,dive,required"`
	CheckType int                  `json:"check_type"` // 0 校验全部 1 校验文本 2校验视频+图片
}

type CheckCreativeJson struct {
	Position    string  `json:"position"`
	Info        []*Info `json:"info" validate:"required,min=1,dive,required"`
	TemplateId  string  `json:"template_id" validate:"required"`
	TemplateKey string  `json:"template_key"`
	UnionKey    string  `json:"union_key"`
}

type Material struct {
	AdvertiserId    int        `json:"advertiser_id"`
	MediaInfo       string     `json:"media_info"`
	Industry        string     `json:"industry"`
	Position        string     `json:"position" validate:"required"`
	CreativeId      string     `json:"creative_id" validate:"required"`
	Name            string     `json:"name" validate:"required"`
	TemplateId      string     `json:"template_id" validate:"required"`
	Action          int        `json:"action" validate:"required,oneof=1 2 3"` //1-打开网页 2-下载 3-deeplink
	LandUrl         string     `json:"land_url"`
	DeeplinkUrl     string     `json:"deeplink_url"`
	MiniProgramId   string     `json:"mini_program_id"`
	MiniProgramPath string     `json:"mini_program_path"`
	StartDate       string     `json:"start_date" time_format:"2006-01-02"` //validate:"required"
	EndDate         string     `json:"end_date" time_format:"2006-01-02"`   //validate:"required,gtefield=StartDate"
	Monitor         []*Monitor `json:"monitor"  validate:"dive,required"`   //validate:"required,min=1,dive,required"`
	Info            []*Info    `json:"info" validate:"required,min=1,dive,required"`
	Cm              []string   `json:"cm" validate:"dive,required"`
	Vm              []string   `json:"vm" validate:"dive,required"`
	Extra           string     `json:"extra"`
}
type Monitor struct {
	T   int    `json:"t"`
	Url string `json:"url" validate:"required,url"`
}
type Info struct {
	AttrName  string `json:"attr_name" validate:"required"`
	AttrValue string `json:"attr_value" validate:"required"`
	Md5       string `json:"md5"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Ext       string `json:"ext"`
	Duration  int    `json:"duration"`
	Size      int    `json:"size"`
}

func (c *Creative) Check(ctx *gin.Context) {

	var (
		checkCreativeJson CheckCreativeJson
		appG              = app.Gin{C: ctx}
		results           map[string]interface{}
		errs              []string
	)
	customer, _ := ctx.MustGet("customer").(*model2.Customer)
	errMsg := app.BindJson(ctx, &checkCreativeJson)
	if len(errMsg) > 0 {
		appG.Response(http.StatusOK, code.INVALID_PARAMS, errMsg)
		return
	}
	var infos []model2.TemplateInfo
	for _, i := range checkCreativeJson.Info {
		infos = append(infos, model2.TemplateInfo{
			AttrName:  i.AttrName,
			AttrValue: i.AttrValue,
			Md5:       i.Md5,
			Width:     i.Width,
			Height:    i.Height,
			Duration:  i.Duration,
			Size:      i.Size,
		})
	}
	creativeService := &service.Creative{
		Position:    checkCreativeJson.Position,
		CustomerId:  customer.Id,
		Info:        infos,
		TemplateId:  checkCreativeJson.TemplateId,
		TemplateKey: checkCreativeJson.TemplateKey,
	}
	results = make(map[string]interface{})
	results["position"] = checkCreativeJson.Position
	results["err_code"] = 0
	results["err_msg"] = ""
	positionInfo, err := cache.GetPositionCache(checkCreativeJson.Position)
	if err != nil {
		errs = append(errs, err.Error())
		results["err_code"] = 1
		results["err_msg"] = errs
		appG.Response(http.StatusOK, code.SUCCESS, results)
		return
	}
	info, errs := creativeService.Check(positionInfo, 0)
	if len(errs) > 0 {
		results["err_code"] = 1
		results["err_msg"] = errs
		appG.Response(http.StatusOK, code.SUCCESS, results)
		return
	}
	template, _ := creativeService.GetTemplate(positionInfo)
	results["info"] = info
	results["extra"] = template.Extra

	appG.Response(http.StatusOK, code.SUCCESS, results)
	return
}
func (c *Creative) BatchCheck(ctx *gin.Context) {

	var (
		checkCreativeJsons BatchCheckJson
		appG               = app.Gin{C: ctx}
		results            map[string]interface{}
		//errs               []string
	)
	customer, _ := ctx.MustGet("customer").(*model2.Customer)
	errMsg := app.BindJson(ctx, &checkCreativeJsons)
	if len(errMsg) > 0 {
		appG.Response(http.StatusOK, code.INVALID_PARAMS, errMsg)
		return
	}
	var lists []map[string]interface{}
	m := sync.Map{}
	wg := sync.WaitGroup{}
	results = make(map[string]interface{})
	results["err_code"] = 0
	results["err_msg"] = ""
	for index, checkCreativeJson := range checkCreativeJsons.Material {

		wg.Add(1)
		go func(checkCreativeJson *CheckCreativeJson, index int, checkType int) {
			defer wg.Done()

			var infos []model2.TemplateInfo
			for _, i := range checkCreativeJson.Info {
				infos = append(infos, model2.TemplateInfo{
					AttrName:  i.AttrName,
					AttrValue: i.AttrValue,
					Md5:       i.Md5,
					Width:     i.Width,
					Height:    i.Height,
					Duration:  i.Duration,
					Size:      i.Size,
				})
			}
			creativeService := &service.Creative{
				Position:    checkCreativeJson.Position,
				CustomerId:  customer.Id,
				Info:        infos,
				TemplateId:  checkCreativeJson.TemplateId,
				TemplateKey: checkCreativeJson.TemplateKey,
			}
			positionInfo, err := cache.GetPositionCache(checkCreativeJson.Position)
			if err != nil {
				errMsg := map[string]interface{}{
					"err_code":  1,
					"err_msg":   []string{err.Error()},
					"union_key": checkCreativeJson.UnionKey,
				}
				m.Store(index, errMsg)
				return
				//appG.Response(http.StatusOK, code.SUCCESS, results)
			}
			info, errs := creativeService.Check(positionInfo, checkType)
			if len(errs) > 0 {
				errMsg := map[string]interface{}{
					"err_code":  1,
					"err_msg":   errs,
					"union_key": checkCreativeJson.UnionKey,
				}
				m.Store(index, errMsg)
				return
				//appG.Response(http.StatusOK, code.SUCCESS, results)
			}
			template, _ := creativeService.GetTemplate(positionInfo)

			list := make(map[string]interface{})
			results["position"] = checkCreativeJson.Position
			list["info"] = info
			list["err_code"] = 0
			list["extra"] = template.Extra
			list["union_key"] = checkCreativeJson.UnionKey
			m.Store(index, list)
		}(checkCreativeJson, index, checkCreativeJsons.CheckType)

	}
	wg.Wait()
	for index, _ := range checkCreativeJsons.Material {
		if list, ok := m.Load(index); ok {
			lists = append(lists, list.(map[string]interface{}))
		}
	}

	results["list"] = lists

	appG.Response(http.StatusOK, code.SUCCESS, results)
	return
}

func (c *Creative) Upload(ctx *gin.Context) {
	start := time.Now()
	requestId := middleware.GetRequestIDFromHeaders(ctx)

	defer func() {
		costs := time.Now().Sub(start).Milliseconds()
		log.GetLogger().With(zap.String("reqId", requestId)).With(zap.Int64("Creative::Upload costs", costs)).Info("creative_upload_costs")
	}()

	var (
		uploadCreativeJson UploadCreativeJson
		appG               = app.Gin{C: ctx}
		results            []uploadResult
	)
	customer, _ := ctx.MustGet("customer").(*model2.Customer)
	errMsg := app.BindJson(ctx, &uploadCreativeJson)
	if len(errMsg) > 0 {
		appG.Response(http.StatusOK, code.INVALID_PARAMS, errMsg)
		return
	}

	var wg = &sync.WaitGroup{}
	var m = &sync.Map{}
	for index, material := range uploadCreativeJson.Material {
		wg.Add(1)
		go uploadOne(wg, m, requestId, customer, index, material)
	}
	wg.Wait()
	for index := range uploadCreativeJson.Material {
		if result, ok := m.Load(index); ok {
			results = append(results, result.(uploadResult))
		}
	}
	appG.Response(http.StatusOK, code.SUCCESS, results)
	return
}

func uploadOne(wg *sync.WaitGroup, m *sync.Map, requestId string, customer *model2.Customer, index int, material *Material) {
	defer wg.Done()
	var monitors []model2.Monitor
	for _, m := range material.Monitor {
		monitors = append(monitors, model2.Monitor{
			T:   m.T,
			Url: m.Url,
		})
	}
	var infos []model2.TemplateInfo
	for _, i := range material.Info {
		infos = append(infos, model2.TemplateInfo{
			AttrName:  i.AttrName,
			AttrValue: i.AttrValue,
			Md5:       i.Md5,
			Width:     i.Width,
			Height:    i.Height,
			Duration:  i.Duration,
		})
	}
	creativeService := &service.Creative{
		CustomerId:      customer.Id,
		AdvertiserId:    material.AdvertiserId,
		MediaInfo:       material.MediaInfo,
		Industry:        material.Industry,
		Position:        material.Position,
		CreativeId:      material.CreativeId,
		Name:            material.Name,
		TemplateId:      material.TemplateId,
		Action:          material.Action,
		LandUrl:         material.LandUrl,
		DeeplinkUrl:     material.DeeplinkUrl,
		MiniProgramId:   material.MiniProgramId,
		MiniProgramPath: material.MiniProgramPath,
		StartDate:       material.StartDate,
		EndDate:         material.EndDate,
		Monitor:         monitors,
		Info:            infos,
		Cm:              material.Cm,
		Vm:              material.Vm,
		Extra:           material.Extra,
		RequestId:       requestId,
	}

	logger := log.GetLogger().
		With(zap.String("method", "Upload")).
		With(zap.String("request_id", requestId)).
		With(zap.String("creative_id", creativeService.CreativeId)).
		With(zap.Int("customer_id", creativeService.CustomerId)).
		With(zap.Int("advertiser_id", creativeService.AdvertiserId)).
		With(zap.String("position", creativeService.Position))
	creative, err := creativeService.Upload()
	if err != nil {
		result := uploadResult{
			Position:   material.Position,
			CreativeId: material.CreativeId,
			ErrCode:    creativeUploadFailed,
			ErrMsg:     err.Error(),
		}
		logger.Info("creative_upload_error", zap.Error(err))
		m.Store(index, result)
		return
	}
	// 添加到送审chan中
	logger.With(zap.Int("id", creative.Id)).With(zap.Int("publisher_id", creative.PublisherId))
	task.CrUploadProducer.Producer(creative, logger)
	result := uploadResult{
		Position:   material.Position,
		CreativeId: material.CreativeId,
		ErrCode:    creativeUploadSuccess,
	}

	logger.Info("creative_upload_success")
	m.Store(index, result)
	return
}

type QueryCreativeJson struct {
	CreativeId []string `json:"creative_id" validate:"required,min=1,dive,min=1"`
}

func (c *Creative) Query(ctx *gin.Context) {

	var (
		queryCreativeJson QueryCreativeJson
		appG              = app.Gin{C: ctx}
		results           []map[string]interface{}
	)
	customer, _ := ctx.MustGet("customer").(*model2.Customer)
	errMsg := app.BindJson(ctx, &queryCreativeJson)
	if len(errMsg) > 0 {
		appG.Response(http.StatusOK, code.INVALID_PARAMS, errMsg)
		return
	}

	creativeService := &service.Creative{
		CreativeIds: queryCreativeJson.CreativeId,
		CustomerId:  customer.Id,
	}
	s := time.Now()
	creatives := creativeService.GetAll()
	costs := time.Now().Sub(s).Milliseconds()

	creativeMaps := make(map[string]model2.Creative)
	for _, creative := range creatives {
		creativeMaps[creative.CreativeId] = *creative
	}
	for _, creativeId := range queryCreativeJson.CreativeId {
		if creative, ok := creativeMaps[creativeId]; ok {
			// 创意审核状态不通过，不需要查询媒体
			if creative.ErrCode == model2.CREATIVE_AUDITING || creative.ErrCode == model2.CREATIVE_QUERY_FAILED {
				logg := log.GetLogger().With(zap.String("request_id", creative.RequestId)).With(zap.Int("creative_id", creative.Id)).With(zap.Int("pub_id", creative.PublisherId)).With(zap.Int64("db costs(ms)", costs))
				task.CrQueryProducer.Producer(creative.PublisherId, creative.Id, logg)
			}

			results = append(results, map[string]interface{}{
				"creative_id":    creativeId,
				"status":         creative.Status,
				"media_cid":      creative.MediaCid,
				"reason":         creative.Reason,
				"extra":          creative.Extra,
				"pub_return_url": creative.PubReturnUrl,
			})
		} else {
			results = append(results, map[string]interface{}{
				"creative_id":    creativeId,
				"status":         -1,
				"media_cid":      "",
				"reason":         "创意不存在",
				"extra":          "",
				"pub_return_url": "",
			})
		}
	}
	appG.Response(http.StatusOK, code.SUCCESS, results)
	return
}
