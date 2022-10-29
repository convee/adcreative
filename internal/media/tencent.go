package media

import (
	"fmt"
	model2 "github.com/convee/adcreative/internal/model"
	"github.com/convee/adcreative/internal/pkg/common"
	"github.com/convee/adcreative/pkg/ding"
	"github.com/convee/adcreative/pkg/utils"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"

	"github.com/convee/adcreative/pkg/httpclient"
	"github.com/convee/adcreative/pkg/md5"
	jsoniter "github.com/json-iterator/go"
)

const (
	TENCENT_RET_SUCCESSED   = 0
	TENCENT_RET_ALL_FAILED  = 1
	TENCENT_RET_PART_FAILED = 2
	TENCENT_RET_AUTH_FAILED = 3
	TENCENT_STATUS_PASSED   = 1
	TENCENT_STATUS_UNPASSED = 2
	TENCENT_STATUS_AUDITING = 3
	TENCENT_TYP_UPLOAD      = 1
	TENCENT_TYP_QUERY       = 2
)

var (
	tencentConnectAgentCode = []int{2, 200, 1115}
	tencentErrCode          = map[int]string{
		200:  "DSP没有权限操作该广告主",
		300:  "必须的参数没有传入",
		301:  "至少必须传入此组参数中的一个",
		302:  "参数必须为整数",
		303:  "参数格式错误",
		304:  "参数不在允许的值的范围内",
		305:  "参数不是正数",
		306:  "参数不是合法的URL",
		307:  "参数太长，超出了允许的长度范围",
		308:  "参数不是合法的YYYY-MM-DD的日期",
		309:  "参数不是合法的JSON数据，无法被解析",
		310:  "可以被解析但是是空的JSON数据",
		311:  "参数太短，小于允许的长度范围",
		312:  "可以被解析的JSON数据，但不是API要求的JSON格式",
		313:  "参数不是0或正数",
		314:  "参数的值太大，超过了允许的最大值",
		315:  "参数的值不允许被修改",
		316:  "页面暂时不开放",
		400:  "广告主ID不存在/没有匹配的广告主信息",
		500:  "广告主名称重复",
		501:  "广告主ID不存在",
		502:  "广告主不能被删除",
		503:  "广告主不能被修改",
		504:  "广告主行业不合法",
		505:  "广告主URL为空或者是URL不合法",
		506:  "广告主vocation_code不合法",
		507:  "广告主area为空",
		508:  "广告主qualification_class不合法",
		510:  "广告主qualification_files不合法",
		511:  "广告主file_name不合法",
		512:  "广告主file_url不合法",
		513:  "不支持的广告主资质文件的格式，目前支持的文件格式：jpg,jpeg,gif,png",
		515:  "广告主name为空",
		516:  "广告主MEMO为空",
		517:  "广告主同步API在一天内相同的广告主资质信息重复调用,5分钟内不允许重复提交完全相同的广告主资质信息内容",
		518:  "品牌gpb不需要自己上传广告主",
		519:  "广告主还未通过审核",
		529:  "广告主投放资质不可修改",
		601:  "文件加载失败",
		602:  "未知的文件格式",
		603:  "不支持的格式",
		604:  "Flv素材获取不到时长信息",
		605:  "URL对应的素材发生了变化，请换一个URL",
		606:  "执行插入过程中发生了错误，请关注是否是同时上传",
		609:  "素材URL为空或者是地址不合法",
		610:  "目标地址为空或者是地址不合法",
		611:  "必须指定广告主",
		612:  "第三方曝光监测地址错误",
		613:  "素材过大，超过素材的大小限制",
		614:  "传入的file_info格式错误，无法解析成数组",
		615:  "URL对应的广告主发生变化，不能上传",
		616:  "同一次请求中，一个素材URL出现了多次",
		617:  "限制素材上传个数",
		618:  "第三方曝光监测地址不在白名单里",
		619:  "第三方曝光监测数目超出限制",
		620:  "DSP侧的素材ID为空或含有不合法字符或长度超过了64",
		625:  "该DSP不支持投放社交化广告",
		626:  "sns_type参数不合法",
		627:  "社交化广告的文件格式不支持，只支持gif和jpg",
		628:  "素材尺寸不支持社交化广告",
		629:  "素材上传个数达到最大日上限",
		630:  "file_vid含有不合法字符或长度超过了255",
		631:  "display_id不合法",
		632:  "monitor_position不合法",
		633:  "monitor_position的个数与monitor_url的个数不一致",
		634:  "DSP侧的素材ID不允许修改",
		635:  "同一个DSP侧的素材ID下的文件有一些参数的值不相等",
		636:  "file_text参数不合法",
		637:  "DSP侧的素材ID在传入file_text参数的情况下必须也同时传入",
		638:  "在多素材的情况下，由于同一个DSP侧的素材ID下的其他的文件有错误，此文件跳过处理",
		639:  "在多素材的情况下，同一个DSP侧的素材ID下的文件个数超过了允许的最大个数",
		640:  "传入的order_info格式错误，无法解析成数组",
		641:  "同一次请求中，一个dsp_order_id出现了多次",
		642:  "无音频",
		643:  "该广告主不允许上传曝光监测点",
		644:  "第三方点击监测地址错误",
		645:  "第三方点击监测数目超出限制",
		646:  "第三方点击监测地址不在白名单里",
		647:  "flv素材尺寸比例错误",
		648:  "flv素材播放时长错误",
		649:  "ott素材仅支持flv格式",
		650:  "ott素材大小限制16：9&&width>1920",
		651:  "ott素材只支持上传一个flv文件",
		652:  "display_id设置无效",
		653:  "display_id和素材不匹配",
		654:  "display_id和素材不匹配,微动图",
		655:  "display_id和素材不匹配,新闪屏",
		656:  "order接口返回的素材没有通用素材",
		657:  "新闪屏中视频时长不符合规格1-5s",
		658:  "文件地址长度超过数据库字段长度1000",
		659:  "目标地址长度超过数据库字段长度1000",
		700:  "广告位不存在",
		801:  "结束时间大于开始时间",
		802:  "超过开始和结束时间限制",
		803:  "开始和结束时间不能跨年",
		900:  "用户不存在",
		901:  "更新数组为空",
		902:  "用户名称重复",
		903:  "dsp帐号只允许改密码",
		904:  "达到最大账号数",
		1101: "广告创意(dsp_order_id)已存在",
		1102: "广告创意(dsp_order_id)不存在",
		1103: "广告创意修改-广告形式id不能被修改",
		1104: "广告创意修改-广告形式id不存在",
		1105: "广告创意ad_content必填",
		1106: "广告创意ad_content必须是数组",
		1107: "广告创意ad_content必填字段",
		1108: "广告创意ad_content.file_text为空",
		1109: "广告创意ad_content.file_url为空",
		1110: "广告创意ad_content.file_text超长>255",
		1111: "广告创意ad_content.file_url超长>1000",
		1112: "广告创意ad_content.file_url不合法",
		1113: "广告创意ad_content.file_md5必须是32位md5值",
		1114: "广告创意ad_content.file_md5未找到该广告主在系统中匹配项创意，请先上传一个文件，再使用该字段",
		1115: "广告创意ad_content内容和display_id不匹配",
		1116: "广告形式微动图特殊处理失败wdt10711-news_client_msg5-新闻广告主端-信息流GIF-外链",
		1117: "广告形式新闪屏特殊处理失败xsp11267-App_News_Splash",
		1118: "没有数据",
		1119: "广告创意修改-广告形式id不开放",
		1120: "monitor_url监测地址，url解析异常",
		1121: "monitor_url监测地址，采用了HTTPS监测",
		1122: "monitor_url监测地址，域名有误",
		1123: "monitor_url监测地址，非腾讯媒体",
		1124: "monitor_url监测地址，宏参数有问题",
		1125: "monitor_url监测地址，内部接口调用失败",
		1126: "广告创意ad_ext.evokeapp.app_name超长 >255",
		1127: "广告创意ad_ext.evokeapp.pkg_name超长 >255",
		1128: "广告创意ad_ext.evokeapp.deep_link超长 >1000",
		1129: "5分钟内有完全重复的相同内容请求调用接口",
		1130: "不支持阿里unidesk监测代码",
		1131: "音频不符合规范 声音标准(ebu r.128-2011标准)",
		1132: "monitor_settle_bill的个数与monitor_url的个数不一致（PDB指定deal_id时）",
		1133: "monitor_settle_bill 必须指定并且只能指定一条（PDB指定deal_id时）",
		1134: "可见曝光监测地址错误",
		1201: "未找到排期计划",
		1202: "deal信息配置失败",
		1301: "广告位未开放",
		1302: "未找到广告位信息",
	}
)

