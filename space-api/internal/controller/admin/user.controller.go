package admin

import (
	"space-api/dto"
	"space-api/internal/service/v1"
	"space-api/middleware/outbound"
	"space-api/util"

	"github.com/gin-gonic/gin"
)

func UseUserController(group *gin.RouterGroup) {
	userGroup := group.Group("/user")

	userService := service.DefaultUserService

	// 查看已经登录的用户的凭据
	{
		userGroup.GET("/session/list", func(ctx *gin.Context) {
			req := &dto.GetLoginUserSessionsReq{}
			if e := ctx.ShouldBindUri(req); e != nil {
				ctx.Error(util.CreateBizErr("参数错误: "+e.Error(), e))
				return
			}
			if resp, err := userService.GetLoginSessions(req, ctx); err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceResponse(resp, ctx)
			}
		})
	}

	// 强制用户下线
	{
		userGroup.DELETE("/session/expire", func(ctx *gin.Context) {
			req := &dto.ExpireUserLoginSessionReq{}
			if e := ctx.ShouldBindBodyWithJSON(req); e != nil {
				ctx.Error(util.CreateBizErr("参数错误: "+e.Error(), e))
				return
			}
			if resp, err := userService.ExpireLoginSessions(req, ctx); err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceResponse(resp, ctx)
			}
		})
	}

	// 更新登录的 oauth2 账户信息
	{
		userGroup.POST("/oauth/update", func(ctx *gin.Context) {
			req := &dto.UpdateOauthUserReq{}
			if e := ctx.ShouldBindBodyWithJSON(req); e != nil {
				ctx.Error(util.CreateBizErr("参数错误: "+e.Error(), e))
				return
			}
			if resp, err := userService.UpdateOauth2User(req, ctx); err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceResponse(resp, ctx)
			}
		})
	}

	// 删除 oauth2 账户信息
	{
		userGroup.DELETE("/oauth/delete", func(ctx *gin.Context) {
			req := &dto.UpdateOauthUserReq{}
			if e := ctx.ShouldBindBodyWithJSON(req); e != nil {
				ctx.Error(util.CreateBizErr("参数错误: "+e.Error(), e))
				return
			}
			if resp, err := userService.DeleteOauth2User(req, ctx); err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceResponse(resp, ctx)
			}
		})
	}

	// 更新本地账户的基本显示信息
	{
		userGroup.POST("/admin/basic/update", func(ctx *gin.Context) {
			req := &dto.UpdateLocalUserBasicReq{}
			if e := ctx.ShouldBindBodyWithJSON(req); e != nil {
				ctx.Error(util.CreateBizErr("参数错误: "+e.Error(), e))
				return
			}
			if resp, err := userService.UpdateLocalUserBasic(req, ctx); err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceResponse(resp, ctx)
			}
		})
	}

	// 更新本地账户的密码信息
	{
		userGroup.POST("/admin/password/password", func(ctx *gin.Context) {
			req := &dto.UpdateLocalUserPassReq{}
			if e := ctx.ShouldBindBodyWithJSON(req); e != nil {
				ctx.Error(util.CreateBizErr("参数错误: "+e.Error(), e))
				return
			}
			if resp, err := userService.UpdateLocalUserPassword(req, ctx); err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceResponse(resp, ctx)
			}
		})
	}

}
