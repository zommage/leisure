/** * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
 *
 * 路由控制器
 * generate by DavidYang 2017.5.23
 *
 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zommage/leisure/controllers/base"
	"github.com/zommage/leisure/controllers/users"
)

var (
	RouteVersionOnePoint = "/leisure/gateway/v1"
)

func ApiRouter(router *gin.Engine) {

	authorized := router.Group("/")
	// 用户鉴权
	authorized.Use(base.AuthRequired())

	version1 := authorized.Group(RouteVersionOnePoint)
	{
		// 用户登录
		version1.POST("/login", users.Login)

		version1.GET("/health", users.HealthCheck)
	}
}
