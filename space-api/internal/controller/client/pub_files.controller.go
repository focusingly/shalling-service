package client

import (
	"space-api/internal/service/v1"

	"github.com/gin-gonic/gin"
)

func UsePubStaticFilesController(group *gin.RouterGroup) {
	fileService := service.DefaultStaticFileService

	// 对外的公共文件访问点
	group.GET("/static/*file", func(ctx *gin.Context) {
		fileService.HandlePubVisit(ctx.Param("file"), ctx)
	})
}
