package common

import (
	"space-api/middleware/outbound"

	"github.com/gin-gonic/gin"
)

func UseHealthCheckController(routeGroup *gin.RouterGroup) {
	routeGroup.GET("/ping", func(ctx *gin.Context) {
		outbound.NotifyProduceResponse("ok", ctx)
	})
}
