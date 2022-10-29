package media

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/convee/adcreative/internal/model"
	"github.com/convee/adcreative/internal/pkg/common"
	"github.com/convee/adcreative/internal/service"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

type Fancy struct {
	*BaseInfo
}

type FancySignSortJson struct {
	DspId        string              `json:"dspId"`
	ExpireTime   int                 `json:"expireTime"`
	MaterialList []FancyMaterialList `json:"materialList"`
	Timestamp    string              `json:"timestamp"`
}

type FancyMaterialList struct {
	AdvertiserId int    `json:"advertiserId"`
	HtmlContent  string `json:"htmlContent"`
	ImgUrl       string `json:"imgUrl"`
	LinkUrl      string `json:"linkUrl"`
	MaterialType int    `json:"materialType"`
	OrderCode    string `json:"orderCode"`
	Platform     int    `json:"platform"`
	TemplateId   int    `json:"templateId"`
	ViewUrl      string `json:"viewUrl"`
}

type FancyUploadResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	AdId    int32  `json:"ad_id"`
}

type FancyAdvertiserUploadResp struct {
	Code         int    `json:"code"`
	Message      string `json:"message"`
	AdvertiserId int32  `json:"advertiser_id"`
}

type FancyResultUpload struct {
	DepositId string `json:"depositId"`
}

type FancyQueryResp struct {
	Code         int            `json:"code"`
	Message      string         `json:"message"`
	Audit        string         `json:"audit"`
	AuditMessage string         `json:"audit_message"`
	Data         FancyQueryData `json:"data"`
}

type FancyAdvertiserQueryResp struct {
	Code         int                      `json:"code"`
	Message      string                   `json:"message"`
	Audit        string                   `json:"audit"`
	AuditMessage string                   `json:"audit_message"`
	Data         FancyAdvertiserQueryData `json:"data"`
}

type FancyAdvertiserQueryData struct {
}

type FancyQueryData struct {
}

func NewFancyHandler(b *BaseInfo) MediaHandler {
	return &Fancy{BaseInfo: b}
}
func (f *Fancy) UploadAdvertiser() Ret {
	var advertiserAuditInfo *service.AdvertiserAuditInfo
	_ = jsoniter.Unmarshal([]byte(f.AdvertiserAudit.Info), &advertiserAuditInfo)
	file := map[string]string{}
	for _, qualification := range advertiserAuditInfo.Qualifications {
		file["file_url"] = qualification.FileUrl
	}
	params := map[string]string{
		"dsp_id":                     f.PublisherAccount.DspId,
		"name":                       f.AdvertiserAudit.AdvertiserName,
		"web_site":                   advertiserAuditInfo.AuthorizeState,
		"industry_id":                advertiserAuditInfo.Industry,
		"business_licence_file_path": file["file_url"],
	}
	str := JoinStringsInASCII(params, "&", false, false, "sign")
	params["signature"] = f.fancySign(str + f.PublisherAccount.Token)
	var ret Ret
	if f.CreativeInfo.IsRsync == 1 {
		ret = f.sendAdvertiserPost(f.AdvertiserUrls.UpdateUrl, params)
	} else {
		ret = f.sendAdvertiserPost(f.AdvertiserUrls.CreateUrl, params)
	}
	return ret
}

