package configs

import (
	"github.com/convee/adcreative/internal/enum"
	"github.com/convee/adcreative/pkg/log"
	"github.com/convee/adcreative/pkg/redis"
	"github.com/convee/adcreative/pkg/storage/orm"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

// Config global config
type Config struct {
	// common
	App  AppConfig
	Cron CronConfig
	// component config
	Logger log.Config
	ORM    orm.Config
	Redis  redis.Config
}

// AppConfig app config
type AppConfig struct {
	Name        string
	Version     string
	Mode        string
	Addr        string
	Host        string
	Resource    string
	FfprobePath string
	Env         string
}

// CronConfig cron config
type CronConfig struct {
	Push bool
}

var (
	// Conf app global config
	Conf                = &Config{}
	UploadBatchSizeConf = make(map[int]int)
	QueryBatchSizeConf  = make(map[int]int)

	UploadChanSizeConf = make(map[int]int)
	QueryChanSizeConf  = make(map[int]int)

	UploadConcurrenceLimitConf = make(map[int]int)
	QueryConcurrenceLimitConf  = make(map[int]int)
)

func Init(configPath string) *Config {
	viper.SetConfigType("yml")
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&Conf); err != nil {
		panic(err)
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		if err := viper.Unmarshal(&Conf); err != nil {
			panic(err)
		}
	})
	return Conf
}

func InitPubConf() {
	// 媒体批量送审和查询数量, 为0时单个送审
	uploadBatchSize := viper.GetStringMapString("pub.UploadBatchSize")
	queryBatchSize := viper.GetStringMapString("pub.QueryBatchSize")

	// 媒体创意协程数上限
	uploadChanSize := viper.GetStringMapString("pub.UploadChanSize")
	queryChanSize := viper.GetStringMapString("pub.QueryChanSize")

	// 媒体并发数
	uploadConcurrenceLimit := viper.GetStringMapString("pub.UploadConcurrenceLimit")
	queryConcurrenceLimit := viper.GetStringMapString("pub.QueryConcurrenceLimit")

	for _, pubId := range enum.PubList {
		// 媒体批量送审和查询数量, 为0时单个送审
		UploadBatchSizeConf[pubId] = cast.ToInt(uploadBatchSize["default"])
		QueryBatchSizeConf[pubId] = cast.ToInt(queryBatchSize["default"])

		// 媒体创意协程数上限
		UploadChanSizeConf[pubId] = cast.ToInt(uploadChanSize["default"])
		QueryChanSizeConf[pubId] = cast.ToInt(queryChanSize["default"])

		// 媒体并发数
		UploadConcurrenceLimitConf[pubId] = cast.ToInt(uploadConcurrenceLimit["default"])
		QueryConcurrenceLimitConf[pubId] = cast.ToInt(queryConcurrenceLimit["default"])

		if ubs, ok := uploadBatchSize[cast.ToString(pubId)]; ok {
			UploadBatchSizeConf[pubId] = cast.ToInt(ubs)
		}
		if qbs, ok := queryBatchSize[cast.ToString(pubId)]; ok {
			QueryBatchSizeConf[pubId] = cast.ToInt(qbs)
		}
		if ucs, ok := uploadChanSize[cast.ToString(pubId)]; ok {
			UploadChanSizeConf[pubId] = cast.ToInt(ucs)
		}
		if qcs, ok := queryChanSize[cast.ToString(pubId)]; ok {
			QueryChanSizeConf[pubId] = cast.ToInt(qcs)
		}
		if ucl, ok := uploadConcurrenceLimit[cast.ToString(pubId)]; ok {
			UploadConcurrenceLimitConf[pubId] = cast.ToInt(ucl)
		}
		if qcl, ok := queryConcurrenceLimit[cast.ToString(pubId)]; ok {
			QueryConcurrenceLimitConf[pubId] = cast.ToInt(qcl)
		}
	}
}
