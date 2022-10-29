package ding

import (
	"github.com/convee/adcreative/configs"
	"testing"
)

func Test_Ding(t *testing.T) {
	configs.Conf.App.Env = "DEV"
	SendAlert("Material-Service", "钉钉预警测试", false)
}
