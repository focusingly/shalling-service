package cmd

import (
	"fmt"
	"os"
	"space-api/conf"
	"space-api/constants"
	"space-api/db"
	"space-api/effect"
	"space-api/internal/controller"
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

	biz.SetDefault(db.GetBizDB())
	extra.SetDefault(db.GetExtraHelperDB())

	gin.SetMode(util.TernaryExpr(useDebug, gin.DebugMode, gin.ReleaseMode))
	gin.ForceConsoleColor()
	engine := gin.New()
	engine.MaxMultipartMemory = int64(constants.MB * 8)

	v := conf.GetProjectViper()
	var appConf conf.AppConf
	if err := v.UnmarshalKey("app", &appConf); err != nil {
		panic(err)
	}

	// TODO 时区设置, 暂不设置, 统一全部直接使用 unix 时间戳; 数据格式化由客户端自己解析
	// 定时任务的时区直接遵循服务器所设置的时区
	// if tz, err := time.LoadLocation(appConf.ServerTimezone); err != nil {
	// 	log.Fatal("获取时区失败: ", err)
	// } else {
	// 	time.Local = tz
	// }

	middlewares := []gin.HandlerFunc{
		outbound.UseErrorHandler(),
		outbound.UseServerResponseHintMiddleware(),
		outbound.UseRestProduceHandler(),
		inbound.UseUseragentParserMiddleware(),
		inbound.UseExtractIPv4Middleware(),
		// auth.UseJwtAuthHandler(),
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
	apiRouteGroup.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, []any{
			inbound.GetUserAgentFromContext(ctx),
			inbound.GetRealIpWithContext(ctx),
		})
	})
	controller.RegisterAllControllers(apiRouteGroup)

	effect.InvokeInit()
	engine.Run(fmt.Sprintf(":%d", appConf.Port))
}
