package client

import (
	"fmt"
	"space-api/constants"
	"space-api/dto"
	"space-api/internal/service/v1/comment"
	"space-api/middleware/auth"
	"space-api/middleware/outbound"
	"space-api/util"
	"space-api/util/performance"
	"time"

	"github.com/gin-gonic/gin"
)

func UseCommentController(routeGroup *gin.RouterGroup) {
	commentsGroup := routeGroup.Group("/comments")
	commentService := comment.DefaultCommentService

	// 根评论分页公开查询
	{
		// 对根分页进行缓存
		var rootCachedGroup = &performance.Group[*dto.GetRootCommentPagesResp]{}
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
		var subCachedGroup = &performance.Group[*dto.GetSubCommentPagesResp]{}
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
		limitCache := performance.NewCache(constants.MB * 2)

		commentsGroup.POST("/",
			func(ctx *gin.Context) {
				// TODO 暂时仅限登录用户进行评论
				loginSession, err := auth.GetCurrentLoginSession(ctx)
				if err != nil {
					ctx.Error(err)
					ctx.Abort()
					return
				}
				// 限制非管理员的言论发表频率
				if loginSession.UserType != constants.Admin {
					if ttl, e := limitCache.GetTTL(fmt.Sprintf("%d", loginSession.ID)); e == nil {
						ctx.Error(util.CreateLimitErr(
							fmt.Sprintf("发言时间限制, 下一条评论发表时间 %d 秒后", ttl),
							fmt.Errorf("post comment limit, left %d seconds", ttl)),
						)
						ctx.Abort()
						return
					}

				}

				ctx.Next()
			},
			func(ctx *gin.Context) {
				req := &dto.CreateCommentReq{}
				if err := ctx.ShouldBindQuery(req); err != nil {
					ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
					return
				}
				if resp, err := commentService.CheckAndCreateComment(req, ctx); err != nil {
					ctx.Error(err)
				} else {
					outbound.NotifyProduceResponse(resp, ctx)
					loginSession, _ := auth.GetCurrentLoginSession(ctx)

					// 设置标记, 限制非管理员评论速率, 一分钟一条
					if loginSession.UserType != constants.Admin {
						limitCache.Set(
							fmt.Sprintf("%d", loginSession.ID),
							&performance.Empty{}, performance.Second(time.Minute*1/time.Second),
						)
					}
				}
			})
	}
}