func (f *Fancy) sendAdvertiserPost(uri string, params map[string]string) Ret {
	var ret Ret
	bodyJson, _ := json.Marshal(params)

	client := &http.Client{}
	req, err := http.NewRequest("POST", uri, strings.NewReader(string(bodyJson)))
	if err != nil {
		ret.ErrCode = model.ADVERTISER_UPLOAD_FAILED
		ret.ErrMsg = err.Error()
		f.Logger.Error("advertiser_upload_failed", zap.Int("err_code", ret.ErrCode), zap.String("err_msg", ret.ErrMsg))
		return ret
	}
	req.Header.Set("Content-type", "application/json")
	resps, err := client.Do(req)
	if err != nil {
		ret.ErrCode = model.ADVERTISER_UPLOAD_FAILED
		ret.ErrMsg = err.Error()
		f.Logger.Error("advertiser_upload_failed", zap.Int("err_code", ret.ErrCode), zap.String("err_msg", ret.ErrMsg))
		return ret
	}
	defer resps.Body.Close()
	response, err := ioutil.ReadAll(resps.Body)
	if err != nil {
		ret.ErrCode = model.ADVERTISER_UPLOAD_FAILED
		ret.ErrMsg = err.Error()
		f.Logger.Error("advertiser_upload_failed", zap.Any("url", uri), zap.Any("req", bodyJson), zap.Any("resp", string(response)), zap.Error(err))
		return ret
	}
	f.Logger.Info("advertiser_upload_info", zap.Any("url", uri), zap.Any("req", string(bodyJson)), zap.Any("resp", string(response)), zap.Error(err))
	var resp *FancyAdvertiserUploadResp
	err = jsoniter.Unmarshal(response, &resp)
	if err != nil {
		ret.ErrCode = model.ADVERTISER_UPLOAD_FAILED
		ret.ErrMsg = err.Error()
		f.Logger.Error("advertiser_upload_failed", zap.Int("err_code", ret.ErrCode), zap.String("err_msg", ret.ErrMsg))
		return ret
	}
	if resp.Code == 0 && resp.Message == "执行成功" {
		ret.ErrCode = model.ADVERTISER_AUDITING
		ret.MediaCid = cast.ToString(resp.AdvertiserId)
		ret.IsRsync = 1
	} else {
		ret.ErrCode = model.ADVERTISER_UPLOAD_UNPASSED
		ret.ErrMsg = fmt.Sprintf("err_code:%d,err_msg:%s", resp.Code, resp.Message)
		f.Logger.Error("advertiser_upload_unpassed", zap.Int("err_code", ret.ErrCode), zap.String("err_msg", ret.ErrMsg))
	}
	return ret
}

func (f *Fancy) QueryAdvertiser() Ret {
	var ret Ret
	url := f.AdvertiserUrls.QueryUrl
	ret = f.sendAdvertiserQueryPost(url)
	return ret
}

func (f *Fancy) sendAdvertiserQueryPost(url string) Ret {
	var ret Ret
	request := map[string]string{
		"dsp_id":        f.PublisherAccount.DspId,
		"advertiser_id": cast.ToString(f.AdvertiserAudit.MediaCid),
	}
	str := JoinStringsInASCII(request, "&", false, false, "sign")
	request["signature"] = f.fancySign(str + f.PublisherAccount.Token)
	bodyJson, _ := jsoniter.Marshal(request)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(string(bodyJson)))
	if err != nil {
		ret.ErrCode = model.ADVERTISER_UPLOAD_FAILED
		ret.ErrMsg = err.Error()
		f.Logger.Error("advertiser_query_failed", zap.Int("err_code", ret.ErrCode), zap.String("err_msg", ret.ErrMsg))
		return ret
	}
	req.Header.Set("Content-type", "application/json")
	resps, err := client.Do(req)
	if err != nil {
		ret.ErrCode = model.ADVERTISER_UPLOAD_FAILED
		ret.ErrMsg = err.Error()
		f.Logger.Error("advertiser_query_failed", zap.Int("err_code", ret.ErrCode), zap.String("err_msg", ret.ErrMsg))
		return ret
	}
	defer resps.Body.Close()
	response, err := ioutil.ReadAll(resps.Body)
	if err != nil {
		ret.ErrCode = model.ADVERTISER_UPLOAD_FAILED
		ret.ErrMsg = err.Error()
		f.Logger.Error("advertiser_query_failed", zap.Any("url", url), zap.Any("req", bodyJson), zap.Any("resp", string(response)), zap.Error(err))
		return ret
	}
	f.Logger.Info("advertiser_query_info", zap.Any("url", url), zap.Any("req", string(bodyJson)), zap.Any("resp", string(response)), zap.Error(err))
	var resp *FancyAdvertiserQueryResp
	err = jsoniter.Unmarshal(response, &resp)
	if err != nil {
		ret.ErrCode = model.ADVERTISER_QUERY_FAILED
		ret.ErrMsg = err.Error()
		f.Logger.Error("advertiser_query_failed", zap.Int("err_code", ret.ErrCode), zap.String("err_msg", ret.ErrMsg))
		return ret
	}
	if resp.Code != 0 {
		ret.ErrCode = model.ADVERTISER_AUDITING
		ret.ErrMsg = fmt.Sprintf("err_code:%d,err_msg:%s", resp.Code, resp.Message)
	} else {
		if resp.Code == 0 && resp.Audit == "pass" {
			ret.ErrCode = model.ADVERTISER_AUDIT_PASSED
		} else if resp.Code == 0 && resp.Audit == "reject" {
			ret.ErrCode = model.ADVERTISER_AUDIT_UNPASSWD
			ret.ErrMsg = fmt.Sprintf("err_code:%d,err_msg:%s", resp.Code, resp.Message)
			f.Logger.Error("advertiser_query_unpassed", zap.Int("err_code", ret.ErrCode), zap.String("err_msg", ret.ErrMsg))
		} else {
			ret.ErrCode = model.ADVERTISER_AUDITING
			ret.ErrMsg = "审核中"
		}
	}
	return ret
}