type Tencent struct {
	*BaseInfo
}

func NewTencentHandler(b *BaseInfo) MediaHandler {
	return &Tencent{BaseInfo: b}
}

type TencentUploadResp struct {
	RetMsg  interface{} `json:"ret_msg"`
	RetCode int         `json:"ret_code"`
	ErrCode int         `json:"error_code,omitempty"`
}

type TencentQueryResp struct {
	RetMsg  interface{} `json:"ret_msg"`
	RetCode int         `json:"ret_code"`
	ErrCode int         `json:"error_code,omitempty"`
}

type TencentExtendParams struct {
	DspOrderId string `json:"dsp_order_id"`
	EndDate    string `json:"end_date"`
}

type TencentUploadRetMsg struct {
	DspOrderId string `json:"dsp_order_id,omitempty"`
	ErrCode    int    `json:"err_code,omitempty"`
	ErrMsg     string `json:"err_msg,omitempty"`
}

type TencentQueryRetMsg struct {
	DspOrderId string `json:"dsp_order_id"`
	Status     int    `json:"status"`
	VInfo      string `json:"vinfo"`
}

type VideoMonitorByTimeUrl struct {
	Url  string `json:"url"`
	Time int    `json:"time"`
}

type TencentQueryExtendRetMsg struct {
	DspOrderId string `json:"dsp_order_id"`
	EndDate    string `json:"end_date"`
	ErrCode    int    `json:"err_code"`
	ErrMsg     string `json:"err_msg"`
}

