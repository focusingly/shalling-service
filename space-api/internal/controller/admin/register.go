package admin

import "github.com/gin-gonic/gin"

func RegisterAllAdminControllers(group *gin.RouterGroup) {
	adminGroup := group.Group("/admin")

	UsePostController(adminGroup)
	UseCommentController(adminGroup)
	UseTagController(adminGroup)
	UseSocialMediaController(adminGroup)
	UseJobController(adminGroup)
	UseFriendLinkController(adminGroup)
	UseMenuController(adminGroup)
	UseMonitorController(adminGroup)
	UseUserController(adminGroup)
	UseLogController(adminGroup)
	UseUploadController(adminGroup)
	UseBlockingIPController(adminGroup)
	UseMailController(adminGroup)
	UseS3Controller(adminGroup)
	UseUVController(adminGroup)
}
