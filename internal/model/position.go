package model

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Position struct {
	Model
	Id                int    `json:"id"`
	PublisherId       int    `json:"publisher_id"`
	Name              string `json:"name"`
	Type              string `json:"type"`
	Position          string `json:"position"`
	MaterialInfo      string `json:"material_info"`
	IsSupportDeeplink *int   `json:"is_support_deeplink"`
	PvLimit           int    `json:"pv_limit"`
	ClLimit           int    `json:"cl_limit"`
}

type MaterialInfo struct {
	List []Template `json:"list"`
}

type Template struct {
	Id        string                 `json:"id"`
	Name      string                 `json:"name"`
	Extra     map[string]string      `json:"extra"`
	DisplayId int                    `json:"display_id"`
	Info      []PositionTemplateInfo `json:"info"`
}

type PositionTemplateInfo struct {
	Key      string            `json:"key"`
	Name     string            `json:"name"`
	Extra    map[string]string `json:"extra"`
	Measure  []string          `json:"measure"`
	Duration string            `json:"duration"`
	Ratio    []string          `json:"ratio"`
	Format   []string          `json:"format"`
	Range    []Range           `json:"range"`
	Size     int               `json:"size"`     // 单位k
	Required int               `json:"required"` // 1必填，0非必填
	Tips     string            `json:"tips"`
	MediaKey string            `json:"media_key"`
}

type Range struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
type PositionModel struct {
}

// TableName sets the insert table name for this struct type
func (c *PositionModel) TableName() string {
	return "position"
}

func (c PositionModel) GetPositionTotal(maps interface{}) (int64, error) {
	var count int64
	if err := GetDB().Model(&Position{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (c PositionModel) GetPositionById(id int) (*Position, error) {
	var position *Position
	err := GetDB().First(&position, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("记录不存在")
		}
		return nil, err
	}
	return position, nil
}

func (c PositionModel) ExistsById(id int) (bool, error) {
	var position *Position
	err := GetDB().First(&position, id).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if position.Id > 0 {
		return true, nil
	}
	return false, nil
}

func (c PositionModel) GetPosition(positionName string) (*Position, error) {
	var position *Position
	err := GetDB().Where("position=?", positionName).First(&position).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("记录不存在")
		}
		return nil, err
	}
	return position, nil
}

func (c PositionModel) GetPositionByMap(maps map[string]interface{}) (*Position, error) {
	var position *Position
	err := GetDB().Where(maps).First(&position).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("记录不存在")
		}
		return nil, err
	}
	return position, nil
}

func (c PositionModel) Exists(positionName string) (bool, error) {
	var position *Position
	err := GetDB().Where("position=?", positionName).First(&position).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if position.Id > 0 {
		return true, nil
	}
	return false, nil
}

func (c PositionModel) GetPositions(page int, perPage int, maps interface{}) ([]*Position, error) {
	var positions []*Position
	offset := (page - 1) * perPage
	err := GetDB().Where(maps).Offset(offset).Limit(perPage).Order("id desc").Find(&positions).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("记录不存在")
		}
		return nil, err
	}
	return positions, nil
}

func (c PositionModel) GetApiPositions(maps interface{}) ([]*Position, error) {
	var positions []*Position
	err := GetDB().Where(maps).Find(&positions).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("记录不存在")
		}
		return nil, err
	}
	return positions, nil
}

func (c PositionModel) GetAllPositions() ([]*Position, error) {
	var positions []*Position
	err := GetDB().Find(&positions).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("记录不存在")
		}
		return nil, err
	}
	return positions, nil
}

func (c PositionModel) CreatePosition(position *Position) (*Position, error) {
	result := GetDB().Create(&position)
	if result.Error != nil {
		return position, result.Error
	}
	return position, nil
}

func (c PositionModel) UpdatePosition(position *Position) (*Position, error) {
	result := GetDB().Updates(position)
	if result.Error != nil {
		return position, result.Error
	}
	return position, nil
}

func (c PositionModel) DeletePosition(position *Position) (*Position, error) {
	GetDB().Delete(&position)
	return position, nil
}
