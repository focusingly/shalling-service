package outbound

import (
	"net/http"
	"space-api/util/rest"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
)

var restInjectMarkKey = "rest:" + uuid.NewString()

func UseRestProduceHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		if val, ok := ctx.Get(restInjectMarkKey); ok {
			handleProduce(
				http.StatusOK,
				rest.RestWithSuccess(val),
				ctx,
			)
		}
	}
}

// NotifyProduceResponse 将要返回给客户端的值注到 gin 的上下文当中, 供注册的中间件统一处理返回
func NotifyProduceResponse[T any](val T, ctx *gin.Context) {
	ctx.Set(restInjectMarkKey, val)
}

// 根据请求头的 Accept 返回给客户端指定的值, 如果没找到匹配项, 那么默认返回 json
func handleProduce[T any](code int, val T, ctx *gin.Context) {
	firstAccept := ctx.Request.Header.Get("Accept")
	switch firstAccept {
	case binding.MIMEXML, binding.MIMEXML2, "application/rss+xml":
		ctx.XML(code, val)
	case binding.MIMETOML:
		ctx.TOML(code, val)
	case binding.MIMEYAML, binding.MIMEYAML2:
		ctx.YAML(code, val)
	case binding.MIMEJSON:
		fallthrough
	default:
		ctx.JSON(code, val)
	}
}
