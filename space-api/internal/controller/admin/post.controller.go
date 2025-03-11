package admin

import (
	"space-api/dto"
	"space-api/internal/service/v1"
	"space-api/middleware/outbound"
	"space-api/util"

	"github.com/gin-gonic/gin"
)

func UsePostController(group *gin.RouterGroup) {
	postService := service.DefaultPostService
	postGroup := group.Group("/posts")

	// 增加或者修改文章信息
	postGroup.POST("/update", func(ctx *gin.Context) {
		updatePostReq := &dto.UpdateOrCreatePostReq{}
		if err := ctx.ShouldBindBodyWithJSON(updatePostReq); err != nil {
			ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
			return
		}
		if resp, err := postService.CreateOrUpdatePost(updatePostReq, ctx); err != nil {
			ctx.Error(util.CreateBizErr(
				"操作失败",
				err,
			))
		} else {
			outbound.NotifyProduceResponse(resp, ctx)
		}
	})

	// 根据条件查询分页列表数据
	postGroup.GET("/list", func(ctx *gin.Context) {
		req := &dto.GetPostPageListReq{}
		if err := ctx.ShouldBindQuery(req); err != nil {
			ctx.Error(util.CreateBizErr("参数错误 :"+err.Error(), err))
			return
		}

		if resp, err := postService.GetAllPostList(req, ctx); err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(resp, ctx)
		}

	})

	// 根据 ID 查询单条文章的详细数据
	postGroup.GET("/detail/:postID", func(ctx *gin.Context) {
		req := &dto.GetPostDetailReq{}
		if err := ctx.ShouldBindUri(req); err != nil {
			ctx.Error(util.CreateBizErr("非法的请求参数: "+err.Error(), err))
			return
		}

		if resp, err := postService.GetAnyPostById(req, ctx); err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(resp, ctx)
		}
	})

	// 根据 ID 列表删除数据
	postGroup.POST("/delete", func(ctx *gin.Context) {
		req := &dto.DeletePostByIdListReq{}
		if err := ctx.BindJSON(req); err != nil {
			ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
			return
		}
		if resp, err := postService.DeletePostByIdList(req, ctx); err != nil {
			ctx.Error(util.CreateBizErr("删除失败: "+err.Error(), err))
			return
		} else {
			outbound.NotifyProduceResponse(resp, ctx)
		}
	})
}
