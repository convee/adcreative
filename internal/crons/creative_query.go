package crons

import (
	"fmt"
	"github.com/convee/adcreative/internal/enum"
	"github.com/convee/adcreative/internal/model"
	"github.com/convee/adcreative/internal/service"
	"time"

	"github.com/convee/adcreative/internal/task"
	logger "github.com/convee/adcreative/pkg/log"
	"go.uber.org/zap"
)

type CreativeQuery struct {
}

func (c *CreativeQuery) Run() {
	for _, p := range enum.PubList {
		// 程序重启时，腾讯媒体加载一次
		c.Query(p)
		go func(mId int) {
			for range time.Tick(time.Minute * 1) {
				c.Query(mId)
			}
		}(p)
	}
}

func (c *CreativeQuery) Query(pubId int) {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("query panic:", pubId, e)
		}
	}()
	errCodes := []int{model.CREATIVE_AUDITING, model.CREATIVE_QUERY_FAILED, model.CREATIVE_AUDIT_UNPASSWD}
	creatives := new(service.Creative).GetAllByErrCodes(pubId, errCodes)
	for _, creative := range creatives {
		logg := logger.GetLogger().With(zap.String("request_id", creative.RequestId)).With(zap.Int("creative_id", creative.Id)).With(zap.Int("pub_id", pubId))
		logg.Info("Query", zap.Int("pubid", creative.PublisherId), zap.String("creativeId", creative.CreativeId), zap.Int("priority", creative.Priority))
		task.CrQueryProducer.Producer(pubId, creative.Id, logg)
	}
}

func (c *CreativeQuery) QueryOne(creativeId string) {
	creativeService := &service.Creative{
		CreativeId: creativeId,
		ErrCode:    5,
	}
	creatives := creativeService.GetAllId()
	if len(creatives) <= 0 {
		fmt.Println("not found creative")
		return
	}
	creative := creatives[0]
	logg := logger.GetLogger().With(zap.String("request_id", creative.RequestId)).With(zap.Int("creative_id", creative.Id)).With(zap.Int("pub_id", creative.PublisherId))
	task.CrQueryProducer.Producer(creative.PublisherId, creative.Id, logg)
}
