package admin

import (
	"space-api/dto"
	"space-api/internal/service/v1"
	"space-api/middleware/outbound"
	"space-api/util"

	"github.com/gin-gonic/gin"
)

func UsePostController(group *gin.RouterGroup) {
	postGroup := group.Group("/posts")
	postService := service.DefaultPostService

	// 增加或者修改文章信息
	{
		postGroup.POST("/", func(ctx *gin.Context) {
			updatePostReq := &dto.UpdateOrCreatePostReq{}
			if err := ctx.BindJSON(updatePostReq); err != nil {
				ctx.Error(&util.BizErr{
					Reason: err,
					Msg:    "参数错误: " + err.Error(),
				})
				return
			}
			if resp, err := postService.CreateOrUpdatePost(updatePostReq, ctx); err != nil {
				ctx.Error(&util.BizErr{
					Reason: err,
					Msg:    "操作失败",
				})
			} else {
				outbound.NotifyProduceResponse(resp, ctx)
			}
		})
	}

	// 根据条件查询分页列表数据
	{
		postGroup.GET("/list", func(ctx *gin.Context) {
			req := &dto.GetPostPageListReq{}
			if err := ctx.ShouldBindQuery(req); err != nil {
				ctx.Error(&util.BizErr{
					Reason: err,
					Msg:    "参数错误 :" + err.Error(),
				})
				return
			}

			if resp, err := postService.GetPostList(req, ctx); err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceResponse(resp, ctx)
			}

		})
	}

	// 根据 ID 查询单条文章的详细数据
	{
		postGroup.GET("/:id", func(ctx *gin.Context) {
			req := &dto.GetPostDetailReq{}
			if err := ctx.ShouldBindUri(req); err != nil {
				ctx.Error(&util.BizErr{
					Msg:    "非法的请求参数: " + err.Error(),
					Reason: err,
				})
				return
			}

			if resp, err := postService.GetPostById(req, ctx); err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceResponse(resp, ctx)
			}
		})
	}

	// 根据 ID 列表删除数据
	{
		postGroup.POST("/deletes", func(ctx *gin.Context) {
			req := &dto.DeletePostByIdListReq{}
			if err := ctx.BindJSON(req); err != nil {
				ctx.Error(&util.BizErr{
					Reason: err,
					Msg:    "参数错误: " + err.Error(),
				})

				return
			}
			if resp, err := postService.DeletePostByIdList(req, ctx); err != nil {
				ctx.Error(&util.BizErr{
					Msg:    "删除失败: " + err.Error(),
					Reason: err,
				})
				return
			} else {
				outbound.NotifyProduceResponse(resp, ctx)
			}
		})
	}
}
