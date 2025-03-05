package admin

import (
	"space-api/middleware/auth"

	"github.com/gin-gonic/gin"
)

func RegisterAllAdminControllers(group *gin.RouterGroup) {
	adminGroup := group.Group(
		"/admin",
		auth.UseAdminAuthMiddleware(),
	)

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
	UseNftController(adminGroup)
	UseMailController(adminGroup)
	UseS3Controller(adminGroup)
	UseUVController(adminGroup)
}
