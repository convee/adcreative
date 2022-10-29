package service

import (
	"context"
	"errors"
	"github.com/convee/adcreative/internal/model"
	"github.com/convee/adcreative/internal/pkg/cache"
	logger "github.com/convee/adcreative/pkg/log"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

type Publisher struct {
	Id                             int
	Name                           string
	IsRsyncAdvertiser              *int
	IsRsyncCreative                *int
	IsPublisherCdn                 *int
	IsCreativeBind                 *int
	MonitorCodeChangeNeedRsync     *int
	LandingChangeNeedRsync         *int
	MonitorPositionChangeNeedRsync *int
	S2sStateInfo                   string
	PubReturnInfo                  string
	PvLimit                        int
	ClLimit                        int
	Nickname                       string
	AdvertiserUrls                 string
	CreativeUrls                   string
	Page                           int // 第几页
	PerPage                        int // 每页显示条数

}

var (
	publisherModel = model.PublisherModel{}
)

func (c *Publisher) GetList() map[string]interface{} {
	data := make(map[string]interface{})

	page := c.Page
	if page == 0 {
		page = 1
	}
	perPage := c.PerPage
	if perPage == 0 {
		perPage = 20
	}

	list, err := publisherModel.GetPublishers(page, perPage, c.getMaps())
	lists := make([]map[string]interface{}, 0)
	for _, val := range list {
		list := make(map[string]interface{})
		list["id"] = val.Id
		list["name"] = val.Name
		list["is_rsync_advertiser"] = val.IsRsyncAdvertiser
		list["is_rsync_creative"] = val.IsRsyncCreative
		list["created_at"] = val.CreatedAt
		list["updated_at"] = val.UpdatedAt
		list["deleted_at"] = val.DeletedAt
		if find := val.AdvertiserUrls == "{}"; find {
			val.AdvertiserUrls = ""
		}
		if find := val.CreativeUrls == "{}"; find {
			val.CreativeUrls = ""
		}
		list["advertiser_urls"] = val.AdvertiserUrls
		list["creative_urls"] = val.CreativeUrls
		list["pv_limit"] = val.PvLimit
		list["cl_limit"] = val.ClLimit
		list["nick_name"] = val.Nickname
		lists = append(lists, list)
	}
	if err != nil {
		logger.Error("publisher get list data err ", zap.Error(err))
		return data
	}
	total, err := publisherModel.GetPublisherTotal(c.getMaps())
	if err != nil {
		logger.Error("publisher get list count err ", zap.Error(err))
		return data
	}
	data["lists"] = lists
	data["total"] = total
	return data
}

func (c *Publisher) GetApiList() map[string]interface{} {
	data := make(map[string]interface{})
	list, err := publisherModel.GetApiPublishers(c.getMaps())
	if err != nil {
		logger.Error("publisher get list data err ", zap.Error(err))
		return data
	}
	lists := make([]map[string]interface{}, 0)
	for _, val := range list {
		list := make(map[string]interface{})
		list["id"] = val.Id
		list["name"] = val.Name
		lists = append(lists, list)
	}
	data["lists"] = lists
	return data
}

func (c *Publisher) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	if len(c.Name) > 0 {
		maps["name"] = c.Name
	}
	if c.Id > 0 {
		maps["id"] = c.Id
	}
	return maps
}

func (c *Publisher) Add() (*model.Publisher, error) {
	return publisherModel.CreatePublisher(&model.Publisher{
		Name:              c.Name,
		IsRsyncAdvertiser: c.IsRsyncAdvertiser,
		IsRsyncCreative:   c.IsRsyncCreative,
		PvLimit:           c.PvLimit,
		ClLimit:           c.ClLimit,
		Nickname:          c.Nickname,
		AdvertiserUrls:    c.AdvertiserUrls,
		CreativeUrls:      c.CreativeUrls,
	})
}

func (c *Publisher) Edit() (*model.Publisher, error) {
	publisher, err := publisherModel.GetPublisherById(c.Id)
	if err != nil {
		return publisher, err
	}
	publisher.Name = c.Name
	publisher.IsRsyncAdvertiser = c.IsRsyncAdvertiser
	publisher.IsRsyncCreative = c.IsRsyncCreative
	publisher.PvLimit = c.PvLimit
	publisher.ClLimit = c.ClLimit
	publisher.Nickname = c.Nickname
	publisher.AdvertiserUrls = c.AdvertiserUrls
	publisher.CreativeUrls = c.CreativeUrls
	err1 := cache.NewUserCache().Del(context.Background(), cache.RedisKeyPublisher+publisher.Name)
	err2 := cache.NewUserCache().Del(context.Background(), cache.RedisKeyPublisher+cast.ToString(publisher.Id))
	if err1 != nil || err2 != nil {
		return nil, errors.New("缓存删除失败")
	}
	return publisherModel.UpdatePublisher(publisher)
}

func (c *Publisher) Delete() (*model.Publisher, error) {
	publisher, err := publisherModel.GetPublisherById(c.Id)
	if err != nil {
		return publisher, err
	}
	err1 := cache.NewUserCache().Del(context.Background(), cache.RedisKeyPublisher+publisher.Name)
	err2 := cache.NewUserCache().Del(context.Background(), cache.RedisKeyPublisher+cast.ToString(publisher.Id))
	if err1 != nil || err2 != nil {
		return nil, errors.New("缓存删除失败")
	}
	return publisherModel.DeletePublisher(publisher)
}

func (c *Publisher) GetPublisherInfo() *model.Publisher {
	publisher, err := publisherModel.GetPublisherById(c.Id)
	if err != nil {
		logger.Error("publisher get info err ", zap.Error(err))
		return publisher
	}
	return publisher
}

func (c *Publisher) GetPublisherByName() (*model.Publisher, error) {
	publisher, err := publisherModel.GetPublisherByName(c.Name)
	if err != nil {
		return nil, err
	}
	return publisher, nil
}
