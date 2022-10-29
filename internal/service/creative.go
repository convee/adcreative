package service

import (
	"fmt"
	model2 "github.com/convee/adcreative/internal/model"
	"github.com/convee/adcreative/internal/pkg/cache"
	"github.com/convee/adcreative/internal/pkg/common"
	"github.com/convee/adcreative/pkg/utils"
	"gorm.io/gorm"
	"math"
	"path/filepath"
	"strings"
	"sync"
	"time"

	logger "github.com/convee/adcreative/pkg/log"
	probe2 "github.com/convee/adcreative/pkg/probe"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

const (
	CHECK_TXT_ONLY  = 1
	CHECK_FILE_ONLY = 2
)

type Creative struct {
	Id                 int
	IsRsync            int
	CustomerId         int
	AdvertiserId       int
	MediaInfo          string
	Industry           string
	PublisherId        int
	PositionId         int
	PublisherAccountId int
	Position           string
	CreativeId         string
	MediaCid           string
	CreativeIds        []string
	Name               string
	TemplateId         string
	Action             int
	LandUrl            string
	DeeplinkUrl        string
	MiniProgramId      string
	MiniProgramPath    string
	StartDate          string
	EndDate            string
	Monitor            []model2.Monitor
	Info               []model2.TemplateInfo
	Cm                 []string
	Vm                 []string
	Status             int
	Reason             string
	ErrCode            int
	ErrCodes           []int
	ErrMsg             string
	Page               int // 第几页
	PerPage            int // 每页显示条数
	TemplateKey        string
	RequestId          string
	Extra              string
	VideoCdnUrl        string
	PicCdnUrl          string
	PubReturnUrl       string
}

var (
	creativeModel = model2.CreativeModel{}
)

func (c *Creative) GetList() map[string]interface{} {
	data := make(map[string]interface{})

	page := c.Page
	if page == 0 {
		page = 1
	}
	perPage := c.PerPage
	if perPage == 0 {
		perPage = 20
	}

	list, err := creativeModel.GetCreatives(page, perPage, c.getMaps())
	if err != nil {
		logger.Error("creative get list data err ", zap.Error(err))
		return data
	}
	total, err := creativeModel.GetCreativeTotal(c.getMaps())
	if err != nil {
		logger.Error("creative get list count err ", zap.Error(err))
		return data
	}
	data["lists"] = list
	data["total"] = total
	return data
}

func (c *Creative) GetAll() []*model2.Creative {
	maps := make(map[string]interface{})

	if c.CustomerId > 0 {
		maps["customer_id"] = c.CustomerId
	}
	if c.PublisherId > 0 {
		maps["publisher_id"] = c.PublisherId
	}
	if c.AdvertiserId > 0 {
		maps["advertiser_id"] = c.AdvertiserId
	}
	if len(c.CreativeId) > 0 {
		maps["creative_id"] = c.CreativeId
	}
	if c.Status >= 0 && c.IsRsync == 1 {
		maps["status"] = c.Status
	}
	list, err := creativeModel.GetAllCreativesByMaps(c.CreativeIds)
	if err != nil {
		logger.Error("creative_get_all_data_err ", zap.Error(err))
		return nil
	}
	return list
}

func (c *Creative) GetPassCreative() []*model2.Creative {
	maps := c.getMaps()
	date := time.Now().Add(-time.Hour * 240).Format("2006-01-02 15:04:05")

	list, err := creativeModel.GetPassCreatives(c.Page, c.PerPage, maps, date)
	if err != nil {
		logger.Error("creative_get_all_data_err ", zap.Error(err))
		return nil
	}
	return list
}

func (c *Creative) GetAllId() []*model2.Creative {
	list, err := creativeModel.GetAllCreativeIdsByMaps(c.getMaps())
	if err != nil {
		logger.Error("creative_get_all_data_err ", zap.Error(err))
		return nil
	}
	return list
}

func (c *Creative) GetAllByErrCodes(publisherId int, errCodes []int) []*model2.Creative {
	list, err := creativeModel.GetAllCreativeByErrCodes(publisherId, errCodes)
	if err != nil {
		logger.Error("creative_get_all_data_err ", zap.Error(err))
		return nil
	}
	return list
}

func (c *Creative) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})

	if c.CustomerId > 0 {
		maps["customer_id"] = c.CustomerId
	}
	if c.PublisherId > 0 {
		maps["publisher_id"] = c.PublisherId
	}
	if c.AdvertiserId > 0 {
		maps["advertiser_id"] = c.AdvertiserId
	}
	if len(c.CreativeId) > 0 {
		maps["creative_id"] = c.CreativeId
	}
	if len(c.CreativeIds) > 0 {
		maps["creative_id"] = c.CreativeIds
	}
	if c.Status >= 0 && c.IsRsync == 1 {
		maps["status"] = c.Status
	}
	if c.ErrCode > 0 {
		maps["err_code"] = c.ErrCode
	}
	return maps
}

