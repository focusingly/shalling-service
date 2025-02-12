package cmd

import (
	"fmt"
	"space-api/conf"
	"space-api/constants"
	"space-api/internal/controller"
	"space-api/util"

	"github.com/gin-gonic/gin"
)

func Run() {
	_, isDebug := conf.GetParsedArgs()
	gin.SetMode(util.TernaryExpr(isDebug, gin.DebugMode, gin.ReleaseMode))
	gin.ForceConsoleColor()
	engine := gin.New()

	engine.MaxMultipartMemory = int64(constants.MB * 32) // 设置表单处理占用的最内存大小
	appConf := conf.ProjectConf.GetAppConf()             // 应用配置

	setTimeZone()                          // 设置时区(如果有指定的话)
	setDataSource()                        // 设置分库数据源
	engine.Use(getMiddlewares(appConf)...) // 应用全局中间件

	// 处理未知请求方法
	engine.NoMethod(func(ctx *gin.Context) {
		ctx.Error(util.CreateBizErr(
			"未知的请求方法: "+ctx.Request.Method,
			fmt.Errorf("unknown request method: %s", ctx.Request.Method),
		))
	})
	// 处理未注册路由
	engine.NoRoute(func(ctx *gin.Context) {
		ctx.Error(util.CreateBizErr(
			"未知的请求资源: "+ctx.Request.RequestURI,
			fmt.Errorf("unknown request uri resource: %s", ctx.Request.RequestURI),
		))
	})

	// 版本控制
	apiRouteGroup := engine.Group("/v1/api")
	// debug 测试示例路由
	useDebugController(engine.Group("/debug"))

	controller.RegisterAllControllers(apiRouteGroup)
	prepareStartup() // 设置项目初始化数据

	engine.Run(fmt.Sprintf(":%d", appConf.Port))
}
