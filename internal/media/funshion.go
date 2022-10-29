package media

import (
	"fmt"
	model2 "github.com/convee/adcreative/internal/model"
	"strings"
	"time"

	"github.com/convee/adcreative/pkg/httpclient"
	"github.com/golang-module/carbon"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cast"
)

const (
	FunshionSuccess = 0

	FunshionQueryStatusAuditing = 0
	FunshionQueryStatusPassed   = 1
	FunshionQueryStatusRefused  = -1
)

type Funshion struct {
	*BaseInfo
}

type FunshionUploadResp struct {
	Result  int         `json:"result"`
	Message interface{} `json:"message"`
}
type FunshionUploadRespMessage struct {
	Crid   string `json:"crid"`
	Result int    `json:"result"`
	Reason string `json:"reason,omitempty"`
}
type FunshionUploadRequest struct {
	DspId    string                   `json:"dspid"`
	Token    string                   `json:"token"`
	Material []FunshionUploadMaterial `json:"material"`
}

type FunshionUploadMaterial struct {
	Crid        string       `json:"crid"`       //创意ID，必填
	Advertiser  string       `json:"advertiser"` //广告主名称
	Adm         string       `json:"adm"`
	Duration    int          `json:"duration"`
	LandingPage string       `json:"landingpage"`
	StartDate   string       `json:"startdate"`
	EndDate     string       `json:"enddate"`
	Framing     string       `json:"framing,omitempty"` //定帧图，暂未使用
	Title       string       `json:"title,omitempty"`
	Description string       `json:"description,omitempty"`
	Pm          []FunshionPm `json:"pm"`
	Cm          []string     `json:"cm"`
	Type        string       `json:"type"` //物料类型：必填，图片（image）、视频（video）、Flash（flash）
}

type FunshionPm struct {
	Point int    `json:"point"`
	Url   string `json:"url"`
}

type FunshionQueryRequest struct {
	DspId string   `json:"dspid"`
	Token string   `json:"token"`
	Crid  []string `json:"crid"`
}

type FunshionQueryResp struct {
	Result  int         `json:"result"`
	Message interface{} `json:"message"`
}

type FunshionQueryRespMessage struct {
	Total   int                              `json:"total"`
	Records []FunshionQueryRespMessageRecord `json:"records"`
}

type FunshionQueryRespMessageRecord struct {
	Crid   string `json:"crid"`
	Result int    `json:"result"`
	Reason string `json:"reason,omitempty"`
}

func NewFunshionHandler(b *BaseInfo) MediaHandler {
	return &Funshion{BaseInfo: b}
}
func (f *Funshion) UploadAdvertiser() Ret {
	return Ret{}
}
func (f *Funshion) QueryAdvertiser() Ret {
	return Ret{}
}

func (f *Funshion) getUploadParams(creativeInfo CreativeInfo, template model2.Template) FunshionUploadMaterial {
	material := FunshionUploadMaterial{
		Crid:        creativeInfo.MediaCid,
		LandingPage: creativeInfo.LandUrl,
		StartDate:   creativeInfo.StartDate,
		EndDate:     carbon.Parse(creativeInfo.EndDate).AddMonths(3).ToDateString(),
		Advertiser:  creativeInfo.MediaInfo,
	}

	if typ, ok := template.Extra["type"]; ok {
		material.Type = typ
	}
	for _, f := range creativeInfo.Info {
		if strings.HasPrefix(f.AttrName, "image") || strings.HasPrefix(f.AttrName, "video") {
			material.Adm = f.AttrValue
			if f.Duration > 0 {
				material.Duration = f.Duration
			}
		}
	}
	if duration, ok := template.Extra["duration"]; ok {
		material.Duration = cast.ToInt(duration)
	}
	if len(creativeInfo.Monitor) > 0 {
		for _, m := range creativeInfo.Monitor {
			material.Pm = append(material.Pm, FunshionPm{
				Point: m.T,
				Url:   m.Url,
			})
		}
	}
	if len(creativeInfo.Cm) > 0 {
		material.Cm = creativeInfo.Cm
	}
	return material
}

func (f *Funshion) UploadCreative() Ret {
	uri := f.CreativeUrls.CreateUrl
	materials := make([]FunshionUploadMaterial, 0)
	material := f.getUploadParams(f.CreativeInfo, f.Template)
	materials = append(materials, material)
	ret := f.sendUploadPost(uri, materials)
	if len(ret.BatchUploadRetMap) > 0 {
		ret.ErrCode = ret.BatchUploadRetMap[material.Crid].ErrCode
		ret.ErrMsg = ret.BatchUploadRetMap[material.Crid].ErrMsg
	}
	return ret
}

func (f *Funshion) BatchQueryCreative() Ret {
	var (
		crid []string
	)
	for _, batch := range f.BatchQuery {
		crid = append(crid, batch.Creative.MediaCid)
	}
	return f.sendQueryPost(f.CreativeUrls.QueryUrl, crid)
}
func (f *Funshion) BatchUploadCreative() Ret {
	var (
		materials []FunshionUploadMaterial
		ret       Ret
	)
	uri := f.CreativeUrls.CreateUrl
	ret.BatchUploadRetMap = make(map[string]BatchUploadRet)
	for _, batch := range f.Batch {
		materials = append(materials, f.getUploadParams(batch.CreativeInfo, batch.Template))
	}
	return f.sendUploadPost(uri, materials)
}
func (f *Funshion) QueryCreative() Ret {
	var (
		ret  Ret
		crid []string
	)
	url := f.CreativeUrls.QueryUrl
	crid = append(crid, f.CreativeInfo.MediaCid)
	ret = f.sendQueryPost(url, crid)
	if len(ret.BatchQueryRet) > 0 {
		ret.ErrCode = ret.BatchQueryRet[0].ErrCode
		ret.ErrMsg = ret.BatchQueryRet[0].ErrMsg
		ret.MediaCid = ret.BatchQueryRet[0].MediaCid
	}
	return ret
}