func (f *Fancy) UploadCreative() Ret {
	uri := f.CreativeUrls.CreateUrl
	VmJson, _ := jsoniter.Marshal(f.CreativeInfo.Vm)
	CmJson, _ := jsoniter.Marshal(f.CreativeInfo.Cm)
	MonitorJson, _ := jsoniter.Marshal(f.CreativeInfo.Monitor)
	var monitors []string
	monitors = append(monitors, "")

	postData := map[string]string{
		"dsp_id":             f.PublisherAccount.DspId,
		"name":               f.CreativeInfo.Name,
		"landing_page":       f.CreativeInfo.LandUrl,
		"monitor_impression": cast.ToString(VmJson),
		"monitor_click":      cast.ToString(CmJson),
		//"monitor_winnotice":  cast.ToString(MonitorJson),
		"industry_id": f.CreativeInfo.MediaInfo,
	}
	if f.CreativeInfo.Monitor == nil {
		MonitorJson, _ := jsoniter.Marshal(monitors)
		postData["monitor_winnotice"] = cast.ToString(MonitorJson)
	} else {
		postData["monitor_winnotice"] = cast.ToString(MonitorJson)
	}
	if advertiserId, ok := f.Template.Extra["advertiser_id"]; ok {
		postData["advertiser_id"] = advertiserId
	}
	if adType, ok := f.Template.Extra["ad_type"]; ok {
		postData["ad_type"] = adType
	}
	if templateId, ok := f.Template.Extra["template_id"]; ok {
		postData["template_id"] = templateId
	}
	if vendorId, ok := f.Template.Extra["vendor_ids"]; ok {
		var vendorIds []string
		vendorIds = append(vendorIds, vendorId)
		vendorIdsJson, _ := jsoniter.Marshal(vendorIds)
		postData["vendor_ids"] = string(vendorIdsJson)
	}
	materials := make([]map[string]interface{}, 0)
	for _, i := range f.CreativeInfo.Info {
		material := make(map[string]interface{})
		// 素材url
		if common.StringsContain(i.AttrName, "image", "image1", "image2", "video") {
			material["file_url"] = i.AttrValue
			material["file_md5"] = i.Md5
			material["width"] = i.Width
			material["height"] = i.Height
			if i.Duration > 0 {
				material["duration"] = i.Duration
			}
			if fileType, ok := f.Template.Extra["file_type"]; ok {
				material["file_type"] = cast.ToInt(fileType)
			}
			materials = append(materials, material)
		}
	}
	deeplink := map[string]interface{}{
		"url": f.CreativeInfo.DeeplinkUrl,
	}
	if deeplinkUrlType, ok := f.Template.Extra["deeplink_url_type"]; ok {
		deeplink["type"] = cast.ToInt(deeplinkUrlType)
	}
	deeplinkJson, _ := jsoniter.Marshal(deeplink)
	materialJson, _ := jsoniter.Marshal(materials)
	if _, ok := f.Template.Extra["deeplink_url_type"]; ok {
		postData["deeplink"] = string(deeplinkJson)
	}
	postData["material"] = string(materialJson)
	str := JoinStringsInASCII(postData, "&", false, false, "sign")
	postData["signature"] = f.fancySign(str + f.PublisherAccount.Token)
	var ret Ret
	ret = f.sendUploadPost(uri, postData)
	return ret
}

