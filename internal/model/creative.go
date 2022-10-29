package model

import (
	"github.com/convee/adcreative/internal/pkg/common"
	"github.com/convee/adcreative/pkg/utils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

const (
	// 0送审返回结果异常、1上传失败、2上传不通过、3待送审、5审核中、6查询失败、7审核通过、8审核不通过、
	CREATIVE_UPDATE_EXCEPTION = 0
	CREATIVE_UPLOAD_FAILED    = 1
	CREATIVE_UPLOAD_UNPASSED  = 2
	CREATIVE_UPLOADING        = 3
	CREATIVE_AUDITING         = 5
	CREATIVE_QUERY_FAILED     = 6
	CREATIVE_AUDIT_PASSED     = 7
	CREATIVE_AUDIT_UNPASSWD   = 8

	// 0待审核，1审核通过，2审核不通过
	CREATIVE_STATUS_UPLOADING = -1
	CREATIVE_STATUS_AUDITING  = 0
	CREATIVE_STATUS_PASSED    = 1
	CREATIVE_STATUS_UNPASSED  = 2
)

var (
	StatusMap = map[int]string{
		CREATIVE_UPDATE_EXCEPTION: "0：送审查询结果更新异常",
		CREATIVE_UPLOAD_FAILED:    "1：上传失败",
		CREATIVE_UPLOAD_UNPASSED:  "2：上传不通过",
		CREATIVE_UPLOADING:        "3：待送审",
		CREATIVE_AUDITING:         "5：审核中",
		CREATIVE_QUERY_FAILED:     "6：查询失败",
		CREATIVE_AUDIT_PASSED:     "7：审核通过",
		CREATIVE_AUDIT_UNPASSWD:   "8：审核不通过",
	}
)

type Creative struct {
	Model
	Id                 int    `json:"id"`
	CreativeId         string `json:"creative_id"`
	MediaCid           string `json:"media_cid"`
	CustomerId         int    `json:"customer_id"`
	PublisherAccountId int    `json:"publisher_account_id"`
	AdvertiserId       int    `json:"advertiser_id"`
	MediaInfo          string `json:"media_info"`
	Industry           string `json:"industry"`
	PublisherId        int    `json:"publisher_id"`
	PositionId         int    `json:"position_id"`
	Name               string `json:"name"`
	TemplateId         string `json:"template_id"`
	Info               string `json:"info"`
	LandUrl            string `json:"land_url"`
	DeeplinkUrl        string `json:"deeplink_url"`
	StartDate          string `json:"start_date"`
	EndDate            string `json:"end_date"`
	Action             int    `json:"action"`
	MiniProgramId      string `json:"mini_program_id"`
	MiniProgramPath    string `json:"mini_program_path"`
	Cm                 string `json:"cm"`
	Vm                 string `json:"vm"`
	Monitor            string `json:"monitor"`
	Status             int    `json:"status"`   // 0待审核，1审核通过，2审核不通过
	IsRsync            int    `json:"is_rsync"` // 是否同步
	Reason             string `json:"reason"`
	ErrCode            int    `json:"err_code"` // 1上传失败、2待送审、3送审失败、4审核中、5审核通过、6审核不通过
	ErrMsg             string `json:"err_msg"`
	RequestId          string `json:"request_id"`     // 唯一请求ID
	Extra              string `json:"extra"`          //额外信息
	VideoCdnUrl        string `json:"video_cdn_url"`  //媒体视频ID
	PicCdnUrl          string `json:"pic_cdn_url"`    //媒体图片ID
	PubReturnUrl       string `json:"pub_return_url"` //媒体素材cdn
	Priority           int    `json:"priority"`       // 优先级
}

type CreativeGroupStat struct {
	PublisherId  string `json:"publisher_id"`
	CustomerId   string `json:"customer_id"`
	AdvertiserId int    `json:"advertiser_id"`
	StatusType   string `json:"status_type"` // 内部错误码 1上传失败、2上传不通过、3待送审、4送审失败、5审核中、6查询失败、7审核通过、8审核不通过
	Count        int    `json:"count"`
}

type CreativeStatusStat struct {
	PublisherName string `json:"publisher_name"`
	StatusType    string `json:"status_type"` // 内部错误码 1上传失败、2上传不通过、3待送审、4送审失败、5审核中、6查询失败、7审核通过、8审核不通过
	Count         int    `json:"count"`
}

type Monitor struct {
	T   int    `json:"t"`
	Url string `json:"url"`
}

type TemplateInfo struct {
	AttrName  string            `json:"attr_name"`
	AttrValue string            `json:"attr_value"`
	Md5       string            `json:"md5"`
	Width     int               `json:"width"`
	Height    int               `json:"height"`
	Ext       string            `json:"ext"`
	Duration  int               `json:"duration"`
	Size      int               `json:"size"`
	MediaKey  string            `json:"media_key"`
	Extra     map[string]string `json:"extra"`
}

type CreativeModel struct {
}

func GetCreativeStatusByErrCode(errCode int) int {
	if errCode == CREATIVE_UPLOAD_FAILED || errCode == CREATIVE_UPLOADING {
		return CREATIVE_STATUS_UPLOADING
	} else if errCode == CREATIVE_AUDITING || errCode == CREATIVE_QUERY_FAILED {
		return CREATIVE_STATUS_AUDITING
	} else if errCode == CREATIVE_AUDIT_PASSED {
		return CREATIVE_STATUS_PASSED
	} else if errCode == CREATIVE_UPLOAD_UNPASSED || errCode == CREATIVE_AUDIT_UNPASSWD {
		return CREATIVE_STATUS_UNPASSED
	}
	return CREATIVE_STATUS_AUDITING
}

// GetGroupCreativeStats
// 状态码 0错误码异常、1上传失败、2上传不通过、3待送审、5审核中、6查询失败、7审核通过、8审核不通过',
func GetGroupCreativeStats() (lst []CreativeGroupStat, err error) {
	ql := `publisher.name as publisher_id,customer.name as customer_id,creative.advertiser_id,case creative.err_code
 when 0 then '错误码异常' 
 when 1 then '送审失败' 
 when 2 then '送审不通过' 
 when 3 then '待送审'
 when 5 then '审核中' 
 when 6 then '审核中' 
 when 7 then '审核通过' 
 when 8 then '审核不通过'
 else creative.err_code end as status_type,count(1) as count`
	err = GetDB().Table("creative").Select(ql).Joins("left join publisher on publisher.id=creative.publisher_id left join customer on customer.id=creative.customer_id").
		Where("creative.deleted_at is null").
		Group("creative.publisher_id,creative.customer_id,creative.advertiser_id,status_type,creative.err_code").Find(&lst).Error
	return
}

// TableName sets the insert table name for this struct type
func (c *CreativeModel) TableName() string {
	return "creative"
}

func (c CreativeModel) GetCreativeTotal(maps interface{}) (int64, error) {
	var count int64
	if err := GetDB().Model(&Creative{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (c CreativeModel) GetCreativeById(id int) (*Creative, error) {
	var creative *Creative
	err := GetDB().First(&creative, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("记录不存在")
		}
		return nil, err
	}
	return creative, nil
}

func (c CreativeModel) GetOneCreativeByMaps(maps interface{}) (*Creative, error) {
	var creative *Creative
	err := GetDB().Where(maps).First(&creative).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return creative, nil
}

func (c CreativeModel) GetAllCreativesByMaps(creativeIds []string) ([]*Creative, error) {
	var creatives []*Creative
	err := GetDB().Where("creative_id in ? ", creativeIds).Find(&creatives).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("记录不存在")
		}
		return nil, err
	}
	return creatives, nil
}

func (c CreativeModel) GetAllCreativeIdsByMaps(maps interface{}) ([]*Creative, error) {
	var creatives []*Creative
	err := GetDB().Where(maps).Find(&creatives).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("记录不存在")
		}
		return nil, err
	}
	return creatives, nil
}

