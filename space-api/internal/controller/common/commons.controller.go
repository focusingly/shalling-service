package common

import (
	"space-api/middleware"
	"space-api/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UseCommonsController(group *gin.RouterGroup) {
	commonGroup := group.Group("/auth")
	// 创建 token
	commonGroup.GET("/", func(ctx *gin.Context) {
		if token, err := util.CreateJwtToken(uuid.NewString()); err != nil {
			ctx.Error(&util.BizErr{
				Msg:    err.Error(),
				Reason: err,
			})
		} else {
			middleware.NotifyRestProducer(token, ctx)
		}
	})
}
