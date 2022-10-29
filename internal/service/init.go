package service

import (
	"github.com/convee/adcreative/internal/model"
)

var (
	pubInfo map[int]string
)

func Init() {
	publishers, err := new(model.PublisherModel).GetAllPublisherNames()
	if err != nil {
		panic(err)
	}
	pubInfo = make(map[int]string)
	for _, pub := range publishers {
		pubInfo[pub.Id] = pub.Name
	}
}

func GetPubName(pubId int) string {
	if pub, ok := pubInfo[pubId]; ok {
		return pub
	}
	return ""
}
