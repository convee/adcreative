package crons

import (
	"github.com/convee/adcreative/internal/service"
	"github.com/convee/adcreative/internal/task"
	logger "github.com/convee/adcreative/pkg/log"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AdvertiserQuery struct {
}

func (t *AdvertiserQuery) Run() {
	AdvertiserAuditService := &service.AdvertiserAudit{
		Status:  0,
		IsRsync: 1,
	}
	advertisers := AdvertiserAuditService.GetAll()
	for _, advertiser := range advertisers {
		logg := logger.GetLogger().With(zap.String("request_id", generateAdvertiserID())).With(zap.Int("advertiser_id", advertiser.Id))
		go task.AdvQueryProducer.Producer(advertiser.Id, logg)
	}
}

// generateID 生成随机字符串，eg: 76d27e8c-a80e-48c8-ad20-e5562e0f67e4
func generateAdvertiserID() string {
	reqID, _ := uuid.NewRandom()
	return reqID.String()
}
