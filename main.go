package main

import (
	"context"
	"github.com/convee/adcreative/internal/crons"
	"github.com/convee/adcreative/internal/enum"
	"github.com/convee/adcreative/internal/model"
	"github.com/convee/adcreative/internal/service"
	"github.com/convee/adcreative/tests"
	"log"
	"net/http"
	"time"

	"github.com/convee/adcreative/configs"
	logger "github.com/convee/adcreative/pkg/log"
	"github.com/convee/adcreative/pkg/redis"
	"github.com/convee/adcreative/pkg/shutdown"

	"github.com/convee/adcreative/internal/routers"
	"github.com/convee/adcreative/internal/task"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
)

var (
	cfgFile = pflag.StringP("config", "c", "./configs/dev.yml", "config file path.")
	version = pflag.BoolP("version", "v", false, "show version info.")
)

func main() {
	pflag.Parse()
	if *version {
		log.Println("version:", "v1.0")
	}
	// init config
	cfg := configs.Init(*cfgFile)
	// init pub config
	configs.InitPubConf()
	// init logger
	logger.Init(&cfg.Logger)
	// init redis
	redis.Init(&cfg.Redis)
	// init mysql
	model.Init(&cfg.ORM)

	// init Publisher info
	service.Init()

	gin.SetMode(cfg.App.Mode)

	log.Println("http server startup", cfg.App.Addr)
	logger.Info("http server startup")

	srv := &http.Server{
		Addr:    cfg.App.Addr,
		Handler: routers.InitRouter(),
	}
	// 初始化创意送审任务
	task.InitCreativeUploadTask()
	task.InitCreativeQueryTask()
	// 初始化广告主送审任务
	task.InitAdvertiserUploadTask()
	task.InitAdvertiserQueryTask()

	// 广告主送审
	log.Println("advertiser consumer work.")
	go task.AdvUploadConsumer.Work()
	// 广告主审核状态查询
	go task.AdvQueryConsumer.Work()

	log.Println("creative consumer work.")
	for _, v := range enum.PubList {
		// 素材送审任务
		go task.CrUploadConsumer[v].Work()
		// 素材查询任务
		go task.CrQueryConsumer[v].Work()
	}

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	//定时执行任务
	crons.Init()

	// 测试脚本
	tests.Init()

	// 优雅关闭
	shutdown.NewHook().Close(
		// 关闭广告主上传和查询任务
		func() {
			task.AdvUploadProducer.Stop()
			task.AdvQueryProducer.Stop()
			log.Printf("advertiser producer has stoped.")
		},
		// 关闭创意上传和查询任务
		func() {
			task.CrUploadProducer.Stop()
			task.CrQueryProducer.Stop()
			log.Printf("creative producer has stoped.")
		},
		// 关闭 http server
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			if err := srv.Shutdown(ctx); err != nil {
				log.Println("http server closed err", err)
			} else {
				log.Println("http server closed")
			}
		},
	)

}
