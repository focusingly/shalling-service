package auth

import (
	"fmt"
	"space-api/util"
	"space-api/util/performance"
	"space-api/util/verify"
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

var _authCacheGroup = performance.DefaultJsonCache.Group("auth")

func GetMiddlewareRelativeAuthCache() *performance.JsonCache {
	return _authCacheGroup
}

const (
	BaseVersion = "/v1/api"

	AdminPath = BaseVersion + "/admin"
	Client    = BaseVersion + "/client"
	Common    = BaseVersion + "/common"
)

// UseJwtAuthExtractMiddleware 提取请求中 JWT 信息, 并设置到当前请求上下文当中
func UseJwtAuthExtractMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		loadTokenAndSetupContext(ctx)

		ctx.Next()
	}
}

// GetCurrentLoginSession 获取当前请求上下文中的用户登录信息
func GetCurrentLoginSession(ctx *gin.Context) (user *model.UserLoginSession, err error) {
	if u, ok := ctx.Get(_jwtRandMark); !ok {
		err = util.CreateAuthErr(
			"无授权凭据, 请先登录",
			fmt.Errorf("oo found login user"))
		return
	} else {
		if t, ok := u.(*model.UserLoginSession); !ok {
			err = util.CreateAuthErr("提取用户标识失败",
				fmt.Errorf("extract user id failed"))
			return
		} else {
			user = t
		}
	}

	return
}

func loadTokenAndSetupContext(ctx *gin.Context) {
	bearerToken := ctx.Request.Header.Get("Authorization")
	// 没有携带 token 的情况下直接跳过设置上下文, 不进行解析
	if bearerToken == "" {
		return
	}

	if !strings.HasPrefix(bearerToken, BearerAuthPrefix) {
		ctx.Error(
			util.CreateAuthErr(
				"非法的授权凭据",
				fmt.Errorf("illegal principal: %s", bearerToken),
			),
		)
		ctx.Abort()
		return
	}

	// 获取 token
	claims, err := verify.VerifyAndGetParsedBizClaims(bearerToken[len(BearerAuthPrefix):])

	// token 本身无效
	if err != nil {
		ctx.Error(util.CreateAuthErr(
			"凭据无效, 请重新登录",
			err,
		))
		ctx.Abort()
	} else {
		// 优先尝试从缓存中获取用户信息
		cacheUUIDKey := claims.UUID
		cachedLoginSession := &model.UserLoginSession{}
		if err := _authCacheGroup.Get(cacheUUIDKey, cachedLoginSession); err == nil {
			// 存在命中情况, 直接返回即可
			ctx.Set(_jwtRandMark, cachedLoginSession)
			return
		}

		// 不存在的情况下, 进行查表进行二次判断
		userId, err := strconv.ParseInt(claims.Jti, 10, 64)
		if err != nil {
			ctx.Error(util.CreateAuthErr("提取用户 ID 失败: "+err.Error(), err))
			ctx.Abort()
			return
		} else {
			loginSessionTx := biz.UserLoginSession
			findLoginSession, err := loginSessionTx.
				WithContext(ctx).
				Where(loginSessionTx.ID.Eq(userId), loginSessionTx.UUID.Eq(cacheUUIDKey)).
				Take()
			if err != nil {
				ctx.Error(util.CreateAuthErr("用户登录会话已失效, 请重新登录", fmt.Errorf("user login session expired, please re-login")))
				ctx.Abort()
				return
			}

			// 设置到缓存
			_authCacheGroup.Set(cacheUUIDKey, findLoginSession, performance.Second(findLoginSession.ExpiredAt-time.Now().Unix()))
			//设置到上下文
			ctx.Set(_jwtRandMark, findLoginSession)
		}
	}
}
