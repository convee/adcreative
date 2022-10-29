package media

import (
	"encoding/json"
	"fmt"
	"github.com/convee/adcreative/internal/model"
	"github.com/convee/adcreative/internal/pkg/common"
	"github.com/convee/adcreative/pkg/httpclient"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cast"
	"time"
)

type UC struct {
	*BaseInfo
}

const (
	UCQueryStatusPassed   = 0
	UCQueryStatusAuditing = 1
	UCQueryStatusUnPassed = 2
)

type UCUploadRequest struct {
	DspId    uint64       `json:"dsp_id"`
	Token    string       `json:"token"`
	Material []UCMaterial `json:"material"`
}
type UCMaterial struct {
	CreativeId      uint64       `json:"creative_id"`       //创意ID
	Creative        []UCCreative `json:"creative"`          //创意内容
	Type            uint32       `json:"type"`              //1：上传；2删除
	MonitorUrl      string       `json:"monitor_url"`       // 曝光监测地址，非必填
	ClickMonitorUrl string       `json:"click_monitor_url"` //点击监测地址，非必填
	TemplateId      string       `json:"template_id"`       // 创意实用模板ID，见template_id枚举
	StartDate       string       `json:"start_date"`        //素材生效日期 格式 YYYY-mm-dd
	EndDate         string       `json:"end_date"`          //素材失效日期 格式 YYYY-mm-dd

}

type UCCreative struct {
	AttrName  string `json:"attr_name"`  //创意属性名
	AttrValue string `json:"attr_value"` //创意属性值
	Width     uint32 `json:"width"`      //图片或视频必填，文本或h5为空
	Height    uint32 `json:"height"`     //图片或视频必填，文本或h5为空
	Ext       string `json:"ext"`        //图片或视频的扩展名或文本类型（txt）
}
type UCUploadResp struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}
type UCQueryRequest struct {
	DspId       uint64   `json:"dsp_id"`
	Token       string   `json:"token"`
	CreativeIds []uint64 `json:"creative_ids"`
}

type UCQueryResp struct {
	Status  int             `json:"status"`
	Message string          `json:"message"`
	Data    UCQueryRespData `json:"data"`
}
type UCQueryRespData struct {
	Result  int                    `json:"result"`
	Total   uint32                 `json:"total"`
	Records []UCResultQueryRecords `json:"records"`
}

type UCResultQueryRecords struct {
	CreativeId uint64 `json:"creative_id"`
	State      int    `json:"state"` //1审核中，2审核拒绝，0审核通过
	Reason     string `json:"reason"`
}

func NewUCHandler(b *BaseInfo) MediaHandler {
	return &UC{BaseInfo: b}
}

func (u *UC) UploadCreative() Ret {
	uri := u.CreativeUrls.CreateUrl

	request := UCUploadRequest{
		DspId: cast.ToUint64(u.PublisherAccount.DspId),
		Token: u.PublisherAccount.Token,
	}
	material := UCMaterial{
		CreativeId: cast.ToUint64(u.CreativeInfo.MediaCid),
		//MonitorUrl:      u.CreativeInfo.Monitor,
		//ClickMonitorUrl: u.CreativeInfo.Cm,
		TemplateId: u.Template.Extra["template_id"],
		Type:       1,
		StartDate:  u.CreativeInfo.StartDate,
		EndDate:    u.CreativeInfo.EndDate,
	}
	mediaKeys := map[string]string{}
	for _, t := range u.Template.Info {
		if len(t.MediaKey) > 0 {
			mediaKeys[t.Key] = t.MediaKey
		}
	}

	for _, i := range u.CreativeInfo.Info {
		if mediaKey, ok := mediaKeys[i.AttrName]; ok {
			if mediaKey == "deeplink" && u.CreativeInfo.DeeplinkUrl != "" {
				material.Creative = append(material.Creative, UCCreative{
					AttrName:  mediaKey,
					AttrValue: u.CreativeInfo.DeeplinkUrl,
					Ext:       "txt",
				})
				continue
			}
			if i.AttrValue == "" {
				continue
			}
			if common.StringsContain(mediaKey, "title", "description", "sub_title", "logo_title", "source", "deeplink", "download_url", "universal_link", "app_name", "app_info_url", "developer", "update_time", "privacy", "permission", "version_name", "related_creative_id", "platform", "package_name") {
				material.Creative = append(material.Creative, UCCreative{
					AttrName:  mediaKey,
					AttrValue: i.AttrValue,
					Ext:       "txt",
				})
			}
			if common.StringsContain(mediaKey, "video", "pic1_img", "pic2_img", "pic3_img", "logo_image") {
				material.Creative = append(material.Creative, UCCreative{
					AttrName:  mediaKey,
					AttrValue: i.AttrValue,
					Width:     cast.ToUint32(i.Width),
					Height:    cast.ToUint32(i.Height),
					Ext:       i.Ext,
				})
			}
		}
	}

	if u.CreativeInfo.LandUrl != "" {
		material.Creative = append(material.Creative, UCCreative{
			AttrName:  "landing_page",
			AttrValue: u.CreativeInfo.LandUrl,
			Ext:       "txt",
		})
	}

	request.Material = append(request.Material, material)
	return u.sendUploadPost(uri, request)
}

