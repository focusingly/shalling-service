package admin

import "github.com/gin-gonic/gin"

func RegisterAllAdminControllers(group *gin.RouterGroup) {
	adminGroup := group.Group("/admin")

	UsePostController(adminGroup)
	UseTagController(adminGroup)
}
