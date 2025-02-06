package outbound

import (
	"net/http"
	"space-api/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var _restInjectMarkKey = "rest:" + uuid.NewString()

func UseRestProduceHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if val, ok := ctx.Get(_restInjectMarkKey); ok {
			ctx.JSON(http.StatusOK, util.RestWithSuccess(val))
		}
	}
}

// NotifyProduceRestJSON 将要返回给客户端的值注到 gin 的上下文当中, 供注册的中间件统一处理返回
func NotifyProduceRestJSON[T any](val T, ctx *gin.Context) {
	ctx.Set(_restInjectMarkKey, val)
}
