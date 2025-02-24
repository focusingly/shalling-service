package admin

import (
	"space-api/dto"
	"space-api/internal/service/v1/comment"
	"space-api/middleware/outbound"
	"space-api/util"

	"github.com/gin-gonic/gin"
)

func UseCommentController(group *gin.RouterGroup) {
	cmtGroup := group.Group("/comment")
	commentService := comment.DefaultCommentService

	// 查询根评论分页
	{
		cmtGroup.GET("/list/root", func(ctx *gin.Context) {
			req := &dto.GetRootCommentPagesReq{}
			if err := ctx.ShouldBindBodyWithJSON(req); err != nil {
				ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
				return
			}
			if resp, err := commentService.GetAnyRootCommentPages(req, ctx); err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceResponse(resp, ctx)
			}
		})
	}

	// 管理员发表后台直接发表评论
	{
		cmtGroup.POST("/pub", func(ctx *gin.Context) {
			req := &dto.CreateCommentReq{}
			if err := ctx.ShouldBindQuery(req); err != nil {
				ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
				return
			}
			if resp, err := commentService.CheckAndCreateComment(req, ctx); err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceResponse(resp, ctx)
			}
		})
	}

	// 修改评论
	{
		cmtGroup.POST("/edit", func(ctx *gin.Context) {
			req := &dto.UpdateCommentReq{}
			if err := ctx.ShouldBindQuery(req); err != nil {
				ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
				return
			}
			if resp, err := commentService.UpdateComment(req, ctx); err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceResponse(resp, ctx)
			}
		})
	}

	// 查询子评论分页
	{
		cmtGroup.GET("/list/sub", func(ctx *gin.Context) {
			req := &dto.GetSubCommentPagesReq{}
			if err := ctx.ShouldBindBodyWithJSON(req); err != nil {
				ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
				return
			}
			if resp, err := commentService.GetAnySubCommentPages(req, ctx); err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceResponse(resp, ctx)
			}
		})
	}

	// 删子除评论
	{
		cmtGroup.DELETE("/sub", func(ctx *gin.Context) {
			req := &dto.DeleteSubCommentReq{}
			if err := ctx.ShouldBindBodyWithJSON(req); err != nil {
				ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
				return
			}
			if resp, err := commentService.DeleteSubComments(req, ctx); err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceResponse(resp, ctx)
			}
		})
	}

	// 删除根评论(连带删除所有的子评论)
	{
		cmtGroup.DELETE("/root", func(ctx *gin.Context) {
			req := &dto.DeleteRootCommentReq{}
			if err := ctx.ShouldBindBodyWithJSON(req); err != nil {
				ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
				return
			}
			if resp, err := commentService.DeleteRootComments(req, ctx); err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceResponse(resp, ctx)
			}
		})
	}

}
