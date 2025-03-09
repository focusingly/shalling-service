package inbound

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"slices"
	"space-api/constants"
	"space-api/internal/service/v1/blocking"
	"strings"

	"space-api/util"
	"space-api/util/performance"
	"time"

	"github.com/gin-gonic/gin"
)

func UseReqRateLimitMiddleware(d time.Duration, maxReq int, limitPath ...string) gin.HandlerFunc {
	cache := performance.DefaultJsonCache.Group("req-limit")
	blockingService := blocking.DefaultIPBlockingService

	if err := blockingService.SyncBlockingRecordInCache(context.Background()); err != nil {
		log.Fatal("init blocking ip list failure", err)
	}

	return func(ctx *gin.Context) {
		// 只限制某些路径需要进行限流检查
		if !slices.ContainsFunc(limitPath, func(p string) bool {
			return strings.HasPrefix(ctx.Request.URL.Path, p)
		}) {
			ctx.Next()

			return
		}

		// 管理员忽略任务访问基数限制
		user, e := GetCurrentLoginSession(ctx)
		if e == nil && user.UserType == constants.Admin {
			ctx.Next()
			return
		}

		ip := GetRealIpWithContext(ctx)
		// 如果在黑名单中, 直接不进行任何操作
		if blockingService.IpInBlockingList(ip) {
			ctx.Abort()
			return
		}

		count, err := cache.GetInt64(ip)
		// 此前不存在访问或者访问已经重置
		if err != nil {
			cache.IncrAndGet(ip, 1, d)
			ctx.Next()
			return
		}

		count++
		cache.Set(ip, count)

		// 游客/非管理员需要进行访问限制
		switch {
		case count == int64(maxReq+1): // 第一次刚超过的时候给个提示
			// 重设过期时间
			cache.SetTTL(ip, d)
			ctx.Error(util.CreateLimitErr(
				"当前 ip 访问过快, 请稍后再试",
				fmt.Errorf("current ip request run out limit: %f/sec", float64(maxReq)/float64(d/time.Second))),
			)
			ctx.Abort()
		case count > int64(maxReq+1): // 后续继续超过的话不响应
			// 重设过期时间
			cache.SetTTL(ip, d)
			ctx.Status(http.StatusTooManyRequests)
			ctx.Abort()
		default:
			ctx.Next()
		}
	}
}
