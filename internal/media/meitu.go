package media

import (
	"fmt"
	model2 "github.com/convee/adcreative/internal/model"
	"github.com/convee/adcreative/pkg/httpclient"
	"github.com/convee/adcreative/pkg/md5"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cast"
	"strconv"
	"time"
)

const (
	MeiTuSuccess           = 0
	MeiTuTokenError        = -1
	MeiTuDataParamsErr     = -2
	MeiTuRequestNetworkErr = -3
	MeiTuFrequencyErr      = -99 //同一个 IP一天只能请求审核 接口送审1000次，查询3000次，第二天自动重置。
	MeiTuOtherError        = -100

	MeiTuQueryStatusAuditing = "待审核"
	MeiTuQueryStatusPassed   = "审核通过"
	MeiTuQueryStatusRefused  = "审核拒绝"
)

type MeiTu struct {
	*BaseInfo
}

type MeiTuUploadResp struct {
	ErrCode int    `json:"error_code"`
	Msg     string `json:"msg"`
}

type MeiTuUploadRequest struct {
	Token     string                `json:"token"`
	Timestamp string                `json:"timestamp"`
	Data      []MeiTuUploadMaterial `json:"data"`
}

type MeiTuUploadMaterial struct {
	AdId             string `json:"ad_id"`                    //创意ID，必填
	OsType           string `json:"os_type"`                  //平台=ios 或者 android，必填
	PositionId       int    `json:"position_id"`              //广告位 ID，注意广告位 ID 不存在的会报错，必填
	LinkInstructions string `json:"link_instructions"`        //跳转链接，必填
	Main             string `json:"main,omitempty"`           //大图 URL
	Icon             string `json:"icon,omitempty"`           //icon,即小图 URL
	Video            string `json:"video,omitempty"`          //视频URL
	Cover            string `json:"cover,omitempty"`          //封面URL
	Title            string `json:"title,omitempty"`          //标题文案
	Desc             string `json:"desc,omitempty"`           //描述文案
	CtaText          string `json:"ctatext,omitempty"`        //按钮文案
	AdNetworkId      string `json:"ad_network_id"`            //广告源标识，注意不存在的广告源会报错，必填
	ThirdTemplate    string `json:"third_template,omitempty"` //模板标识，摇一摇 =mt_shake_splash
}

type MeiTuQueryRequest struct {
	Token     string               `json:"token"`
	Timestamp string               `json:"timestamp"`
	Data      []MeiTuQueryMaterial `json:"data"`
}
type MeiTuQueryMaterial struct {
	AdId        string `json:"ad_id"`
	AdNetworkId string `json:"ad_network_id"`
}

type MeiTuQueryResp struct {
	ErrCode int                `json:"error_code"`
	Msg     string             `json:"msg"`
	Data    []MeiTuResultQuery `json:"data"`
}

type MeiTuResultQuery struct {
	AdId      string `json:"ad_id"`
	AdNetwork string `json:"ad_network"`
	Status    string `json:"status,omitempty"`
	Reason    string `json:"reason,omitempty"`
}

func NewMeiTuHandler(b *BaseInfo) MediaHandler {
	return &MeiTu{BaseInfo: b}
}
func (m *MeiTu) UploadAdvertiser() Ret {
	return Ret{}
}
func (m *MeiTu) QueryAdvertiser() Ret {
	return Ret{}
}

func (m *MeiTu) getUploadParams(creativeInfo CreativeInfo, template model2.Template) MeiTuUploadMaterial {
	material := MeiTuUploadMaterial{
		AdId:             creativeInfo.MediaCid,
		LinkInstructions: creativeInfo.LandUrl,
		AdNetworkId:      creativeInfo.MediaInfo,
	}

	if adslot, ok := template.Extra["adslot"]; ok {
		material.PositionId = cast.ToInt(adslot)
	}
	if osType, ok := template.Extra["os_type"]; ok {
		material.OsType = osType
	}
	if thirdTemplate, ok := template.Extra["third_template"]; ok {
		material.ThirdTemplate = thirdTemplate
	}
	for _, f := range creativeInfo.Info {
		if f.AttrName == "image" {
			material.Main = f.AttrValue
		}
		if f.AttrName == "video" {
			material.Video = f.AttrValue
		}
		if f.AttrName == "cover" {
			material.Cover = f.AttrValue
		}
		if f.AttrName == "icon" {
			material.Icon = f.AttrValue
		}
		if f.AttrName == "title" {
			material.Title = f.AttrValue
		}
		if f.AttrName == "desc" {
			material.Desc = f.AttrValue
		}
		if f.AttrName == "ctatext" {
			material.CtaText = f.AttrValue
		}

	}
	return material
}

func (m *MeiTu) UploadCreative() Ret {
	uri := m.CreativeUrls.CreateUrl
	materials := make([]MeiTuUploadMaterial, 0)
	material := m.getUploadParams(m.CreativeInfo, m.Template)
	materials = append(materials, material)
	return m.sendUploadPost(uri, materials)
}

