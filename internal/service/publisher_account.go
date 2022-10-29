package service

import (
	"context"
	"github.com/convee/adcreative/internal/model"
	"github.com/convee/adcreative/internal/pkg/cache"
	logger "github.com/convee/adcreative/pkg/log"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

type PublisherAccount struct {
	Id          int
	CustomerId  int
	PublisherId int
	DspId       string
	Token       string
	Remark      string
	CallbackUrl string
	Page        int // 第几页
	PerPage     int // 每页显示条数

}

var (
	publisherAccountModel = model.PublisherAccountModel{}
)

func (pa *PublisherAccount) GetList() map[string]interface{} {
	data := make(map[string]interface{})

	page := pa.Page
	if page == 0 {
		page = 1
	}
	perPage := pa.PerPage
	if perPage == 0 {
		perPage = 20
	}

	list, err := publisherAccountModel.GetPublisherAccounts(page, perPage, pa.getMaps())
	if err != nil {
		logger.Error("publisherAccount get list data err ", zap.Error(err))
		return data
	}
	total, err := publisherAccountModel.GetPublisherAccountTotal(pa.getMaps())
	if err != nil {
		logger.Error("publisherAccount get list count err ", zap.Error(err))
		return data
	}
	data["lists"] = list
	data["total"] = total
	return data
}

func (pa *PublisherAccount) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})

	if pa.CustomerId > 0 {
		maps["customer_id"] = pa.CustomerId
	}
	if pa.PublisherId > 0 {
		maps["publisher_id"] = pa.PublisherId
	}
	return maps
}

func (pa *PublisherAccount) Add() (*model.PublisherAccount, error) {
	return publisherAccountModel.CreatePublisherAccount(&model.PublisherAccount{
		CustomerId:  pa.CustomerId,
		Token:       pa.Token,
		DspId:       pa.DspId,
		PublisherId: pa.PublisherId,
		CallbackUrl: pa.CallbackUrl,
		Remark:      pa.Remark,
	})
}

func (pa *PublisherAccount) GetOnePublisherAccountByMaps() (*model.PublisherAccount, error) {
	account, err := publisherAccountModel.GetOnePublisherAccountByMaps(pa.getMaps())
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (pa *PublisherAccount) CreateOrUpdate() (*model.PublisherAccount, error) {
	exists, err := publisherAccountModel.GetOnePublisherAccountByMaps(pa.getMaps())
	if err != nil {
		return nil, err
	}
	publisherCount := &model.PublisherAccount{
		CustomerId:  pa.CustomerId,
		Token:       pa.Token,
		DspId:       pa.DspId,
		PublisherId: pa.PublisherId,
		CallbackUrl: pa.CallbackUrl,
		Remark:      pa.Remark,
	}
	err = cache.NewUserCache().Del(context.Background(), cache.RedisKeyPublisherAccount+cast.ToString(exists.CustomerId)+cast.ToString(exists.PublisherId))
	if err != nil {
		return nil, err
	}
	if exists.Id > 0 {
		publisherCount.Id = exists.Id
		return publisherAccountModel.UpdatePublisherAccount(publisherCount)
	} else {
		return publisherAccountModel.CreatePublisherAccount(publisherCount)
	}
}

func (pa *PublisherAccount) Edit() (*model.PublisherAccount, error) {
	publisherAccount, err := publisherAccountModel.GetPublisherAccountById(pa.Id)
	if err != nil {
		return publisherAccount, err
	}
	publisherAccount.Token = pa.Token
	publisherAccount.CustomerId = pa.CustomerId
	publisherAccount.DspId = pa.DspId
	publisherAccount.PublisherId = pa.PublisherId
	publisherAccount.CallbackUrl = pa.CallbackUrl
	publisherAccount.Remark = pa.Remark
	err = cache.NewUserCache().Del(context.Background(), cache.RedisKeyPublisherAccount+cast.ToString(publisherAccount.CustomerId)+cast.ToString(publisherAccount.PublisherId))
	if err != nil {
		return nil, err
	}
	return publisherAccountModel.UpdatePublisherAccount(publisherAccount)
}

func (pa *PublisherAccount) Delete(id int) (*model.PublisherAccount, error) {
	publisherAccount, err := publisherAccountModel.GetPublisherAccountById(id)
	if err != nil {
		return publisherAccount, err
	}
	err = cache.NewUserCache().Del(context.Background(), cache.RedisKeyPublisherAccount+cast.ToString(publisherAccount.CustomerId)+cast.ToString(publisherAccount.PublisherId))
	if err != nil {
		return nil, err
	}
	return publisherAccountModel.DeletePublisherAccount(publisherAccount)
}

func (pa *PublisherAccount) GetPublisherAccountInfo() *model.PublisherAccount {
	account, err := publisherAccountModel.GetPublisherAccountById(pa.Id)
	if err != nil {
		logger.Error("account get info err ", zap.Error(err))
		return account
	}
	return account
}
