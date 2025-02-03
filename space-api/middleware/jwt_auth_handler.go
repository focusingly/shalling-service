package middleware

import (
	"space-api/util"
	"space-domain/model"

	"github.com/gin-gonic/gin"
)

const AuthPrefix = "Bearer "

func UseJwtAuthHandler() gin.HandlerFunc {
	user := model.LoginUser{
		UserType:    "Admin",
		DisplayName: "Admin",
		PlatformId:  1232625235235,
		Link:        "https://www.shalling.me",
		Email:       "shalling@shalling.me",
	}

	return func(ctx *gin.Context) {
		ctx.Set("current-member", &user)

		ctx.Next()
	}
}

// GetCurrentLoginUser 获取当前的凭据
func GetCurrentLoginUser(ctx *gin.Context) (user *model.LoginUser, err error) {
	if u, exits := ctx.Get("current-member"); !exits {
		err = &util.BizErr{
			Msg:    "未登录的用户",
			Reason: "No Found Login User",
		}
		return
	} else {
		if t, ok := u.(*model.LoginUser); !ok {
			err = &util.BizErr{
				Msg:    "未登录的用户",
				Reason: "No Found Login User",
			}
			return
		} else {
			user = t
		}
	}

	return
}
