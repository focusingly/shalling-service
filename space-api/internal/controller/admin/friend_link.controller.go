package admin

import (
	"space-api/dto"
	"space-api/internal/service/v1"
	"space-api/middleware/outbound"
	"space-api/util"

	"github.com/gin-gonic/gin"
)

func UseFriendLinkController(group *gin.RouterGroup) {
	friendLinkService := service.DefaultFriendLinkService
	friendLinkGroup := group.Group("/friend")

	// 获取所有的友情信息列表信息
	friendLinkGroup.GET("/", func(ctx *gin.Context) {
		req := &dto.GetFriendLinksReq{}
		if err := ctx.ShouldBindQuery(req); err != nil {
			ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
			return
		}
		if resp, err := friendLinkService.GetVisibleFriendLinks(req, ctx); err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(resp, ctx)
		}
	})

	friendLinkGroup.POST("/", func(ctx *gin.Context) {
		req := &dto.CreateOrUpdateFriendLinkReq{}

		if err := ctx.ShouldBindBodyWithJSON(req); err != nil {
			ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
			return
		}

		if resp, err := friendLinkService.CreateOrUpdateFriendLink(req, ctx); err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(resp, ctx)
		}
	})

	friendLinkGroup.DELETE("/", func(ctx *gin.Context) {
		req := &dto.DeleteFriendLinkReq{}
		if err := ctx.ShouldBindBodyWithJSON(req); err != nil {
			ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
			return
		}

		resp, err := friendLinkService.DeleteFriendLinkByIDList(req, ctx)
		if err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(resp, ctx)
		}
	})
}
