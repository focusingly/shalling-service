package middleware

import (
	"net/http"
	"space-api/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var RestInjectMarkKey = uuid.New().String()

func UseRestProduceHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if val, ok := ctx.Get(RestInjectMarkKey); ok {
			ctx.JSON(http.StatusOK, util.RestWithSuccess(val))
		}
	}
}

// NotifyRestProducer 将要返回给客户端的值注到 gin 的上下文当中, 供中间件统一处理返回
func NotifyRestProducer[T any](val T, ctx *gin.Context) {
	ctx.Set(RestInjectMarkKey, val)
}
