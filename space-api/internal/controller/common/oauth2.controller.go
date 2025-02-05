package common

import (
	"net/http"
	"space-api/internal/service/v1"
	"space-api/util"

	"github.com/gin-gonic/gin"
)

func UseOauth2Controllers(routeGroup *gin.RouterGroup) {
	oauth2Group := routeGroup.Group("/oauth2")

	{
		githubRouting := oauth2Group.Group("/github")
		githubRouting.GET("/login", func(ctx *gin.Context) {
			if url, err := service.GetGithubLoginURL(ctx); err != nil {
				ctx.Error(&util.BizErr{
					Reason: err,
					Msg:    "获取认证链接失败",
				})
			} else {
				ctx.JSON(http.StatusOK, util.RestWithSuccess(url))
			}
		})
		githubRouting.GET("/callback", func(ctx *gin.Context) {
			if val, err := service.GithubCallbackHandler(ctx); err != nil {
				ctx.Error(&util.BizErr{
					Reason: err,
					Msg:    "登录失败",
				})
			} else {
				ctx.JSON(http.StatusOK, val)
			}
		})
	}

	{
		googleRouting := oauth2Group.Group("/google")
		googleRouting.GET("/login", func(ctx *gin.Context) {
			if url, err := service.GetGoogleLoginURL(ctx); err != nil {
				ctx.Error(&util.BizErr{
					Reason: err,
					Msg:    "获取认证链接失败",
				})
			} else {
				ctx.JSON(http.StatusOK, util.RestWithSuccess(url))
			}
		})
		googleRouting.GET("/callback", func(ctx *gin.Context) {
			if val, err := service.GoogleCallbackHandler(ctx); err != nil {
				ctx.Error(&util.BizErr{
					Reason: err,
					Msg:    "登录失败",
				})
			} else {
				ctx.JSON(http.StatusOK, val)
			}
		})
	}
}
