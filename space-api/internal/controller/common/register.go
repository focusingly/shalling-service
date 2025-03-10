package common

import (
	"github.com/gin-gonic/gin"
)

func RegisterAllCommonControllers(group *gin.RouterGroup) {
	commonGroup := group.Group("/common")
	UseHealthCheckController(commonGroup)
	UseAuthController(commonGroup)
}
