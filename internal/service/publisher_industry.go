package service

import (
	"github.com/convee/adcreative/internal/model"
	logger "github.com/convee/adcreative/pkg/log"
	"go.uber.org/zap"
	"strconv"
)

type PublisherIndustry struct {
	Id        int
	Name      string
	Pid       int
	Level     int
	Sort      int
	Publisher string
}

var (
	publisherIndustryModel = model.PublisherIndustryModel{}
)

func (pi *PublisherIndustry) GetList() map[string]interface{} {
	data := make(map[string]interface{})
	list, err := publisherIndustryModel.GetPublisherIndustrys(pi.getMaps())
	lists := make([]map[string]string, 0)
	for _, val := range list {
		info := make(map[string]string)
		info["id"] = strconv.Itoa(val.TypeId)
		info["value"] = val.Name
		lists = append(lists, info)
	}
	if err != nil {
		logger.Error("publisherIndustry get list data err ", zap.Error(err))
		return data
	}
	data["lists"] = lists
	return data
}

func GetTree(list []*model.PublisherIndustry, pid, level, maxLevel int) []interface{} {
	var results []interface{}
	var values model.TreeIndustry
	for _, value := range list {
		if value.Pid == pid {
			values.Id = value.Id
			values.Name = value.Name
			values.Pid = value.Pid
			values.Level = value.Level
			child := GetTree(list, value.Id, level+1, maxLevel)
			if len(child) > 0 {
				values.Leaf = false
				values.Child = child
			} else {
				values.Leaf = true
			}
			results = append(results, values)
		}
	}
	return results
}

func (pi *PublisherIndustry) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	if len(pi.Publisher) > 0 {
		maps["publisher"] = pi.Publisher
	}
	return maps
}
