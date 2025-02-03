package admin

import (
	"net/http"
	"space-api/dto"
	"space-api/internal/service/v1"
	"space-api/util"
	"space-domain/dao/biz"
	"space-domain/model"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gen"
)

func UsePostController(group *gin.RouterGroup) {
	postGroup := group.Group("/posts")

	// 增加或者修改文章信息
	{
		postGroup.POST("/", func(ctx *gin.Context) {
			updatePostReq := &dto.UpdatePostReq{}
			if err := ctx.BindJSON(updatePostReq); err != nil {
				ctx.Error(&util.BizErr{
					Reason: err,
					Msg:    "参数错误: " + err.Error(),
				})
				return
			}
			if p, err := service.UpdateOrCreatePost(updatePostReq, ctx); err != nil {
				ctx.Error(&util.BizErr{
					Reason: err,
					Msg:    "操作失败",
				})
			} else {
				ctx.JSON(http.StatusOK, util.RestWithSuccess(p))
			}

		})
	}

	// 根据条件查询分页列表数据
	{
		type FindCond struct {
			Page int
			Size int
		}

		postGroup.GET("/list", func(ctx *gin.Context) {
			cond := new(FindCond)
			if err := ctx.ShouldBindQuery(cond); err != nil {
				ctx.Error(&util.BizErr{
					Reason: err,
					Msg:    err.Error(),
				})

				return
			}

			condList := []gen.Condition{}
			biz.Post.WithContext(ctx).
				Where(condList...).
				FindByPage(cond.Size*(cond.Page-1), cond.Size)
		})
	}

	// 根据 ID 查询单条文章数据
	{
		type PostId struct {
			Id int64 `uri:"id" binding:"required"`
		}

		postGroup.GET("/:id", func(ctx *gin.Context) {
			tmpId := new(PostId)
			if err := ctx.ShouldBindUri(tmpId); err != nil {
				ctx.Error(&util.BizErr{
					Msg:    "非法的请求参数: " + err.Error(),
					Reason: err,
				})

				return
			}

			query := biz.Post.WithContext(ctx)
			if result, err := query.Where(biz.Post.Id.Eq(tmpId.Id)).First(); err != nil {
				ctx.Error(&util.BizErr{
					Msg:    "获取数据失败: " + err.Error(),
					Reason: err,
				})
			} else {
				ctx.JSON(http.StatusOK, util.RestWithSuccess(result))
			}
		})
	}

	// 根据 ID 删除数据
	{
		type DeleteOp struct {
			Deletes []int64 `json:"deletes"`
		}

		postGroup.POST("/deletes", func(ctx *gin.Context) {
			tmp := new(DeleteOp)
			if err := ctx.MustBindWith(tmp, binding.JSON); err != nil {
				ctx.Error(&util.BizErr{
					Reason: err,
					Msg:    "参数错误: " + err.Error(),
				})

				return
			}
			deletes := []*model.Post{}
			for _, val := range tmp.Deletes {
				deletes = append(deletes, &model.Post{
					BaseColumn: model.BaseColumn{
						Id: val,
					},
				})
			}
			if v, err := biz.Post.WithContext(ctx).Delete(deletes...); err != nil {
				ctx.Error(&util.BizErr{
					Msg:    "删除失败: " + err.Error(),
					Reason: err,
				})
				return
			} else {
				ctx.JSON(http.StatusOK, util.RestWithSuccess(v.RowsAffected))
			}
		})
	}
}
