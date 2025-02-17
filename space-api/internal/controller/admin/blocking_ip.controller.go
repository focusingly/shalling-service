package admin

import (
	"context"
	"space-api/dto"
	"space-api/internal/service/v1/blocking"
	"space-api/middleware/outbound"
	"space-api/util"
	"space-api/util/performance"

	"github.com/gin-gonic/gin"
)

func UseBlockingIPController(group *gin.RouterGroup) {
	blockingGroup := group.Group("/ip")
	blockingService := blocking.DefaultIPBlockingService

	// 获取阻塞列表
	{
		blockingGroup.GET("/list", func(ctx *gin.Context) {
			req := &dto.GetBlockingPagesReq{}

			if err := ctx.ShouldBindQuery(req); err != nil {
				ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
				return
			}
			if resp, err := blockingService.GetBlockingPages(req, ctx); err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceResponse(resp, ctx)
			}
		})
	}

	// 删除阻塞数据
	{
		blockingGroup.DELETE("/delete", func(ctx *gin.Context) {
			req := &dto.DeleteBlockingRecordReq{}

			if err := ctx.ShouldBindQuery(req); err != nil {
				ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
				return
			}
			if resp, err := blockingService.DeleteBlockingRecord(req, ctx); err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceResponse(resp, ctx)
			}
		})
	}

	// 添加阻塞数据
	{
		blockingGroup.POST("/", func(ctx *gin.Context) {
			req := &dto.AddBlockingIPReq{}

			if err := ctx.ShouldBindQuery(req); err != nil {
				ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
				return
			}
			if resp, err := blockingService.AddBlockingIP(req, ctx); err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceResponse(resp, ctx)
			}
		})
	}

	// 同步数据
	{
		blockingGroup.POST("/sync", func(ctx *gin.Context) {
			if err := blockingService.SyncBlockingRecordInCache(context.TODO()); err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceResponse(&performance.Empty{}, ctx)
			}
		})
	}
}
