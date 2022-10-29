package service

import (
	"context"
	"encoding/json"
	"fmt"
	model2 "github.com/convee/adcreative/internal/model"
	"github.com/convee/adcreative/internal/pkg/cache"
	"github.com/convee/adcreative/internal/pkg/common"
	logger "github.com/convee/adcreative/pkg/log"
	probe2 "github.com/convee/adcreative/pkg/probe"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"path/filepath"
	"unicode/utf8"
)

type AdvertiserAudit struct {
	IsOnlyValid        int
	CheckRules         bool
	Id                 int
	AdvertiserName     string
	AdvertiserIds      []int
	Publisher          string
	PublisherId        int
	CustomerId         int
	PublisherAccountId int
	Status             int
	Info               map[string]interface{}
	IsRsync            int
	MediaCid           string
	ErrCode            int
	ErrMsg             string
	Extra              string
	Page               int // 第几页
	PerPage            int // 每页显示条数
	AdvertiserId       int
}

type AdvertiserAuditInfo struct {
	CompanyName      string                 `json:"company_name"`
	CompanySummary   string                 `json:"company_summary"`
	WebsiteName      string                 `json:"website_name"`
	WebsiteAddress   string                 `json:"website_address"`
	WebsiteNumber    string                 `json:"website_number"`
	BusinessLicenser string                 `json:"business_licenser"`
	AuthorizeState   string                 `json:"authorize_state"`
	Industry         string                 `json:"Industry"`
	Extra            map[string]interface{} `json:"Extra"`
	Qualifications   []model2.Qualification
}

type Qualification struct {
	FileName string
	FileUrl  string
}

var (
	advertiserAuditModel = model2.AdvertiserAuditModel{}
)

func (pa *AdvertiserAudit) GetList() map[string]interface{} {
	data := make(map[string]interface{})

	page := pa.Page
	if page == 0 {
		page = 1
	}
	perPage := pa.PerPage
	if perPage == 0 {
		perPage = 20
	}

	list, err := advertiserAuditModel.GetAdvertiserAudits(page, perPage, pa.getMaps())
	if err != nil {
		logger.Error("advertiserAudit get list data err ", zap.Error(err))
		return data
	}
	total, err := advertiserAuditModel.GetAdvertiserAuditTotal(pa.getMaps())
	if err != nil {
		logger.Error("advertiserAudit get list count err ", zap.Error(err))
		return data
	}
	data["lists"] = list
	data["total"] = total
	return data
}

func (aa *AdvertiserAudit) GetAll() []*model2.AdvertiserAudit {

	list, err := advertiserAuditModel.GetAllAdvertiserAuditsByMaps(aa.getMaps())
	if err != nil {
		logger.Error("advertiser_audit get all data err ", zap.Error(err))
		return []*model2.AdvertiserAudit{}
	}
	return list
}

func (aa *AdvertiserAudit) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	if aa.CustomerId > 0 {
		maps["customer_id"] = aa.CustomerId
	}
	if len(aa.AdvertiserName) > 0 {
		maps["advertiser_name"] = aa.AdvertiserName
	}
	if len(aa.Publisher) > 0 {
		maps["publisher_id"] = aa.PublisherId
	}
	if len(aa.AdvertiserIds) > 0 {
		maps["id"] = aa.AdvertiserIds
	}
	if aa.Status >= 0 && aa.IsRsync == 1 {
		maps["status"] = aa.Status
	}
	return maps
}

func (aa *AdvertiserAudit) CreateOrUpdate() (*model2.AdvertiserAudit, error) {
	publisherInfo, err := cache.GetPublisherCacheByName(aa.Publisher)
	if err != nil {
		return nil, errors.New("媒体不存在")
	}
	if *publisherInfo.IsRsyncAdvertiser == 0 {
		return nil, errors.New("该媒体方不支持送审")
	}
	aa.PublisherId = publisherInfo.Id
	exists, err := advertiserAuditModel.GetOneAdvertiserAuditByMaps(aa.getMaps())
	if err != nil {
		return nil, err
	}
	publisherAccount, err := cache.GetPublisherAccount(aa.CustomerId, publisherInfo.Id)
	if err != nil {
		return nil, errors.New("媒体账号不存在")
	}
	advertiserAuditService := AdvertiserRules{
		PublisherId: &publisherInfo.Id,
	}
	advertiserAuditInfo, err := advertiserAuditService.GetAdvertiserRulesByPublisherId()
	if err := aa.Check(advertiserAuditInfo); err != nil {
		return nil, errors.New(err[0])
	}
	if aa.IsOnlyValid == 1 {
		return exists, nil
	}

	info, _ := json.Marshal(aa.Info)
	advertiserAudit := &model2.AdvertiserAudit{
		AdvertiserName:     aa.AdvertiserName,
		PublisherId:        publisherInfo.Id,
		CustomerId:         aa.CustomerId,
		PublisherAccountId: publisherAccount.Id,
		Info:               string(info),
		Extra:              aa.Extra,
		IsRsync:            aa.IsRsync,
		ErrMsg:             aa.ErrMsg,
		ErrCode:            aa.ErrCode,
		AdvertiserId:       aa.AdvertiserId,
		Status:             aa.Status,
	}
	err = cache.NewUserCache().Del(context.Background(), cache.RedisKeyPublisherAccount+cast.ToString(advertiserAudit.CustomerId)+cast.ToString(advertiserAudit.AdvertiserId))
	if err != nil {
		logger.Error("advertiser_audit_del_cache_err", zap.Error(err))
		return nil, errors.New("数据更新失败，请重试")
	}
	if exists.Id > 0 {
		advertiserAudit.Id = exists.Id
		return advertiserAuditModel.UpdateAdvertiserAudit(advertiserAudit)
	} else {
		return advertiserAuditModel.CreateAdvertiserAudit(advertiserAudit)
	}
}

