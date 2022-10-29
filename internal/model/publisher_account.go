package model

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type PublisherAccount struct {
	Model
	Id          int    `json:"id"`
	CustomerId  int    `json:"customer_id"`
	PublisherId int    `json:"publisher_id"`
	DspId       string `json:"dsp_id"`
	Token       string `json:"token"`
	Remark      string `json:"remark"`
	CallbackUrl string `json:"callback_url"`
	Extra       string `json:"extra"`
}

type PublisherAccountModel struct {
}

// TableName sets the insert table name for this struct type
func (c *PublisherAccountModel) TableName() string {
	return "publisher_account"
}

func (c PublisherAccountModel) GetPublisherAccountTotal(maps interface{}) (int64, error) {
	var count int64
	if err := GetDB().Model(&PublisherAccount{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (c PublisherAccountModel) GetPublisherAccountById(id int) (*PublisherAccount, error) {
	var publisherAccount *PublisherAccount
	err := GetDB().First(&publisherAccount, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("记录不存在")
		}
		return nil, err
	}
	return publisherAccount, nil
}

func (c PublisherAccountModel) GetOnePublisherAccountByMaps(maps interface{}) (*PublisherAccount, error) {
	var publisherAccount *PublisherAccount
	err := GetDB().Where(maps).First(&publisherAccount).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("记录不存在")
		}
		return nil, err
	}
	return publisherAccount, nil
}

func (c PublisherAccountModel) GetPublisherAccounts(page int, perPage int, maps interface{}) ([]*PublisherAccount, error) {
	var publisherAccounts []*PublisherAccount
	offset := (page - 1) * perPage
	err := GetDB().Where(maps).Offset(offset).Limit(perPage).Order("id desc").Find(&publisherAccounts).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("记录不存在")
		}
		return nil, err
	}
	return publisherAccounts, nil
}

func (c PublisherAccountModel) CreatePublisherAccount(publisherAccount *PublisherAccount) (*PublisherAccount, error) {
	result := GetDB().Create(&publisherAccount)
	if result.Error != nil {
		return publisherAccount, result.Error
	}
	return publisherAccount, nil
}

func (c PublisherAccountModel) UpdatePublisherAccount(publisherAccount *PublisherAccount) (*PublisherAccount, error) {
	result := GetDB().Updates(publisherAccount)
	if result.Error != nil {
		return publisherAccount, result.Error
	}
	return publisherAccount, nil
}

func (c PublisherAccountModel) DeletePublisherAccount(publisherAccount *PublisherAccount) (*PublisherAccount, error) {
	GetDB().Delete(&publisherAccount)
	return publisherAccount, nil
}
