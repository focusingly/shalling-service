package controller

import (
	"space-api/internal/controller/admin"
	"space-api/internal/controller/common"

	"github.com/gin-gonic/gin"
)

func RegisterAllControllers(routeGroup *gin.RouterGroup) {
	admin.RegisterAllAdminControllers(routeGroup)
	common.RegisterAllCommonControllers(routeGroup)
}