func (m *Tencent) getUploadParams(creativeInfo CreativeInfo, template model2.Template) map[string]interface{} {
	params := map[string]interface{}{
		"dsp_order_id":    creativeInfo.MediaCid,
		"advertiser_name": creativeInfo.MediaInfo,
		"display_id":      template.DisplayId,
		"targeting_url":   creativeInfo.LandUrl,
		"end_date":        creativeInfo.EndDate,
	}
	monitorPosition, monitorUrl, videoMonitorByTimeUrl := getMonitorInfo(creativeInfo)
	if len(monitorUrl) > 0 {
		params["monitor_url"] = monitorUrl

	}
	if len(videoMonitorByTimeUrl) > 0 {
		params["video_monitor_by_time_url"] = videoMonitorByTimeUrl

	}
	if len(monitorPosition) > 0 {
		params["monitor_position"] = monitorPosition
	}
	if len(creativeInfo.Vm) > 0 {
		params["visible_monitor_url"] = creativeInfo.Vm
	}
	if len(creativeInfo.Cm) > 0 {
		params["click_monitor_url"] = creativeInfo.Cm
	}
	if len(creativeInfo.MiniProgramId) > 0 && len(creativeInfo.MiniProgramPath) > 0 {
		params["mini_program_id"] = creativeInfo.MiniProgramId
		params["mini_program_path"] = creativeInfo.MiniProgramPath
	}
	if len(creativeInfo.DeeplinkUrl) > 0 {
		params["ad_ext"] = map[string]interface{}{
			"app_info": map[string]interface{}{
				"app_id":    1,
				"deep_link": creativeInfo.DeeplinkUrl,
			},
		}
	}
	params["ad_content"] = getAdContent(creativeInfo, template)
	return params
}

