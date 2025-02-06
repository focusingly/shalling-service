package cmd

import (
	"fmt"
	"os"
	"space-api/conf"
	"space-api/constants"
	"space-api/db"
	"space-api/internal/controller"
	"space-api/middleware/auth"
	"space-api/middleware/inbound"
	"space-api/middleware/outbound"
	"space-api/util"
	"space-domain/dao/biz"
	"space-domain/dao/extra"

	"github.com/gin-gonic/gin"
)

func Run() {
	useDebug := false
	for _, arg := range os.Args {
		if arg == "-use-debug" {
			useDebug = true
			continue
		}
	}

	gin.SetMode(util.TernaryExpr(useDebug, gin.DebugMode, gin.ReleaseMode))
	biz.SetDefault(db.GetBizDB())
	extra.SetDefault(db.GetExtraHelperDB())
	gin.ForceConsoleColor()
	engine := gin.New()
	engine.MaxMultipartMemory = int64(constants.MB * 8)
	v := conf.GetProjectViper()
	var appConf conf.AppConf

	if err := v.UnmarshalKey("app", &appConf); err != nil {
		panic(err)
	}

	middlewares := []gin.HandlerFunc{
		outbound.UseErrorHandler(),
		outbound.UseServerResponseHintMiddleware(),
		outbound.UseRestProduceHandler(),
		inbound.UseExtractIPv4Middleware(),
		auth.UseJwtAuthHandler(),
	}

	if useDebug {
		middlewares = append(middlewares, gin.Logger())
	}
	engine.Use(middlewares...)

	engine.NoMethod(func(ctx *gin.Context) {
		ctx.Error(&util.BizErr{Msg: "不存在的方法: " + ctx.Request.Method})
	})
	engine.NoRoute(func(ctx *gin.Context) {
		ctx.Error(&util.BizErr{Msg: "不存在的请求资源: " + ctx.Request.URL.String()})
	})

	apiRouteGroup := engine.Group("/v1/api")
	controller.RegisterAllControllers(apiRouteGroup)
	engine.Run(fmt.Sprintf(":%d", appConf.Port))
}
