package common

import (
	"space-api/internal/service/v1"
	"space-api/middleware"
	"space-api/util"

	"github.com/gin-gonic/gin"
)

func UseAdminAuthController(group *gin.RouterGroup) {
	authGroup := group.Group("/auth/admin")

	// 系统管理员/站主 登录接口
	{
		authGroup.POST("/login", func(ctx *gin.Context) {

		})
	}

}

func UseOauth2Controller(group *gin.RouterGroup) {
	oauth2Group := group.Group("/oauth2")

	// Github 登录认证
	{
		githubRoute := oauth2Group.Group("/github")
		// 获取登录链接
		githubRoute.GET("/login", func(ctx *gin.Context) {
			if url, err := service.GetGithubLoginURL(ctx); err != nil {
				ctx.Error(&util.BizErr{
					Reason: err,
					Msg:    "获取认证链接失败",
				})
			} else {
				middleware.NotifyRestProducer(url, ctx)
			}
		})

		// 验证回调信息
		githubRoute.GET("/callback", func(ctx *gin.Context) {
			if resp, err := service.VerifyGithubCallback(ctx); err != nil {
				ctx.Error(&util.BizErr{
					Reason: err,
					Msg:    "登录失败",
				})
			} else {
				middleware.NotifyRestProducer(resp, ctx)
			}
		})
	}

	// 谷歌登录认证
	{
		googleRouting := oauth2Group.Group("/google")
		googleRouting.GET("/login", func(ctx *gin.Context) {
			if url, err := service.GetGoogleLoginURL(ctx); err != nil {
				ctx.Error(&util.BizErr{
					Reason: err,
					Msg:    "获取认证链接失败",
				})
			} else {
				middleware.NotifyRestProducer(url, ctx)
			}
		})
		googleRouting.GET("/callback", func(ctx *gin.Context) {
			if val, err := service.VerifyGoogleCallback(ctx); err != nil {
				ctx.Error(&util.BizErr{
					Reason: err,
					Msg:    "登录失败",
				})
			} else {
				middleware.NotifyRestProducer(val, ctx)
			}
		})
	}
}

func UseLogoutController(group *gin.RouterGroup) {
	// 退出登录
	{
		group.GET("/logout", func(ctx *gin.Context) {
			_, err := middleware.GetCurrentLoginUser(ctx)
			if err != nil {
				ctx.Error(err)
				return
			}

			middleware.NotifyRestProducer(new(struct{}), ctx)
		})
	}
}
