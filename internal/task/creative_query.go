package task

import (
	"context"
	"fmt"
	"github.com/convee/adcreative/configs"
	"github.com/convee/adcreative/internal/enum"
	"github.com/convee/adcreative/internal/media"
	model2 "github.com/convee/adcreative/internal/model"
	"github.com/convee/adcreative/internal/service"
	"github.com/convee/adcreative/pkg/log"
	"runtime/debug"
	"time"

	"github.com/convee/adcreative/internal/pkg/stats"
	"golang.org/x/time/rate"

	"github.com/convee/adcreative/internal/pkg/cache"
	"github.com/convee/adcreative/pkg/ding"
	"github.com/convee/adcreative/pkg/httpclient"
	"github.com/convee/adcreative/pkg/utils"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

var (
	queryCreativeChan map[int]chan QueryCreativeChan
	CrQueryProducer   *CreativeQueryProducer
	CrQueryConsumer   map[int]*CreativeQueryConsumer
)

type QueryCreativeChan struct {
	CreativeId int
	MediaCid   string
	Logger     *zap.Logger
	CacheKey   string
}

type CreativeQueryProducer struct {
}
type CreativeQueryConsumer struct {
	PubId       int
	LimitChan   chan struct{}
	RateLimiter *rate.Limiter
	BatchSize   int
}

func InitCreativeQueryTask() {
	CrQueryProducer = &CreativeQueryProducer{}
	CrQueryConsumer = make(map[int]*CreativeQueryConsumer)
	queryCreativeChan = make(map[int]chan QueryCreativeChan)
	for _, pubId := range enum.PubList {
		CrQueryConsumer[pubId] = &CreativeQueryConsumer{
			PubId:     pubId,
			LimitChan: make(chan struct{}, configs.QueryConcurrenceLimitConf[pubId]),
			BatchSize: configs.QueryBatchSizeConf[pubId],
		}

		queryCreativeChan[pubId] = make(chan QueryCreativeChan, configs.QueryChanSizeConf[pubId])
	}
}

/*
	creativeId ： 素材服务里面的自增 ID
*/
func (p *CreativeQueryProducer) Producer(publisherId, creativeId int, logger *zap.Logger) {
	key := fmt.Sprintf("query_%d", creativeId)
	if cache.QueryCache.IsExist(key) {
		return
	}
	select {
	case queryCreativeChan[publisherId] <- QueryCreativeChan{
		CreativeId: creativeId,
		Logger:     logger,
		CacheKey:   key,
	}:
		cache.QueryCache.Set(key)
		logger.With(zap.Int("chan_query", 1)).With(zap.Int("chan_len", len(queryCreativeChan[publisherId]))).Info("producer_query_creative")

	case <-time.After(100 * time.Millisecond):
		logger.With(zap.Int("chan_query", 2)).With(zap.Int("chan_len", len(queryCreativeChan[publisherId]))).Info("producer_query_creative")
	}
}

func (p *CreativeQueryProducer) Stop() {
	for _, pubId := range enum.PubList {
		queryCreativeChan[pubId] <- QueryCreativeChan{
			CreativeId: 0,
		}
	}
}

func (p *CreativeQueryConsumer) Work() {
	t := time.NewTimer(time.Second)
	if !t.Stop() {
		<-t.C
	}
	var batch []*QueryCreativeChan
	for {
		select {
		case <-t.C:
			// 定时器到时间了，执行 batch 里面的内容
			if len(batch) > 0 {
				if p.PubId == enum.PUB_TENCENT && p.RateLimiter != nil {
					ctx := context.Background()
					_ = p.RateLimiter.Wait(ctx)
				}
				p.LimitChan <- struct{}{}
				go p.queryCreativeBatch(batch)
				batch = nil
			}

		case creativeChan := <-queryCreativeChan[p.PubId]:
			if creativeChan.CreativeId == 0 {
				return
			}
			creativeChan.Logger.Info("consumer_query_creative", zap.Int("limit_chan_len", len(p.LimitChan)))
			if p.BatchSize == 0 {
				if p.PubId == enum.PUB_TENCENT && p.RateLimiter != nil {
					ctx := context.Background()
					_ = p.RateLimiter.Wait(ctx)
				}
				p.LimitChan <- struct{}{}
				go p.queryCreative(creativeChan)
				continue
			}
			batch = append(batch, &creativeChan)
			if len(batch) == 1 {
				if !t.Stop() {
					select {
					case <-t.C:
					default:
					}
				}
				t.Reset(time.Second)
			}
			if len(batch) == p.BatchSize {
				if !t.Stop() {
					select {
					case <-t.C:
					default:
					}
				}
				if p.PubId == enum.PUB_TENCENT && p.RateLimiter != nil {
					ctx := context.Background()
					_ = p.RateLimiter.Wait(ctx)
				}
				p.LimitChan <- struct{}{}
				go p.queryCreativeBatch(batch)
				batch = nil
			}

		}
	}
}
func (p *CreativeQueryConsumer) queryCreativeBatch(batch []*QueryCreativeChan) {
	defer func() {
		for _, queryCreativeChan := range batch {
			// 清理缓存
			cache.QueryCache.Del(queryCreativeChan.CacheKey)
		}

		<-p.LimitChan // 释放管道

		if e := recover(); e != nil {
			dingMsg := map[string]interface{}{
				"publisher_id": p.PubId,
				"err":          e,
			}
			log.Error("query_creative_batch_panic", zap.Any("msg", dingMsg))
			ding.SendAlert("批量查询创意 panic...", dingMsg, false)
			log.Error("批量查询创意 panic...", zap.String("stack", string(debug.Stack())))
		}
	}()

	var batchCreative []media.BatchCreativeQuery
	var batchCreativeMap = make(map[string]media.BatchCreativeQuery)

	for _, queryCreativeChan := range batch {
		logger := queryCreativeChan.Logger
		// 获取创意信息
		creative, err := cache.GetCreativeCacheById(queryCreativeChan.CreativeId)
		if err != nil {
			dingMsg := map[string]interface{}{
				"publisher": p.PubId,
				"method":    "queryCreativeBatch",
				"id":        queryCreativeChan.CreativeId,
				"source":    utils.GetHostname(),
				"err":       err.Error(),
			}
			logger.Error("get_creative_error", zap.Any("msg", dingMsg))
			ding.SendAlert("创意查询-获取创意失败预警", dingMsg, false)
			continue
		}
		queryCreativeChan.MediaCid = creative.MediaCid
		// 获取客户信息
		customer, err := cache.GetCustomerCacheById(creative.CustomerId)
		if err != nil {
			dingMsg := map[string]interface{}{
				"publisher":   creative.PublisherId,
				"method":      "queryCreative",
				"creative_id": creative.CreativeId,
				"customer_id": creative.CustomerId,
				"source":      utils.GetHostname(),
				"err":         err.Error(),
			}
			logger.Error("get_customer_error", zap.Any("msg", dingMsg))
			ding.SendAlert("创意查询-获取客户信息失败预警", dingMsg, false)
			continue
		}
		batchCreative = append(batchCreative, media.BatchCreativeQuery{
			MediaCid: creative.MediaCid,
			Creative: creative,
			Customer: customer,
		})
		batchCreativeMap[creative.MediaCid] = media.BatchCreativeQuery{
			MediaCid: creative.MediaCid,
			Creative: creative,
			Customer: customer,
		}
	}

	l := log.GetLogger()
	l.With(zap.String("method", "batchQueryCreativeBatch"))

	// customerId = 0 目前媒体账号客户ID都是用0
	handler, err := media.GetBatchCreativeHandler(batchCreative, p.PubId, 0, l)
	if err != nil {
		ret := media.Ret{ErrCode: model2.CREATIVE_QUERY_FAILED, ErrMsg: err.Error()}
		for _, queryCreativeChan := range batch {
			if creativeInfo, ok := batchCreativeMap[queryCreativeChan.MediaCid]; ok {
				p.UpdateCreativeStatus(creativeInfo.Creative, creativeInfo.Customer, ret, queryCreativeChan.Logger)
			}
		}
		return
	}

	// 执行实际的查询
	start := time.Now()
	ret := handler.BatchQueryCreative()
	ret.MediaCosts = time.Now().Sub(start).Milliseconds()
	stats.PublisherApiObserve(p.PubId, model2.HANDLER_METHOD_QUERY, ret.MediaCosts)
	rets := map[string]media.Ret{}
	for _, r := range ret.BatchQueryRet {
		//ret.ErrCode = r.ErrCode
		//ret.ErrMsg = r.ErrMsg
		//ret.MediaCid = r.MediaCid
		var singleRet = ret
		singleRet.ErrCode = r.ErrCode
		singleRet.ErrMsg = r.ErrMsg
		singleRet.MediaCid = r.MediaCid
		rets[r.MediaCid] = singleRet
	}

	for _, queryCreativeChan := range batch {
		if batchInfo, ok := rets[queryCreativeChan.MediaCid]; ok {
			creativeInfo := batchCreativeMap[queryCreativeChan.MediaCid]
			p.UpdateCreativeStatus(creativeInfo.Creative, creativeInfo.Customer, batchInfo, queryCreativeChan.Logger)
		} else {
			creativeInfo := batchCreativeMap[queryCreativeChan.MediaCid]
			p.UpdateCreativeStatus(creativeInfo.Creative, creativeInfo.Customer, ret, queryCreativeChan.Logger)
		}
	}

	return
}
func (p *CreativeQueryConsumer) queryCreative(creativeChan QueryCreativeChan) {
	logger := creativeChan.Logger
	logger.With(zap.String("method", "queryCreative"))

	defer func() {
		<-p.LimitChan // 释放管道
		// 清理缓存
		cache.QueryCache.Del(creativeChan.CacheKey)
		if e := recover(); e != nil {
			dingMsg := map[string]interface{}{
				"creative_id":  creativeChan.CreativeId,
				"publisher_id": p.PubId,
				"err":          e,
			}
			logger.Error("query_creative_panic", zap.Any("msg", dingMsg))
			ding.SendAlert("创意查询 panic...", dingMsg, false)
			log.Error("创意查询 panic...", zap.String("stack", string(debug.Stack())))
		}
	}()

	creativeId := creativeChan.CreativeId
	// 获取创意信息
	creative, err := cache.GetCreativeCacheById(creativeId)
	if err != nil {

		dingMsg := map[string]interface{}{
			"publisher": p.PubId,
			"method":    "queryCreative",
			"id":        creativeChan.CreativeId,
			"source":    utils.GetHostname(),
			"err":       err.Error(),
		}
		logger.Error("get_creative_error", zap.Any("msg", dingMsg))
		ding.SendAlert("创意查询-获取创意失败预警", dingMsg, false)
		return
	}

	// 获取客户信息
	customer, err := cache.GetCustomerCacheById(creative.CustomerId)
	if err != nil {
		dingMsg := map[string]interface{}{
			"publisher":   creative.PublisherId,
			"method":      "queryCreative",
			"creative_id": creative.CreativeId,
			"customer_id": creative.CustomerId,
			"source":      utils.GetHostname(),
			"err":         err.Error(),
		}
		logger.Error("get_customer_error", zap.Any("msg", dingMsg))
		ding.SendAlert("创意查询-获取客户信息失败预警", dingMsg, false)
		return
	}
	mediaHandler, err := media.GetCreativeHandler(creative, customer, logger)
	if err != nil {
		logger.Error("get_creative_handler_error", zap.Error(err))
		ret := media.Ret{ErrCode: model2.CREATIVE_QUERY_FAILED, ErrMsg: err.Error()}
		dingMsg := map[string]interface{}{
			"publisher":     creative.PublisherId,
			"method":        "queryCreative",
			"customer_id":   creative.CustomerId,
			"advertiser_id": creative.AdvertiserId,
			"creative_id":   creative.CreativeId,
			"source":        utils.GetHostname(),
			"err_code":      ret.ErrCode,
			"err_msg":       ret.ErrMsg,
		}
		ding.SendAlert("创意查询-获取Handler失败预警", dingMsg, false)
		p.UpdateCreativeStatus(creative, customer, ret, logger)
		return
	}

	// 执行实际的查询
	start := time.Now()
	ret := mediaHandler.QueryCreative()
	ret.MediaCosts = time.Now().Sub(start).Milliseconds()
	stats.PublisherApiObserve(creative.PublisherId, model2.HANDLER_METHOD_QUERY, ret.MediaCosts)
	p.UpdateCreativeStatus(creative, customer, ret, logger)

	return

}

func (p *CreativeQueryConsumer) UpdateCreativeStatus(creative *model2.Creative, customer *model2.Customer, ret media.Ret, logger *zap.Logger) {
	dingMsg := map[string]interface{}{
		"id":             creative.Id,
		"publisher":      creative.PublisherId,
		"creative_id":    creative.CreativeId,
		"media_cid":      creative.MediaCid,
		"customer_id":    creative.CustomerId,
		"advertiser_id":  creative.AdvertiserId,
		"source":         utils.GetHostname(),
		"err_code":       ret.ErrCode,
		"err_msg":        model2.StatusMap[ret.ErrCode] + ":" + ret.ErrMsg,
		"url":            ret.Url,
		"req":            ret.Req,
		"header":         ret.Header,
		"resp":           ret.Resp,
		"media_costs":    ret.MediaCosts,
		"batch_creative": ret.BatchQueryRet,
	}
	if ret.ErrCode == model2.CREATIVE_QUERY_FAILED {
		ding.SendAlert("创意查询失败预警", dingMsg, false)
	} else if ret.ErrCode == model2.CREATIVE_AUDIT_UNPASSWD {
		ding.SendAlert("创意查询审核不通过预警", dingMsg, false)
	} else if ret.ErrCode == model2.CREATIVE_UPDATE_EXCEPTION {
		ret.ErrCode = creative.ErrCode
		ding.SendAlert("创意查询更新异常预警", dingMsg, false)
	}
	logger.Info("query_creative_info", zap.Any("msg", dingMsg))
	status := model2.GetCreativeStatusByErrCode(ret.ErrCode)
	if ret.MediaCid == "" {
		ret.MediaCid = creative.MediaCid
	}
	creativeService := &service.Creative{
		Id:       creative.Id,
		Status:   status,
		MediaCid: ret.MediaCid,
		Extra:    ret.Extra,
		ErrCode:  ret.ErrCode,
		Reason:   ret.ErrMsg,
		ErrMsg:   ret.ErrMsg,
	}
	if ret.ErrCode == model2.CREATIVE_QUERY_FAILED {
		creativeService.Reason = ""
	}
	err := creativeService.UpdateCreativeByMaps()
	if err != nil {
		logger.Error("creative_update_error", zap.Error(err))
		dingMsg := map[string]interface{}{
			"publisher":   creative.PublisherId,
			"method":      "UpdateCreativeStatus",
			"creative_id": creative.CreativeId,
			"source":      utils.GetHostname(),
			"err":         err.Error(),
		}
		ding.SendAlert("创意查询-更新创意失败预警", dingMsg, false)
	}
	if len(customer.CreativeCallbackUrl) > 0 {
		p.CreativeQueryCallBack(customer.CreativeCallbackUrl, creative.CreativeId, ret.MediaCid, status, creativeService.Reason, ret.Extra, logger, creative.PubReturnUrl)
	}
}

// CreativeQueryCallBack 将物料审核的结果回调给送审方，只调一次，不论成功失败
// status:0待审核，1审核通过，2审核不通过
func (p *CreativeQueryConsumer) CreativeQueryCallBack(uri string, creativeId string, mediaInfo string, status int, reason string, extra string, logger *zap.Logger, pubReturnUrl string) {
	material := map[string]interface{}{
		"creative_id":    creativeId,
		"media_info":     mediaInfo,
		"status":         status,
		"reason":         reason,
		"extra":          extra,
		"pub_return_url": pubReturnUrl,
	}
	var data []map[string]interface{}
	data = append(data, material)
	request := map[string]interface{}{
		"material": data,
	}
	//uri := configs.Conf.App.CallbackUrl
	bodyJson, _ := jsoniter.Marshal(request)
	response, err := httpclient.PostJSON(uri, bodyJson, httpclient.WithTTL(time.Second*100))
	logger.Info("creative_query_callback_info", zap.Any("url", uri), zap.Any("req", string(bodyJson)), zap.Any("resp", string(response)), zap.Error(err))
	if err != nil {
		logger.Error("creative_query_callback_info_error", zap.Error(err))
	}
}
