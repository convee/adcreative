package cache

import (
	"context"
	model2 "github.com/convee/adcreative/internal/model"
	"time"

	"github.com/convee/adcreative/pkg/encoding"
	"github.com/convee/adcreative/pkg/redis"
	"github.com/spf13/cast"
)

const (
	RedisKeyPosition           = "position:"
	RedisKeyPublisher          = "publisher:"
	RedisKeyCustomer           = "customer:"
	RedisKeyCreative           = "creative:"
	RedisKeyPublisherAccount   = "publisher_account:"
	RedisKeyAdvertiserAudit    = "advertiser_audit:"
	RedisKeyAdvertiserRule     = "advertiser_rule:"
	RedisKeyCreativeUploadLock = "creative_upload_lock:"
	RedisKeyCreativeQueryLock  = "creative_query_lock:"
	RedisTTL                   = 600 * time.Second
)

// NewUserCache new一个用户cache
func NewUserCache() Cache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := "ms"
	return NewRedisCache(redis.RedisClient, cachePrefix, jsonEncoding)
}

func GetPublisherCacheByName(name string) (*model2.Publisher, error) {
	key := RedisKeyPublisher + name
	var publisherFromCache *model2.Publisher
	err := NewUserCache().Get(context.Background(), key, &publisherFromCache)
	if err != nil {
		return nil, err
	}
	if publisherFromCache != nil {
		return publisherFromCache, nil
	}
	publisherModel := model2.PublisherModel{}
	publisherFromDb, err := publisherModel.GetPublisherByName(name)
	if err != nil {
		return nil, err
	}
	err = NewUserCache().Set(context.Background(), key, publisherFromDb, RedisTTL)
	if err != nil {
		return nil, err
	}
	return publisherFromDb, nil
}

func GetPublisherAccount(customerId int, publisherId int) (*model2.PublisherAccount, error) {
	if customerId != 0 {
		customer, err := new(model2.CustomerModel).GetCustomerById(customerId)
		if err != nil {
			return nil, err
		}
		// 非私有化账号使用默认账号
		if *customer.IsPrivate == 0 {
			customerId = 0
		}
	}
	key := RedisKeyPublisherAccount + cast.ToString(customerId) + cast.ToString(publisherId)
	var publisherAccountFromCache *model2.PublisherAccount
	err := NewUserCache().Get(context.Background(), key, &publisherAccountFromCache)
	if err != nil {
		return nil, err
	}
	if publisherAccountFromCache != nil {
		return publisherAccountFromCache, nil
	}
	maps := make(map[string]interface{})
	maps["customer_id"] = customerId
	maps["publisher_id"] = publisherId

	publisherAccountModel := model2.PublisherAccountModel{}
	publisherAccountFromDb, err := publisherAccountModel.GetOnePublisherAccountByMaps(maps)
	if err != nil {
		return nil, err
	}
	err = NewUserCache().Set(context.Background(), key, publisherAccountFromDb, RedisTTL)
	if err != nil {
		return nil, err
	}
	return publisherAccountFromDb, nil
}

func GetPublisherCacheById(id int) (*model2.Publisher, error) {
	key := RedisKeyPublisher + cast.ToString(id)
	var publisherFromCache *model2.Publisher
	err := NewUserCache().Get(context.Background(), key, &publisherFromCache)
	if err != nil {
		return nil, err
	}
	if publisherFromCache != nil {
		return publisherFromCache, nil
	}
	publisherModel := model2.PublisherModel{}
	publisherFromDb, err := publisherModel.GetPublisherById(id)
	if err != nil {
		return nil, err
	}
	err = NewUserCache().Set(context.Background(), key, publisherFromDb, RedisTTL)
	if err != nil {
		return nil, err
	}
	return publisherFromDb, nil

}

func GetAdvertiserAuditCache(customerId int, advertiserId int, publisherId int) (*model2.AdvertiserAudit, error) {
	key := RedisKeyAdvertiserAudit + cast.ToString(customerId) + cast.ToString(advertiserId) + cast.ToString(publisherId)
	var advertiserAuditFromCache *model2.AdvertiserAudit
	err := NewUserCache().Get(context.Background(), key, &advertiserAuditFromCache)
	if err != nil {
		return nil, err
	}
	if advertiserAuditFromCache != nil {
		return advertiserAuditFromCache, nil
	}
	maps := make(map[string]interface{})
	maps["customer_id"] = customerId
	maps["advertiser_id"] = advertiserId
	maps["publisher_id"] = publisherId
	advertiserAuditModel := model2.AdvertiserAuditModel{}
	advertiserAuditFromDB, err := advertiserAuditModel.GetOneAdvertiserAudit(maps)
	if err != nil {
		return nil, err
	}
	err = NewUserCache().Set(context.Background(), key, advertiserAuditFromDB, RedisTTL)
	if err != nil {
		return nil, err
	}
	return advertiserAuditFromDB, nil

}

