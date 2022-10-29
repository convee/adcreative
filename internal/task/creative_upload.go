package task

import (
	"context"
	"fmt"
	"github.com/convee/adcreative/internal/enum"
	"github.com/convee/adcreative/internal/media"
	model2 "github.com/convee/adcreative/internal/model"
	"github.com/convee/adcreative/internal/service"
	"runtime/debug"
	"strings"
	"time"

	"github.com/convee/adcreative/configs"
	"github.com/convee/adcreative/internal/pkg/common"

	"github.com/convee/adcreative/internal/pkg/cache"
	"github.com/convee/adcreative/internal/pkg/stats"
	"github.com/convee/adcreative/pkg/ding"
	"github.com/convee/adcreative/pkg/httpclient"
	"github.com/convee/adcreative/pkg/log"
	"github.com/convee/adcreative/pkg/utils"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

var (
	uploadCreativeChan map[int]chan UploadCreativeChan
	CrUploadProducer   *CreativeUploadProducer
	CrUploadConsumer   map[int]*CreativeUploadConsumer
)

type UploadCreativeChan struct {
	IsRsync    int //1已同步，0未同步，为1时走更新接口
	CreativeId int
	MediaCid   string
	Logger     *zap.Logger
	CacheKey   string
}

type CreativeUploadProducer struct {
}
type CreativeUploadConsumer struct {
	PubId       int
	LimitChan   chan struct{}
	RateLimiter *rate.Limiter
	BatchSize   int
}

func InitCreativeUploadTask() {
	CrUploadProducer = &CreativeUploadProducer{}
	CrUploadConsumer = make(map[int]*CreativeUploadConsumer)
	uploadCreativeChan = make(map[int]chan UploadCreativeChan)
	for _, pubId := range enum.PubList {
		CrUploadConsumer[pubId] = &CreativeUploadConsumer{
			PubId:     pubId,
			LimitChan: make(chan struct{}, configs.UploadConcurrenceLimitConf[pubId]),
			BatchSize: configs.UploadBatchSizeConf[pubId],
		}
		uploadCreativeChan[pubId] = make(chan UploadCreativeChan, configs.UploadChanSizeConf[pubId])
	}
}

func (p *CreativeUploadProducer) Producer(creative *model2.Creative, logg *zap.Logger) {
	creativeId := creative.Id
	publisherId := creative.PublisherId
	key := fmt.Sprintf("upload_%d", creativeId)

	if cache.UploadCache.IsExist(key) {
		return
	}

	select {
	case uploadCreativeChan[publisherId] <- UploadCreativeChan{
		CreativeId: creativeId,
		Logger:     logg,
		IsRsync:    creative.IsRsync,
		CacheKey:   key,
	}:
		cache.UploadCache.Set(key)
		logg.With(zap.Int("chan_a", 1)).With(zap.Int("chan_len", len(uploadCreativeChan[publisherId]))).Info("producer_upload_creative")
	case <-time.After(100 * time.Millisecond):
		logg.With(zap.Int("chan_a", 2)).With(zap.Int("chan_len", len(uploadCreativeChan[publisherId]))).Info("producer_upload_creative")
	}

}

func (p *CreativeUploadProducer) Stop() {
	for _, pubId := range enum.PubList {
		uploadCreativeChan[pubId] <- UploadCreativeChan{
			CreativeId: 0,
		}
	}
}

func (p *CreativeUploadConsumer) Work() {
	t_add := time.NewTimer(time.Second)
	if !t_add.Stop() {
		<-t_add.C
	}
	t_update := time.NewTimer(time.Second)
	if !t_update.Stop() {
		<-t_update.C
	}

	var batchAdd []*UploadCreativeChan
	var batchUpdate []*UploadCreativeChan
	for {
		select {
		case <-t_add.C:
			if len(batchAdd) > 0 {
				if p.PubId == enum.PUB_TENCENT && p.RateLimiter != nil {
					ctx := context.Background()
					_ = p.RateLimiter.Wait(ctx)
				}
				p.LimitChan <- struct{}{}
				go p.uploadCreativeBatch(batchAdd, false)
				batchAdd = nil
			}
		case <-t_update.C:
			if len(batchUpdate) > 0 {
				if p.PubId == enum.PUB_TENCENT && p.RateLimiter != nil {
					ctx := context.Background()
					_ = p.RateLimiter.Wait(ctx)
				}
				p.LimitChan <- struct{}{}
				go p.uploadCreativeBatch(batchUpdate, true)
				batchUpdate = nil
			}
		case creativeChan := <-uploadCreativeChan[p.PubId]:
			if creativeChan.CreativeId == 0 {
				return
			}
			creativeChan.Logger.Info("consumer_upload_creative", zap.Int("limit_chan_len", len(p.LimitChan)))
			if p.BatchSize == 0 {
				if p.PubId == enum.PUB_TENCENT && p.RateLimiter != nil {
					ctx := context.Background()
					_ = p.RateLimiter.Wait(ctx)
				}
				p.LimitChan <- struct{}{}
				go p.uploadCreative(creativeChan)
				continue
			}
			// 创意更新
			if creativeChan.IsRsync == 1 {
				batchUpdate = append(batchUpdate, &creativeChan)
				if len(batchUpdate) == 1 {
					if !t_update.Stop() {
						select {
						case <-t_update.C:
						default:
						}
					}
					t_update.Reset(time.Second)
				}
				if len(batchUpdate) == p.BatchSize {
					if !t_update.Stop() {
						select {
						case <-t_update.C:
						default:
						}
					}
					if p.PubId == enum.PUB_TENCENT && p.RateLimiter != nil {
						ctx := context.Background()
						_ = p.RateLimiter.Wait(ctx)
					}
					p.LimitChan <- struct{}{}
					go p.uploadCreativeBatch(batchUpdate, true)
					batchUpdate = nil
				}
			} else {
				//创意新增
				batchAdd = append(batchAdd, &creativeChan)
				if len(batchAdd) == 1 {
					if !t_add.Stop() {
						select {
						case <-t_add.C:
						default:
						}
					}
					t_add.Reset(time.Second)
				}
				if len(batchAdd) == p.BatchSize {
					if !t_add.Stop() {
						select {
						case <-t_add.C:
						default:
						}
					}
					if p.PubId == enum.PUB_TENCENT && p.RateLimiter != nil {
						ctx := context.Background()
						_ = p.RateLimiter.Wait(ctx)
					}
					p.LimitChan <- struct{}{}
					go p.uploadCreativeBatch(batchAdd, false)
					batchAdd = nil
				}

			}

		}
	}
}

func (p *CreativeUploadConsumer) uploadCreativeBatch(batch []*UploadCreativeChan, isUpdate bool) {
	var batchCreative []media.BatchCreativeUpload

	defer func() {
		for _, uploadCreativeChan := range batch {
			// 清理缓存
			cache.UploadCache.Del(uploadCreativeChan.CacheKey)
		}
		<-p.LimitChan // 释放管道
		if e := recover(); e != nil {
			dingMsg := map[string]interface{}{
				"publisher_id": p.PubId,
				"err":          e,
			}
			log.Error("upload_creative_batch_panic", zap.Any("msg", dingMsg))
			ding.SendAlert("创意批量送审 panic...", dingMsg, false)
			log.Error("创意批量送审 panic...", zap.String("stack", string(debug.Stack())))
		}
	}()

	var (
		errCreatives []int
	)

	for _, uploadCreativeChan := range batch {
		logger := uploadCreativeChan.Logger
		// 获取创意信息
		creative, err := cache.GetCreativeCacheById(uploadCreativeChan.CreativeId)
		if err != nil {
			dingMsg := map[string]interface{}{
				"publisher": p.PubId,
				"method":    "uploadCreativeBatch",
				"id":        uploadCreativeChan.CreativeId,
				"source":    utils.GetHostname(),
				"err":       err.Error(),
			}
			logger.Error("get_creative_error", zap.Any("msg", dingMsg))
			ding.SendAlert("创意查询-获取创意失败预警", dingMsg, false)
			errCreatives = append(errCreatives, uploadCreativeChan.CreativeId)
			continue
		}
		uploadCreativeChan.MediaCid = creative.MediaCid
		// 获取客户信息
		customer, err := cache.GetCustomerCacheById(creative.CustomerId)
		if err != nil {

			dingMsg := map[string]interface{}{
				"publisher":   creative.PublisherId,
				"method":      "uploadCreative",
				"id":          creative.Id,
				"customer_id": creative.CustomerId,
				"source":      utils.GetHostname(),
				"err":         err.Error(),
			}
			logger.Error("get_customer_error", zap.Any("msg", dingMsg))
			ding.SendAlert("创意查询-获取客户信息失败预警", dingMsg, false)
			errCreatives = append(errCreatives, uploadCreativeChan.CreativeId)
			continue
		}
		batchCreative = append(batchCreative, media.BatchCreativeUpload{
			MediaCid: creative.MediaCid,
			Creative: creative,
			Customer: customer,
		})

	}

	// customerId = 0 目前媒体账号客户ID都是用0
	logg := log.GetLogger().With(zap.String("method", "uploadCreativeBatch"))
	handler, errs := media.GetBatchHandler(batchCreative, p.PubId, 0, logg, isUpdate)
	if errs != nil {
		for i, uploadCreativeChan := range batch {
			if common.IntContain(uploadCreativeChan.CreativeId, errCreatives...) {
				continue
			}
			if err, ok := errs[uploadCreativeChan.CreativeId]; ok {
				errCreatives = append(errCreatives, uploadCreativeChan.CreativeId)
				ret := media.Ret{ErrCode: model2.CREATIVE_UPLOAD_UNPASSED, ErrMsg: err.Error()}
				p.UpdateCreativeStatus(batchCreative[i].Creative, batchCreative[i].Customer, ret, uploadCreativeChan.Logger)
			}
		}
	}

	// 全部错误，立即返回
	if len(errCreatives) == len(batch) {
		return
	}
	// 执行实际的查询
	start := time.Now()
	ret := handler.BatchUploadCreative()
	ret.MediaCosts = time.Now().Sub(start).Milliseconds()
	stats.PublisherApiObserve(p.PubId, model2.HANDLER_METHOD_UPLOAD, ret.MediaCosts)

	// 统计媒体处理结果
	rets := ret.BatchUploadRetMap

	// 更新数据库，并回调客户端，同步状态
	for i, uploadCreativeChan := range batch {
		if common.IntContain(i, errCreatives...) {
			continue
		}
		var singleRet = ret
		if r, ok := rets[uploadCreativeChan.MediaCid]; ok {
			singleRet.ErrCode = r.ErrCode
			singleRet.ErrMsg = r.ErrMsg
			singleRet.IsRsync = r.IsRsync
		}
		p.UpdateCreativeStatus(batchCreative[i].Creative, batchCreative[i].Customer, singleRet, uploadCreativeChan.Logger)

	}
	return

}

// err_code:1上传失败、2待送审、3送审失败、4审核中、5审核通过、6审核不通过
// status:0待审核，1审核通过，2审核不通过
// 2,4=0 5=1 1,3,6=2
func (p *CreativeUploadConsumer) uploadCreative(creativeChan UploadCreativeChan) {
	defer func() {
		<-p.LimitChan // 释放管道
		// 清理缓存
		cache.UploadCache.Del(creativeChan.CacheKey)
		if e := recover(); e != nil {
			dingMsg := map[string]interface{}{
				"creative_id":  creativeChan.CreativeId,
				"publisher_id": p.PubId,
				"err":          e,
			}
			log.Error("upload_creative_panic", zap.Any("msg", dingMsg))
			ding.SendAlert("创意送审 panic...", dingMsg, false)
			log.Error("创意送审 panic...", zap.String("stack", string(debug.Stack())))
		}
	}()
	logger := creativeChan.Logger
	logger.With(zap.String("method", "uploadCreative"))
	creative, err := cache.GetCreativeCacheById(creativeChan.CreativeId)

	if err != nil {
		dingMsg := map[string]interface{}{
			"publisher":   creative.PublisherId,
			"method":      "uploadCreative",
			"creative_id": creative.CreativeId,
			"source":      utils.GetHostname(),
			"err":         err.Error(),
		}
		logger.Error("get_creative_error", zap.Any("msg", dingMsg))
		ding.SendAlert("创意上传-获取创意失败预警", dingMsg, false)
		return
	}

	customer, err := cache.GetCustomerCacheById(creative.CustomerId)
	if err != nil {
		dingMsg := map[string]interface{}{
			"publisher":   creative.PublisherId,
			"method":      "uploadCreative",
			"creative_id": creative.CreativeId,
			"customer_id": creative.CustomerId,
			"source":      utils.GetHostname(),
			"err":         err.Error(),
		}
		logger.Error("get_customer_error", zap.Any("msg", dingMsg))
		ding.SendAlert("创意上传-获取客户信息失败预警", dingMsg, false)
		return
	}

	var infos []model2.TemplateInfo
	_ = jsoniter.Unmarshal([]byte(creative.Info), &infos)
	creativeService := &service.Creative{
		TemplateId: creative.TemplateId,
		Info:       infos,
	}
	positionInfo, err := cache.GetPositionCacheById(creative.PositionId)
	if err != nil {
		logger.Error("get_position_error", zap.Error(err))
		dingMsg := map[string]interface{}{
			"publisher":   creative.PublisherId,
			"method":      "uploadCreative",
			"creative_id": creative.CreativeId,
			"source":      utils.GetHostname(),
			"err":         err.Error(),
		}
		ding.SendAlert("创意上传-获取广告位失败预警", dingMsg, false)
		return
	}
	info, errs := creativeService.Check(positionInfo, 0)
	if len(errs) > 0 {
		ret := media.Ret{ErrCode: model2.CREATIVE_UPLOAD_UNPASSED, ErrMsg: strings.Join(errs, ",")}
		p.UpdateCreativeStatus(creative, customer, ret, logger)
		return
	}
	infoJson, _ := jsoniter.Marshal(info)
	creative.Info = string(infoJson)

	mediaHandler, err := media.GetCreativeHandler(creative, customer, logger)
	if err != nil {
		logger.Error("get_creative_handler_error", zap.Error(err))
		ret := media.Ret{ErrCode: model2.CREATIVE_UPLOAD_FAILED, ErrMsg: err.Error()}
		p.UpdateCreativeStatus(creative, customer, ret, logger)
		return
	}
	start := time.Now()
	ret := mediaHandler.UploadCreative()
	ret.MediaCosts = time.Now().Sub(start).Milliseconds()
	stats.PublisherApiObserve(p.PubId, model2.HANDLER_METHOD_UPLOAD, ret.MediaCosts)

	p.UpdateCreativeStatus(creative, customer, ret, logger)
}

func (p *CreativeUploadConsumer) UpdateCreativeStatus(creative *model2.Creative, customer *model2.Customer, ret media.Ret, logger *zap.Logger) {
	msg := map[string]interface{}{
		"id":            creative.Id,
		"publisher":     creative.PublisherId,
		"creative_id":   creative.CreativeId,
		"media_cid":     creative.MediaCid,
		"customer_id":   creative.CustomerId,
		"advertiser_id": creative.AdvertiserId,
		"source":        utils.GetHostname(),
		"err_code":      ret.ErrCode,
		"err_msg":       model2.StatusMap[ret.ErrCode] + ":" + ret.ErrMsg,
		"url":           ret.Url,
		"header":        ret.Header,
		"req":           ret.Req,
		"resp":          ret.Resp,
		"media_costs":   ret.MediaCosts,
	}
	if ret.ErrCode == model2.CREATIVE_UPLOAD_FAILED {
		ding.SendAlert("创意上传失败", msg, false)
	} else if ret.ErrCode == model2.CREATIVE_UPLOAD_UNPASSED {
		ding.SendAlert("创意上传不通过预警", msg, false)
	} else if ret.ErrCode == model2.CREATIVE_UPDATE_EXCEPTION { //创意更新异常，没有对应到错误码，不更新
		ret.ErrCode = creative.ErrCode
		ding.SendAlert("创意上传更新异常预警", msg, false)
	}
	logger.Info("upload_creative_info", zap.Any("msg", msg))
	status := model2.GetCreativeStatusByErrCode(ret.ErrCode)
	isRsync := 0
	if ret.IsRsync == 1 {
		isRsync = 1
	} else if ret.ErrCode == model2.CREATIVE_AUDITING {
		isRsync = 1
	}
	if ret.MediaCid == "" {
		ret.MediaCid = creative.MediaCid
	}
	creativeService := &service.Creative{
		Id:           creative.Id,
		Status:       status,
		ErrCode:      ret.ErrCode,
		Reason:       ret.ErrMsg,
		ErrMsg:       ret.ErrMsg,
		MediaCid:     ret.MediaCid,
		Extra:        ret.Extra,
		IsRsync:      isRsync,
		RequestId:    creative.RequestId,
		VideoCdnUrl:  ret.VideoCdnUrl,
		PicCdnUrl:    ret.PicCdnUrl,
		PubReturnUrl: ret.PubReturnUrl,
	}

	if ret.ErrCode == model2.CREATIVE_UPLOAD_FAILED {
		creativeService.Reason = ""
	}
	err := creativeService.UpdateCreativeByMaps()
	if err != nil {
		logger.Error("creative_update_error", zap.Error(err))
		dingMsg := map[string]interface{}{
			"publisher":     creative.PublisherId,
			"method":        "UpdateCreativeStatus",
			"creative_id":   creative.CreativeId,
			"customer_id":   creative.CustomerId,
			"advertiser_id": creative.AdvertiserId,
			"source":        utils.GetHostname(),
			"err":           err.Error(),
		}
		ding.SendAlert("创意上传-更新创意失败预警", dingMsg, false)
	}
	if len(customer.CreativeCallbackUrl) > 0 {
		p.CreativeUploadCallBack(customer.CreativeCallbackUrl, creative, ret.MediaCid, status, creativeService.Reason, logger, ret.PubReturnUrl)
	}
}

// CreativeUploadCallBack 将物料审核的结果回调给送审方，只调一次，不论成功失败
// status:0待审核，1审核通过，2审核不通过
func (p *CreativeUploadConsumer) CreativeUploadCallBack(uri string, creative *model2.Creative, mediaInfo string, status int, reason string, logger *zap.Logger, pubReturnUrl string) {
	material := map[string]interface{}{
		"creative_id":    creative.CreativeId,
		"media_info":     mediaInfo,
		"status":         status,
		"reason":         reason,
		"pub_return_url": pubReturnUrl,
	}
	var data []map[string]interface{}
	data = append(data, material)
	request := map[string]interface{}{
		"material": data,
	}
	bodyJson, _ := jsoniter.Marshal(request)
	response, err := httpclient.PostJSON(uri, bodyJson, httpclient.WithTTL(time.Second*100))
	if err != nil {
		logger.Info("creative_upload_callback_error", zap.Error(err))
		dingMsg := map[string]interface{}{
			"method":        "CreativeUploadCallBack",
			"creative_id":   creative.CreativeId,
			"customer_id":   creative.CustomerId,
			"advertiser_id": creative.AdvertiserId,
			"publisher":     creative.PublisherId,
			"source":        utils.GetHostname(),
			"err":           err.Error(),
		}
		ding.SendAlert("创意上传-回调失败预警", dingMsg, false)
	} else {
		logger.Info("creative_upload_callback_info", zap.Any("url", uri), zap.Any("req", string(bodyJson)), zap.Any("resp", string(response)), zap.Error(err))
	}
}
