package admin

import (
	"space-api/internal/service/v1/monitor"
	"space-api/middleware/outbound"

	"github.com/gin-gonic/gin"
)

func UseMonitorController(group *gin.RouterGroup) {
	perfGroup := group.Group("/monitor")
	pefService := monitor.DefaultMonitorService

	// 查看当前系统的负载情况
	perfGroup.GET("/info", func(ctx *gin.Context) {
		if resp, err := pefService.GetStatus(); err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(resp, ctx)
		}
	})
}
