package admin

import (
	"space-api/dto"
	"space-api/internal/service/v1/nft"
	"space-api/middleware/outbound"
	"space-api/util"

	"github.com/gin-gonic/gin"
)

func UseNftController(group *gin.RouterGroup) {
	nftService := nft.DefaultNftService
	nftGroup := group.Group("/nft")

	// 获取封禁列表
	nftGroup.GET("/list", func(ctx *gin.Context) {
		if resp, err := nftService.GetBlockList(); err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(resp, ctx)
		}
	})

	// 添加封禁 IP
	nftGroup.POST("/update", func(ctx *gin.Context) {
		req := &dto.AddNftBanIPReq{}
		if err := ctx.ShouldBindBodyWithJSON(req); err != nil {
			ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
			return
		}

		if err := nftService.AddBlockIPList(req.IPList); err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(&dto.AddNftBanIPResp{}, ctx)
		}
	})

	// 解禁 IP
	nftGroup.POST("/delete", func(ctx *gin.Context) {
		req := &dto.UnbindNftIPReq{}
		if err := ctx.ShouldBindBodyWithJSON(req); err != nil {
			ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
			return
		}
		if err := nftService.RemoveBlockIPList(req.IPList); err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(&dto.AddNftBanIPResp{}, ctx)
		}
	})
}
