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

type blockingController struct {
	blockingService blocking.INetBlockingService
}

func (s *blockingController) listBLockingList() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := &dto.GetBlockingPagesReq{}
		if err := ctx.ShouldBindQuery(req); err != nil {
			ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
			return
		}
		if resp, err := s.blockingService.GetBlockingPages(req, ctx); err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(resp, ctx)
		}
	}
}

func (s *blockingController) deleteBlockingRecords() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := &dto.DeleteBlockingRecordReq{}
		if err := ctx.ShouldBindQuery(req); err != nil {
			ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
			return
		}
		if resp, err := s.blockingService.DeleteBlockingRecord(req, ctx); err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(resp, ctx)
		}
	}
}

func (s *blockingController) updateBlockingRecord() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := &dto.AddBlockingIPReq{}
		if err := ctx.ShouldBindQuery(req); err != nil {
			ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
			return
		}
		if resp, err := s.blockingService.AddBlockingIP(req, ctx); err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(resp, ctx)
		}
	}
}

func (s *blockingController) syncBlockingRecords() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := s.blockingService.SyncBlockingRecordInCache(context.TODO()); err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(&performance.Empty{}, ctx)
		}
	}
}

func UseBlockingIPController(group *gin.RouterGroup) {
	blockingGroup := group.Group("/ip")
	blockingService := blocking.DefaultIPBlockingService
	controller := &blockingController{
		blockingService: blockingService,
	}

	// 获取 IP 黑名单列表
	blockingGroup.GET("/list", controller.listBLockingList())
	// 删除黑名单中的 IP 记录
	blockingGroup.POST("/delete", controller.deleteBlockingRecords())
	// 添加 IP 黑名单记录
	blockingGroup.POST("/update", controller.updateBlockingRecord())
	// 同步数据
	blockingGroup.POST("/sync", controller.syncBlockingRecords())
}
