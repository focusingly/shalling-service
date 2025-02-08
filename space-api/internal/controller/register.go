package controller

import (
	"space-api/internal/controller/admin"
	"space-api/internal/controller/client"
	"space-api/internal/controller/common"

	"github.com/gin-gonic/gin"
)

func RegisterAllControllers(routeGroup *gin.RouterGroup) {
	admin.RegisterAllAdminControllers(routeGroup)
	client.RegisterAllClientComments(routeGroup)
	common.RegisterAllCommonControllers(routeGroup)
}