func GetAdvertiserRuleCache(id int) (*model2.AdvertiserRules, error) {
	key := RedisKeyAdvertiserRule + cast.ToString(id)
	var advertiserRuleFromCache *model2.AdvertiserRules
	err := NewUserCache().Get(context.Background(), key, &advertiserRuleFromCache)
	if err != nil {
		return nil, err
	}
	if advertiserRuleFromCache != nil {
		return advertiserRuleFromCache, nil
	}
	advertiserAuditModel := model2.AdvertiserRulesModel{}
	advertiserRuleFromDB, err := advertiserAuditModel.GetAdvertiserRulesById(id)
	if err != nil {
		return nil, err
	}
	err = NewUserCache().Set(context.Background(), key, advertiserRuleFromDB, RedisTTL)
	if err != nil {
		return nil, err
	}
	return advertiserRuleFromDB, nil

}

func GetPositionCache(position string) (*model2.Position, error) {

	key := RedisKeyPosition + position
	var positionFromCache *model2.Position
	err := NewUserCache().Get(context.Background(), key, &positionFromCache)
	if err != nil {
		return nil, err
	}
	if positionFromCache != nil {
		return positionFromCache, nil
	}
	positionModel := model2.PositionModel{}
	positionFromDB, err := positionModel.GetPosition(position)
	if err != nil {
		return nil, err
	}
	err = NewUserCache().Set(context.Background(), key, positionFromDB, RedisTTL)
	if err != nil {
		return nil, err
	}
	return positionFromDB, nil

}

func GetPositionCacheById(positionId int) (*model2.Position, error) {

	key := RedisKeyPosition + cast.ToString(positionId)
	var positionFromCache *model2.Position
	err := NewUserCache().Get(context.Background(), key, &positionFromCache)
	if err != nil {
		return nil, err
	}
	if positionFromCache != nil {
		return positionFromCache, nil
	}
	positionModel := model2.PositionModel{}
	positionFromDB, err := positionModel.GetPositionById(positionId)
	if err != nil {
		return nil, err
	}
	err = NewUserCache().Set(context.Background(), key, positionFromDB, RedisTTL)
	if err != nil {
		return nil, err
	}
	return positionFromDB, nil

}

func GetCustomerCache(appId string) (*model2.Customer, error) {

	key := RedisKeyCustomer + appId
	var customerFromCache *model2.Customer
	err := NewUserCache().Get(context.Background(), key, &customerFromCache)
	if err != nil {
		return nil, err
	}
	if customerFromCache != nil {
		return customerFromCache, err
	}
	customerModel := model2.CustomerModel{}
	customerFromDB, err := customerModel.GetCustomerByAppid(appId)
	if err != nil {
		return nil, err
	}
	err = NewUserCache().Set(context.Background(), key, customerFromDB, RedisTTL)
	if err != nil {
		return nil, err
	}
	return customerFromDB, nil
}

func GetCustomerCacheById(customerId int) (*model2.Customer, error) {

	key := RedisKeyCustomer + cast.ToString(customerId)
	var customerFromCache *model2.Customer
	err := NewUserCache().Get(context.Background(), key, &customerFromCache)
	if err != nil {
		return nil, err
	}
	if customerFromCache != nil {
		return customerFromCache, err
	}
	customerModel := model2.CustomerModel{}
	customerFromDB, err := customerModel.GetCustomerById(customerId)
	if err != nil {
		return nil, err
	}
	err = NewUserCache().Set(context.Background(), key, customerFromDB, RedisTTL)
	if err != nil {
		return nil, err
	}
	return customerFromDB, nil
}

func GetCreativeCacheById(id int) (*model2.Creative, error) {

	// 获取创意信息
	creativeModel := model2.CreativeModel{}
	creativeFromDB, err := creativeModel.GetCreativeById(id)
	if err != nil {
		return nil, err
	}

	return creativeFromDB, nil
}

func LockUploadCreative(creativeId int) (bool, error) {
	key := RedisKeyCreativeUploadLock + cast.ToString(creativeId)
	return NewUserCache().SetNX(context.Background(), key, 1, 60*time.Second)
}

func UnLockUploadCreative(creativeId int) error {
	key := RedisKeyCreativeUploadLock + cast.ToString(creativeId)
	return NewUserCache().Del(context.Background(), key)
}

func LockUploadUrlUnique(urlUnique string, val string) (bool, error) {
	return NewUserCache().SetNX(context.Background(), urlUnique, val, 120*time.Second)
}

func GetUploadUrlUnique(urlUnique string, val interface{}) error {
	return NewUserCache().Get(context.Background(), urlUnique, val)
}

func SetUploadUrlUnique(urlUnique string, val interface{}) error {
	return NewUserCache().Set(context.Background(), urlUnique, val, 86400*time.Second)
}

func DelUploadUrlUnique(urlUnique string) error {
	return NewUserCache().Del(context.Background(), urlUnique)
}
