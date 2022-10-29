package model

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type AdvertiserRules struct {
	Model
	Id          int	   `json:"id"`
	PublisherId *int   `json:"publisher_id"`
	Info        string `json:"info"`
}

type ApiAdvertiserRules struct {
	PublisherId *int        `json:"publisher_id"`
	Rules       interface{} `json:"rules"`
}

type RuleTemplate struct {
	Key      string   `json:"key"`
	Desc     string   `json:"desc"`
	Size     int      `json:"size"`
	Type     string   `json:"type"`
	Limit    int      `json:"limit"`
	Format   []string `json:"format"`
	Measure  []string `json:"measure"`
	Required int      `json:"required"`
}

type AdvertiserRulesModel struct {
}

// TableName sets the insert table name for this struct type
func (c *AdvertiserRulesModel) TableName() string {
	return "advertiser_rules"
}

func (c AdvertiserRulesModel) GetAdvertiserRulesTotal(maps interface{}) (int64, error) {
	var count int64
	if err := GetDB().Model(&AdvertiserRules{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (c AdvertiserRulesModel) GetAdvertiserRulesById(id int) (*AdvertiserRules, error) {
	var advertiserRules *AdvertiserRules
	err := GetDB().First(&advertiserRules, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("记录不存在")
		}
		return nil, err
	}
	return advertiserRules, nil
}

func (c AdvertiserRulesModel) ExistsById(id int) (bool, error) {
	var advertiserRules *AdvertiserRules
	err := GetDB().First(&advertiserRules, id).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if advertiserRules.Id > 0 {
		return true, nil
	}
	return false, nil
}

func (c AdvertiserRulesModel) GetAdvertiserRuless(page int, perPage int, maps interface{}) ([]*AdvertiserRules, error) {
	var advertiserRuless []*AdvertiserRules
	offset := (page - 1) * perPage
	err := GetDB().Where(maps).Offset(offset).Limit(perPage).Order("id desc").Find(&advertiserRuless).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("记录不存在")
		}
		return nil, err
	}
	return advertiserRuless, nil
}

func (c AdvertiserRulesModel) GetApiAdvertiserRuless(publisherId int) ([]*AdvertiserRules, error) {
	var advertiserRuless []*AdvertiserRules
	err := GetDB().Where("publisher_id=?", publisherId).First(&advertiserRuless).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("记录不存在")
		}
		return nil, err
	}
	return advertiserRuless, nil
}

func (c AdvertiserRulesModel) CreateAdvertiserRules(advertiserRules *AdvertiserRules) (*AdvertiserRules, error) {
	result := GetDB().Create(&advertiserRules)
	if result.Error != nil {
		return advertiserRules, result.Error
	}
	return advertiserRules, nil
}

func (c AdvertiserRulesModel) UpdateAdvertiserRules(advertiserRules *AdvertiserRules) (*AdvertiserRules, error) {
	result := GetDB().Updates(advertiserRules)
	if result.Error != nil {
		return advertiserRules, result.Error
	}
	return advertiserRules, nil
}

func (c AdvertiserRulesModel) DeleteAdvertiserRules(advertiserRules *AdvertiserRules) (*AdvertiserRules, error) {
	GetDB().Delete(&advertiserRules)
	return advertiserRules, nil
}

func (c AdvertiserRulesModel) GetAdvertiserRulesByPublisherId(publisherId int) (*AdvertiserRules, error) {
	var advertiserRules *AdvertiserRules
	err := GetDB().Where("publisher_id=?", publisherId).First(&advertiserRules).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("记录不存在")
		}
		return nil, err
	}
	return advertiserRules, nil
}
