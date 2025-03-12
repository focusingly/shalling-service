package admin

import (
	"fmt"
	"space-api/constants"
	"space-api/dto"
	"space-api/internal/service/v1"
	"space-api/middleware/outbound"
	"space-api/util"
	"space-api/util/performance"

	"github.com/gin-gonic/gin"
)

func UseUploadController(group *gin.RouterGroup) {
	uploadFileService := service.DefaultUploadService
	uploadGroup := group.Group("/local-file")
	localFileVisitService := service.DefaultStaticFileService

	// 获取本地文件分页数据
	uploadGroup.GET("/list", func(ctx *gin.Context) {
		req := &dto.GetLocalFilesPaginationReq{}
		if err := ctx.ShouldBindQuery(req); err != nil {
			ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
			return
		}
		if resp, err := uploadFileService.GetLocalFilesByPagination(req, ctx); err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(resp, ctx)
		}
	})

	// 上传普通的文件, 如果携带 MD5 并匹配通过, 那么跳过重复保存
	uploadGroup.POST("/upload", func(ctx *gin.Context) {
		if resp, err := uploadFileService.
			// 临时设置最大上传限制
			Upload(ctx, constants.MB*1024); err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(resp, ctx)
		}
	})

	// 本地文件访问服务
	uploadGroup.GET("/static/visit",
		func(ctx *gin.Context) {
			filename, ok := ctx.GetQuery("file")
			if !ok {
				ctx.Error(util.CreateBizErr("未指定访问的文件", fmt.Errorf("not define any request file name")))
				return
			}

			localFileVisitService.HandleAnyVisit(filename, ctx)
		})

	// 设置图象转码任务的并行度, 防止 webp 转码过度消耗 CPU
	tokenLimitChan := make(chan performance.Empty, 1)
	// 上传图片文件并转码为 webp
	uploadGroup.POST("/upload/webp",
		func(ctx *gin.Context) {
			tokenLimitChan <- performance.Empty{}
			defer func() {
				<-tokenLimitChan
			}()
			ctx.Next()
		},
		func(ctx *gin.Context) {
			if resp, err := uploadFileService.UploadImage2Webp(ctx); err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceResponse(resp, ctx)
			}
		},
	)
}
