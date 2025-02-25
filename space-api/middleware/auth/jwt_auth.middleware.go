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

var authSpaceCache = performance.DefaultJsonCache.Group("auth")

func GetMiddlewareRelativeAuthCache() performance.CacheGroupInf {
	return authSpaceCache
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

// 加载请求头中的 token 解析并设置到当前请求上下文当中
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
			err.Error(),
			err,
		))
		ctx.Abort()
	} else {
		sessionID := claims.Jti
		parsedSessionID, err := strconv.ParseInt(sessionID, 10, 64)
		if err != nil {
			err = util.CreateAuthErr(err.Error(), err)
			return
		}
		// 优先尝试从缓存当中获取凭据
		cachedLoginSession := &model.UserLoginSession{}
		if fetchErr := authSpaceCache.Get(sessionID, cachedLoginSession); fetchErr == nil {
			// 存在命中情况, 直接返回即可
			ctx.Set(_jwtRandMark, cachedLoginSession)
			return
		}

		// 缓存不存在的情况下, 查表进行二次判断
		loginSessionTx := biz.UserLoginSession
		findLoginSession, err := loginSessionTx.
			WithContext(ctx).
			Where(loginSessionTx.ID.Eq(parsedSessionID)).
			Take()
		if err != nil {
			ctx.Error(util.CreateAuthErr("用户登录会话已失效, 请重新登录", fmt.Errorf("user login session expired, please re-login")))
			ctx.Abort()
			return
		}

		// 设置到缓存
		authSpaceCache.Set(sessionID, findLoginSession, time.UnixMilli(findLoginSession.ExpiredAt).Sub(time.Now()))
		//设置到上下文
		ctx.Set(_jwtRandMark, findLoginSession)
	}
}
