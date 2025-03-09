package admin

import (
	"space-api/dto"
	"space-api/internal/service/v1/mail"
	"space-api/middleware/outbound"
	"space-api/util"

	"github.com/gin-gonic/gin"
)

func UseMailController(group *gin.RouterGroup) {
	emailService := mail.DefaultMailService
	emailGroup := group.Group("/mail")

	// 获取邮件配置列表
	emailGroup.GET("/list", func(ctx *gin.Context) {
		resp := emailService.GetConfList(&dto.GetMailConfListReq{})
		outbound.NotifyProduceResponse(resp, ctx)
	})

	// 使用主邮件配置发送邮件
	emailGroup.POST("/primary/send", func(ctx *gin.Context) {
		req := &dto.SendMailReq{}
		if err := ctx.ShouldBindBodyWithJSON(req); err != nil {
			ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
			return
		}
		if resp, err := emailService.SendEmailByPrimary(req); err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(resp, ctx)
		}
	})

	// 使用自选的邮件配置发送邮件
	emailGroup.POST("/chose/send", func(ctx *gin.Context) {
		req := &dto.SendMailWithSelectionReq{}
		if err := ctx.ShouldBindBodyWithJSON(req); err != nil {
			ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
			return
		}
		if resp, err := emailService.SendEmailBySelection(req); err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(resp, ctx)
		}
	})
}
