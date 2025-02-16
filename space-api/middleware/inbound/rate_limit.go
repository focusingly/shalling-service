package inbound

import (
	"fmt"
	"space-api/constants"
	"space-api/middleware/auth"
	"space-api/util"
	"space-api/util/performance"
	"time"

	"github.com/gin-gonic/gin"
)

func UseReqRateLimitMiddleware(d time.Duration, maxReq int) gin.HandlerFunc {
	cache := performance.NewCache(constants.MB * 4)

	return func(ctx *gin.Context) {
		ip := GetRealIpWithContext(ctx)
		count, err := cache.GetInt64(ip)

		// 此前不存在访问或者访问已经重置
		if err != nil {
			cache.IncAndGet(ip, 1, performance.Second(d/time.Second))
			ctx.Next()
			return
		}

		// 管理员忽略任务访问基数限制
		user, e := auth.GetCurrentLoginSession(ctx)
		if e == nil && user.UserType == constants.Admin {
			ctx.Next()
			return
		}

		// 游客/非管理员需要进行访问限制
		switch {
		case count == int64(maxReq):
			ctx.Error(util.CreateLimitErr(
				"当前 ip 访问过快, 请稍后再试",
				fmt.Errorf("current ip request run out limit: %f/sec", float64(maxReq)/float64(d/time.Second))),
			)
			ctx.Abort()
		case count > int64(maxReq):
			// ctx.Status(http.StatusTooManyRequests)
			ctx.Abort()
		default:
			cache.GetAndIncr(ip, 1)
			ctx.Next()
		}
	}
}
