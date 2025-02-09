package admin

import "github.com/gin-gonic/gin"

func RegisterAllAdminControllers(group *gin.RouterGroup) {
	adminGroup := group.Group("/admin")

	UsePostController(adminGroup)
	UseTagController(adminGroup)
	UseSocialMediaController(adminGroup)
	UseJobController(adminGroup)
	UseFriendLinkController(adminGroup)
	UseMenuController(adminGroup)
	UsePerformanceMonitorController(adminGroup)
}
