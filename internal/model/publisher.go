package model

import (
	"errors"
	"gorm.io/gorm"
)

type Publisher struct {
	Model
	Id                int    `json:"id"`
	Name              string `json:"name"`
	IsRsyncAdvertiser *int   `json:"is_rsync_advertiser"`
	IsRsyncCreative   *int   `json:"is_rsync_creative"`
	PvLimit           int    `json:"pv_limit"`
	ClLimit           int    `json:"cl_limit"`
	Nickname          string `json:"nickname"`
	AdvertiserUrls    string `json:"advertiser_urls"`
	CreativeUrls      string `json:"creative_urls"`
}
type AdvertiserUrls struct {
	CreateUrl string `json:"create_url"`
	QueryUrl  string `json:"query_url"`
	UpdateUrl string `json:"update_url"`
}

type CreativeUrls struct {
	CreateUrl              string            `json:"create_url"`
	QueryUrl               string            `json:"query_url"`
	ExtendUrl              string            `json:"extend_url"`
	KpQueryUrl             string            `json:"kp_query_url"`
	UpdateUrl              string            `json:"update_url"`
	CdnUrl                 string            `json:"cdn_url"`
	DeleteUrl              string            `json:"delete_url"`
	PicUploadUrl           string            `json:"pic_upload_url"`
	InitMediaUploadUrl     string            `json:"init_media_upload_url"`
	CheckMediaUploadUrl    string            `json:"check_media_upload_url"`
	PartMediaUploadUrl     string            `json:"part_media_upload_url"`
	CompleteMediaUploadUrl string            `json:"complete_media_upload_url"`
	MultiUrl               map[string]string `json:"multi_url"`
}
type PublisherModel struct {
}

// TableName sets the insert table name for this struct type
func (c *PublisherModel) TableName() string {
	return "publisher"
}

func (c PublisherModel) GetPublisherTotal(maps interface{}) (int64, error) {
	var count int64
	if err := GetDB().Model(&Publisher{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (c PublisherModel) GetPublisherById(id int) (*Publisher, error) {
	var publisher *Publisher
	err := GetDB().First(&publisher, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("记录不存在")
		}
		return nil, err
	}
	return publisher, nil
}

func (c PublisherModel) ExistsById(id int) (bool, error) {
	var publisher *Publisher
	err := GetDB().First(&publisher, id).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if publisher.Id > 0 {
		return true, nil
	}
	return false, nil
}

func (c PublisherModel) GetPublisherByName(name string) (*Publisher, error) {
	var publisher *Publisher
	err := GetDB().Where("name=?", name).First(&publisher).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("记录不存在")
		}
		return nil, err
	}
	return publisher, nil
}

func (c PublisherModel) ExistsByName(name string) (bool, error) {
	var publisher *Publisher
	err := GetDB().Where("name=?", name).First(&publisher).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if publisher.Id > 0 {
		return true, nil
	}
	return false, nil
}

func (c PublisherModel) GetPublishers(page int, perPage int, maps interface{}) ([]*Publisher, error) {
	var publishers []*Publisher
	offset := (page - 1) * perPage
	err := GetDB().Where(maps).Offset(offset).Limit(perPage).Order("id desc").Find(&publishers).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("记录不存在")
		}
		return nil, err
	}
	return publishers, nil
}

func (c PublisherModel) GetApiPublishers(maps interface{}) ([]*Publisher, error) {
	var publishers []*Publisher
	err := GetDB().Where(maps).Find(&publishers).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("记录不存在")
		}
		return nil, err
	}
	return publishers, nil
}

func (c PublisherModel) GetAllPublishers() ([]*Publisher, error) {
	var publishers []*Publisher
	err := GetDB().Find(&publishers).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("记录不存在")
		}
		return nil, err
	}
	return publishers, nil
}

func (c PublisherModel) GetAllPublisherNames() ([]*Publisher, error) {
	var publishers []*Publisher
	err := GetDB().Select("id,name").Find(&publishers).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("记录不存在")
		}
		return nil, err
	}
	return publishers, nil
}

func (c PublisherModel) CreatePublisher(publisher *Publisher) (*Publisher, error) {
	result := GetDB().Create(&publisher)
	if result.Error != nil {
		return publisher, result.Error
	}
	return publisher, nil
}

func (c PublisherModel) UpdatePublisher(publisher *Publisher) (*Publisher, error) {
	result := GetDB().Updates(publisher)
	if result.Error != nil {
		return publisher, result.Error
	}
	return publisher, nil
}

func (c PublisherModel) DeletePublisher(publisher *Publisher) (*Publisher, error) {
	GetDB().Delete(&publisher)
	return publisher, nil
}
