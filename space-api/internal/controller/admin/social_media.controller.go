package admin

import (
	"space-api/dto"
	"space-api/internal/service/v1"
	"space-api/middleware/outbound"

	"space-api/util"

	"github.com/gin-gonic/gin"
)

func UseSocialMediaController(group *gin.RouterGroup) {
	mediaGroup := group.Group("/social-media")
	mediaService := service.DefaultMediaService

	{
		// 获取单个标签信息
		mediaGroup.GET("/:id", func(ctx *gin.Context) {
			req := &dto.GetSocialMediaDetailReq{}
			if err := ctx.ShouldBindUri(req); err != nil {
				ctx.Error(util.CreateBizErr("参数校验错误: "+err.Error(), err))
				return
			}
			resp, err := mediaService.GetMediaTagDetailById(req, ctx)
			if err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceRestJSON(resp, ctx)
			}
		})

		// 获取列表信息
		mediaGroup.GET("/list", func(ctx *gin.Context) {
			req := &dto.GetSocialMediaPageListReq{}
			if err := ctx.ShouldBindQuery(req); err != nil {
				ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
				return
			}

			resp, err := mediaService.GetMediaTagPages(req, ctx)
			if err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceRestJSON(resp, ctx)
			}
		})

		// 创建/更新 公开的媒体标签
		mediaGroup.POST("/", func(ctx *gin.Context) {
			req := &dto.CreateOrUpdateSocialMediaReq{}
			if err := ctx.ShouldBindBodyWithJSON(req); err != nil {
				ctx.Error(util.CreateBizErr("参数校验失败: "+err.Error(), err))
				return
			}

			if resp, err := mediaService.CreateOrUpdateMediaTag(req, ctx); err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceRestJSON(resp, ctx)
			}

		})

		// 删除标签
		mediaGroup.DELETE("/", func(ctx *gin.Context) {
			req := &dto.DeleteSocialMediaByIdListReq{}
			err := ctx.ShouldBindBodyWithJSON(req)
			if err != nil {
				ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
				return
			}

			resp, err := mediaService.DeleteMediaTagByIdList(req, ctx)
			if err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceRestJSON(resp, ctx)
			}
		})
	}
}