func (m *Tencent) UploadCreative() Ret {
	var ret Ret
	params := m.getUploadParams(m.CreativeInfo, m.Template)
	if m.CreativeInfo.IsRsync == 1 {
		ret = m.sendUploadPost(m.CreativeUrls.UpdateUrl, params)
	} else {
		ret = m.sendUploadPost(m.CreativeUrls.CreateUrl, params)
	}
	return ret
}

func (m *Tencent) QueryCreative() Ret {
	params := map[string]interface{}{
		"dsp_order_id": m.CreativeInfo.MediaCid,
	}
	return m.sendQueryPost(m.CreativeUrls.QueryUrl, params)
}

func (m *Tencent) sendUploadPost(uri string, params map[string]interface{}) Ret {
	var ret Ret
	timestamp := time.Now().Unix()
	var data []map[string]interface{}
	data = append(data, params)
	request := map[string]interface{}{
		"dsp_id": m.PublisherAccount.DspId,
		"token":  m.PublisherAccount.Token,
		"time":   timestamp,
		"sig":    m.sign(timestamp),
		"data":   data,
	}
	bodyJson, _ := jsoniter.Marshal(request)
	response, err := httpclient.PostJSON(uri, bodyJson, httpclient.WithTTL(time.Second*180))
	ret.Url = uri
	ret.Req = string(bodyJson)
	ret.Resp = string(response)
	if err != nil {
		ret.ErrCode = model2.CREATIVE_UPLOAD_FAILED
		ret.ErrMsg = err.Error()
		return ret
	}
	var resp *TencentUploadResp
	err = jsoniter.Unmarshal(response, &resp)
	if err != nil {
		ret.ErrCode = model2.CREATIVE_UPLOAD_FAILED
		ret.ErrMsg = err.Error()
		return ret
	}
	if resp.ErrCode == 316 {
		ret.ErrCode = model2.CREATIVE_UPLOADING
		dingMsg := map[string]interface{}{
			"publisher":   "Tencent",
			"method":      "uploadCreative",
			"customer_id": m.CustomerId,
			"creative_id": m.CreativeInfo.CreativeId,
			"media_cid":   m.CreativeInfo.MediaCid,
			"source":      utils.GetHostname(),
			"err_msg":     string(response),
		}
		ding.SendAlert("腾讯送审异常预警", dingMsg, false)
		return ret
	}
	retMsgJson, _ := jsoniter.Marshal(resp.RetMsg)
	var retMsg []TencentUploadRetMsg
	err = jsoniter.Unmarshal(retMsgJson, &retMsg)
	for _, r := range retMsg {
		if r.DspOrderId == params["dsp_order_id"] {
			if r.ErrCode == 0 {
				ret.ErrCode = model2.CREATIVE_AUDITING
			} else if r.ErrCode == 608 || r.ErrCode == 1129 {
				// adcreative add failed7.Pls retry later.
				ret.ErrCode = model2.CREATIVE_UPLOAD_FAILED
				ret.ErrMsg = r.ErrMsg
			} else {
				ret.ErrCode = model2.CREATIVE_UPLOAD_UNPASSED
				ret.ErrMsg = m.errMsgCovert(r.ErrCode, r.ErrMsg, TENCENT_TYP_UPLOAD)
			}
			if resp.RetCode == 0 || resp.RetCode == 2 || (resp.RetCode == 1 && r.ErrCode == 1101) {
				if r.ErrCode == 1101 {
					// 1129 :5分钟内有完全重复的相同内容请求调用接口
					// 1101 广告创意(dsp_order_id)已存在，并发环境下，媒体虽然给出了限流的错误提示，但是好像还是接受到了创意，
					ret.ErrCode = model2.CREATIVE_UPLOAD_FAILED
					ret.ErrMsg = r.ErrMsg
				}

				ret.IsRsync = 1
			}
			return ret
		}
	}
	return ret

}

