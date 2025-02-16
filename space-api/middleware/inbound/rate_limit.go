package inbound

import (
	"fmt"
	"space-api/constants"
	"space-api/internal/service/v1/blocking"
	"space-api/middleware/auth"

	"space-api/util"
	"space-api/util/performance"
	"time"

	"github.com/gin-gonic/gin"
)

func UseReqRateLimitMiddleware(d time.Duration, maxReq int) gin.HandlerFunc {
	cache := performance.NewCache(constants.MB * 4)
	blockingService := blocking.DefaultIPBlockingService

	return func(ctx *gin.Context) {
		// 管理员忽略任务访问基数限制
		user, e := auth.GetCurrentLoginSession(ctx)
		if e == nil && user.UserType == constants.Admin {
			ctx.Next()
			return
		}

		ip := GetRealIpWithContext(ctx)
		// 如果在黑名单中, 直接不进行任何操作, 返回
		if blockingService.IpInBlockingList(ip) {
			ctx.Abort()
			return
		}

		count, err := cache.GetInt64(ip)

		// 此前不存在访问或者访问已经重置
		if err != nil {
			cache.IncAndGet(ip, 1, performance.Second(d/time.Second))
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