func (u *UC) QueryCreative() Ret {
	url := u.CreativeUrls.QueryUrl
	request := UCQueryRequest{
		DspId:       cast.ToUint64(u.PublisherAccount.DspId),
		Token:       u.PublisherAccount.Token,
		CreativeIds: []uint64{cast.ToUint64(u.CreativeInfo.MediaCid)},
	}
	return u.sendQueryPost(url, request)
}

func (u *UC) sendUploadPost(uri string, request UCUploadRequest) Ret {
	var ret Ret

	bodyJson, _ := json.Marshal(request)
	var option []httpclient.Option
	option = append(option, httpclient.WithTTL(time.Second*100))
	response, err := httpclient.PostJSON(uri, bodyJson, option...)
	ret.Url = uri
	ret.Req = string(bodyJson)
	ret.Resp = string(response)
	if err != nil {
		ret.ErrCode = model.CREATIVE_UPLOAD_FAILED
		ret.ErrMsg = err.Error()
		return ret
	}
	var resp *UCUploadResp
	err = jsoniter.Unmarshal(response, &resp)
	if err != nil {
		ret.ErrCode = model.CREATIVE_UPLOAD_FAILED
		ret.ErrMsg = err.Error()
		return ret
	}
	if resp.Status == 0 {
		ret.ErrCode = model.CREATIVE_AUDITING
		ret.ErrMsg = "审核中"
	} else {
		ret.ErrCode = model.CREATIVE_UPLOAD_UNPASSED
		ret.ErrMsg = fmt.Sprintf("媒体返回：%s", resp.Message)
	}
	return ret
}

func (u *UC) sendQueryPost(uri string, request UCQueryRequest) Ret {
	var ret Ret
	bodyJson, _ := jsoniter.Marshal(request)
	var option []httpclient.Option
	option = append(option, httpclient.WithTTL(time.Second*100))
	response, err := httpclient.PostJSON(uri, bodyJson, option...)
	ret.Url = uri
	ret.Req = string(bodyJson)
	ret.Resp = string(response)
	if err != nil {
		ret.ErrCode = model.CREATIVE_QUERY_FAILED
		ret.ErrMsg = err.Error()
		return ret
	}
	var resp *UCQueryResp
	err = jsoniter.Unmarshal(response, &resp)
	if err != nil {
		ret.ErrCode = model.CREATIVE_QUERY_FAILED
		ret.ErrMsg = err.Error()
		return ret
	}
	if resp.Status != 0 || resp.Data.Result != 0 || len(resp.Data.Records) <= 0 {
		ret.ErrCode = model.CREATIVE_QUERY_FAILED
		ret.ErrMsg = resp.Message
		return ret
	}

	for _, response := range resp.Data.Records {
		if response.State == UCQueryStatusPassed {
			ret.ErrCode = model.CREATIVE_AUDIT_PASSED
			ret.ErrMsg = "审核通过"
		} else if response.State == UCQueryStatusUnPassed {
			ret.ErrCode = model.CREATIVE_AUDIT_UNPASSWD
			ret.ErrMsg = fmt.Sprintf("媒体返回：%s", response.Reason)
		} else {
			ret.ErrCode = model.CREATIVE_AUDITING
			ret.ErrMsg = "审核中"
		}
	}
	return ret
}
