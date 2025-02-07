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

	{
		// 对根分页进行缓存
		var rootCachedGroup = &util.Group[*dto.GetRootCommentPagesResp]{}
		// 获取评论的根分页信息
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
					value, err = commentService.GetRootCommentPages(req, ctx)
					return
				},
				time.Millisecond*800,
			)

			if err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceRestJSON(resp, ctx)
			}
		})
	}

	{
		// 对子分页进行缓存
		var subCachedGroup = &util.Group[*dto.GetSubCommentPagesResp]{}
		// 获取评论的根分页信息
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
					value, err = commentService.GetSubCommentPages(req, ctx)
					return
				},
				time.Millisecond*800,
			)

			if err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceRestJSON(resp, ctx)
			}
		})
	}

	{
		// 创建评论
		commentsGroup.POST("/",
			// TODO 暂时仅限登录用户
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
					outbound.NotifyProduceRestJSON(resp, ctx)
				}
			})
	}
}
