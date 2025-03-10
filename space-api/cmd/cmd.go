package cmd

import (
	"fmt"
	"net/http"

	"space-api/cmd/adapter"
	"space-api/conf"
	"space-api/constants"
	"space-api/internal/controller"
	"space-api/internal/controller/common"
	"space-api/pack"
	"space-api/util"
	"time"

	"github.com/gin-gonic/gin"
)

func Run() {
	_, isDebug := conf.GetParsedArgs()
	gin.SetMode(util.TernaryExpr(isDebug, gin.DebugMode, gin.ReleaseMode))
	gin.ForceConsoleColor()
	engine := gin.New()
	engine.RedirectTrailingSlash = false
	engine.RemoveExtraSlash = true
	engine.MaxMultipartMemory = int64(constants.MB * 16) // 设置较小的表单内存
	appConf := conf.ProjectConf.GetAppConf()
	setTimeZoneIfRequire()                 // 设置时区(如果有指定的话)
	setDataSource()                        // 设置分库数据源
	engine.Use(getMiddlewares(appConf)...) // 应用全局中间件

	// 处理未知请求方法
	engine.NoMethod(func(ctx *gin.Context) {
		ctx.Error(util.CreateNotMethodOrResourceErr(
			"未知的请求方法: "+ctx.Request.Method,
			fmt.Errorf("unknown request method: %s", ctx.Request.Method),
		))
	})

	handleSpa := common.CreateEmbedSpaAppHandler(
		"/",
		"static/dist",
		appConf.ApiPrefix,
		&pack.SpaResource,
		time.Hour*24*15,
	)
	// 处理未注册路由
	engine.NoRoute(func(ctx *gin.Context) {
		// 只处理 GET 和 Head 请求到 SPA
		if ctx.Request.Method == http.MethodGet || ctx.Request.Method == http.MethodHead {
			handleSpa(ctx)
		} else {
			ctx.Error(util.CreateNotMethodOrResourceErr(
				"未知的请求资源: "+ctx.Request.RequestURI,
				fmt.Errorf("unknown request uri resource: %s", ctx.Request.RequestURI),
			))
		}
	})
	// 版本控制
	apiRouteGroup := engine.Group(appConf.ApiPrefix)
	// debug 测试示例路由
	useDebugController(engine.Group("/debug"))
	controller.RegisterAllControllers(apiRouteGroup)
	prepareStartup() // 设置项目初始化数据

	adapter.RunAndServe(engine, appConf)
}