func (c CreativeModel) GetAllCreativeByErrCodes(publisherId int, errCodes []int) ([]*Creative, error) {
	var (
		creatives         []*Creative
		unpassedCreatives []*Creative

		errCodeOk    []int      // 正常的错误码
		unpassedCode int   = -1 // 标记 审核不通过的状态码
		err          error
	)

	for i := range errCodes {
		if errCodes[i] != CREATIVE_AUDIT_UNPASSWD {
			errCodeOk = append(errCodeOk, errCodes[i])
		} else {
			unpassedCode = CREATIVE_AUDIT_UNPASSWD
		}
	}

	if unpassedCode != -1 {
		_ = GetDB().Where("publisher_id=? and err_code = ? and updated_at >= ?", publisherId, unpassedCode, utils.GetLastTenDaysTimestamp()).Order("priority desc").Find(&unpassedCreatives).Error
	}

	err = GetDB().Where("publisher_id=? and err_code in ?", publisherId, errCodeOk).Order("priority desc").Find(&creatives).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("记录不存在")
		}
		return nil, err
	}

	if len(unpassedCreatives) != 0 {
		creatives = append(creatives, unpassedCreatives...)
	}

	return creatives, nil
}

func (c CreativeModel) GetCreatives(page int, perPage int, maps interface{}) ([]*Creative, error) {
	var creatives []*Creative
	offset := (page - 1) * perPage
	err := GetDB().Where(maps).Offset(offset).Limit(perPage).Order("id desc").Find(&creatives).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("记录不存在")
		}
		return nil, err
	}
	return creatives, nil
}