func (c *Creative) getMaps2() map[string]interface{} {
	maps := make(map[string]interface{})

	if c.CustomerId > 0 {
		maps["customer_id"] = c.CustomerId
	}
	if c.PublisherId > 0 {
		maps["publisher_id"] = c.PublisherId
	}
	if c.AdvertiserId > 0 {
		maps["advertiser_id"] = c.AdvertiserId
	}
	if len(c.CreativeId) > 0 {
		maps["creative_id"] = c.CreativeId
	}
	return maps
}
func (c *Creative) Upload() (*model2.Creative, error) {

	exists, err := creativeModel.GetOneCreativeByMaps(c.getMaps2())
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	positionInfo, err := cache.GetPositionCache(c.Position)
	if err != nil {
		return nil, errors.New("广告位有误")
	}
	publisherInfo, err := cache.GetPublisherCacheById(positionInfo.PublisherId)
	if err != nil {
		return nil, errors.New("媒体不存在")
	}
	publisherAccount, err := cache.GetPublisherAccount(c.CustomerId, publisherInfo.Id)
	if err != nil {
		return nil, errors.New("媒体账号不存在")
	}

	if *positionInfo.IsSupportDeeplink == 0 && len(c.DeeplinkUrl) > 0 {
		return nil, errors.New("该广告位不支持 Deeplink")
	}

	if *publisherInfo.IsRsyncCreative == 0 {
		return nil, errors.New("该媒体不需要送审")
	}

	if *publisherInfo.IsRsyncAdvertiser == 1 {
		advertiserAudit, err := cache.GetAdvertiserAuditCache(c.CustomerId, c.AdvertiserId, publisherInfo.Id)
		if err != nil {
			logger.Error("广告主信息获取失败", zap.Any("custom", c.CustomerId), zap.Any("adv", c.AdvertiserId), zap.Any("publisher", publisherInfo), zap.Error(err))
			return nil, errors.New("广告主信息获取失败")
		}
		if advertiserAudit.Status != 1 {
			return nil, errors.New("广告主未审核通过")
		}

	}
	if common.StringsContain(publisherInfo.Name, "Tencent", "B612", "Weibo", "KuaiShou") && len(c.MediaInfo) <= 0 {
		return nil, errors.New("媒体方信息不能为空")
	}

	if len(c.EndDate) > 0 && c.EndDate < time.Now().Format("2006-01-02") {
		return nil, errors.New("素材已失效")
	}

	monitor, _ := jsoniter.Marshal(c.Monitor)
	info, _ := jsoniter.Marshal(c.Info)
	cm, _ := jsoniter.Marshal(c.Cm)
	vm, _ := jsoniter.Marshal(c.Vm)
	if exists.Id > 0 {
		exists.CustomerId = c.CustomerId
		exists.PublisherId = publisherInfo.Id
		exists.PublisherAccountId = publisherAccount.Id
		exists.AdvertiserId = c.AdvertiserId
		exists.MediaInfo = c.MediaInfo
		exists.Industry = c.Industry
		exists.PositionId = positionInfo.Id
		exists.CreativeId = c.CreativeId
		exists.Name = c.Name
		exists.TemplateId = c.TemplateId
		exists.Action = c.Action
		exists.LandUrl = c.LandUrl
		exists.DeeplinkUrl = c.DeeplinkUrl
		exists.MiniProgramId = c.MiniProgramId
		exists.MiniProgramPath = c.MiniProgramPath
		exists.StartDate = c.StartDate
		exists.EndDate = c.EndDate
		exists.Monitor = string(monitor)
		exists.Info = string(info)
		exists.Cm = string(cm)
		exists.Vm = string(vm)
		exists.Status = model2.CREATIVE_STATUS_UPLOADING
		exists.ErrMsg = ""
		exists.Reason = ""
		exists.ErrCode = model2.CREATIVE_UPLOADING
		exists.RequestId = c.RequestId
		if len(c.Extra) > 0 {
			exists.Extra = c.Extra
		}
		return creativeModel.UpdateCreative(exists)
	} else {
		creative := &model2.Creative{
			CustomerId:         c.CustomerId,
			PublisherId:        publisherInfo.Id,
			PublisherAccountId: publisherAccount.Id,
			AdvertiserId:       c.AdvertiserId,
			MediaInfo:          c.MediaInfo,
			Industry:           c.Industry,
			PositionId:         positionInfo.Id,
			CreativeId:         c.CreativeId,
			Name:               c.Name,
			TemplateId:         c.TemplateId,
			Action:             c.Action,
			LandUrl:            c.LandUrl,
			DeeplinkUrl:        c.DeeplinkUrl,
			MiniProgramId:      c.MiniProgramId,
			MiniProgramPath:    c.MiniProgramPath,
			StartDate:          c.StartDate,
			EndDate:            c.EndDate,
			Monitor:            string(monitor),
			Info:               string(info),
			Cm:                 string(cm),
			Vm:                 string(vm),
			RequestId:          c.RequestId,
			Status:             model2.CREATIVE_STATUS_UPLOADING,
			ErrCode:            model2.CREATIVE_UPLOADING,
		}
		if len(c.Extra) > 0 {
			creative.Extra = c.Extra
		}
		return creativeModel.CreateCreative(creative)
	}
}

