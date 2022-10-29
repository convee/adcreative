package middleware

import (
	"github.com/gin-gonic/gin"
)

func BackendAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		// todo 登录信息验证
		var userInfo interface{}

		c.Set("userinfo", userInfo)

		c.Next()
		return
	}
}
