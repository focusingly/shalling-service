package middleware

import (
	"fmt"
	"space-api/constants"
	"space-api/util"
	"space-domain/dao/biz"
	"space-domain/model"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const BearerAuthPrefix = "Bearer "

var _jwtRandMark = uuid.NewString()

var AuthCacheGroup = util.DefaultJsonCache.Group("/auth")

const (
	BaseVersion = "/v1/api"

	AdminPath = BaseVersion + "/admin"
	Client    = BaseVersion + "/client"
	Common    = BaseVersion + "/common"
)

func UseJwtAuthHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		loadTokenAndStoreToContext(ctx)
		p := ctx.Request.URL.Path

		switch {
		// 所有的后台服务都要求使用管理员权限
		case strings.HasPrefix(p, AdminPath):
			if user, err := GetCurrentLoginUser(ctx); err != nil {
				ctx.Error(&util.VerifyErr{
					BizErr: util.BizErr{
						Msg:    "获取用户凭据失败, 请先登录",
						Reason: err,
					},
				})
				ctx.Abort()
				return
			} else {
				if user.UserType != constants.LocalUser {
					ctx.Error(&util.VerifyErr{
						BizErr: util.BizErr{
							Msg:    "仅限管理员用户进行登录",
							Reason: fmt.Errorf("un-support user, want%s, but current is:%s", constants.LocalUser, user.UserType),
						},
					})
					ctx.Abort()
					return
				}
			}
		default:
			// pass....
		}

		ctx.Next()
	}
}

// GetCurrentLoginUser 获取当前的凭据
func GetCurrentLoginUser(ctx *gin.Context) (user *model.LoginUser, err *util.VerifyErr) {
	if u, exits := ctx.Get(_jwtRandMark); !exits {
		err = &util.VerifyErr{
			BizErr: util.BizErr{
				Msg:    "无授权凭据, 请先登录",
				Reason: fmt.Errorf("oo found login user"),
			},
		}
		return
	} else {
		if t, ok := u.(*model.LoginUser); !ok {
			err = &util.VerifyErr{
				BizErr: util.BizErr{
					Msg:    "提取用户标识失败",
					Reason: fmt.Errorf("extract user id failed"),
				},
			}
			return
		} else {
			user = t
		}
	}

	return
}

func loadTokenAndStoreToContext(ctx *gin.Context) {
	bearerToken := ctx.Request.Header.Get("Authorization")
	// 没有携带 token 的情况下直接跳过设置上下文
	if bearerToken == "" {
		return
	}

	if !strings.HasPrefix(bearerToken, BearerAuthPrefix) {
		ctx.Error(&util.BizErr{
			Msg:    "非法的授权凭据",
			Reason: fmt.Errorf("illegal principal: %s", bearerToken),
		})
		ctx.Abort()
		return
	}

	// 获取 token
	claims, err := util.VerifyAndGetClaims(bearerToken[len(BearerAuthPrefix):])

	// token 本身无效
	if err != nil {
		ctx.Error(&util.BizErr{
			Msg:    "凭据验证失败: " + err.Error(),
			Reason: err,
		})
	} else {
		userId, err := strconv.ParseInt((claims["jti"]).(string), 10, 64)
		if err != nil {
			ctx.Error(&util.BizErr{
				Msg:    "提取用户 ID 失败: " + err.Error(),
				Reason: err,
			})
			ctx.Abort()
			return
		} else {
			// 从缓存中获取凭据
			user := new(model.LoginUser)
			if e := AuthCacheGroup.GetById(fmt.Sprintf("%d", userId), user); e != nil {
				// 缓存当中不存在数据, 那么查找数据库并进行设置
				user, err = biz.LoginUser.WithContext(ctx).Where(biz.LoginUser.Id.Eq(userId)).Take()
				if err != nil {
					ctx.Error(&util.BizErr{
						Msg:    "提取用户 ID 失败: " + err.Error(),
						Reason: e,
					})
					ctx.Abort()
					return
				} else {
					exp, _ := claims.GetExpirationTime()
					sec := time.Until(exp.Time).Seconds() / time.Hour.Seconds()
					AuthCacheGroup.SetWith(fmt.Sprintf("%d", userId), user, util.Second(sec))
					ctx.Set(_jwtRandMark, user)
				}
			} else {
				// 缓存当中存在, 直接设置到上下文
				ctx.Set(_jwtRandMark, user)
			}
		}
	}
}
