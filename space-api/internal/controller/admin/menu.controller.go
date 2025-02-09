package admin

import (
	"space-api/dto"
	"space-api/internal/service/v1"
	"space-api/middleware/outbound"
	"space-api/util"

	"github.com/gin-gonic/gin"
)

func UseMenuController(group *gin.RouterGroup) {
	menuService := service.DefaultMenuService
	menuGroup := group.Group("/menu")

	// 获取所有的列表信息
	menuGroup.GET("/", func(ctx *gin.Context) {
		req := &dto.GetMenusReq{}
		if err := ctx.ShouldBindQuery(req); err != nil {
			ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
			return
		}
		if resp, err := menuService.GetAllMenus(req, ctx); err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(resp, ctx)
		}
	})

	menuGroup.POST("/", func(ctx *gin.Context) {
		req := &dto.CreateOrUpdateMenuReq{}

		if err := ctx.ShouldBindBodyWithJSON(req); err != nil {
			ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
			return
		}

		if resp, err := menuService.CreateOrUpdateMenu(req, ctx); err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(resp, ctx)
		}
	})

	menuGroup.DELETE("/", func(ctx *gin.Context) {
		req := &dto.DeleteMenuGroupsReq{}
		if err := ctx.ShouldBindBodyWithJSON(req); err != nil {
			ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
			return
		}

		resp, err := menuService.DeleteMenuGroupByIDList(req, ctx)
		if err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(resp, ctx)
		}
	})
}