func (m *Tencent) sendQueryPost(uri string, params map[string]interface{}) Ret {
	var ret Ret
	timestamp := time.Now().Unix()
	var data []map[string]interface{}
	data = append(data, params)
	request := map[string]interface{}{
		"dsp_id": m.PublisherAccount.DspId,
		"token":  m.PublisherAccount.Token,
		"time":   timestamp,
		"sig":    m.sign(timestamp),
		"data":   data,
	}
	bodyJson, _ := jsoniter.Marshal(request)
	response, err := httpclient.PostJSON(uri, bodyJson, httpclient.WithTTL(time.Second*100))
	ret.Url = uri
	ret.Req = string(bodyJson)
	ret.Resp = string(response)
	if err != nil {
		ret.ErrCode = model2.CREATIVE_QUERY_FAILED
		ret.ErrMsg = err.Error()
		return ret
	}
	var resp *TencentQueryResp
	err = jsoniter.Unmarshal(response, &resp)
	if err != nil {
		ret.ErrCode = model2.CREATIVE_QUERY_FAILED
		ret.ErrMsg = err.Error()
		return ret
	}
	if resp.ErrCode != 0 || resp.RetCode != 0 {
		ret.ErrCode = model2.CREATIVE_QUERY_FAILED
		ret.ErrMsg = string(response)
		return ret
	}

	retMsgJson, _ := jsoniter.Marshal(resp.RetMsg)
	var retMsg []TencentQueryRetMsg
	err = jsoniter.Unmarshal(retMsgJson, &retMsg)
	for _, r := range retMsg {
		if r.DspOrderId == params["dsp_order_id"] {
			if r.Status == 1 { //审核通过
				ret.ErrCode = model2.CREATIVE_AUDIT_PASSED
				ret.ErrMsg = "审核通过"
			} else if r.Status == 2 { // 审核未通过
				ret.ErrCode = model2.CREATIVE_AUDIT_UNPASSWD
				ret.ErrMsg = m.errMsgCovert(r.Status, r.VInfo, TENCENT_TYP_QUERY)
			} else {
				ret.ErrCode = model2.CREATIVE_AUDITING
				ret.ErrMsg = "审核中"
			}
			return ret
		}
	}
	return ret

}

func (m *Tencent) sign(timestamp int64) string {
	return md5.New().Encrypt(fmt.Sprintf("%s%s%d", m.PublisherAccount.DspId, m.PublisherAccount.Token, timestamp))
}

func getMonitorInfo(info CreativeInfo) ([]string, []string, []VideoMonitorByTimeUrl) {
	var monitorPosition []string
	var monitorUrl []string
	var videoMonitorByTimeUrl []VideoMonitorByTimeUrl
	if len(info.Monitor) > 0 {
		for _, m := range info.Monitor {
			// 2s 监测
			if m.T == -2 {
				videoMonitorByTimeUrl = append(videoMonitorByTimeUrl, VideoMonitorByTimeUrl{
					Url:  m.Url,
					Time: 2,
				})
				// play 监测
			} else if m.T == -1 {
				videoMonitorByTimeUrl = append(videoMonitorByTimeUrl, VideoMonitorByTimeUrl{
					Url:  m.Url,
					Time: 0,
				})
			} else {
				monitorPosition = append(monitorPosition, strconv.Itoa(m.T))
				monitorUrl = append(monitorUrl, m.Url)
			}
		}
	}
	return monitorPosition, monitorUrl, videoMonitorByTimeUrl

}

func getAdContent(creativeInfo CreativeInfo, template model2.Template) []interface{} {
	var adContent []interface{}
	creativeByKey := make(map[string]model2.TemplateInfo)
	for _, item := range creativeInfo.Info {
		creativeByKey[item.AttrName] = item
	}
	for _, item := range template.Info {
		content := make(map[string]string)
		if file, ok := creativeByKey[item.Key]; ok && !common.StringsContain(item.Key, "deeplink_url", "mini_program_id", "mini_program_path") {
			if strings.HasPrefix(item.Key, "image") || strings.HasPrefix(item.Key, "video") || strings.HasPrefix(item.Key, "icon") || strings.HasPrefix(item.Key, "cover") {
				content["file_url"] = file.AttrValue
			} else {
				content["file_text"] = file.AttrValue
			}
			adContent = append(adContent, content)
		}

	}
	return adContent
}

