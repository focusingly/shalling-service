package admin

import (
	"space-api/internal/service/v1"
	"space-api/middleware/outbound"

	"github.com/gin-gonic/gin"
)

func UsePerformanceMonitorController(group *gin.RouterGroup) {
	perfGroup := group.Group("/performance")

	pefService := service.DefaultPerformanceService
	{
		perfGroup.GET("/", func(ctx *gin.Context) {
			if resp, err := pefService.GetStatus(); err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceRestJSON(resp, ctx)
			}
		})
	}
}
