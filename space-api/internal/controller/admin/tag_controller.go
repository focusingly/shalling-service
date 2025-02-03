package admin

import (
	"encoding/json"
	"net/http"
	"space-api/internal/service/v1"
	"space-api/util"
	"space-domain/dao/biz"

	"github.com/gin-gonic/gin"
)

func UseTagController(group *gin.RouterGroup) {
	tagGroup := group.GET("/tag")

	// 获取标签列表
	{
		tagGroup.GET("/list", func(ctx *gin.Context) {
			if val, err := service.SelectTagListByPage(ctx); err != nil {
				ctx.Error(&util.BizErr{
					Msg:    "未知错误",
					Reason: err,
				})
			} else {
				ctx.JSON(http.StatusOK, val)
			}
		})
	}

	// 根据 ID 获取单独的标签信息
	{
		type UriId struct {
			Id int64
		}
		tagGroup.GET("/:id", func(ctx *gin.Context) {
			uidStr := new(UriId)
			if err := ctx.ShouldBindUri(uidStr); err != nil {
				ctx.Error(&util.BizErr{
					Reason: err,
					Msg:    "参数错误: " + err.Error(),
				})

				return
			}

			if val, err := biz.Post.WithContext(ctx).Where(biz.Post.Id.Eq(uidStr.Id)).First(); err != nil {
				ctx.Error(&util.BizErr{
					Msg:    "未知错误",
					Reason: err,
				})

				return
			} else {
				ctx.JSON(http.StatusOK, util.RestWithSuccess(val))
			}
		})
	}

	// 更新或者删除已经存在的标签信息
	{
		tagGroup.POST("/", func(ctx *gin.Context) {
			biz.Q.Transaction(func(tx *biz.Query) error {
				postOp := tx.Post
				if p, err := postOp.WithContext(ctx).Where(postOp.Id.Eq(12)).First(); err != nil {
					ctx.Error(&util.BizErr{
						Reason: err,
						Msg:    "操作失败: " + err.Error(),
					})

					return err
				} else {
					ctx.JSON(http.StatusOK, util.RestWithSuccess(p))
				}

				return nil
			})
		})
	}

	// 删除文章标签
	{
		tagGroup.DELETE("/", func(ctx *gin.Context) {
			shouldDeletes := make([]int64, 8)
			defer ctx.Request.Body.Close()
			if err := json.NewDecoder(ctx.Request.Body).Decode(&shouldDeletes); err != nil {
				ctx.Error(&util.BizErr{
					Msg:    "参数序列化错误",
					Reason: err,
				})

				return
			}

			biz.Q.Transaction(func(tx *biz.Query) error {
				if _, err := tx.Post.WithContext(ctx).Where(tx.Post.Id.In(shouldDeletes...)).Delete(); err != nil {
					ctx.Error(&util.BizErr{
						Msg:    "删除数据失败: " + err.Error(),
						Reason: err,
					})
				} else {
					ctx.JSON(http.StatusOK, util.RestWithSuccess[any](nil))
				}

				return nil
			})

		})
	}
}