type TencentQueryCreative struct {
	DspId string                     `json:"dsp_id"`
	Token string                     `json:"token"`
	Time  int64                      `json:"time"`
	Sig   string                     `json:"sig"`
	Data  []TencentQueryCreativeData `json:"data"`
}

type TencentUploadCreative struct {
	DspId string                   `json:"dsp_id"`
	Token string                   `json:"token"`
	Time  int64                    `json:"time"`
	Sig   string                   `json:"sig"`
	Data  []map[string]interface{} `json:"data"`
}

type TencentQueryCreativeData struct {
	DspOrderId string `json:"dsp_order_id"`
}

func (m *Tencent) BatchQueryCreative() Ret {
	var ret Ret
	timestamp := time.Now().Unix()

	request := TencentQueryCreative{
		DspId: m.PublisherAccount.DspId,
		Token: m.PublisherAccount.Token,
		Time:  timestamp,
		Sig:   m.sign(timestamp),
	}

	for _, v := range m.BatchQuery {
		request.Data = append(request.Data, TencentQueryCreativeData{
			DspOrderId: v.Creative.MediaCid,
		})
	}

	bodyJson, _ := jsoniter.Marshal(request)
	response, err := httpclient.PostJSON(m.CreativeUrls.QueryUrl, bodyJson, httpclient.WithTTL(time.Second*100))
	ret.Url = m.CreativeUrls.QueryUrl
	ret.Req = string(bodyJson)
	ret.Resp = string(response)
	if err != nil {
		ret.ErrCode = model2.CREATIVE_QUERY_FAILED
		ret.ErrMsg = err.Error()
		return ret
	}
	var resp *TencentQueryResp
	err = jsoniter.Unmarshal(response, &resp)
	if err != nil {
		ret.ErrCode = model2.CREATIVE_QUERY_FAILED
		ret.ErrMsg = err.Error()
		return ret
	}
	if resp.ErrCode != TENCENT_RET_SUCCESSED || resp.RetCode != TENCENT_RET_SUCCESSED {
		ret.ErrCode = model2.CREATIVE_QUERY_FAILED
		ret.ErrMsg = string(response)
		return ret
	}
	retMsgJson, _ := jsoniter.Marshal(resp.RetMsg)
	var retMsg []TencentQueryRetMsg

	err = jsoniter.Unmarshal(retMsgJson, &retMsg)
	for _, r := range retMsg {
		var queryRet BatchQueryRet
		queryRet.MediaCid = r.DspOrderId
		if r.Status == TENCENT_STATUS_PASSED {
			queryRet.ErrCode = model2.CREATIVE_AUDIT_PASSED
			queryRet.ErrMsg = "审核通过"
		} else if r.Status == TENCENT_STATUS_UNPASSED {
			queryRet.ErrCode = model2.CREATIVE_AUDIT_UNPASSWD
			queryRet.ErrMsg = m.errMsgCovert(r.Status, r.VInfo, TENCENT_TYP_QUERY)
		} else {
			queryRet.ErrCode = model2.CREATIVE_AUDITING
			queryRet.ErrMsg = "审核中"
		}
		ret.BatchQueryRet = append(ret.BatchQueryRet, queryRet)
	}
	return ret
}

