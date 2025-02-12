package admin

import (
	"space-api/internal/service/v1"
	"space-api/middleware/outbound"

	"github.com/gin-gonic/gin"
)

func UseUploadController(group *gin.RouterGroup) {
	uploadService := service.DefaultUploadService
	uploadGroup := group.Group("/upload")

	uploadGroup.POST("/", func(ctx *gin.Context) {
		if resp, err := uploadService.Upload(ctx); err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(resp, ctx)
		}
	})

	uploadGroup.POST("/webp", func(ctx *gin.Context) {
		if resp, err := uploadService.UploadImage2Webp(ctx); err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(resp, ctx)
		}
	})

}
