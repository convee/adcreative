package task

import (
	"github.com/convee/adcreative/internal/media"
	model2 "github.com/convee/adcreative/internal/model"
	"github.com/convee/adcreative/internal/service"
	"github.com/convee/adcreative/pkg/httpclient"
	logger "github.com/convee/adcreative/pkg/log"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"time"
)

var (
	uploadAdvertiserChan chan UploadAdvertiseChan
	AdvUploadProducer    *AdvertiserUploadProducer
	AdvUploadConsumer    *AdvertiserUploadConsumer
)

type UploadAdvertiseChan struct {
	AdvertiserAuditId int
	Logger            *zap.Logger
}

type AdvertiserUploadProducer struct {
}
type AdvertiserUploadConsumer struct {
	PubId int
}

func InitAdvertiserUploadTask() {
	AdvUploadProducer = &AdvertiserUploadProducer{}
	AdvUploadConsumer = &AdvertiserUploadConsumer{}
	uploadAdvertiserChan = make(chan UploadAdvertiseChan, 1024)
}

func (p *AdvertiserUploadProducer) Producer(advertiserAuditId int, requestId string) {
	logg := logger.GetLogger().With(zap.String("request_id", requestId)).With(zap.Int("advertiser_id", advertiserAuditId))
	uploadAdvertiserChan <- UploadAdvertiseChan{
		AdvertiserAuditId: advertiserAuditId,
		Logger:            logg,
	}
	logg.Info("producer_advertiser_upload")
}

func (p *AdvertiserUploadProducer) Stop() {
	uploadAdvertiserChan <- UploadAdvertiseChan{
		AdvertiserAuditId: 0,
	}
}

func (p *AdvertiserUploadConsumer) Work() {
	for {
		select {
		case advertiserAudit := <-uploadAdvertiserChan:
			if advertiserAudit.AdvertiserAuditId == 0 {
				return
			}
			advertiserAudit.Logger.Info("consumer_upload_advertiser")
			go p.uploadAdvertiser(advertiserAudit)
		}
	}
}

func (p *AdvertiserUploadConsumer) uploadAdvertiser(advertiserAuditChan UploadAdvertiseChan) {

	advAuditService := service.AdvertiserAudit{
		Id: advertiserAuditChan.AdvertiserAuditId,
	}
	logger := advertiserAuditChan.Logger
	advAudit, err := advAuditService.GetAdvertiserAuditInfo()
	if err != nil {
		logger.Error("get_advertiser_audit_err", zap.Error(err))
		return
	}
	customer, err := new(model2.CustomerModel).GetCustomerById(advAudit.CustomerId)
	if err != nil {
		logger.Error("get_advertiser_error")
		return
	}
	mediaHandler, err := media.GetAdvertiserHandler(advAudit, logger)
	if err != nil {
		logger.Error("get_advertiser_handler_error")
		return
	}
	ret := mediaHandler.UploadAdvertiser()
	status := model2.GetAdvertiserStatusByErrCode(ret.ErrCode)
	advertiserAuditService := &service.AdvertiserAudit{
		Id:       advertiserAuditChan.AdvertiserAuditId,
		Status:   status,
		IsRsync:  1,
		ErrCode:  ret.ErrCode,
		ErrMsg:   ret.ErrMsg,
		MediaCid: ret.MediaCid,
	}
	advertiserAuditService.UpdateAdvertiserAuditByMaps()
	if len(customer.AdvertiserCallbackUrl) > 0 {
		go p.AdvertiserUploadCallBack(customer.AdvertiserCallbackUrl, advertiserAuditChan.AdvertiserAuditId, status, ret.ErrMsg, logger)
	}
	// 创意状态查询协程
	if status == 0 {
		//go AdvQueryProducer.Producer(advertiserAuditChan.AdvertiserAuditId, advertiserAuditChan.Logger)
	}
}

// AdvertiserUploadCallBack 将广告主审核的结果回调给送审方，只调一次，不论成功失败
// status:0待审核，1审核通过，2审核不通过
func (p *AdvertiserUploadConsumer) AdvertiserUploadCallBack(uri string, advertiserId int, status int, reason string, logger *zap.Logger) {
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
	logger.Info("advertiser_upload_callback_info", zap.Any("url", uri), zap.Any("req", string(bodyJson)), zap.Any("resp", string(response)), zap.Error(err))
	if err != nil {
		logger.Error("advertiser_upload_callback_info error", zap.Error(err))
	}
}
