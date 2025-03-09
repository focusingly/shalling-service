package admin

import (
	"space-api/dto"
	"space-api/internal/service/v1"
	"space-api/middleware/outbound"

	"space-api/util"

	"github.com/gin-gonic/gin"
)

func UseSocialMediaController(group *gin.RouterGroup) {
	mediaService := service.DefaultMediaService
	mediaGroup := group.Group("/pub-media")

	// 创建/更新 媒体标签
	mediaGroup.POST("/update", func(ctx *gin.Context) {
		req := &dto.CreateOrUpdateSocialMediaReq{}
		if err := ctx.ShouldBindBodyWithJSON(req); err != nil {
			ctx.Error(util.CreateBizErr("参数校验失败: "+err.Error(), err))
			return
		}

		if resp, err := mediaService.CreateOrUpdateMediaTag(req, ctx); err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(resp, ctx)
		}

	})

	// 获取所有媒体标签
	mediaGroup.GET("/list", func(ctx *gin.Context) {
		req := &dto.CreateOrUpdateSocialMediaReq{}
		if err := ctx.ShouldBindBodyWithJSON(req); err != nil {
			ctx.Error(util.CreateBizErr("参数校验失败: "+err.Error(), err))
			return
		}
		if resp, err := mediaService.CreateOrUpdateMediaTag(req, ctx); err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(resp, ctx)
		}
	})

	// 删除标签
	mediaGroup.POST("/delete", func(ctx *gin.Context) {
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
			outbound.NotifyProduceResponse(resp, ctx)
		}
	})

}
