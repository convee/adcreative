package model

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type PublisherIndustry struct {
	Model
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Pid    int    `json:"pid"`
	Level  int    `json:"level"`
	Sort   int    `json:"sort"`
	TypeId int    `json:"type_id"`
}

type TreeIndustry struct {
	Id    int         `json:"id"`
	Name  string      `json:"name"`
	Pid   int         `json:"pid"`
	Level int         `json:"level"`
	Leaf  bool        `json:"leaf"`
	Child interface{} `json:"child"`
}

type PublisherIndustryModel struct {
}

// TableName sets the insert table name for this struct type
func (c *PublisherIndustryModel) TableName() string {
	return "publisher_industry"
}

func (c PublisherIndustryModel) GetPublisherIndustrys(maps interface{}) ([]*PublisherIndustry, error) {
	var publisherIndustrys []*PublisherIndustry
	err := GetDB().Where(maps).Order("sort asc").Order("id asc").Find(&publisherIndustrys).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("记录不存在")
		}
		return nil, err
	}
	return publisherIndustrys, nil
}
