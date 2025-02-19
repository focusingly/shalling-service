package auth

import (
	"fmt"
	"space-api/constants"
	"space-api/util"

	"github.com/gin-gonic/gin"
)

func UseAdminAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session, err := GetCurrentLoginSession(ctx)
		switch {
		case err != nil:
			ctx.Error(util.CreateAuthErr(
				"获取用户凭据失败, 请先登录",
				err,
			))
			ctx.Abort()
			return
		case session.UserType == constants.LocalUser:
			ctx.Error(util.CreateAuthErr(
				"当前用户类型不支持此操作",
				fmt.Errorf("un-support user, want%s, but current is:%s", constants.LocalUser, session.UserType),
			))
			// 不需要后续流程
			ctx.Abort()
			return
		case session.UserType == constants.Admin:
			ctx.Next()
		default:
			ctx.Error(util.CreateAuthErr(
				"未知的用户类型: "+session.UserType,
				fmt.Errorf("unknown user: %s", session.UserType),
			))
			ctx.Abort()
		}

	}
}
