package common

import (
	"net/http"
	"space-api/util"

	"github.com/gin-gonic/gin"
)

func UseHealthCheckController(routeGroup *gin.RouterGroup) {
	routeGroup.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, util.RestWithSuccess("pong"))
	})
}
