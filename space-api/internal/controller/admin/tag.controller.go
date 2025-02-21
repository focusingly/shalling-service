package admin

import (
	"space-api/dto"
	"space-api/internal/service/v1"
	"space-api/middleware/outbound"
	"space-api/util"

	"github.com/gin-gonic/gin"
)

func UseTagController(group *gin.RouterGroup) {
	tagGroup := group.Group("/tag")
	service := service.DefaultTagService

	// 获取标签分页列表
	{
		tagGroup.GET("/list", func(ctx *gin.Context) {
			req := &dto.GetTagPageListReq{}
			err := ctx.BindQuery(req)
			if err != nil {
				ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
				return
			}
			if resp, err := service.GetAnyTagPages(req, ctx); err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceResponse(resp, ctx)
			}
		})
	}

	// 根据 ID 获取单独的标签信息
	{
		tagGroup.GET("/:id", func(ctx *gin.Context) {
			req := &dto.GetTagDetailReq{}
			if err := ctx.ShouldBindUri(req); err != nil {
				ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
				return
			}

			if resp, err := service.GetTagDetailById(req, ctx); err != nil {
				ctx.Error(err)
				return
			} else {
				outbound.NotifyProduceResponse(resp, ctx)
			}
		})
	}

	// 更新/创建标签
	{
		tagGroup.POST("/", func(ctx *gin.Context) {
			req := &dto.CreateOrUpdateTagReq{}
			if err := ctx.BindJSON(req); err != nil {
				ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
				return
			}

			if resp, err := service.CreateOrUpdateTag(req, ctx); err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceResponse(resp, ctx)
			}
		})
	}

	// 根据指定的 ID 列表删除标签
	{
		tagGroup.DELETE("/", func(ctx *gin.Context) {
			req := &dto.DeleteTagByIdListReq{}

			if err := ctx.ShouldBindBodyWithJSON(req); err != nil {
				ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
				return
			}
			if resp, err := service.DeleteTagByIdList(req, ctx); err != nil {
				ctx.Error(err)
				return
			} else {
				outbound.NotifyProduceResponse(resp, ctx)
			}
		})
	}
}
