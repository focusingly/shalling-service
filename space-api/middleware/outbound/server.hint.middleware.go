package outbound

import (
	"space-api/conf"

	"github.com/gin-gonic/gin"
)

func UseServerResponseHintMiddleware() gin.HandlerFunc {
	hint := conf.ProjectConf.GetAppConf().ServerHint

	return func(ctx *gin.Context) {
		ctx.Header("Server", hint)
		ctx.Next()
	}
}
