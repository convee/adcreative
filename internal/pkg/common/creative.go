package common

import (
	"github.com/convee/adcreative/configs"
	"github.com/convee/adcreative/internal/enum"
	"github.com/spf13/cast"
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	// MediaCidPubList 由素材服务生成素材id并且送审到媒体
	MediaCidPubList = []int{
		enum.PUB_TENCENT,
		enum.PUB_UC,
	}
	PubTxtLength = []int{ // 字符长度按照2个英文字符算一个中文字符计算媒体列表
		enum.PUB_TENCENT,
	}
)

// GenMediaCid 带前缀避免重复
func GenMediaCid(pubId, id int) string {
	if configs.Conf.App.Env == "prod" {
		return GenProdGenMediaCid(pubId, id)
	}
	return GenUatGenMediaCid(pubId, id)
}

func GenUatGenMediaCid(pubId, id int) string {
	switch pubId {
	default:
		return "conveeuat" + cast.ToString(id)
	}
}

func GenProdGenMediaCid(pubId, id int) string {
	switch pubId {
	case enum.PUB_UC:
		return cast.ToString(1000000000000 + id)

	default:
		return "convee" + cast.ToString(id)
	}
}

// StringLength 字符长度计算
func StringLength(s string, pubId int) float32 {
	// 中文算一个，英文数字算半个
	if IntContain(pubId, PubTxtLength...) {
		l := cast.ToFloat32(utf8.RuneCountInString(s))
		li := cast.ToFloat32(len(s))
		return l - (l-(li-l)/2)/2
		// 中文2个字符，英文数字1个
	} else if false { //todo 中文2个字符，英文数字1个
		re, _ := regexp.Compile("[^\\x00-\\xff]")
		r := re.ReplaceAllString(s, "AA")
		return cast.ToFloat32(len(r))
	}
	// 中文、数字、英文都算1个
	return cast.ToFloat32(utf8.RuneCountInString(s))
}

func GetActualExt(ext string) string {
	var actualExt string
	if strings.Contains(ext[1:], "?") {
		lowerExt := strings.ToLower(ext[1:])
		comma := strings.Index(lowerExt, "?")
		actualExt = lowerExt[:comma]
	} else if strings.Contains(ext[1:], "#") {
		lowerExt := strings.ToLower(ext[1:])
		comma := strings.Index(lowerExt, "#")
		actualExt = lowerExt[:comma]
	} else {
		actualExt = strings.ToLower(ext[1:])
	}
	return actualExt
}
