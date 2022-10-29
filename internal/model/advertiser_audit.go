package model

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const (
	// 1上传失败、2上传不通过、3待送审、4送审失败、5审核中、6查询失败、7审核通过、8审核不通过
	ADVERTISER_UPLOAD_FAILED   = 1
	ADVERTISER_UPLOAD_UNPASSED = 2
	ADVERTISER_UPLOADING       = 3
	ADVERTISER_AUDIT_FAILED    = 4
	ADVERTISER_AUDITING        = 5
	ADVERTISER_QUERY_FAILED    = 6
	ADVERTISER_AUDIT_PASSED    = 7
	ADVERTISER_AUDIT_UNPASSWD  = 8

	// 0待审核，1审核通过，2审核不通过
	ADVERTISER_STATUS_AUDITING = 0
	ADVERTISER_STATUS_PASSED   = 1
	ADVERTISER_STATUS_UNPASSED = 2
)

type AdvertiserAudit struct {
	Model
	Id                 int    `json:"id"`
	AdvertiserName     string `json:"advertiser_name"`
	AdvertiserId       int    `json:"advertiser_id"`
	PublisherId        int    `json:"publisher_id"`
	CustomerId         int    `json:"customer_id"`
	PublisherAccountId int    `json:"publisher_account_id"`
	Status             int    `json:"status"`
	Info               string `json:"info"`
	ErrCode            int    `json:"err_code"`
	ErrMsg             string `json:"err_msg"`
	IsRsync            int    `json:"is_rsync"`
	MediaCid           string `json:"media_cid"`
	Extra              string `json:"extra"`
}

type AdvertiserAuditInfo struct {
	CompanyName      string
	CompanySummary   string
	WebsiteName      string
	WebsiteAddress   string
	WebsiteNumber    string
	BusinessLicenser string
	AuthorizeState   string
	Industry         string
	Qualifications   []Qualification
	Extra            map[string]interface{}
}

type Qualification struct {
	FileName string
	FileUrl  string
}

type AdvertiserAuditModel struct {
}

// TableName sets the insert table name for this struct type
func (c *AdvertiserAuditModel) TableName() string {
	return "advertiser_audit"
}

func (c AdvertiserAuditModel) GetAdvertiserAuditTotal(maps interface{}) (int64, error) {
	var count int64
	if err := GetDB().Model(&AdvertiserAudit{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (c AdvertiserAuditModel) GetAllAdvertiserAuditsByMaps(maps interface{}) ([]*AdvertiserAudit, error) {
	var advertiserAudits []*AdvertiserAudit
	err := GetDB().Where(maps).Find(&advertiserAudits).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("记录不存在")
		}
		return nil, err
	}
	return advertiserAudits, nil
}

func (c AdvertiserAuditModel) GetAdvertiserAuditById(id int) (*AdvertiserAudit, error) {
	var advertiserAudit *AdvertiserAudit
	err := GetDB().First(&advertiserAudit, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("记录不存在")
		}
		return nil, err
	}
	return advertiserAudit, nil
}

func (c AdvertiserAuditModel) GetOneAdvertiserAuditByMaps(maps interface{}) (*AdvertiserAudit, error) {
	var advertiserAudit *AdvertiserAudit
	err := GetDB().Where(maps).First(&advertiserAudit).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return advertiserAudit, nil
}

func (c AdvertiserAuditModel) GetOneAdvertiserAudit(maps interface{}) (*AdvertiserAudit, error) {
	var advertiserAudit *AdvertiserAudit
	err := GetDB().Where(maps).First(&advertiserAudit).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("记录不存在")
		}
		return nil, err
	}
	return advertiserAudit, nil
}

func (c AdvertiserAuditModel) GetAdvertiserAudits(page int, perPage int, maps interface{}) ([]*AdvertiserAudit, error) {
	var advertiserAudits []*AdvertiserAudit
	offset := (page - 1) * perPage
	err := GetDB().Where(maps).Offset(offset).Limit(perPage).Order("id desc").Find(&advertiserAudits).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("记录不存在")
		}
		return nil, err
	}
	return advertiserAudits, nil
}

func (c AdvertiserAuditModel) CreateAdvertiserAudit(advertiserAudit *AdvertiserAudit) (*AdvertiserAudit, error) {
	result := GetDB().Create(&advertiserAudit)
	if result.Error != nil {
		return advertiserAudit, result.Error
	}
	if advertiserAudit.AdvertiserId <= 0 {
		//创建数据后马上更新广告主id
		maps := make(map[string]interface{})
		maps["advertiser_id"] = advertiserAudit.Id
		resultUpdate := GetDB().Model(&advertiserAudit).Updates(maps)
		if resultUpdate.Error != nil {
			return advertiserAudit, result.Error
		}
	}
	return advertiserAudit, nil
}

func (c AdvertiserAuditModel) UpdateAdvertiserAudit(advertiserAudit *AdvertiserAudit) (*AdvertiserAudit, error) {
	result := GetDB().Updates(advertiserAudit)
	if result.Error != nil {
		return advertiserAudit, result.Error
	}
	return advertiserAudit, nil
}

func (c AdvertiserAuditModel) UpdateCreativeByMap(id int, maps map[string]interface{}) (AdvertiserAudit, error) {
	var advertiserAudit AdvertiserAudit
	advertiserAudit.Id = id
	result := GetDB().Model(&advertiserAudit).Updates(maps)
	if result.Error != nil {
		return advertiserAudit, result.Error
	}
	return advertiserAudit, nil
}

func (c AdvertiserAuditModel) DeleteAdvertiserAudit(advertiserAudit *AdvertiserAudit) (*AdvertiserAudit, error) {
	GetDB().Delete(&advertiserAudit)
	return advertiserAudit, nil
}

func GetAdvertiserStatusByErrCode(errCode int) int {
	if errCode == ADVERTISER_UPLOADING || errCode == ADVERTISER_AUDITING || errCode == ADVERTISER_QUERY_FAILED {
		return ADVERTISER_STATUS_AUDITING
	} else if errCode == ADVERTISER_AUDIT_PASSED {
		return ADVERTISER_STATUS_PASSED
	} else if errCode == ADVERTISER_UPLOAD_FAILED || errCode == ADVERTISER_UPLOAD_UNPASSED || errCode == ADVERTISER_AUDIT_FAILED || errCode == ADVERTISER_AUDIT_UNPASSWD {
		return ADVERTISER_STATUS_UNPASSED
	}
	return ADVERTISER_STATUS_AUDITING
}
