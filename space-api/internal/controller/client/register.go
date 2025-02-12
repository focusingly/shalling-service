package client

import (
	"github.com/gin-gonic/gin"
)

func RegisterAllClientComments(group *gin.RouterGroup) {
	clientGroup := group.Group("/client")

	IndexController(clientGroup)
	UseCommentController(clientGroup)
	UsePubStaticFilesController(clientGroup)
}
