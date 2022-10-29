package tests

import (
	"github.com/convee/adcreative/configs"
)

func Init() {
	if configs.Conf.App.Env == "local" {
		//new(crons.CreativeQuery).QueryOne("41084511")
		//new(crons.CreativeUpload).UploadOne("23428")
	}
}