func (m *MeiTu) BatchQueryCreative() Ret {
	var (
		materials []MeiTuQueryMaterial
	)
	for _, batch := range m.BatchQuery {
		materials = append(materials, MeiTuQueryMaterial{
			AdId:        batch.Creative.MediaCid,
			AdNetworkId: batch.Creative.MediaInfo,
		})
	}
	return m.sendQueryPost(m.CreativeUrls.QueryUrl, materials)
}
func (m *MeiTu) BatchUploadCreative() Ret {
	var (
		materials []MeiTuUploadMaterial
		ret       Ret
	)
	uri := m.CreativeUrls.CreateUrl
	ret.BatchUploadRetMap = make(map[string]BatchUploadRet)
	for _, batch := range m.Batch {
		materials = append(materials, m.getUploadParams(batch.CreativeInfo, batch.Template))
	}
	return m.sendUploadPost(uri, materials)
}
func (m *MeiTu) QueryCreative() Ret {
	var (
		ret       Ret
		materials []MeiTuQueryMaterial
	)
	url := m.CreativeUrls.QueryUrl
	materials = append(materials, MeiTuQueryMaterial{
		AdId:        m.CreativeInfo.MediaCid,
		AdNetworkId: m.CreativeInfo.MediaInfo,
	})
	ret = m.sendQueryPost(url, materials)
	return ret
}

func (m *MeiTu) sendUploadPost(uri string, materials []MeiTuUploadMaterial) Ret {
	var ret Ret
	timestamp := time.Now().Unix()
	request := MeiTuUploadRequest{
		Token:     m.sign(timestamp),
		Timestamp: cast.ToString(timestamp),
		Data:      materials,
	}
	bodyJson, _ := jsoniter.Marshal(request)
	ret.Url = uri
	ret.Req = string(bodyJson)
	response, err := httpclient.PostJSON(uri, bodyJson, httpclient.WithTTL(time.Second*180))
	ret.Resp = string(response)
	if err != nil {
		ret.ErrCode = model2.CREATIVE_UPLOAD_FAILED
		ret.ErrMsg = err.Error()
		return ret
	}
	var resp *MeiTuUploadResp
	err = jsoniter.Unmarshal(response, &resp)
	if err != nil {
		ret.ErrCode = model2.CREATIVE_UPLOAD_FAILED
		ret.ErrMsg = err.Error()
		return ret
	}
	if resp.ErrCode != MeiTuSuccess {
		ret.ErrCode = model2.CREATIVE_UPLOAD_FAILED
		ret.ErrMsg = resp.Msg
		return ret
	}
	ret.BatchUploadRetMap = make(map[string]BatchUploadRet)
	for _, batch := range m.Batch {
		ret.BatchUploadRetMap[batch.CreativeInfo.MediaCid] = BatchUploadRet{
			ErrCode: model2.CREATIVE_AUDITING,
			ErrMsg:  "审核中",
		}
	}
	return ret
}

func (m *MeiTu) sendQueryPost(uri string, materials []MeiTuQueryMaterial) Ret {
	var ret Ret
	timestamp := time.Now().Unix()
	request := MeiTuQueryRequest{
		Token:     m.sign(timestamp),
		Timestamp: cast.ToString(timestamp),
		Data:      materials,
	}
	bodyJson, _ := jsoniter.Marshal(request)
	ret.Url = uri
	ret.Req = string(bodyJson)
	response, err := httpclient.PostJSON(uri, bodyJson, httpclient.WithTTL(time.Second*180))
	ret.Resp = string(response)
	if err != nil {
		ret.ErrCode = model2.CREATIVE_QUERY_FAILED
		ret.ErrMsg = err.Error()
		return ret
	}
	var resp *MeiTuQueryResp
	err = jsoniter.Unmarshal(response, &resp)
	if err != nil {
		ret.ErrCode = model2.CREATIVE_QUERY_FAILED
		ret.ErrMsg = err.Error()
		return ret
	}
	if resp.ErrCode != MeiTuSuccess || len(resp.Data) <= 0 {
		ret.ErrCode = model2.CREATIVE_QUERY_FAILED
		ret.ErrMsg = string(response)
		return ret
	}
	for _, r := range resp.Data {
		var queryRet BatchQueryRet
		queryRet.MediaCid = r.AdId
		if r.Status == MeiTuQueryStatusPassed {
			queryRet.ErrCode = model2.CREATIVE_AUDIT_PASSED
			queryRet.ErrMsg = "审核通过"
		} else if r.Status == MeiTuQueryStatusRefused {
			queryRet.ErrCode = model2.CREATIVE_AUDIT_UNPASSWD
			queryRet.ErrMsg = fmt.Sprintf("媒体返回：%s", r.Reason)
		} else {
			queryRet.ErrCode = model2.CREATIVE_AUDITING
			queryRet.ErrMsg = "审核中"
		}
		ret.BatchQueryRet = append(ret.BatchQueryRet, queryRet)
	}

	return ret
}

func (m *MeiTu) sign(timestamp int64) string {
	return md5.New().Encrypt(m.PublisherAccount.DspId + "*" + m.PublisherAccount.Token + "&" + strconv.FormatInt(timestamp, 10))
}
