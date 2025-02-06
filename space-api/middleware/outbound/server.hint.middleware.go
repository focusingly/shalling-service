package outbound

import (
	"space-api/conf"

	"github.com/gin-gonic/gin"
)

func UseServerResponseHintMiddleware() gin.HandlerFunc {
	v := conf.GetProjectViper()
	appConf := &conf.AppConf{}
	if err := v.UnmarshalKey("app", appConf); err != nil {
		panic(err)
	}

	return func(ctx *gin.Context) {
		ctx.Header("Server", appConf.ServerHint)
		ctx.Next()
	}
}
