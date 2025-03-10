package admin

import (
	"space-api/constants"
	"space-api/internal/service/v1"
	"space-api/middleware/outbound"
	"space-api/util/performance"

	"github.com/gin-gonic/gin"
)

func UseUploadController(group *gin.RouterGroup) {
	uploadService := service.DefaultUploadService
	uploadGroup := group.Group("/upload")

	// 上传普通的文件, 如果携带 MD5 并匹配通过, 那么跳过重复保存
	uploadGroup.POST("/common", func(ctx *gin.Context) {
		if resp, err := uploadService.
			// 临时设置最大上传限制
			Upload(ctx, constants.MB*1024); err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(resp, ctx)
		}
	})

	// 设置图象转码任务的并行度, 防止 webp 转码过度消耗 CPU
	tokenLimitChan := make(chan performance.Empty, 1)
	// 上传图片文件并转码为 webp
	uploadGroup.POST(
		"/webp",
		func(ctx *gin.Context) {
			tokenLimitChan <- performance.Empty{}
			defer func() {
				<-tokenLimitChan
			}()
			ctx.Next()
		},
		func(ctx *gin.Context) {
			if resp, err := uploadService.UploadImage2Webp(ctx); err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceResponse(resp, ctx)
			}
		},
	)
}
