package admin

import (
	"log"
	"space-api/dto"
	"space-api/internal/service/v1"
	"space-api/middleware/outbound"
	"space-api/util"

	"github.com/gin-gonic/gin"
)

func UseJobController(group *gin.RouterGroup) {
	taskService := service.DefaultTaskService
	taskGroup := group.Group("/task")

	// 恢复记录
	if err := taskService.ResumeFromDatabase(); err != nil {
		log.Fatal("从数据库中恢复任务失败: ", err)
	}

	// 添加/更新 定时任务
	{
		taskGroup.POST("/update", func(ctx *gin.Context) {
			req := &dto.CreateOrUpdateJobReq{}
			if err := ctx.ShouldBindBodyWithJSON(req); err != nil {
				ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
				return
			}
			if resp, err := taskService.CreateOrUpdateNewJob(req, ctx); err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceRestJSON(resp, ctx)
			}
		})
	}

	// 获取可用的任务列表/已定义的任务
	{
		taskGroup.GET("/available/list", func(ctx *gin.Context) {
			req := &dto.GetAvailableJobListReq{}
			if err := ctx.ShouldBindUri(req); err != nil {
				ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
				return
			}
			if resp, err := taskService.GetAvailableJobList(req, ctx); err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceRestJSON(resp, ctx)
			}
		})
	}

	// 获取已经添加到数据库中的任务
	{
		taskGroup.GET("/db/list", func(ctx *gin.Context) {
			req := &dto.GetRunningJobListReq{}
			if err := ctx.ShouldBindUri(req); err != nil {
				ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
				return
			}

			if resp, err := taskService.GetRunningJobs(req, ctx); err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceRestJSON(resp, ctx)
			}

		})
	}

	// 立即执行一个任务
	{
		taskGroup.POST("/execute", func(ctx *gin.Context) {
			req := &dto.RunJobReq{}
			if err := ctx.ShouldBindBodyWithJSON(req); err != nil {
				ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
				return
			}
			if resp, err := taskService.RunJobImmediately(req, ctx); err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceRestJSON(resp, ctx)
			}

		})
	}

	// 删除任务
	{
		taskGroup.DELETE("/", func(ctx *gin.Context) {
			req := &dto.DeleteRunningJobListReq{}
			if err := ctx.ShouldBindBodyWithJSON(req); err != nil {
				ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
				return
			}

			if resp, err := taskService.DeleteRunningJobs(req, ctx); err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceRestJSON(resp, ctx)
			}
		})
	}

}
