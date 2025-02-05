package admin

import "github.com/gin-gonic/gin"

func UseAdminAuthController(group *gin.RouterGroup) {
	authGroup := group.Group("/auth")
	// 登录
	{
		authGroup.GET("/login", func(ctx *gin.Context) {

		})
	}

	// 注销
	{
		authGroup.GET("/logout", func(ctx *gin.Context) {

		})
	}
}
