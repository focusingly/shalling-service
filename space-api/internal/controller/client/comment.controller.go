package client

import (
	"space-api/dto"
	"space-api/internal/service/v1"
	"space-api/middleware/auth"
	"space-api/middleware/outbound"
	"space-api/util"
	"time"

	"github.com/gin-gonic/gin"
)

func UseCommentController(routeGroup *gin.RouterGroup) {
	commentsGroup := routeGroup.Group("/comments")
	commentService := service.DefaultCommentService

	// 根评论分页公开查询
	{
		// 对根分页进行缓存
		var rootCachedGroup = &util.Group[*dto.GetRootCommentPagesResp]{}
		commentsGroup.GET("/list", func(ctx *gin.Context) {
			cachedKey := ctx.Request.RequestURI

			resp, _, err := rootCachedGroup.Do(
				cachedKey,
				func() (value *dto.GetRootCommentPagesResp, err error) {
					req := &dto.GetRootCommentPagesReq{}
					if err = ctx.ShouldBindQuery(req); err != nil {
						ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
						return
					}
					value, err = commentService.GetVisibleRootCommentPages(req, ctx)
					return
				},
				time.Millisecond*800,
			)

			if err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceResponse(resp, ctx)
			}
		})
	}

	// 子分页评论公开查询
	{
		// 对子分页进行缓存
		var subCachedGroup = &util.Group[*dto.GetSubCommentPagesResp]{}
		commentsGroup.GET("/list/sub", func(ctx *gin.Context) {
			cachedKey := ctx.Request.RequestURI

			resp, _, err := subCachedGroup.Do(
				cachedKey,
				func() (value *dto.GetSubCommentPagesResp, err error) {
					req := &dto.GetSubCommentPagesReq{}
					if err = ctx.ShouldBindQuery(req); err != nil {
						ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
						return
					}
					value, err = commentService.GetVisibleSubCommentPages(req, ctx)
					return
				},
				time.Millisecond*800,
			)

			if err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceResponse(resp, ctx)
			}
		})
	}

	// 创建评论
	{
		commentsGroup.POST("/",
			// TODO 暂时仅限登录用户进行评论
			func(ctx *gin.Context) {
				if _, err := auth.GetCurrentLoginSession(ctx); err != nil {
					ctx.Error(err)
					ctx.Abort()
					return
				}
				ctx.Next()
			},
			func(ctx *gin.Context) {
				req := &dto.CreateCommentReq{}
				if err := ctx.ShouldBindQuery(req); err != nil {
					ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
					return
				}
				if resp, err := commentService.SimpleVerifyAndCreateComment(req, ctx); err != nil {
					ctx.Error(err)
				} else {
					outbound.NotifyProduceResponse(resp, ctx)
				}
			})
	}
}
