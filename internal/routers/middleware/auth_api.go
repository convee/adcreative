package middleware

import (
	"fmt"
	"github.com/convee/adcreative/internal/pkg/cache"
	logger "github.com/convee/adcreative/pkg/log"
	"go.uber.org/zap"
	"net/http"

	"github.com/convee/adcreative/internal/pkg/app"
	"github.com/convee/adcreative/internal/pkg/code"
	"github.com/convee/adcreative/pkg/md5"
	"github.com/gin-gonic/gin"
)

type AuthBind struct {
	Appid     string `header:"appid" validate:"required"`
	Secret    string `header:"secret" validate:"required"`
	Sign      string `header:"sign" validate:"required"`
	Timestamp int    `header:"timestamp" validate:"required"`
}

func ApiAuth() (g gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			appG     = app.Gin{C: c}
			authBind AuthBind
		)
		validateErr := app.BindHeader(c, &authBind)
		if len(validateErr) > 0 {
			appG.Response(http.StatusOK, code.INVALID_PARAMS, validateErr)
			c.Abort()
			return
		}
		str := md5.New().Encrypt(fmt.Sprintf("%s%s%d", authBind.Appid, authBind.Secret, authBind.Timestamp))
		customer, err := cache.GetCustomerCache(authBind.Appid)
		if err != nil {
			// 没有找到客户记录
			logger.Error("empty_customer", zap.Error(err))
			appG.Response(http.StatusOK, code.CUSTOMER_EMPTY, nil)
			c.Abort()
			return
		}
		if customer.Secret != authBind.Secret || str != authBind.Sign {
			// 验证签名失败
			appG.Response(http.StatusOK, code.API_SIGN_FAILED, nil)
			c.Abort()
			return
		}
		c.Set("customer", customer)
		c.Next()
		return
	}
}
