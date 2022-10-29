package routers

import (
	"github.com/convee/adcreative/internal/handler/api"
	"github.com/convee/adcreative/internal/handler/backend"
	"github.com/convee/adcreative/internal/routers/middleware"
	"github.com/convee/adcreative/pkg/utils"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	backendSystemHandler           = backend.System{}
	backendCustomerHandler         = backend.Customer{}
	backendPublisherAccountHandler = backend.PublisherAccount{}
	backendPublisherHandler        = backend.Publisher{}
	backendPositionHandler         = backend.Position{}
	backendAdvertiserRulesHandler  = backend.AdvertiserRules{}
	backendAdvertiserAuditHandler  = backend.AdvertiserAudit{}
	backendCreativeHandler         = backend.Creative{}
	apiPublisherAccountHandler     = api.PublisherAccount{}
	apiPublisherIndustryHandler    = api.PublisherIndustry{}
	apiAdvertiserRulesHandler      = api.AdvertiserRules{}
	apiCreativeHandler             = api.Creative{}
	apiPublisherHandler            = api.Publisher{}
	apiPositionHandler             = api.Position{}
	apiAdvertiserAuditHandler      = api.AdvertiserAudit{}
	apiStatisticsHandler           = api.Statistics{}
)

type healthCheckResponse struct {
	Status   string `json:"status"`
	Hostname string `json:"hostname"`
}

// HealthCheck will return OK if the underlying BoltDB is healthy. At least healthy enough for demoing purposes.
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, healthCheckResponse{Status: "UP", Hostname: utils.GetHostname()})
}

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()

	// pprof router 性能分析路由
	// 默认关闭，开发环境下可以打开
	// 访问方式: HOST/debug/pprof
	// 通过 HOST/debug/pprof/profile 生成profile
	// 查看分析图 go tool pprof -http=:5000 profile (安装graphviz: brew install graphviz)
	// see: https://github.com/gin-contrib/pprof
	pprof.Register(r)
	r.Use(middleware.RequestID())
	r.Use(middleware.Logging())
	r.Use(middleware.Metrics(nil))
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	// HealthCheck 健康检查路由
	r.GET("/health", HealthCheck)
	// metrics router 可以在 prometheus 中进行监控
	// 通过 grafana 可视化查看 prometheus 的监控数据，使用插件6671查看
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	backendGroup := r.Group("/backend")
	backendGroup.Use(middleware.BackendAuth())
	{
		// 获取用户登录信息
		backendGroup.GET("/system/login_info", backendSystemHandler.LoginInfo)
		// 退出登录
		backendGroup.POST("/system/logout", backendSystemHandler.Logout)

		// 获取客户列表
		backendGroup.GET("/customer/list", backendCustomerHandler.List)
		// 添加客户
		backendGroup.POST("/customer/add", backendCustomerHandler.Add)
		// 编辑客户
		backendGroup.POST("/customer/edit", backendCustomerHandler.Edit)
		// 删除客户
		backendGroup.POST("/customer/delete", backendCustomerHandler.Delete)

		// 获取媒体账号列表
		backendGroup.GET("/publisher_account/list", backendPublisherAccountHandler.List)
		// 添加媒体账号
		backendGroup.POST("/publisher_account/add", backendPublisherAccountHandler.Add)
		// 编辑媒体账号
		backendGroup.POST("/publisher_account/edit", backendPublisherAccountHandler.Edit)
		// 删除媒体账号
		backendGroup.POST("/publisher_account/delete", backendPublisherAccountHandler.Delete)

		// 获取媒体列表
		backendGroup.GET("/publisher/list", backendPublisherHandler.List)
		// 添加媒体
		backendGroup.POST("/publisher/add", backendPublisherHandler.Add)
		// 编辑媒体
		backendGroup.POST("/publisher/edit", backendPublisherHandler.Edit)
		// 删除媒体
		backendGroup.POST("/publisher/delete", backendPublisherHandler.Delete)

		// 获取广告位列表
		backendGroup.GET("/position/list", backendPositionHandler.List)
		// 添加广告位
		backendGroup.POST("/position/add", backendPositionHandler.Add)
		// 编辑广告位
		backendGroup.POST("/position/edit", backendPositionHandler.Edit)
		// 删除广告位
		backendGroup.POST("/position/delete", backendPositionHandler.Delete)

		// 获取广告主规则
		backendGroup.GET("/advertiser_rules/list", backendAdvertiserRulesHandler.List)
		// 添加广告主规则
		backendGroup.POST("/advertiser_rules/add", backendAdvertiserRulesHandler.Add)
		// 编辑广告主规则
		backendGroup.POST("/advertiser_rules/edit", backendAdvertiserRulesHandler.Edit)
		// 删除广告主规则
		backendGroup.POST("/advertiser_rules/delete", backendAdvertiserRulesHandler.Delete)

		// 获取广告主审核列表
		backendGroup.GET("/advertiser_audit/list", backendAdvertiserAuditHandler.List)
		// 获取创意列表
		backendGroup.GET("/creative/list", backendCreativeHandler.List)

	}

	apiGroup := r.Group("/api")
	apiGroup.Use(middleware.ApiAuth())
	{
		// 媒体列表
		apiGroup.GET("/publisher/list", apiPublisherHandler.List)
		// 广告位列表
		apiGroup.GET("/position/list", apiPositionHandler.List)
		// 批量获取广告位列表
		apiGroup.POST("/position/batch", apiPositionHandler.Batch)
		// 广告位json模版
		apiGroup.GET("/position/material", apiPositionHandler.Material)
		// 添加媒体账号
		apiGroup.POST("/publisher_account/add", apiPublisherAccountHandler.Add)
		// 创意上传
		apiGroup.POST("/creative/upload", apiCreativeHandler.Upload)
		// 创意校验
		apiGroup.POST("/creative/check", apiCreativeHandler.Check)
		// 创意校验
		apiGroup.POST("/creative/batch_check", apiCreativeHandler.BatchCheck)
		// 创意状态查询
		apiGroup.POST("/creative/query", apiCreativeHandler.Query)
		// 媒体行业ID列表
		apiGroup.GET("/industry/list", apiPublisherIndustryHandler.List)
		// 广告主列表
		apiGroup.GET("/advertiser/rules", apiAdvertiserRulesHandler.List)
		// 广告主信息保存，用于同步广告主信息
		apiGroup.POST("/advertiser_audit/save", apiAdvertiserAuditHandler.Save)
		// 广告主上传
		apiGroup.POST("/advertiser_audit/upload", apiAdvertiserAuditHandler.Upload)
		// 广告主查询
		apiGroup.POST("/advertiser_audit/query", apiAdvertiserAuditHandler.Query)
		// 素材状态统计
		apiGroup.POST("/material/status", apiStatisticsHandler.Status)
	}

	return r
}
