package admin

import (
	"space-api/dto"
	"space-api/internal/service/v1/uv"
	"space-api/middleware/outbound"
	"space-api/util"

	"github.com/gin-gonic/gin"
)

func UseUVController(group *gin.RouterGroup) {
	uvService := uv.DefaultUVService
	uvGroup := group.Group("/uv")

	// 获取指定日期的访问数
	uvGroup.GET("/daily", func(ctx *gin.Context) {
		req := &dto.GetDailyCountReq{}
		if err := ctx.ShouldBindQuery(req); err != nil {
			ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
			return
		}
		resp, err := uvService.GetDailyUVCount(req, ctx)
		if err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(resp, ctx)
		}
	})

	// 获取某几日前到今天的趋势
	uvGroup.GET("/trend", func(ctx *gin.Context) {
		req := &dto.GetUVTrendReq{}
		if err := ctx.ShouldBindQuery(req); err != nil {
			ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
			return
		}
		resp, err := uvService.GetUVTrend(req, ctx)
		if err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(resp, ctx)
		}
	})

	// 获取分页
	uvGroup.POST("/list", func(ctx *gin.Context) {
		req := &dto.GetUvPagesReq{}
		if err := ctx.ShouldBindBodyWithJSON(req); err != nil {
			ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
			return
		}
		resp, err := uvService.GetUvPages(req, ctx)
		if err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(resp, ctx)
		}
	})

	// 获取指定条件的独立访问数
	uvGroup.POST("/query/count", func(ctx *gin.Context) {
		req := &dto.QueryUvCountReq{}
		if err := ctx.ShouldBindBodyWithJSON(req); err != nil {
			ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
			return
		}
		resp, err := uvService.QueryRangeUV(req, ctx)
		if err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(resp, ctx)
		}
	})

	// 删除指定条件的独立访问数
	uvGroup.POST("/delete", func(ctx *gin.Context) {
		req := &dto.DeleteUVReq{}
		if err := ctx.ShouldBindBodyWithJSON(req); err != nil {
			ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
			return
		}
		resp, err := uvService.DeleteUVRecord(req, ctx)
		if err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(resp, ctx)
		}
	})
}
