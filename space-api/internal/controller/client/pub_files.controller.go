package client

import (
	"space-api/internal/service/v1"

	"github.com/gin-gonic/gin"
)

func UsePubStaticFilesController(group *gin.RouterGroup) {
	fileService := service.DefaultStaticFileService
	fileGroup := group.Group("/static")
	// 对外的公共文件访问点
	fileGroup.GET("*file", func(ctx *gin.Context) {
		fileService.HandlePubVisit(ctx.Param("file"), ctx)
	})
}
