package media

import (
	"github.com/convee/adcreative/internal/enum"
	model2 "github.com/convee/adcreative/internal/model"
	"github.com/convee/adcreative/internal/service"
	"github.com/json-iterator/go/extra"
	"strings"

	"github.com/convee/adcreative/internal/pkg/cache"
	"github.com/convee/adcreative/internal/pkg/common"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type MediaHandler interface {
	UploadAdvertiser() Ret
	QueryAdvertiser() Ret
	UploadCreative() Ret
	QueryCreative() Ret
	BatchQueryCreative() Ret
	BatchUploadCreative() Ret
}

type Base struct {
	CreativeId      int
	AdvertiserAudit *model2.AdvertiserAudit
	CreativeInfo    CreativeInfo
	PositionInfo    *model2.Position
	Template        model2.Template
	OriginUrl       string
}

type BaseInfo struct {
	Logger           *zap.Logger
	PublisherId      int
	CustomerId       int
	CreativeId       int
	CustomerInfo     *model2.Customer
	AdvertiserUrls   model2.AdvertiserUrls
	CreativeUrls     model2.CreativeUrls
	AdvertiserAudit  *model2.AdvertiserAudit
	PublisherInfo    *model2.Publisher
	CreativeInfo     CreativeInfo
	PositionInfo     *model2.Position
	Template         model2.Template
	PublisherAccount *model2.PublisherAccount
	OriginUrl        string
	BatchQuery       []BatchCreativeQuery
	Batch            []Base // 批量送审、查询使用
	IsUpdate         bool
	UrlUnique        string
}

type BatchCreativeQuery struct {
	MediaCid string
	Customer *model2.Customer
	Creative *model2.Creative
}

type BatchCreativeUpload struct {
	MediaCid string
	Customer *model2.Customer
	Creative *model2.Creative
}

func init() {
	// jsoniter 容忍字符串和数字互转、容忍空数组作为对象等操作
	// extra.RegisterFuzzyDecoders()时，注意不要在多个协程中重复设置，在init函数执行一次即可，否则会报错：fatal error: concurrent map writes
	extra.RegisterFuzzyDecoders()
}

func (b *BaseInfo) UploadAdvertiser() Ret {
	return Ret{}
}
func (b *BaseInfo) QueryAdvertiser() Ret {
	return Ret{}
}
func (b *BaseInfo) UploadCreative() Ret {
	return Ret{}
}
func (b *BaseInfo) QueryCreative() Ret {
	return Ret{}
}
func (b *BaseInfo) BatchQueryCreative() Ret {
	return Ret{}
}
func (b *BaseInfo) BatchUploadCreative() Ret {
	return Ret{}
}

type CreativeInfo struct {
	IsRsync            int
	CustomerId         int
	AdvertiserId       int
	AdvertiserName     string
	MediaInfo          string
	Industry           string
	PublisherId        int
	PositionId         int
	Publisher          string
	PublisherAccountId int
	Position           string
	CreativeId         string
	MaterialId         string //素材服务生成的创意id
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
	MediaCid           string
	Monitor            []model2.Monitor
	Info               []model2.TemplateInfo
	Cm                 []string
	Vm                 []string
	Extra              string
	RequestId          string
	VideoCdnUrl        string
	PicCdnUrl          string
	ErrCode            int
	CreateTime         int64
}

type Ret struct {
	// 公共参数
	Url    string `json:"url,omitempty"`
	Header string `json:"header,omitempty"`
	Req    string `json:"req,omitempty"`
	Resp   string `json:"resp,omitempty"`

	// 非批量接口参数
	Extra        string `json:"extra,omitempty"`
	VideoCdnUrl  string `json:"video_cdn_url,omitempty"`
	PicCdnUrl    string `json:"pic_cdn_url,omitempty"`
	PubReturnUrl string `json:"pub_return_url,omitempty"`
	MediaCosts   int64  `json:"media_costs"`
	UpLoadCosts  int64  `json:"up_load_costs"`
	CreativeId   int    `json:"creative_id"`
	ErrCode      int    `json:"err_code"`
	MediaCid     string `json:"media_cid"`
	ErrMsg       string `json:"err_msg"`
	IsRsync      int    `json:"is_rsync"`

	// 批量的参数
	BatchQueryRet     []BatchQueryRet           `json:"batch_query_ret"`
	BatchUploadRetMap map[string]BatchUploadRet `json:"batch_upload_ret_map"`
}

type BatchQueryRet struct {
	ErrCode  int    `json:"err_code"`
	MediaCid string `json:"media_cid"`
	ErrMsg   string `json:"err_msg"`
}

type BatchUploadRet struct {
	Extra       string `json:"extra,omitempty"`
	VideoCdnUrl string `json:"video_cdn_url,omitempty"`
	PicCdnUrl   string `json:"pic_cdn_url,omitempty"`
	MediaCosts  int64  `json:"media_costs"`
	CreativeId  int    `json:"creative_id"`
	ErrCode     int    `json:"err_code"`
	MediaCid    string `json:"media_cid"`
	ErrMsg      string `json:"err_msg"`
	IsRsync     int    `json:"is_rsync"`
}

func GetAdvertiserHandler(advertiserAudit *model2.AdvertiserAudit, logger *zap.Logger) (mediaHandler MediaHandler, err error) {

	publisherAccount, err := cache.GetPublisherAccount(advertiserAudit.CustomerId, advertiserAudit.PublisherId)
	if err != nil {
		return nil, err
	}
	publisherInfo, err := cache.GetPublisherCacheById(advertiserAudit.PublisherId)
	if err != nil {
		return nil, err
	}
	var advertiserUrls model2.AdvertiserUrls
	_ = jsoniter.Unmarshal([]byte(publisherInfo.AdvertiserUrls), &advertiserUrls)
	logg := logger.With(zap.String("advertiser_name", advertiserAudit.AdvertiserName))
	baseInfo := &BaseInfo{
		Logger:           logg,
		CustomerId:       advertiserAudit.CustomerId,
		PublisherId:      advertiserAudit.PublisherId,
		AdvertiserAudit:  advertiserAudit,
		PublisherAccount: publisherAccount,
		PublisherInfo:    publisherInfo,
		AdvertiserUrls:   advertiserUrls,
	}
	switch advertiserAudit.PublisherId {
	case enum.PUB_FANCY:
		mediaHandler = NewFancyHandler(baseInfo)
		break
	}
	return
}

func GetCreativeHandler(creative *model2.Creative, customer *model2.Customer, logger *zap.Logger) (mediaHandler MediaHandler, err error) {
	publisherId := creative.PublisherId
	// 获取媒体信息
	publisherInfo, err := cache.GetPublisherCacheById(publisherId)
	if err != nil {
		return nil, errors.Wrap(err, "媒体信息获取失败")
	}

	// 获取广告位信息
	positionInfo, err := cache.GetPositionCacheById(creative.PositionId)
	if err != nil {
		return nil, errors.Wrap(err, "广告位信息获取")
	}

	// 获取媒体账号信息
	publisherAccount, err := cache.GetPublisherAccount(creative.CustomerId, publisherId)
	if err != nil {
		return nil, errors.Wrap(err, "媒体账号不存在")
	}

	var advertiserAudit *model2.AdvertiserAudit
	// 广告主送审信息
	advertiserAudit, _ = cache.GetAdvertiserAuditCache(creative.CustomerId, creative.AdvertiserId, publisherId)

	// 获取创意信息
	creativeInfo := getCreativeInfo(creative)

	// 获取媒体素材送审url
	var creativeUrls model2.CreativeUrls
	err = jsoniter.Unmarshal([]byte(publisherInfo.CreativeUrls), &creativeUrls)
	if err != nil {
		return nil, errors.Wrap(err, "媒体送审域名获取失败")
	}

	// 广告位模板信息
	template := getTemplate(positionInfo, creativeInfo)
	logg := logger.With(zap.String("pub", publisherInfo.Name))
	baseInfo := &BaseInfo{
		CreativeId:       creative.Id,
		CustomerInfo:     customer,
		Logger:           logg,
		CustomerId:       creativeInfo.CustomerId,
		PublisherId:      creativeInfo.PublisherId,
		PublisherAccount: publisherAccount,
		PublisherInfo:    publisherInfo,
		CreativeInfo:     creativeInfo,
		PositionInfo:     positionInfo,
		Template:         template,
		CreativeUrls:     creativeUrls,
		AdvertiserAudit:  advertiserAudit,
	}

	return getHandler(publisherId, baseInfo), err
}

func getHandler(publisherId int, baseInfo *BaseInfo) (mediaHandler MediaHandler) {
	switch publisherId {
	case enum.PUB_TENCENT:
		mediaHandler = NewTencentHandler(baseInfo)
		break
	case enum.PUB_FANCY:
		mediaHandler = NewFancyHandler(baseInfo)
		break
	case enum.PUB_UC:
		mediaHandler = NewUCHandler(baseInfo)
		break

	}
	return mediaHandler
}
func GetBatchCreativeHandler(batchCreative []BatchCreativeQuery, publisherId int, customerId int, logger *zap.Logger) (mediaHandler MediaHandler, err error) {
	// 获取媒体信息
	publisherInfo, err := cache.GetPublisherCacheById(publisherId)
	if err != nil {
		return nil, errors.Wrap(err, "媒体信息获取失败")
	}

	// 获取媒体账号信息
	publisherAccount, err := cache.GetPublisherAccount(customerId, publisherId)
	if err != nil {
		return nil, errors.Wrap(err, "媒体账号不存在")
	}

	// 获取媒体素材送审url
	var creativeUrls model2.CreativeUrls
	err = jsoniter.Unmarshal([]byte(publisherInfo.CreativeUrls), &creativeUrls)
	if err != nil {
		return nil, errors.Wrap(err, "媒体送审域名获取失败")
	}
	logg := logger.With(zap.String("pub", publisherInfo.Name))
	baseInfo := &BaseInfo{
		PublisherAccount: publisherAccount,
		PublisherInfo:    publisherInfo,
		CreativeUrls:     creativeUrls,
		BatchQuery:       batchCreative,
		Logger:           logg,
	}
	return getHandler(publisherId, baseInfo), err
}

func GetBatchHandler(batchCreative []BatchCreativeUpload, publisherId int, customerId int, logger *zap.Logger, isUpdate bool) (mediaHandler MediaHandler, errs map[int]error) {
	var (
		batch []Base
	)
	errs = make(map[int]error)
	// 获取媒体信息
	publisherInfo, err := cache.GetPublisherCacheById(publisherId)
	if err != nil {
		for _, b := range batchCreative {
			errs[b.Creative.Id] = errors.New("媒体信息获取失败")
		}
		return nil, errs
	}
	// 获取媒体素材送审url
	var creativeUrls model2.CreativeUrls
	err = jsoniter.Unmarshal([]byte(publisherInfo.CreativeUrls), &creativeUrls)
	if err != nil {
		for _, b := range batchCreative {
			errs[b.Creative.Id] = errors.New("媒体送审域名获取失败")
		}
		return nil, errs
	}
	// 获取媒体账号信息
	publisherAccount, err := cache.GetPublisherAccount(customerId, publisherId)
	if err != nil {
		for _, b := range batchCreative {
			errs[b.Creative.Id] = errors.New("媒体账号不存在")
		}
		return nil, errs
	}
	for _, b := range batchCreative {
		creative := b.Creative

		// 获取广告位信息
		positionInfo, err := cache.GetPositionCacheById(creative.PositionId)
		if err != nil {
			errs[creative.Id] = errors.New("广告位信息获取")
			continue
		}
		var infos []model2.TemplateInfo
		_ = jsoniter.Unmarshal([]byte(creative.Info), &infos)
		creativeService := &service.Creative{
			TemplateId: creative.TemplateId,
			Info:       infos,
		}
		info, checkErrs := creativeService.Check(positionInfo, 0)
		if len(checkErrs) > 0 {
			errs[creative.Id] = errors.New(strings.Join(checkErrs, ","))
			continue
		}
		infoJson, _ := jsoniter.Marshal(info)
		creative.Info = string(infoJson)
		var advertiserAudit *model2.AdvertiserAudit
		if *publisherInfo.IsRsyncAdvertiser == 1 {
			// 广告主送审信息
			advertiserAudit, err = cache.GetAdvertiserAuditCache(creative.CustomerId, creative.AdvertiserId, publisherId)
			if err != nil {
				errs[creative.Id] = errors.New("广告主未审核")
				continue
			}
		}

		// 获取创意信息
		creativeInfo := getCreativeInfo(creative)

		// 广告位模板信息
		template := getTemplate(positionInfo, creativeInfo)
		batch = append(batch, Base{
			CreativeId:      creative.Id,
			CreativeInfo:    creativeInfo,
			PositionInfo:    positionInfo,
			Template:        template,
			AdvertiserAudit: advertiserAudit,
		})

	}
	baseInfo := &BaseInfo{
		IsUpdate:         isUpdate,
		CustomerId:       customerId,
		PublisherAccount: publisherAccount,
		PublisherId:      publisherId,
		Logger:           logger,
		Batch:            batch,
		CreativeUrls:     creativeUrls,
		PublisherInfo:    publisherInfo,
	}
	return getHandler(publisherId, baseInfo), errs
}

func getCreativeInfo(creative *model2.Creative) CreativeInfo {
	var (
		monitors []model2.Monitor
		infos    []model2.TemplateInfo
		cm       []string
		vm       []string
	)
	// 获取广告位模板
	_ = jsoniter.Unmarshal([]byte(creative.Monitor), &monitors)
	_ = jsoniter.Unmarshal([]byte(creative.Info), &infos)
	_ = jsoniter.Unmarshal([]byte(creative.Vm), &vm)
	_ = jsoniter.Unmarshal([]byte(creative.Cm), &cm)

	return CreativeInfo{
		IsRsync:            creative.IsRsync,
		CustomerId:         creative.CustomerId,
		AdvertiserId:       creative.AdvertiserId,
		PublisherAccountId: creative.PublisherAccountId,
		MediaInfo:          creative.MediaInfo,
		Industry:           creative.Industry,
		PublisherId:        creative.PublisherId,
		PositionId:         creative.PositionId,
		CreativeId:         creative.CreativeId,
		MaterialId:         common.GenMediaCid(creative.PublisherId, creative.Id),
		MediaCid:           creative.MediaCid,
		Name:               creative.Name,
		TemplateId:         creative.TemplateId,
		Action:             creative.Action,
		LandUrl:            creative.LandUrl,
		DeeplinkUrl:        creative.DeeplinkUrl,
		MiniProgramId:      creative.MiniProgramId,
		MiniProgramPath:    creative.MiniProgramPath,
		StartDate:          creative.StartDate,
		EndDate:            creative.EndDate,
		Monitor:            monitors,
		Info:               infos,
		Cm:                 cm,
		Vm:                 vm,
		RequestId:          creative.RequestId,
		Extra:              creative.Extra,
		VideoCdnUrl:        creative.VideoCdnUrl,
		PicCdnUrl:          creative.PicCdnUrl,
		ErrCode:            creative.ErrCode,
		CreateTime:         creative.CreatedAt.Unix(),
	}
}

func getTemplate(positionInfo *model2.Position, creativeInfo CreativeInfo) model2.Template {
	var materialInfo *model2.MaterialInfo
	_ = jsoniter.Unmarshal([]byte(positionInfo.MaterialInfo), &materialInfo)
	// 根据模板ID获取广告位模板信息
	var template model2.Template
	for _, list := range materialInfo.List {
		if list.Id == creativeInfo.TemplateId {
			template = list
		}
	}
	return template
}
