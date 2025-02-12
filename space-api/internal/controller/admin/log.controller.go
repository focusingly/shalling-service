package admin

import (
	"space-api/dto"
	"space-api/internal/service/v1"
	"space-api/middleware/outbound"
	"space-api/util"

	"github.com/gin-gonic/gin"
)

func UseLogController(group *gin.RouterGroup) {
	logGroup := group.Group("/log")
	logService := service.DefaultLogService

	// 查看日志分页信息
	{
		logGroup.POST("/list", func(ctx *gin.Context) {
			req := &dto.GetLogPagesReq{}
			if e := ctx.ShouldBindBodyWithJSON(req); e != nil {
				ctx.Error(util.CreateBizErr("参数错误: "+e.Error(), e))
				return
			}
			if resp, err := logService.GetLogPages(req, ctx); err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceResponse(resp, ctx)
			}
		})
	}

	// 删除日志分页信息
	{
		logGroup.DELETE("/", func(ctx *gin.Context) {
			req := &dto.DeleteLogReq{}
			if e := ctx.ShouldBindBodyWithJSON(req); e != nil {
				ctx.Error(util.CreateBizErr("参数错误: "+e.Error(), e))
				return
			}
			if resp, err := logService.DeleteLogsByCondition(req, ctx); err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceResponse(resp, ctx)
			}
		})
	}

	// 导出日志
	{
		logGroup.GET("/dump", func(ctx *gin.Context) {
			req := &dto.DumpLogReq{}
			if e := ctx.ShouldBindQuery(req); e != nil {
				ctx.Error(util.CreateBizErr("参数错误: "+e.Error(), e))
				return
			}

			logService.DumLogsStream(req, ctx)
		})
	}
}
