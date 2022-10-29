package task

import (
	"github.com/convee/adcreative/internal/media"
	model2 "github.com/convee/adcreative/internal/model"
	"github.com/convee/adcreative/internal/service"
	"github.com/convee/adcreative/pkg/httpclient"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"time"
)

var (
	queryAdvertiserChan chan QueryAdvertiseChan
	AdvQueryProducer    *AdvertiserQueryProducer
	AdvQueryConsumer    *AdvertiserQueryConsumer
	Logg                *zap.Logger
)

type QueryAdvertiseChan struct {
	AdvertiserAuditId int
	Logger            *zap.Logger
}

type AdvertiserQueryProducer struct {
}
type AdvertiserQueryConsumer struct {
	PubId int
}

func InitAdvertiserQueryTask() {
	AdvQueryProducer = &AdvertiserQueryProducer{}
	AdvQueryConsumer = &AdvertiserQueryConsumer{}
	queryAdvertiserChan = make(chan QueryAdvertiseChan, 1024)
}

func (p *AdvertiserQueryProducer) Producer(advertiserAuditId int, logger *zap.Logger) {
	queryAdvertiserChan <- QueryAdvertiseChan{
		AdvertiserAuditId: advertiserAuditId,
		Logger:            logger,
	}
	logger.Info("producer_query_advertiser")
}

func (p *AdvertiserQueryProducer) Stop() {
	queryAdvertiserChan <- QueryAdvertiseChan{
		AdvertiserAuditId: 0,
	}
}

func (p *AdvertiserQueryConsumer) Work() {
	for {
		select {
		case advertiserAudit := <-queryAdvertiserChan:
			if advertiserAudit.AdvertiserAuditId == 0 {
				return
			}
			advertiserAudit.Logger.Info("consumer_query_advertiser")
			go p.queryAdvertiser(advertiserAudit)
		}
	}
}

func (p *AdvertiserQueryConsumer) queryAdvertiser(advertiserAuditChan QueryAdvertiseChan) {
	advAuditService := service.AdvertiserAudit{
		Id: advertiserAuditChan.AdvertiserAuditId,
	}
	logger := advertiserAuditChan.Logger
	advAudit, err := advAuditService.GetAdvertiserAuditInfo()
	if err != nil {
		logger.Error("get_advertiser_audit_err", zap.Error(err))
		return
	}
	var (
		customerModel = model2.CustomerModel{}
	)
	customer, err := customerModel.GetCustomerById(advAudit.CustomerId)
	if err != nil {
		logger.Error("get_advertiser_error")
		return
	}
	mediaHandler, err := media.GetAdvertiserHandler(advAudit, logger)
	if err != nil {
		logger.Error("get_advertiser_handler_error")
		return
	}
	ret := mediaHandler.QueryAdvertiser()
	status := model2.GetAdvertiserStatusByErrCode(ret.ErrCode)
	advertiserAuditService := &service.AdvertiserAudit{
		Id:       advertiserAuditChan.AdvertiserAuditId,
		Status:   status,
		ErrCode:  ret.ErrCode,
		ErrMsg:   ret.ErrMsg,
		MediaCid: ret.MediaCid,
	}
	advertiserAuditService.UpdateAdvertiserAuditByMaps()
	if len(customer.AdvertiserCallbackUrl) > 0 {
		go p.AdvertiserQueryCallBack(customer.AdvertiserCallbackUrl, advertiserAuditChan.AdvertiserAuditId, status, ret.ErrMsg, logger)
	}
}

// AdvertiserQueryCallBack 将广告主审核的结果回调给送审方，只调一次，不论成功失败
// status:0待审核，1审核通过，2审核不通过
func (p *AdvertiserQueryConsumer) AdvertiserQueryCallBack(uri string, advertiserId int, status int, reason string, logger *zap.Logger) {
	advertiser := map[string]interface{}{
		"advertiser_id": advertiserId,
		"status":        status,
		"reason":        reason,
	}
	var data []map[string]interface{}
	data = append(data, advertiser)
	request := map[string]interface{}{
		"advertiser": data,
	}
	bodyJson, _ := jsoniter.Marshal(request)
	response, err := httpclient.PostJSON(uri, bodyJson, httpclient.WithTTL(time.Second*100))
	logger.Info("advertiser_query_callback_info", zap.Any("url", uri), zap.Any("req", string(bodyJson)), zap.Any("resp", string(response)), zap.Error(err))
	if err != nil {
		logger.Error("advertiser_query_callback_info error", zap.Error(err))
	}
}
