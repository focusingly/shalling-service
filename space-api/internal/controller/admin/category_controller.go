package admin

import "github.com/gin-gonic/gin"

func UseCategoryController(group *gin.RouterGroup) {
	categoryGroup := group.Group("/category")
	// 创建或者更新已有的分类
	{
		categoryGroup.POST("/", func(ctx *gin.Context) {})
	}

	// 查询分类列表信息
	{

		categoryGroup.GET("/list", func(ctx *gin.Context) {})
	}

	// 获取单个分类信息
	{
		categoryGroup.GET("/:id", func(ctx *gin.Context) {})
	}

	// 删除分类
	{
		categoryGroup.DELETE("/", func(ctx *gin.Context) {})
	}

}
