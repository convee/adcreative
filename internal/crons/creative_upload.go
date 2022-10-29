package crons

import (
	"fmt"
	"github.com/convee/adcreative/internal/enum"
	"github.com/convee/adcreative/internal/model"
	service2 "github.com/convee/adcreative/internal/service"
	"github.com/convee/adcreative/internal/task"
	"github.com/convee/adcreative/pkg/log"
	"go.uber.org/zap"
	"time"
)

type CreativeUpload struct {
}

func (c *CreativeUpload) Run() {
	for _, p := range enum.PubList {
		// 程序重启时,媒体加载一次
		c.Upload(p)

		go func(mediaId int) {
			t := time.Duration(2)
			if mediaId == enum.PUB_TENCENT {
				t = time.Duration(5)
			}

			for range time.Tick(time.Minute * t) {
				c.Upload(mediaId)
			}
		}(p)

	}
}

func (c *CreativeUpload) Upload(pubId int) {
	start := time.Now()
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("upload panic:", pubId, e)
		}
	}()
	errCodes := []int{model.CREATIVE_UPLOADING, model.CREATIVE_UPLOAD_FAILED}
	creatives := new(service2.Creative).GetAllByErrCodes(pubId, errCodes)
	costs := time.Now().Sub(start).Milliseconds()

	for _, creative := range creatives {
		if creative == nil {
			continue
		}
		log.Info("Upload", zap.Int("pubid", creative.PublisherId), zap.String("creativeId", creative.CreativeId), zap.Int("priority", creative.Priority))
		positionService := &service2.Position{Id: creative.PositionId}
		positionInfo := positionService.GetPositionInfo()
		if positionInfo == nil {
			continue
		}

		logger := log.GetLogger().
			With(zap.String("method", "Upload")).
			With(zap.String("request_id", creative.RequestId)).
			With(zap.String("creative_id", creative.CreativeId)).
			With(zap.Int("customer_id", creative.CustomerId)).
			With(zap.Int("advertiser_id", creative.AdvertiserId)).
			With(zap.String("position", positionInfo.Position))

		// 添加到送审chan中
		logger.With(zap.Int("id", creative.Id)).With(zap.Int("publisher_id", creative.PublisherId)).With(zap.Int64("CreativeUpload getallId", costs))
		task.CrUploadProducer.Producer(creative, logger)
	}

}

func (c *CreativeUpload) UploadOne(creativeId string) {
	creativeService := &service2.Creative{
		CreativeId: creativeId,
	}
	creatives := creativeService.GetAllId()
	if len(creatives) <= 0 {
		fmt.Println("not found creative")
		return
	}
	creative := creatives[0]
	logger := log.GetLogger().
		With(zap.String("method", "Upload")).
		With(zap.String("request_id", creative.RequestId)).
		With(zap.String("creative_id", creative.CreativeId)).
		With(zap.Int("id", creative.Id)).
		With(zap.Int("publisher_id", creative.PublisherId))
	// 添加到送审chan中
	task.CrUploadProducer.Producer(creative, logger)
}