func (aa *AdvertiserAudit) UpdateAdvertiserAuditByMaps() model2.AdvertiserAudit {
	maps := make(map[string]interface{})
	maps["status"] = aa.Status
	maps["err_code"] = aa.ErrCode
	maps["err_msg"] = aa.ErrMsg
	if aa.IsRsync > 0 {
		maps["is_rsync"] = aa.IsRsync
	}
	if len(aa.MediaCid) > 0 {
		maps["media_cid"] = aa.MediaCid
	}
	creative, err := advertiserAuditModel.UpdateCreativeByMap(aa.Id, maps)
	if err != nil {
		logger.Error("广告主状态更新失败", zap.Error(err))
	}
	return creative
}

func (aa *AdvertiserAudit) GetAdvertiserAuditInfo() (*model2.AdvertiserAudit, error) {
	advertiserAudit, err := advertiserAuditModel.GetAdvertiserAuditById(aa.Id)
	if err != nil {
		return nil, err
	}
	return advertiserAudit, nil
}

func (aa *AdvertiserAudit) Check(advertiserRules *model2.AdvertiserRules) []string {
	var ruleTemplate []model2.RuleTemplate
	_ = json.Unmarshal([]byte(advertiserRules.Info), &ruleTemplate)
	var errs []string
	var attrNames []string
	var QualificationsLen int
	var FileUrl string
	for key, val := range aa.Info {
		if key == "qualifications" {
			if vv, ok := val.([]model2.Qualification); ok {
				QualificationsLen = len(vv)
				FileUrl = vv[0].FileUrl
			}
		}
		attrNames = append(attrNames, key)
	}
	for _, ai := range ruleTemplate {
		//1、必填参数校验
		//if ai.Required == 1 && !utils.StringsContain(ai.Key, attrNames...) {
		//	return errors.New(fmt.Sprintf("%s必填", ai.Desc))
		//}
		if info, ok := aa.Info[ai.Key]; ok {
			if ai.Measure == nil {
				continue
			}
			Type := common.GetFileType(ai.Type)
			switch Type {
			case "txt":
				if len(ai.Measure) == 2 {
					min := cast.ToInt(ai.Measure[0])
					max := cast.ToInt(ai.Measure[1])
					length := utf8.RuneCountInString(info.(string))
					//1、长度验证
					if length > max || length < min {
						errs = append(errs, fmt.Sprintf("%s长度有误", ai.Desc))
					}
				}
				break
			case "file":
				if aa.PublisherId == 2 {
					if QualificationsLen > ai.Limit {
						errs = append(errs, fmt.Sprintf("%s元素个数有误", ai.Desc))
					}
					if len(FileUrl) <= 0 {
						errs = append(errs, fmt.Sprintf("%s不能为空", ai.Desc))
						return errs
					}
					ext := filepath.Ext(FileUrl)
					if len(ext) < 2 {
						errs = append(errs, fmt.Sprintf("%s后缀有误", ai.Desc))
						continue
					}
					if !common.StringsContain(ext[1:], ai.Format...) {
						errs = append(errs, fmt.Sprintf("%s扩展有误", ai.Desc))
						continue
					}
					storePath := common.GetTmpFilePath()
					probe, err := probe2.ImgDownload(FileUrl, storePath)
					if err != nil {
						errs = append(errs, fmt.Sprintf("%s信息获取有误", ai.Desc))
						continue
					}
					fileSize, _ := probe.FilesizeByUnit("KB")
					if ai.Size > 0 && cast.ToInt(fileSize) > ai.Size {
						errs = append(errs, fmt.Sprintf("%s 大小不能超过%dKB", ai.Desc, ai.Size))
					}
				}
				break
			}
		}
	}
	return errs
}
