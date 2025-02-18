package cmd

import (
	"space-api/conf"
	"space-api/constants"
	"space-api/middleware/auth"
	"space-api/middleware/inbound"
	"space-api/middleware/outbound"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func getMiddlewares(appConf *conf.AppConf) []gin.HandlerFunc {
	_, useDebug := conf.GetParsedArgs()
	store := cookie.NewStore([]byte(appConf.Salt))

	middlewares := []gin.HandlerFunc{
		outbound.UseErrorHandler(),
		outbound.UseServerResponseHintMiddleware(),
		outbound.UseRestProduceHandler(),

		inbound.UseUploadFileLimitMiddleware(constants.MemoryByteSize(appConf.ParsedUploadSize)),
		inbound.UseUseragentParserMiddleware(),
		inbound.UseExtractIPv4Middleware(),
		auth.UseJwtAuthHandler(),
		inbound.UseReqRateLimitMiddleware(time.Second*16, 32),

		sessions.Sessions("shalling.space", store),
		inbound.
			NewUVManager("/favicon.ico").
			CreateUVMiddleware(),
	}

	if useDebug {
		middlewares = append(middlewares, gin.Logger())
	}

	return middlewares
}
