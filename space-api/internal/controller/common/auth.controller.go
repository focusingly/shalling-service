package common

import (
	"space-api/dto"
	"space-api/internal/service/v1"
	"space-api/middleware/outbound"
	"space-api/util"

	"github.com/gin-gonic/gin"
)

func UseAuthController(group *gin.RouterGroup) {
	authService := service.DefaultAuthService

	// 系统管理员/站主/本地用户进行 登录
	{
		authGroup := group.Group("/local/login")
		authGroup.POST("/", func(ctx *gin.Context) {
			req := &dto.AdminLoginReq{}
			if err := ctx.ShouldBindBodyWithJSON(req); err != nil {
				ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
				return
			}
			if resp, err := authService.AdminLogin(req, ctx); err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceResponse(resp, ctx)
			}
		})
	}

	// 处理 oauth2 相关的登录/认证
	{
		oauth2LoginGroup := group.Group("/oauth2/url")
		// 获取登录链接
		oauth2LoginGroup.GET("/url", func(ctx *gin.Context) {
			req := &dto.GetLoginURLReq{}
			if err := ctx.ShouldBindQuery(req); err != nil {
				ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
				return
			}

			if resp, err := authService.GetOauth2LoginGrantURL(req, ctx); err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceResponse(resp, ctx)
			}
		})

		// 解析 Oauth 用户登录的回调信息
		oauth2LoginGroup.POST("/login", func(ctx *gin.Context) {
			req := &dto.OauthLoginCallbackReq{}

			if e := ctx.ShouldBindBodyWithJSON(req); e != nil {
				ctx.Error(util.CreateBizErr("参数错误: "+e.Error(), e))
				return
			}

			if resp, err := authService.HandleOauthLogin(req, ctx); err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceResponse(resp, ctx)
			}
		})
	}

	// 退出登录
	{
		group.Group("/logout").GET("/", func(ctx *gin.Context) {
			if resp, err := authService.CurrentUserLogout(ctx); err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceResponse(resp, ctx)
			}
		})
	}
}
