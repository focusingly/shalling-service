package admin

import (
	"net/http"
	"space-api/dto"
	"space-api/internal/service/v1"
	"space-api/util"

	"github.com/gin-gonic/gin"
)

func UseTagController(group *gin.RouterGroup) {
	tagGroup := group.Group("/tag")
	service := service.DefaultTagService

	// 获取标签分页列表
	{
		tagGroup.GET("/list", func(ctx *gin.Context) {
			req := &dto.GetTagPageListReq{}
			err := ctx.BindQuery(req)
			if err != nil {
				ctx.Error(&util.BizErr{
					Msg:    "参数错误: " + err.Error(),
					Reason: err,
				})

				return
			}
			if val, err := service.GetTagPageList(req, ctx); err != nil {
				ctx.Error(err)
			} else {
				ctx.JSON(http.StatusOK, util.RestWithSuccess(val))
			}
		})
	}

	// 根据 ID 获取单独的标签信息
	{
		tagGroup.GET("/:id", func(ctx *gin.Context) {
			req := &dto.GetTagDetailReq{}
			if err := ctx.ShouldBindUri(req); err != nil {
				ctx.Error(&util.BizErr{
					Reason: err,
					Msg:    "参数错误: " + err.Error(),
				})

				return
			}

			if resp, err := service.GetTagDetailById(req, ctx); err != nil {
				ctx.Error(err)
				return
			} else {
				ctx.JSON(http.StatusOK, util.RestWithSuccess(resp))
			}
		})
	}

	// 更新/创建标签
	{
		tagGroup.POST("/", func(ctx *gin.Context) {
			req := &dto.CreateOrUpdateTagReq{}
			if err := ctx.BindJSON(req); err != nil {
				ctx.Error(&util.BizErr{
					Msg:    "参数错误: " + err.Error(),
					Reason: err,
				})
				return
			}

			if resp, err := service.CreateOrUpdateTag(req, ctx); err != nil {
				ctx.Error(err)

			} else {
				ctx.JSON(http.StatusOK, util.RestWithSuccess(resp))
			}
		})
	}

	// 根据指定的 ID 列表删除标签
	{
		tagGroup.DELETE("/", func(ctx *gin.Context) {
			req := &dto.DeleteTagByIdListReq{}

			if err := ctx.ShouldBindBodyWithJSON(req); err != nil {
				ctx.Error(&util.BizErr{
					Reason: err,
					Msg:    "参数错误: " + err.Error(),
				})
				return
			}
			if val, err := service.DeleteTagByIdList(req, ctx); err != nil {
				ctx.Error(err)
				return
			} else {
				ctx.JSON(http.StatusOK, util.RestWithSuccess(val))
			}
		})
	}
}
