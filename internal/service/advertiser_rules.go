package service

import (
	"context"
	"github.com/convee/adcreative/internal/model"
	"github.com/convee/adcreative/internal/pkg/cache"
	logger "github.com/convee/adcreative/pkg/log"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

type AdvertiserRules struct {
	Id          int
	PublisherId *int
	Publisher   string
	Info        string
	Page        int // 第几页
	PerPage     int // 每页显示条数

}

var (
	advertiserRulesModel = model.AdvertiserRulesModel{}
)

func (c *AdvertiserRules) GetList() map[string]interface{} {
	data := make(map[string]interface{})

	page := c.Page
	if page == 0 {
		page = 1
	}
	perPage := c.PerPage
	if perPage == 0 {
		perPage = 20
	}

	list, err := advertiserRulesModel.GetAdvertiserRuless(page, perPage, c.getMaps())
	if err != nil {
		logger.Error("AdvertiserRules get list data err ", zap.Error(err))
		return data
	}
	total, err := advertiserRulesModel.GetAdvertiserRulesTotal(c.getMaps())
	if err != nil {
		logger.Error("AdvertiserRules get list count err ", zap.Error(err))
		return data
	}
	data["lists"] = list
	data["total"] = total
	return data
}

func (c *AdvertiserRules) GetApiList() map[string]interface{} {
	var results []interface{}
	data := make(map[string]interface{})
	publisherService := Publisher{
		Name: c.Publisher,
	}
	PublisherList := make(map[string]int, 9)
	PublisherList["IQIYI"] = 2
	PublisherList["B612"] = 8
	PublisherList["YouTu"] = 29
	PublisherList["Sina"] = 32
	PublisherList["Weibo"] = 33
	PublisherList["Fancy"] = 34
	if _, ok := PublisherList[c.Publisher]; ok {

	} else {
		logger.Error("PublisherName put err ", zap.Error(nil))
		return data
	}
	publisherInfo, err := publisherService.GetPublisherByName()
	if err != nil {
		logger.Error("PublisherName get data err ", zap.Error(err))
		return data
	}
	mslice := make([]map[string]interface{}, 0)
	var values model.ApiAdvertiserRules
	list, err := advertiserRulesModel.GetApiAdvertiserRuless(publisherInfo.Id)
	for _, value := range list {
		values.PublisherId = value.PublisherId
		err := jsoniter.Unmarshal([]byte(value.Info), &mslice)
		values.Rules = mslice
		if err != nil {
			logger.Error("Json Deserialization err ", zap.Error(err))
			return data
		}
		results = append(results, values)
	}
	if err != nil {
		logger.Error("AdvertiserRules get list data err ", zap.Error(err))
		return data
	}
	//total, err := advertiserRulesModel.GetAdvertiserRulesTotal(c.getMaps())
	//if err != nil {
	//	logger.Error("AdvertiserRules get list count err ", zap.Error(err))
	//	return data
	//}
	data["lists"] = results
	//data["total"] = total
	return data
}

func (c *AdvertiserRules) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	return maps
}

func (c *AdvertiserRules) Add() (*model.AdvertiserRules, error) {
	return advertiserRulesModel.CreateAdvertiserRules(&model.AdvertiserRules{
		PublisherId: c.PublisherId,
		Info:        c.Info,
	})
}

func (c *AdvertiserRules) Edit() (*model.AdvertiserRules, error) {
	advertiserRules, err := advertiserRulesModel.GetAdvertiserRulesById(c.Id)
	if err != nil {
		return advertiserRules, err
	}
	err = cache.NewUserCache().Del(context.Background(), cache.RedisKeyAdvertiserRule+cast.ToString(c.Id))
	if err != nil {
		return nil, err
	}
	advertiserRules.PublisherId = c.PublisherId
	advertiserRules.Info = c.Info
	return advertiserRulesModel.UpdateAdvertiserRules(advertiserRules)
}

func (c *AdvertiserRules) Delete(id int) (*model.AdvertiserRules, error) {
	advertiserRules, err := advertiserRulesModel.GetAdvertiserRulesById(id)
	if err != nil {
		return advertiserRules, err
	}
	err = cache.NewUserCache().Del(context.Background(), cache.RedisKeyAdvertiserRule+cast.ToString(id))
	if err != nil {
		return nil, err
	}
	return advertiserRulesModel.DeleteAdvertiserRules(advertiserRules)
}

func (c *AdvertiserRules) GetAdvertiserRulesByPublisherId() (*model.AdvertiserRules, error) {
	advertiserRules, err := advertiserRulesModel.GetAdvertiserRulesByPublisherId(*c.PublisherId)
	if err != nil {
		return nil, err
	}
	return advertiserRules, nil
}
