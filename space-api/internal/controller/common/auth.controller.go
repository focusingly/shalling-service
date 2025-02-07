package common

import (
	"space-api/dto"
	"space-api/internal/service/v1"
	"space-api/middleware/outbound"
	"space-api/util"

	"github.com/gin-gonic/gin"
)

func UseAuthController(group *gin.RouterGroup) {

	authService := service.DefaultOauth2Service

	// 系统管理员/站主 登录接口
	{
		authGroup := group.Group("/auth/admin")
		authGroup.POST("/login", func(ctx *gin.Context) {
			req := &dto.AdminLoginReq{}
			if err := ctx.ShouldBindBodyWithJSON(req); err != nil {
				ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
				return
			}
			if resp, err := authService.AdminLogin(req, ctx); err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceRestJSON(resp, ctx)
			}
		})
	}

	// 处理 oauth2 相关的登录/认证
	{
		oauth2LoginGroup := group.Group("/auth/login")

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
				outbound.NotifyProduceRestJSON(resp, ctx)
			}
		})
	}
}
