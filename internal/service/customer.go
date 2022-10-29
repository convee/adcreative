package service

import (
	"context"
	"fmt"
	"github.com/convee/adcreative/internal/model"
	"github.com/convee/adcreative/internal/pkg/cache"
	"github.com/convee/adcreative/internal/pkg/common"
	logger "github.com/convee/adcreative/pkg/log"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"strings"
)

type Customer struct {
	Id                    int
	CustomerId            int
	PublisherId           int
	IsPrivate             *int
	Name                  string // 名称
	CreativeCallbackUrl   string // 创意回调地址
	AdvertiserCallbackUrl string // 广告主回调地址
	Page                  int    // 第几页
	PerPage               int    // 每页显示条数

}

var (
	customerModel = model.CustomerModel{}
)

func (c *Customer) GetList() map[string]interface{} {
	data := make(map[string]interface{})

	page := c.Page
	if page == 0 {
		page = 1
	}
	perPage := c.PerPage
	if perPage == 0 {
		perPage = 20
	}

	list, err := customerModel.GetCustomers(page, perPage, c.getMaps())
	if err != nil {
		logger.Error("customer get list data err ", zap.Error(err))
		return data
	}
	total, err := customerModel.GetCustomerTotal(c.getMaps())
	if err != nil {
		logger.Error("customer get list count err ", zap.Error(err))
		return data
	}
	data["lists"] = list
	data["total"] = total
	return data
}

func (c *Customer) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	if len(c.Name) > 0 {
		maps["name"] = c.Name
	}
	return maps
}

func (c *Customer) Add() (*model.Customer, error) {
	return customerModel.CreateCustomer(&model.Customer{
		Name:                  c.Name,
		Appid:                 fmt.Sprintf("YZ%s", common.GetRandomInt(8)),
		Secret:                strings.ToUpper(uuid.NewV4().String()),
		IsPrivate:             c.IsPrivate,
		CreativeCallbackUrl:   c.CreativeCallbackUrl,
		AdvertiserCallbackUrl: c.AdvertiserCallbackUrl,
	})
}

func (c *Customer) Edit() (*model.Customer, error) {
	customer, err := customerModel.GetCustomerById(c.Id)
	if err != nil {
		return customer, err
	}
	customer.Name = c.Name
	customer.IsPrivate = c.IsPrivate
	customer.CreativeCallbackUrl = c.CreativeCallbackUrl
	customer.AdvertiserCallbackUrl = c.AdvertiserCallbackUrl
	err1 := cache.NewUserCache().Del(context.Background(), cache.RedisKeyCustomer+customer.Appid)
	err2 := cache.NewUserCache().Del(context.Background(), cache.RedisKeyCustomer+cast.ToString(customer.Id))
	if err1 != nil || err2 != nil {
		return nil, errors.New("缓存删除失败")
	}
	return customerModel.UpdateCustomer(customer)
}

func (c *Customer) Delete(id int) (*model.Customer, error) {
	customer, err := customerModel.GetCustomerById(id)
	if err != nil {
		return customer, err
	}
	err1 := cache.NewUserCache().Del(context.Background(), cache.RedisKeyCustomer+customer.Appid)
	err2 := cache.NewUserCache().Del(context.Background(), cache.RedisKeyCustomer+cast.ToString(customer.Id))
	if err1 != nil || err2 != nil {
		return nil, errors.New("缓存删除失败")
	}
	return customerModel.DeleteCustomer(customer)
}