func (c *Creative) GetTemplate(positionInfo *model2.Position) (model2.Template, error) {
	var (
		materialInfo *model2.MaterialInfo
		template     model2.Template
	)
	err := jsoniter.Unmarshal([]byte(positionInfo.MaterialInfo), &materialInfo)
	if err != nil {
		return template, err
	}
	for _, list := range materialInfo.List {
		if list.Id == c.TemplateId {
			template = list
		}
	}
	if template.Id == "" {
		return template, errors.New("模板不存在")
	}
	return template, nil
}

func (c *Creative) Check(positionInfo *model2.Position, checkType int) ([]*model2.TemplateInfo, []string) {
	errs := make([]string, 0)
	template, err := c.GetTemplate(positionInfo)
	if err != nil {
		errs = append(errs, err.Error())
		return nil, errs
	}
	var attrNames []string
	infoByKey := make(map[string]*model2.TemplateInfo)
	for _, cr := range c.Info {
		attrNames = append(attrNames, cr.AttrName)
		crea := cr
		infoByKey[cr.AttrName] = &crea
	}
	var templateKeys []string
	for _, ti := range template.Info {
		templateKeys = append(templateKeys, ti.Key)
	}
	if len(c.TemplateKey) > 0 && !common.StringsContain(c.TemplateKey, templateKeys...) {
		errs = append(errs, "素材规格不正确")
		return nil, errs
	}
	wg := sync.WaitGroup{}
	var m sync.Map
	var i sync.Map

	for _, ti := range template.Info {
		// 校验单个Key
		if len(c.TemplateKey) > 0 && ti.Key != c.TemplateKey {
			continue
		}
		wg.Add(1)
		go func(ti model2.PositionTemplateInfo, infoByKey map[string]*model2.TemplateInfo, checkType, publisherId int) {
			defer wg.Done()
			// 1、必填参数校验

			format := common.GetFileType(ti.Format[0])
			if checkType == CHECK_TXT_ONLY && c.isMaterial(ti.Key) {
				return
			}
			if checkType == CHECK_FILE_ONLY && !c.isMaterial(ti.Key) {
				return
			}

			if ti.Required == 1 && !common.StringsContain(ti.Key, attrNames...) {
				m.Store(ti.Key, fmt.Sprintf("%s必填", ti.Name))
				return
			}
			if info, ok := infoByKey[ti.Key]; ok {
				info.Extra = ti.Extra
				info.MediaKey = ti.MediaKey

				switch format {
				case "image":
					// 1、扩展名验证
					// 2、size 验证
					// 3、宽高验证
					// 4、md5 验证
					ext := filepath.Ext(info.AttrValue)
					if len(ext) < 2 {
						m.Store(ti.Key, fmt.Sprintf("%s后缀有误", ti.Name))
						return
					}
					actualExt := common.GetActualExt(ext)
					if !common.StringsContain(actualExt, ti.Format...) {
						m.Store(ti.Key, fmt.Sprintf("%s扩展有误,格式要求后缀为%s,实际后缀为%s", ti.Name, ti.Format, actualExt))
						return
					}
					path, err := common.DownloadFile(info.AttrValue)
					defer utils.RemoveFile(path)
					if err != nil {
						logger.Error("down_load_file_error", zap.Error(err))
						m.Store(ti.Key, fmt.Sprintf("%s信息下载有误", ti.Name))
						return
					}
					probe := probe2.New(path)

					fileSize, err := probe.FilesizeByUnit("KB")
					if err != nil {
						logger.Error("get_image_filesize_error", zap.Error(err))
						m.Store(ti.Key, fmt.Sprintf("%s大小获取失败", ti.Name))
						return
					}
					size, err := probe.GetSize()
					if err != nil {
						logger.Error("get_image_size_error", zap.Error(err))
						m.Store(ti.Key, fmt.Sprintf("%s宽高获取失败", ti.Name))
						return
					}
					width := size.X
					height := size.Y
					info.Width = width
					info.Height = height
					info.Ext = strings.ToLower(actualExt)

					md5, err := common.GetFileMd5(path)
					if err != nil {
						logger.Error("get_file_md5_error", zap.Error(err))
						m.Store(ti.Key, fmt.Sprintf("%s信息获取md5有误", ti.Name))
						return
					}
					info.Md5 = md5
					info.Size = cast.ToInt(math.Ceil(fileSize))
					measure := fmt.Sprintf("%d*%d", width, height)
					measures := strings.Join(ti.Measure, "")
					ratios := strings.Join(ti.Ratio, "")
					if len(ti.Measure) > 0 {
						if strings.Index(measures, ":") >= 0 {
							strArrayNew := strings.Split(measures, ":")
							if cast.ToFloat32(strArrayNew[0])/cast.ToFloat32(strArrayNew[1]) != cast.ToFloat32(width)/cast.ToFloat32(height) {
								m.Store(ti.Key, fmt.Sprintf("%s尺寸比例有误,模板尺寸为%s,实际尺寸为%s", ti.Name, ti.Measure, measure))
								return
							}
						} else if strings.Index(measures, ">") >= 0 {
							if strings.Index(measures, ">=") >= 0 {
								strArray := strings.Split(measures, ">=")
								strArrayNew := strings.Split(strArray[1], "*")
								if cast.ToInt(strArrayNew[0]) > width || cast.ToInt(strArrayNew[1]) > height {
									m.Store(ti.Key, fmt.Sprintf("%s尺寸范围有误,模板尺寸为%s,实际尺寸为%s", ti.Name, ti.Measure, measure))
									return
								}
							} else {
								strArray := strings.Split(measures, ">")
								strArrayNew := strings.Split(strArray[1], "*")
								if cast.ToInt(strArrayNew[0]) >= width || cast.ToInt(strArrayNew[1]) >= height {
									m.Store(ti.Key, fmt.Sprintf("%s尺寸范围有误,模板尺寸为%s,实际尺寸为%s", ti.Name, ti.Measure, measure))
									return
								}
							}
						} else if strings.Index(measures, "<") >= 0 {
							if strings.Index(measures, "<=") >= 0 {
								strArray := strings.Split(measures, "<=")
								strArrayNew := strings.Split(strArray[1], "*")
								if cast.ToInt(strArrayNew[0]) < width || cast.ToInt(strArrayNew[1]) < height {
									m.Store(ti.Key, fmt.Sprintf("%s尺寸范围有误,模板尺寸为%s,实际尺寸为%s", ti.Name, ti.Measure, measure))
									return
								}
							} else {
								strArray := strings.Split(measures, "<")
								strArrayNew := strings.Split(strArray[1], "*")
								if cast.ToInt(strArrayNew[0]) <= width || cast.ToInt(strArrayNew[1]) <= height {
									m.Store(ti.Key, fmt.Sprintf("%s尺寸范围有误,模板尺寸为%s,实际尺寸为%s", ti.Name, ti.Measure, measure))
									return
								}
							}
						} else if !common.StringsContain(measure, ti.Measure...) {
							m.Store(ti.Key, fmt.Sprintf("%s尺寸有误,模板尺寸为%s,实际尺寸为%s", ti.Name, ti.Measure, measure))
							return
						}
					}

					// Ratio 仅支持配置比例，与measure尺寸取交集
					if len(ti.Ratio) > 0 {
						if strings.Index(ratios, ":") >= 0 {
							strArrayNew := strings.Split(ratios, ":")
							if cast.ToFloat32(strArrayNew[0])/cast.ToFloat32(strArrayNew[1]) != cast.ToFloat32(width)/cast.ToFloat32(height) {
								m.Store(ti.Key, fmt.Sprintf("%s尺寸比例有误,模板尺寸为%s,实际尺寸为%s", ti.Name, ti.Ratio, measure))
								return
							}
						}
					}
					if ti.Size > 0 && cast.ToInt(math.Ceil(fileSize)) > ti.Size {
						m.Store(ti.Key, fmt.Sprintf("%s 大小不能超过%dKB", ti.Name, ti.Size))
						return
					}
					i.Store(ti.Key, info)
					break
				case "video":
					// 1、扩展名验证
					// 2、size 验证
					// 3、宽高验证（或视频比例验证）
					// 4、md5 验证
					// 5、时长验证
					ext := filepath.Ext(info.AttrValue)
					if len(ext) < 2 {
						m.Store(ti.Key, fmt.Sprintf("%s后缀有误", ti.Name))
						return
					}

					actualExt := common.GetActualExt(ext)
					if !common.StringsContain(actualExt, ti.Format...) {
						m.Store(ti.Key, fmt.Sprintf("%s扩展有误,格式要求后缀为%s,实际后缀为%s", ti.Name, ti.Format, actualExt))
						return
					}
					path, err := common.DownloadFile(info.AttrValue)
					defer utils.RemoveFile(path)
					if err != nil {
						logger.Error("get_video_path_error", zap.Error(err))
						m.Store(ti.Key, fmt.Sprintf("%s信息获取有误", ti.Name))
						return
					}
					videoInfo, err := common.GetVideoInfo(path)
					if err != nil {
						logger.Error("get_video_error", zap.Error(err))
						m.Store(ti.Key, fmt.Sprintf("%s信息获取有误", ti.Name))
						return
					}
					width := videoInfo.FirstVideoStream().Width
					height := videoInfo.FirstVideoStream().Height
					duration := videoInfo.Format.Duration
					//durationStr := strings.Split(cast.ToString(duration()), ".")
					if find := strings.Contains(cast.ToString(duration()), "m"); find {
						durationM := strings.Split(cast.ToString(duration()), "m")
						durationInt := cast.ToInt(durationM[0]) * 60
						info.Duration = durationInt
					} else {
						durationStr := strings.Replace(cast.ToString(duration()), "s", "", -1)
						durationInt := math.Floor(cast.ToFloat64(durationStr))
						info.Duration = cast.ToInt(durationInt)
					}
					info.Width = cast.ToInt(width)
					info.Height = cast.ToInt(height)
					info.Ext = strings.ToLower(actualExt)
					md5, err := common.GetFileMd5(path)
					if err != nil {
						logger.Error("get_file_md5_error", zap.Error(err))
						m.Store(ti.Key, fmt.Sprintf("%s信息获取md5有误", ti.Name))
						return
					}
					info.Md5 = md5
					info.Size = cast.ToInt(math.Ceil(cast.ToFloat64(videoInfo.Format.Size) / 1024))
					measure := fmt.Sprintf("%d*%d", width, height)
					measures := strings.Join(ti.Measure, "")
					ratios := strings.Join(ti.Ratio, "")
					if len(ti.Measure) > 0 {
						if strings.Index(measures, ":") >= 0 {
							strArrayNew := strings.Split(measures, ":")
							if cast.ToFloat32(strArrayNew[0])/cast.ToFloat32(strArrayNew[1]) != cast.ToFloat32(width)/cast.ToFloat32(height) {
								m.Store(ti.Key, fmt.Sprintf("%s尺寸比例有误,模板尺寸为%s,实际尺寸为%s", ti.Name, ti.Measure, measure))
								return
							}
						} else if strings.Index(measures, ">") >= 0 {
							if strings.Index(measures, ">=") >= 0 {
								strArray := strings.Split(measures, ">=")
								strArrayNew := strings.Split(strArray[1], "*")
								if cast.ToInt(strArrayNew[0]) > width || cast.ToInt(strArrayNew[1]) > height {
									m.Store(ti.Key, fmt.Sprintf("%s尺寸范围有误,模板尺寸为%s,实际尺寸为%s", ti.Name, ti.Measure, measure))
									return
								}
							} else {
								strArray := strings.Split(measures, ">")
								strArrayNew := strings.Split(strArray[1], "*")
								if cast.ToInt(strArrayNew[0]) >= width || cast.ToInt(strArrayNew[1]) >= height {
									m.Store(ti.Key, fmt.Sprintf("%s尺寸范围有误,模板尺寸为%s,实际尺寸为%s", ti.Name, ti.Measure, measure))
									return
								}
							}
						} else if strings.Index(measures, "<") >= 0 {
							if strings.Index(measures, "<=") >= 0 {
								strArray := strings.Split(measures, "<=")
								strArrayNew := strings.Split(strArray[1], "*")
								if cast.ToInt(strArrayNew[0]) < width || cast.ToInt(strArrayNew[1]) < height {
									m.Store(ti.Key, fmt.Sprintf("%s尺寸范围有误,模板尺寸为%s,实际尺寸为%s", ti.Name, ti.Measure, measure))
									return
								}
							} else {
								strArray := strings.Split(measures, "<")
								strArrayNew := strings.Split(strArray[1], "*")
								if cast.ToInt(strArrayNew[0]) <= width || cast.ToInt(strArrayNew[1]) <= height {
									m.Store(ti.Key, fmt.Sprintf("%s尺寸范围有误,模板尺寸为%s,实际尺寸为%s", ti.Name, ti.Measure, measure))
									return
								}
							}
						} else if !common.StringsContain(measure, ti.Measure...) {
							m.Store(ti.Key, fmt.Sprintf("%s尺寸有误,模板尺寸为%s,实际尺寸为%s", ti.Name, ti.Measure, measure))
							return
						}
					}

					// Ratio 仅支持配置比例，与measure尺寸取交集
					if len(ti.Ratio) > 0 {
						if strings.Index(ratios, ":") >= 0 {
							strArrayNew := strings.Split(ratios, ":")
							if cast.ToFloat32(strArrayNew[0])/cast.ToFloat32(strArrayNew[1]) != cast.ToFloat32(width)/cast.ToFloat32(height) {
								m.Store(ti.Key, fmt.Sprintf("%s尺寸比例有误,模板尺寸为%s,实际尺寸为%s", ti.Name, ti.Ratio, measure))
								return
							}
						}
					}
					if len(ti.Duration) > 0 && info.Duration != cast.ToInt(ti.Duration) {
						m.Store(ti.Key, fmt.Sprintf("%s时长有误,模板时长为%s,实际时长为%d", ti.Name, ti.Duration, info.Duration))
						return
					}
					if ti.Size > 0 && cast.ToInt(videoInfo.Format.Size) > ti.Size*1024 {
						m.Store(ti.Key, fmt.Sprintf("%s大小不能超过%dKB", ti.Name, ti.Size))
						return
					}
					i.Store(ti.Key, info)
					break
				case "txt":
					if len(ti.Measure) == 2 {
						min := cast.ToFloat32(ti.Measure[0])
						max := cast.ToFloat32(ti.Measure[1])
						length := common.StringLength(info.AttrValue, publisherId)
						// 1、长度验证
						if length > max || length < min {
							maxI := cast.ToInt(max)
							minI := cast.ToInt(min)

							m.Store(ti.Key, fmt.Sprintf("%s长度有误,模板长度为%d-%d,实际长度为%v", ti.Name, minI, maxI, length))
							return
						}
					}
					i.Store(ti.Key, info)
					break
				case "select":
					// 1、枚举值验证
					var rangeKeys []string
					for _, r := range ti.Range {
						rangeKeys = append(rangeKeys, r.Key)
					}
					if !common.StringsContain(info.AttrValue, rangeKeys...) {
						m.Store(ti.Key, fmt.Sprintf("%s值有误", ti.Name))
						return
					}
					i.Store(ti.Key, info)
					break
				case "url":
					path, err := common.DownloadFile(info.AttrValue)
					defer utils.RemoveFile(path)
					if err != nil {
						logger.Error("down_load_url_error", zap.Error(err))
						m.Store(ti.Key, fmt.Sprintf("%s信息下载有误", ti.Name))
						return
					}
					md5, err := common.GetFileMd5(path)
					if err != nil {
						logger.Error("get_url_md5_error", zap.Error(err))
						m.Store(ti.Key, fmt.Sprintf("%s信息获取md5有误", ti.Name))
						return
					}
					info.Md5 = md5
					ext := filepath.Ext(info.AttrValue)
					if len(ext) < 2 {
						m.Store(ti.Key, fmt.Sprintf("%s后缀有误", ti.Name))
						logger.Error("can not get_file_ext")
						return
					}
					actualExt := common.GetActualExt(ext)
					info.Ext = strings.ToLower(actualExt)
					extType := common.GetFileType(info.Ext)

					if extType != "image" && extType != "video" {
						if strings.Index(ti.Key, "video") > -1 || strings.Index(ti.Key, "image") > -1 || strings.Index(ti.Key, "cover") > -1 {
							m.Store(ti.Key, fmt.Sprintf("%s格式有误", ti.Name))
							return
						}
					}
					if extType == "image" {
						probe := probe2.New(path)
						fileSize, err := probe.FilesizeByUnit("KB")
						if err != nil {
							logger.Error("get_image_filesize_error", zap.Error(err))
							m.Store(ti.Key, fmt.Sprintf("%s大小获取失败", ti.Name))
							return
						}
						size, err := probe.GetSize()
						if err != nil {
							logger.Error("get_image_size_error", zap.Error(err))
							m.Store(ti.Key, fmt.Sprintf("%s宽高获取失败", ti.Name))
							return
						}
						width := size.X
						height := size.Y
						info.Width = width
						info.Height = height
						info.Size = cast.ToInt(math.Ceil(fileSize))
						measure := fmt.Sprintf("%d*%d", width, height)
						measures := strings.Join(ti.Measure, "")
						ratios := strings.Join(ti.Ratio, "")
						if len(ti.Measure) > 0 {
							if strings.Index(measures, ":") >= 0 {
								strArrayNew := strings.Split(measures, ":")
								if cast.ToFloat32(strArrayNew[0])/cast.ToFloat32(strArrayNew[1]) != cast.ToFloat32(width)/cast.ToFloat32(height) {
									m.Store(ti.Key, fmt.Sprintf("%s尺寸比例有误,模板尺寸为%s,实际尺寸为%s", ti.Name, ti.Measure, measure))
									return
								}
							} else if strings.Index(measures, ">") >= 0 {
								if strings.Index(measures, ">=") >= 0 {
									strArray := strings.Split(measures, ">=")
									strArrayNew := strings.Split(strArray[1], "*")
									if cast.ToInt(strArrayNew[0]) > width || cast.ToInt(strArrayNew[1]) > height {
										m.Store(ti.Key, fmt.Sprintf("%s尺寸范围有误,模板尺寸为%s,实际尺寸为%s", ti.Name, ti.Measure, measure))
										return
									}
								} else {
									strArray := strings.Split(measures, ">")
									strArrayNew := strings.Split(strArray[1], "*")
									if cast.ToInt(strArrayNew[0]) >= width || cast.ToInt(strArrayNew[1]) >= height {
										m.Store(ti.Key, fmt.Sprintf("%s尺寸范围有误,模板尺寸为%s,实际尺寸为%s", ti.Name, ti.Measure, measure))
										return
									}
								}
							} else if strings.Index(measures, "<") >= 0 {
								if strings.Index(measures, "<=") >= 0 {
									strArray := strings.Split(measures, "<=")
									strArrayNew := strings.Split(strArray[1], "*")
									if cast.ToInt(strArrayNew[0]) < width || cast.ToInt(strArrayNew[1]) < height {
										m.Store(ti.Key, fmt.Sprintf("%s尺寸范围有误,模板尺寸为%s,实际尺寸为%s", ti.Name, ti.Measure, measure))
										return
									}
								} else {
									strArray := strings.Split(measures, "<")
									strArrayNew := strings.Split(strArray[1], "*")
									if cast.ToInt(strArrayNew[0]) <= width || cast.ToInt(strArrayNew[1]) <= height {
										m.Store(ti.Key, fmt.Sprintf("%s尺寸范围有误,模板尺寸为%s,实际尺寸为%s", ti.Name, ti.Measure, measure))
										return
									}
								}
							} else if !common.StringsContain(measure, ti.Measure...) {
								m.Store(ti.Key, fmt.Sprintf("%s尺寸有误,模板尺寸为%s,实际尺寸为%s", ti.Name, ti.Measure, measure))
								return
							}
						}
						// Ratio 仅支持配置比例，与measure尺寸取交集
						if len(ti.Ratio) > 0 {
							if strings.Index(ratios, ":") >= 0 {
								strArrayNew := strings.Split(ratios, ":")
								if cast.ToFloat32(strArrayNew[0])/cast.ToFloat32(strArrayNew[1]) != cast.ToFloat32(width)/cast.ToFloat32(height) {
									m.Store(ti.Key, fmt.Sprintf("%s尺寸比例有误,模板尺寸为%s,实际尺寸为%s", ti.Name, ti.Ratio, measure))
									return
								}
							}
						}
						if urlFormat, ok := ti.Extra["format"]; ok && urlFormat != "" {
							urlFormatArr := strings.Split(urlFormat, ",")

							if !common.StringsContain(info.Ext, urlFormatArr...) {
								m.Store(ti.Key, fmt.Sprintf("%s扩展有误,格式要求后缀为%s,实际后缀为%s", ti.Name, urlFormatArr, info.Ext))
								return
							}
						}
						if ti.Size > 0 && cast.ToInt(math.Ceil(fileSize)) > ti.Size {
							m.Store(ti.Key, fmt.Sprintf("%s 大小不能超过%dKB", ti.Name, ti.Size))
							return
						}
					}
					if extType == "video" {
						videoInfo, err := common.GetVideoInfo(info.AttrValue)
						if err != nil {
							logger.Error("get_video_error", zap.Error(err))
							m.Store(ti.Key, fmt.Sprintf("%s信息获取有误", ti.Name))
							return
						}
						width := videoInfo.FirstVideoStream().Width
						height := videoInfo.FirstVideoStream().Height
						measure := fmt.Sprintf("%d*%d", width, height)
						measures := strings.Join(ti.Measure, "")
						ratios := strings.Join(ti.Ratio, "")
						if len(ti.Measure) > 0 {
							if strings.Index(measures, ":") >= 0 {
								strArrayNew := strings.Split(measures, ":")
								if cast.ToFloat32(strArrayNew[0])/cast.ToFloat32(strArrayNew[1]) != cast.ToFloat32(width)/cast.ToFloat32(height) {
									m.Store(ti.Key, fmt.Sprintf("%s尺寸比例有误,模板尺寸为%s,实际尺寸为%s", ti.Name, ti.Measure, measure))
									return
								}
							} else if strings.Index(measures, ">") >= 0 {
								if strings.Index(measures, ">=") >= 0 {
									strArray := strings.Split(measures, ">=")
									strArrayNew := strings.Split(strArray[1], "*")
									if cast.ToInt(strArrayNew[0]) > width || cast.ToInt(strArrayNew[1]) > height {
										m.Store(ti.Key, fmt.Sprintf("%s尺寸范围有误,模板尺寸为%s,实际尺寸为%s", ti.Name, ti.Measure, measure))
										return
									}
								} else {
									strArray := strings.Split(measures, ">")
									strArrayNew := strings.Split(strArray[1], "*")
									if cast.ToInt(strArrayNew[0]) >= width || cast.ToInt(strArrayNew[1]) >= height {
										m.Store(ti.Key, fmt.Sprintf("%s尺寸范围有误,模板尺寸为%s,实际尺寸为%s", ti.Name, ti.Measure, measure))
										return
									}
								}
							} else if strings.Index(measures, "<") >= 0 {
								if strings.Index(measures, "<=") >= 0 {
									strArray := strings.Split(measures, "<=")
									strArrayNew := strings.Split(strArray[1], "*")
									if cast.ToInt(strArrayNew[0]) < width || cast.ToInt(strArrayNew[1]) < height {
										m.Store(ti.Key, fmt.Sprintf("%s尺寸范围有误,模板尺寸为%s,实际尺寸为%s", ti.Name, ti.Measure, measure))
										return
									}
								} else {
									strArray := strings.Split(measures, "<")
									strArrayNew := strings.Split(strArray[1], "*")
									if cast.ToInt(strArrayNew[0]) <= width || cast.ToInt(strArrayNew[1]) <= height {
										m.Store(ti.Key, fmt.Sprintf("%s尺寸范围有误,模板尺寸为%s,实际尺寸为%s", ti.Name, ti.Measure, measure))
										return
									}
								}
							} else if !common.StringsContain(measure, ti.Measure...) {
								m.Store(ti.Key, fmt.Sprintf("%s尺寸有误,模板尺寸为%s,实际尺寸为%s", ti.Name, ti.Measure, measure))
								return
							}
						}
						// Ratio 仅支持配置比例，与measure尺寸取交集
						if len(ti.Ratio) > 0 {
							if strings.Index(ratios, ":") >= 0 {
								strArrayNew := strings.Split(ratios, ":")
								if cast.ToFloat32(strArrayNew[0])/cast.ToFloat32(strArrayNew[1]) != cast.ToFloat32(width)/cast.ToFloat32(height) {
									m.Store(ti.Key, fmt.Sprintf("%s尺寸比例有误,模板尺寸为%s,实际尺寸为%s", ti.Name, ti.Ratio, measure))
									return
								}
							}
						}
						info.Width = cast.ToInt(width)
						info.Height = cast.ToInt(height)
						duration := videoInfo.Format.Duration
						if find := strings.Contains(cast.ToString(duration()), "m"); find {
							durationM := strings.Split(cast.ToString(duration()), "m")
							durationInt := cast.ToInt(durationM[0]) * 60
							info.Duration = durationInt
						} else {
							durationStr := strings.Replace(cast.ToString(duration()), "s", "", -1)
							durationInt := math.Floor(cast.ToFloat64(durationStr))
							info.Duration = cast.ToInt(durationInt)
						}
						if len(ti.Duration) > 0 && info.Duration != cast.ToInt(ti.Duration) {
							m.Store(ti.Key, fmt.Sprintf("%s时长有误,模板时长为%s,实际时长为%d", ti.Name, ti.Duration, info.Duration))
							return
						}
						if urlFormat, ok := ti.Extra["format"]; ok && urlFormat != "" {
							urlFormatArr := strings.Split(urlFormat, ",")

							if !common.StringsContain(info.Ext, urlFormatArr...) {
								m.Store(ti.Key, fmt.Sprintf("%s扩展有误,格式要求后缀为%s,实际后缀为%s", ti.Name, urlFormatArr, info.Ext))
								return
							}
						}
						if ti.Size > 0 && cast.ToInt(videoInfo.Format.Size) > ti.Size*1024 {
							m.Store(ti.Key, fmt.Sprintf("%s大小不能超过%dKB", ti.Name, ti.Size))
							return
						}
					}
					i.Store(ti.Key, info)
					break
					// 1、url 验证
					// 暂不验证
				}
			}
		}(ti, infoByKey, checkType, positionInfo.PublisherId)
	}

	wg.Wait()

	for _, ti := range template.Info {
		if err, ok := m.Load(ti.Key); ok {
			errs = append(errs, err.(string))
		}
	}

	var adContent []*model2.TemplateInfo
	for _, val := range infoByKey {
		if v, ok := i.Load(val.AttrName); ok {
			adContent = append(adContent, v.(*model2.TemplateInfo))
		} else {
			adContent = append(adContent, val)
		}
	}
	return adContent, errs
}
func (c *Creative) UpdateCreativeByMaps() error {
	maps := make(map[string]interface{})
	maps["status"] = c.Status
	maps["err_code"] = c.ErrCode
	maps["err_msg"] = c.ErrMsg
	if c.IsRsync > 0 {
		maps["is_rsync"] = c.IsRsync
	}
	if len(c.Reason) > 0 {
		maps["reason"] = c.Reason
	}
	if len(c.MediaCid) > 0 {
		maps["media_cid"] = c.MediaCid
	}
	if len(c.Extra) > 0 {
		maps["extra"] = c.Extra
	}
	if c.VideoCdnUrl != "" {
		maps["video_cdn_url"] = c.VideoCdnUrl
	}
	if c.PicCdnUrl != "" {
		maps["pic_cdn_url"] = c.PicCdnUrl
	}
	if c.PubReturnUrl != "" {
		maps["pub_return_url"] = c.PubReturnUrl
	}

	_, err := creativeModel.UpdateCreativeByMap(c.Id, maps)
	if err != nil {
		return err
	}
	return nil
}

func (c *Creative) GetCreativeById() (*model2.Creative, error) {
	creative, err := creativeModel.GetCreativeById(c.Id)
	if err != nil {
		return nil, err
	}
	return creative, nil
}

func (c *Creative) isMaterial(key string) bool {
	if strings.Index(key, "image") >= 0 || strings.Index(key, "video") >= 0 || strings.Index(key, "icon") >= 0 {
		return true
	}
	return false
}
