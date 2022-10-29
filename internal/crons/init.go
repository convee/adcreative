package crons

import (
	"github.com/convee/adcreative/internal/model"
	"log"

	"github.com/convee/adcreative/configs"
	"github.com/convee/adcreative/internal/pkg/stats"
	"github.com/convee/adcreative/internal/task"
	logger "github.com/convee/adcreative/pkg/log"
	"github.com/robfig/cron"
	"go.uber.org/zap"
)

func Init() {
	// second     = field(fields[0], seconds)
	// minute     = field(fields[1], minutes)
	// hour       = field(fields[2], hours)
	// dayofmonth = field(fields[3], dom)
	// month      = field(fields[4], months)
	// dayofweek  = field(fields[5], dow)
	c := cron.New()
	if configs.Conf.Cron.Push {

		new(CreativeQuery).Run()
		new(CreativeUpload).Run()
		_ = c.AddFunc("0 */5 * * * *", func() {
			new(AdvertiserQuery).Run()
		})
		c.AddFunc("*/15 * * * * *", func() {
			collectCreativeStats()
		})
	}

	c.Start()
	log.Println("cron work starting...")

}

func collectCreativeStats() {
	task.CollectChanStats()
	lst, err := model.GetGroupCreativeStats()
	if err != nil {
		logger.GetLogger().Error("GetGroupCreativeStats failed", zap.Error(err))
		return
	}
	stats.CollectMaterialTaskReset()
	for _, n := range lst {
		stats.CollectMaterialTaskStatus(n.CustomerId, n.AdvertiserId, n.PublisherId, n.StatusType, n.Count)
	}
}
