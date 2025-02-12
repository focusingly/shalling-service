package admin

import (
	"space-api/dto"
	"space-api/internal/service/v1"
	"space-api/middleware/outbound"
	"space-api/util"

	"github.com/gin-gonic/gin"
)

func UseCategoryController(group *gin.RouterGroup) {
	categoryGroup := group.Group("/category")
	categoryService := service.DefaultCategoryService

	// 创建或者更新已有的分类
	{
		categoryGroup.POST("/", func(ctx *gin.Context) {
			req := &dto.CreateOrUpdateCategoryReq{}
			if err := ctx.ShouldBindBodyWithJSON(req); err != nil {
				ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
				return
			}
			resp, err := categoryService.
				CreateOrUpdateCategory(req, ctx)
			if err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceResponse(resp, ctx)
			}
		})
	}

	// 获取分类信息和其相关联的文章信息
	{
		categoryGroup.GET("/:catID", func(ctx *gin.Context) {
			req := &dto.GetCategoryWithPostsReq{}
			if err := ctx.ShouldBindUri(req); err != nil {
				ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
				return
			}
			resp, err := categoryService.GetCategoryWithAllPosts(req, ctx)
			if err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceResponse(resp, ctx)
			}
		})
	}

	// 获取所有的分类列表信息
	{
		categoryGroup.GET("/list", func(ctx *gin.Context) {
			req := &dto.GetCategoryListReq{}
			if err := ctx.ShouldBindQuery(req); err != nil {
				ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
				return
			}
			resp, err := categoryService.GetAllCategories(req, ctx)
			if err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceResponse(resp, ctx)
			}
		})
	}

	// 删除分类
	{
		categoryGroup.DELETE("/", func(ctx *gin.Context) {
			req := &dto.DeleteCategoryReq{}
			if err := ctx.ShouldBindBodyWithJSON(req); err != nil {
				ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
				return
			}
			if resp, err := categoryService.DeleteCategoryByIDList(req, ctx); err != nil {
				ctx.Error(err)
			} else {
				outbound.NotifyProduceResponse(resp, ctx)
			}
		})
	}
}