func (f *Funshion) sendUploadPost(uri string, materials []FunshionUploadMaterial) Ret {
	var ret Ret
	request := FunshionUploadRequest{
		DspId:    f.PublisherAccount.DspId,
		Token:    f.PublisherAccount.Token,
		Material: materials,
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
	var resp *FunshionUploadResp
	err = jsoniter.Unmarshal(response, &resp)
	if err != nil {
		ret.ErrCode = model2.CREATIVE_UPLOAD_FAILED
		ret.ErrMsg = err.Error()
		return ret
	}
	if resp.Result != FunshionSuccess {
		ret.ErrCode = model2.CREATIVE_UPLOAD_FAILED
		ret.ErrMsg = string(response)
		return ret
	}
	messageByte, err := jsoniter.Marshal(resp.Message)
	if err != nil {
		ret.ErrCode = model2.CREATIVE_UPLOAD_FAILED
		ret.ErrMsg = err.Error()
		return ret
	}
	var messageSlice []FunshionUploadRespMessage
	err = jsoniter.Unmarshal(messageByte, &messageSlice)
	if err != nil {
		ret.ErrCode = model2.CREATIVE_UPLOAD_FAILED
		ret.ErrMsg = err.Error()
		return ret
	}
	if len(messageSlice) <= 0 {
		ret.ErrCode = model2.CREATIVE_UPLOAD_FAILED
		ret.ErrMsg = "媒体返回为空"
		return ret
	}
	var messageByCrid = map[string]FunshionUploadRespMessage{}
	for _, m := range messageSlice {
		messageByCrid[m.Crid] = m
	}
	ret.BatchUploadRetMap = make(map[string]BatchUploadRet)
	for _, batch := range f.Batch {
		if a, ok := messageByCrid[batch.CreativeInfo.MediaCid]; ok {
			if a.Result != 0 {
				ret.BatchUploadRetMap[batch.CreativeInfo.MediaCid] = BatchUploadRet{
					ErrCode: model2.CREATIVE_UPLOAD_UNPASSED,
					ErrMsg:  a.Reason,
				}
			} else {
				ret.BatchUploadRetMap[batch.CreativeInfo.MediaCid] = BatchUploadRet{
					ErrCode: model2.CREATIVE_AUDITING,
					ErrMsg:  "审核中",
				}
			}
		} else {
			ret.BatchUploadRetMap[batch.CreativeInfo.MediaCid] = BatchUploadRet{
				ErrCode: model2.CREATIVE_UPLOAD_FAILED,
				ErrMsg:  "媒体未返回该创意",
			}
		}

	}
	return ret
}

func (f *Funshion) sendQueryPost(uri string, crid []string) Ret {
	var ret Ret
	request := FunshionQueryRequest{
		DspId: f.PublisherAccount.DspId,
		Token: f.PublisherAccount.Token,
		Crid:  crid,
	}
	bodyJson, _ := jsoniter.Marshal(request)
	ret.Url = uri
	ret.Req = string(bodyJson)
	response, err := httpclient.PostJSON(uri, bodyJson, httpclient.WithTTL(time.Second*100))
	ret.Resp = string(response)
	if err != nil {
		ret.ErrCode = model2.CREATIVE_QUERY_FAILED
		ret.ErrMsg = err.Error()
		return ret
	}
	var resp *FunshionQueryResp
	err = jsoniter.Unmarshal(response, &resp)
	if err != nil {
		ret.ErrCode = model2.CREATIVE_QUERY_FAILED
		ret.ErrMsg = err.Error()
		return ret
	}
	if resp.Result != FunshionSuccess {
		ret.ErrCode = model2.CREATIVE_QUERY_FAILED
		ret.ErrMsg = string(response)
		return ret
	}
	messageByte, err := jsoniter.Marshal(resp.Message)
	if err != nil {
		ret.ErrCode = model2.CREATIVE_QUERY_FAILED
		ret.ErrMsg = err.Error()
		return ret
	}
	var message FunshionQueryRespMessage
	err = jsoniter.Unmarshal(messageByte, &message)
	if err != nil {
		ret.ErrCode = model2.CREATIVE_UPLOAD_FAILED
		ret.ErrMsg = err.Error()
		return ret
	}
	if len(message.Records) <= 0 {
		ret.ErrCode = model2.CREATIVE_UPLOAD_FAILED
		ret.ErrMsg = "媒体返回为空"
		return ret
	}
	for _, r := range message.Records {
		var queryRet BatchQueryRet
		queryRet.MediaCid = r.Crid
		if r.Result == FunshionQueryStatusPassed {
			queryRet.ErrCode = model2.CREATIVE_AUDIT_PASSED
			queryRet.ErrMsg = "审核通过"
		} else if r.Result == FunshionQueryStatusRefused {
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