func (f *Fancy) QueryCreative() Ret {
	var ret Ret
	url := f.CreativeUrls.QueryUrl
	ret = f.sendQueryPost(url)
	return ret
}

func (f *Fancy) sendUploadPost(uri string, postData map[string]string) Ret {
	var ret Ret
	bodyJson, _ := json.Marshal(postData)

	client := &http.Client{}
	req, err := http.NewRequest("POST", uri, strings.NewReader(string(bodyJson)))
	if err != nil {
		ret.ErrCode = model.CREATIVE_UPLOAD_FAILED
		ret.ErrMsg = err.Error()
		f.Logger.Error("creative_upload_failed", zap.Int("err_code", ret.ErrCode), zap.String("err_msg", ret.ErrMsg))
		return ret
	}
	req.Header.Set("Content-type", "application/json")
	resps, err := client.Do(req)
	if err != nil {
		ret.ErrCode = model.CREATIVE_UPLOAD_FAILED
		ret.ErrMsg = err.Error()
		f.Logger.Error("creative_upload_failed", zap.Int("err_code", ret.ErrCode), zap.String("err_msg", ret.ErrMsg))
		return ret
	}
	defer resps.Body.Close()
	response, err := ioutil.ReadAll(resps.Body)
	if err != nil {
		ret.ErrCode = model.CREATIVE_UPLOAD_FAILED
		ret.ErrMsg = err.Error()
		f.Logger.Error("creative_upload_failed", zap.Any("url", uri), zap.Any("req", bodyJson), zap.Any("resp", string(response)), zap.Error(err))
		return ret
	}
	f.Logger.Info("creative_upload_info", zap.Any("url", uri), zap.Any("req", string(bodyJson)), zap.Any("resp", string(response)), zap.Error(err))
	var resp *FancyUploadResp
	err = jsoniter.Unmarshal(response, &resp)
	if err != nil {
		ret.ErrCode = model.CREATIVE_UPLOAD_FAILED
		ret.ErrMsg = err.Error()
		f.Logger.Error("creative_upload_failed", zap.Int("err_code", ret.ErrCode), zap.String("err_msg", ret.ErrMsg))
		return ret
	}
	if resp.Code == 0 && resp.Message == "执行成功" {
		ret.ErrCode = model.CREATIVE_AUDITING
		ret.MediaCid = cast.ToString(resp.AdId)
	} else {
		ret.ErrCode = model.CREATIVE_UPLOAD_UNPASSED
		ret.ErrMsg = fmt.Sprintf("err_code:%d,err_msg:%s", resp.Code, resp.Message)
		f.Logger.Error("creative_upload_unpassed", zap.Int("err_code", ret.ErrCode), zap.String("err_msg", ret.ErrMsg))
	}
	return ret
}

