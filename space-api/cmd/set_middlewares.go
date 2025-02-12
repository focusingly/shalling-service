package cmd

import (
	"space-api/conf"
	"space-api/constants"
	"space-api/middleware/auth"
	"space-api/middleware/inbound"
	"space-api/middleware/outbound"

	"github.com/gin-gonic/gin"
)

func getMiddlewares(appConf *conf.AppConf) []gin.HandlerFunc {
	_, useDebug := conf.GetParsedArgs()

	middlewares := []gin.HandlerFunc{
		outbound.UseErrorHandler(),
		inbound.UseUploadFileLimitMiddleware(constants.MemoryByteSize(appConf.ParsedUploadSize)),
		outbound.UseServerResponseHintMiddleware(),
		outbound.UseRestProduceHandler(),
		inbound.UseUseragentParserMiddleware(),
		inbound.UseExtractIPv4Middleware(),
		auth.UseJwtAuthHandler(),
	}

	if useDebug {
		middlewares = append(middlewares, gin.Logger())
	}

	return middlewares
}
