package cmd

import (
	"space-api/conf"
	"space-api/dto"
	"space-api/internal/service/v1"
	"space-api/util"

	"github.com/gin-gonic/gin"
)

func useDebugController(group *gin.RouterGroup) {
	_, isDebugMode := conf.GetParsedArgs()
	if !isDebugMode {
		return
	}

	group.POST("/upload", func(ctx *gin.Context) {
		f, e := service.DefaultUploadService.UploadImage2Webp(ctx)
		if e != nil {
			ctx.JSON(200, e)
		} else {
			ctx.JSON(200, f)
		}
	})

	group.GET("/logs", func(ctx *gin.Context) {
		req := &dto.DumpLogReq{}
		if e := ctx.ShouldBindQuery(req); e != nil {
			ctx.Error(util.CreateBizErr("参数错误: "+e.Error(), e))
			return
		}

		service.DefaultLogService.DumLogsStream(req, ctx)
	})
}
