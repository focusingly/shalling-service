package admin

import (
	"space-api/constants"
	"space-api/internal/service/v1"
	"space-api/middleware/outbound"

	"github.com/gin-gonic/gin"
)

func UseUploadController(group *gin.RouterGroup) {
	uploadService := service.DefaultUploadService
	uploadGroup := group.Group("/upload")

	// 上传普通的文件, 如果携带 MD5 并匹配通过, 那么跳过重复保存
	uploadGroup.POST("/", func(ctx *gin.Context) {
		if resp, err := uploadService.
			Upload(ctx, constants.MB*1024); err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(resp, ctx)
		}
	})

	// 上传图片文件并转码为 webp
	uploadGroup.POST("/webp", func(ctx *gin.Context) {
		if resp, err := uploadService.UploadImage2Webp(ctx); err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(resp, ctx)
		}
	})

}