func (f *Fancy) sendQueryPost(url string) Ret {
	var ret Ret
	request := map[string]string{
		"dsp_id": f.PublisherAccount.DspId,
		"ad_id":  f.CreativeInfo.MediaCid,
	}
	str := JoinStringsInASCII(request, "&", false, false, "sign")
	request["signature"] = f.fancySign(str + f.PublisherAccount.Token)
	bodyJson, _ := jsoniter.Marshal(request)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(string(bodyJson)))
	if err != nil {
		ret.ErrCode = model.CREATIVE_QUERY_FAILED
		ret.ErrMsg = err.Error()
		f.Logger.Error("creative_query_failed", zap.Int("err_code", ret.ErrCode), zap.String("err_msg", ret.ErrMsg))
		return ret
	}
	req.Header.Set("Content-type", "application/json")
	resps, err := client.Do(req)
	if err != nil {
		ret.ErrCode = model.CREATIVE_QUERY_FAILED
		ret.ErrMsg = err.Error()
		f.Logger.Error("creative_query_failed", zap.Int("err_code", ret.ErrCode), zap.String("err_msg", ret.ErrMsg))
		return ret
	}
	defer resps.Body.Close()
	response, err := ioutil.ReadAll(resps.Body)
	if err != nil {
		ret.ErrCode = model.CREATIVE_QUERY_FAILED
		ret.ErrMsg = err.Error()
		f.Logger.Error("creative_query_failed", zap.Any("url", url), zap.Any("req", bodyJson), zap.Any("resp", string(response)), zap.Error(err))
		return ret
	}
	f.Logger.Info("creative_query_info", zap.Any("url", url), zap.Any("req", string(bodyJson)), zap.Any("resp", string(response)), zap.Error(err))
	var resp *FancyQueryResp
	err = jsoniter.Unmarshal(response, &resp)
	if err != nil {
		ret.ErrCode = model.CREATIVE_QUERY_FAILED
		ret.ErrMsg = err.Error()
		f.Logger.Error("creative_query_failed", zap.Int("err_code", ret.ErrCode), zap.String("err_msg", ret.ErrMsg))
		return ret
	}
	if resp.Code != 0 {
		ret.ErrCode = model.CREATIVE_AUDITING
		ret.ErrMsg = fmt.Sprintf("err_code:%d,err_msg:%s", resp.Code, resp.Message)
	} else {
		if resp.Code == 0 && resp.Audit == "pass" {
			ret.ErrCode = model.CREATIVE_AUDIT_PASSED
		} else if resp.Code == 0 && resp.Audit == "reject" {
			ret.ErrCode = model.CREATIVE_AUDIT_UNPASSWD
			ret.ErrMsg = fmt.Sprintf("err_code:%d,err_msg:%s", resp.Code, resp.Message)
			f.Logger.Error("creative_query_unpassed", zap.Int("err_code", ret.ErrCode), zap.String("err_msg", ret.ErrMsg))
		} else {
			ret.ErrCode = model.CREATIVE_AUDITING
			ret.ErrMsg = "审核中"
		}
	}
	return ret
}

//JoinStringsInASCII 按照规则，参数名ASCII码从小到大排序后拼接
//data 待拼接的数据
//sep 连接符
//onlyValues 是否只包含参数值，true则不包含参数名，否则参数名和参数值均有
//includeEmpty 是否包含空值，true则包含空值，否则不包含，注意此参数不影响参数名的存在
//exceptKeys 被排除的参数名，不参与排序及拼接
func JoinStringsInASCII(data map[string]string, sep string, onlyValues, includeEmpty bool, exceptKeys ...string) string {
	var list []string
	var keyList []string
	m := make(map[string]int)
	if len(exceptKeys) > 0 {
		for _, except := range exceptKeys {
			m[except] = 1
		}
	}
	for k := range data {
		if _, ok := m[k]; ok {
			continue
		}
		value := data[k]
		if !includeEmpty && value == "" {
			continue
		}
		if onlyValues {
			keyList = append(keyList, k)
		} else {
			list = append(list, fmt.Sprintf("%s=%s", k, value))
		}
	}
	if onlyValues {
		sort.Strings(keyList)
		for _, v := range keyList {
			list = append(list, data[v])
		}
	} else {
		sort.Strings(list)
	}
	return strings.Join(list, sep)
}

func (f *Fancy) fancySign(str string) string {
	h := sha1.New()
	h.Write([]byte(str))
	l := fmt.Sprintf("%x", h.Sum(nil))
	return l
}
