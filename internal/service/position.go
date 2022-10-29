package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/convee/adcreative/internal/model"
	"github.com/convee/adcreative/internal/pkg/cache"
	logger "github.com/convee/adcreative/pkg/log"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

type Position struct {
	Id                             int
	PublisherId                    int
	Publisher                      string
	Name                           string
	Type                           string
	Position                       string
	MaterialInfo                   string
	MediaType                      string
	AdFormat                       *int
	IsSupportDeeplink              *int
	LandingChangeNeedRsync         *int
	MonitorCodeChangeNeedRsync     *int
	MonitorPositionChangeNeedRsync *int
	IsCreativeBind                 *int
	PvLimit                        int
	ClLimit                        int
	Page                           int // 第几页
	PerPage                        int // 每页显示条数

}

var (
	positionModel = model.PositionModel{}
)

func (c *Position) GetList() map[string]interface{} {
	data := make(map[string]interface{})

	page := c.Page
	if page == 0 {
		page = 1
	}
	perPage := c.PerPage
	if perPage == 0 {
		perPage = 20
	}

	list, err := positionModel.GetPositions(page, perPage, c.getMaps())
	if err != nil {
		logger.Error("position get list data err ", zap.Error(err))
		return data
	}
	total, err := positionModel.GetPositionTotal(c.getMaps())
	if err != nil {
		logger.Error("position get list count err ", zap.Error(err))
		return data
	}
	data["lists"] = list
	data["total"] = total
	return data
}

func (p *Position) GetApiList() map[string]interface{} {
	data := make(map[string]interface{})

	position, err := cache.GetPositionCache(p.Position)
	if err != nil {
		logger.Error("get_position_cache_err", zap.Error(err))
		return data
	}
	lists := make([]map[string]interface{}, 0)
	list := make(map[string]interface{})
	var materialInfo *model.MaterialInfo
	_ = json.Unmarshal([]byte(position.MaterialInfo), &materialInfo)
	list["position"] = position.Position
	list["is_support_deeplink"] = position.IsSupportDeeplink
	templates := make([]map[string]interface{}, 0)
	for _, value := range materialInfo.List {
		template := make(map[string]interface{})
		info := make([]map[string]interface{}, 0)
		for _, attr := range value.Info {
			attrs := make(map[string]interface{})
			attrs["attr_name"] = attr.Key
			attrs["attr_desc"] = attr.Name
			attrs["extra"] = attr.Extra
			attrs["format"] = attr.Format
			attrs["measure"] = attr.Measure
			attrs["ratio"] = attr.Ratio
			attrs["range"] = attr.Range
			attrs["size"] = attr.Size
			attrs["duration"] = attr.Duration
			attrs["required"] = attr.Required
			attrs["tips"] = attr.Tips
			info = append(info, attrs)
		}
		template["template_id"] = value.Id
		template["template_name"] = value.Name
		template["info"] = info
		template["extra"] = value.Extra
		templates = append(templates, template)
	}
	list["templates"] = templates
	lists = append(lists, list)
	data["lists"] = lists
	return data
}

func (p *Position) GetApiMaterial() map[string]string {
	data := make(map[string]string)

	position, _ := cache.GetPositionCache(p.Position)
	data[position.Position] = position.MaterialInfo
	return data
}
func (p *Position) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	if p.PublisherId > 0 {
		maps["publisher_id"] = p.PublisherId
	}
	if len(p.Position) > 0 {
		maps["position"] = p.Position
	}
	if len(p.Name) > 0 {
		maps["name"] = p.Name
	}
	if p.Id > 0 {
		maps["id"] = p.Id
	}
	return maps
}

func (c *Position) Add() (*model.Position, error) {
	return positionModel.CreatePosition(&model.Position{
		PublisherId:       c.PublisherId,
		Name:              c.Name,
		Type:              c.Type,
		Position:          c.Position,
		MaterialInfo:      c.MaterialInfo,
		IsSupportDeeplink: c.IsSupportDeeplink,
		PvLimit:           c.PvLimit,
		ClLimit:           c.ClLimit,
	})
}

func (c *Position) Edit() (*model.Position, error) {
	position, err := positionModel.GetPositionById(c.Id)
	if err != nil {
		return position, err
	}
	position.PublisherId = c.PublisherId
	position.Name = c.Name
	position.Type = c.Type
	position.Position = c.Position
	position.MaterialInfo = c.MaterialInfo
	position.IsSupportDeeplink = c.IsSupportDeeplink
	position.PvLimit = c.PvLimit
	position.ClLimit = c.ClLimit
	err1 := cache.NewUserCache().Del(context.Background(), cache.RedisKeyPosition+position.Position)
	err2 := cache.NewUserCache().Del(context.Background(), cache.RedisKeyPosition+cast.ToString(position.Id))
	if err1 != nil || err2 != nil {
		return nil, errors.New("删除缓存失败")
	}
	return positionModel.UpdatePosition(position)
}

func (c *Position) Delete() (*model.Position, error) {
	position, err := positionModel.GetPositionById(c.Id)
	if err != nil {
		return position, err
	}
	err1 := cache.NewUserCache().Del(context.Background(), cache.RedisKeyPosition+position.Position)
	err2 := cache.NewUserCache().Del(context.Background(), cache.RedisKeyPosition+cast.ToString(position.Id))
	if err1 != nil || err2 != nil {
		return nil, errors.New("删除缓存失败")
	}
	return positionModel.DeletePosition(position)
}

func (c *Position) GetPositionInfo() *model.Position {
	position, err := positionModel.GetPositionById(c.Id)
	if err != nil {
		logger.Error("position get info err ", zap.Error(err))
		return position
	}
	return position
}

func (c *Position) GetPosition() (*model.Position, error) {
	position, err := positionModel.GetPosition(c.Position)
	if err != nil {
		return nil, err
	}
	return position, nil
}

func (c *Position) GetPositionByMap() (*model.Position, error) {
	maps := make(map[string]interface{})
	if c.PublisherId > 0 {
		maps["publisher_id"] = c.PublisherId
	}
	if len(c.Position) > 0 {
		maps["position"] = c.Position
	}
	position, err := positionModel.GetPositionByMap(maps)
	if err != nil {
		return nil, err
	}
	return position, nil
}