func (c CreativeModel) GetPassCreatives(page int, perPage int, maps interface{}, date string) ([]*Creative, error) {
	var creatives []*Creative
	offset := (page - 1) * perPage
	err := GetDB().Where(maps).Where("updated_at > ?", date).Offset(offset).Limit(perPage).Order("id desc").Find(&creatives).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("记录不存在")
		}
		return nil, err
	}
	return creatives, nil
}

// CreateCreative 事务添加创意，同时更新媒体侧创意ID
func (c CreativeModel) CreateCreative(creative *Creative) (*Creative, error) {
	GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&creative).Error; err != nil {
			//返回任何错误都会回滚事务
			return err
		}
		if utils.IsInSlice(creative.PublisherId, common.MediaCidPubList) {
			creative.MediaCid = common.GenMediaCid(creative.PublisherId, creative.Id)
			if err := tx.Save(&creative).Error; err != nil {
				return err
			}
		}
		// 返回 nil 提交事务
		return nil
	})
	return creative, nil
}

func (c CreativeModel) UpdateCreative(creative *Creative) (*Creative, error) {
	result := GetDB().Save(creative)
	if result.Error != nil {
		return creative, result.Error
	}
	return creative, nil
}

func (c CreativeModel) UpdateCreativeByMap(id int, maps map[string]interface{}) (Creative, error) {
	var creative Creative
	creative.Id = id
	result := GetDB().Model(&creative).Updates(maps)
	if result.Error != nil {
		return creative, result.Error
	}
	return creative, nil
}

func (c CreativeModel) DeleteCreative(creative *Creative) (*Creative, error) {
	GetDB().Delete(&creative)
	return creative, nil
}

// GetStatisticsCreativeStats
// 状态码 0错误码异常、1上传失败、2上传不通过、3待送审、5审核中、6查询失败、7审核通过、8审核不通过',
func GetStatisticsCreativeStats(maps map[string]interface{}) (lst []CreativeStatusStat, err error) {
	ql := `publisher.name as publisher_name,err_code as status_type,count(1) as count`

	db := GetDB().Table("creative").Select(ql).Joins("left join publisher on publisher.id=creative.publisher_id").Where(maps)

	// 不是根据RM创意id查询，则查询最近90天数据
	if _, ok := maps[""]; !ok {
		dateTime := time.Now().AddDate(0, 0, -90)
		date := dateTime.Format("2006-01-02 15:04:05")
		db.Where("creative.created_at > ?", date)
	}
	err = db.Where("creative.deleted_at is null").
		Group("creative.publisher_id,creative.customer_id,creative.err_code").Find(&lst).Error
	return lst, err
}