func (m *Tencent) BatchUploadCreative() Ret {
	var (
		params []map[string]interface{}
		ret    Ret
	)
	timestamp := time.Now().Unix()

	ret.BatchUploadRetMap = make(map[string]BatchUploadRet)
	for _, b := range m.Batch {
		params = append(params, m.getUploadParams(b.CreativeInfo, b.Template))
	}

	request := TencentUploadCreative{
		DspId: m.PublisherAccount.DspId,
		Token: m.PublisherAccount.Token,
		Time:  timestamp,
		Sig:   m.sign(timestamp),
		Data:  params,
	}

	bodyJson, _ := jsoniter.Marshal(request)
	var uri string
	if m.IsUpdate {
		uri = m.CreativeUrls.UpdateUrl
	} else {
		uri = m.CreativeUrls.CreateUrl
	}
	response, err := httpclient.PostJSON(uri, bodyJson, httpclient.WithTTL(time.Second*100))
	ret.Url = uri
	ret.Req = string(bodyJson)
	ret.Resp = string(response)
	if err != nil {
		ret.ErrCode = model2.CREATIVE_UPLOAD_FAILED
		ret.ErrMsg = err.Error()

		return ret
	}

	var resp *TencentUploadResp
	err = jsoniter.Unmarshal(response, &resp)
	if err != nil {
		ret.ErrCode = model2.CREATIVE_UPLOAD_FAILED
		ret.ErrMsg = err.Error()
		return ret
	}
	if resp.ErrCode == 316 {
		ret.ErrCode = model2.CREATIVE_UPLOADING
		dingMsg := map[string]interface{}{
			"publisher":   "Tencent",
			"method":      "uploadCreative",
			"customer_id": m.CustomerId,
			"source":      utils.GetHostname(),
			"err_msg":     string(response),
		}
		ding.SendAlert("腾讯送审异常预警", dingMsg, false)
		m.Logger.Error("qq_limit", zap.Any("ret", ret))
		return ret
	}

	retMsgJson, _ := jsoniter.Marshal(resp.RetMsg)
	var retMsg []TencentUploadRetMsg
	err = jsoniter.Unmarshal(retMsgJson, &retMsg)

	for _, r := range retMsg {

		uploadRet := BatchUploadRet{}
		if r.ErrCode == 0 {
			uploadRet.ErrCode = model2.CREATIVE_AUDITING
		} else {
			uploadRet.ErrCode = model2.CREATIVE_UPLOAD_UNPASSED
			//ret.ErrMsg = fmt.Sprintf("dsp_order_id:%s,err_code:%d,err_msg:%s", r.DspOrderId, r.ErrCode, r.ErrMsg)
			uploadRet.ErrMsg = m.errMsgCovert(r.ErrCode, r.ErrMsg, TENCENT_TYP_UPLOAD)
		}
		if resp.RetCode == 0 || resp.RetCode == 2 || (resp.RetCode == 1 && r.ErrCode == 1101) {
			if r.ErrCode == 1101 {
				// TODO
				// 并发环境下，媒体虽然给出了限流的错误提示，但是好像还是接受到了创意，
				// 因此这里需要兼容处理，后续可以继续关注该 case
				// 是否需要增加 创意 状态是 待上传 才这样操作 ？？？
				uploadRet.ErrCode = model2.CREATIVE_UPLOADING
				uploadRet.ErrMsg = ""
				dingMsg := map[string]interface{}{
					"publisher":   "Tencent",
					"method":      "uploadCreative",
					"customer_id": m.CustomerId,
					"source":      utils.GetHostname(),
					"ret":         ret,
				}
				ding.SendAlert("dsp_order_id already exists", dingMsg, false)
				m.Logger.Error("dsp_order_id already exists.Pls use update API.", zap.Any("ret", ret))
			}

			uploadRet.IsRsync = 1
		}
		ret.BatchUploadRetMap[r.DspOrderId] = uploadRet
	}

	return ret
}
func (m *Tencent) errMsgCovert(code int, msg string, typ int) string {
	var errMsg = "媒体返回："
	if codeMsg, ok := tencentErrCode[code]; ok && typ == TENCENT_TYP_UPLOAD {
		errMsg += "素材上传失败，" + codeMsg
	} else {
		errMsg += msg
	}

	if common.IntContain(code, tencentConnectAgentCode...) {
		errMsg += "，请联系代理调整素材"
	} else {
		errMsg += "，请联系产品"
	}
	return errMsg

}
