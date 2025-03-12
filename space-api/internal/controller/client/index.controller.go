package client

import (
	"fmt"
	"space-api/dto"
	"space-api/internal/service/v1"
	"space-api/internal/service/v1/user"
	"space-api/middleware/outbound"
	"space-api/util"
	"space-api/util/performance"
	"time"

	"github.com/gin-gonic/gin"
)

// 客户端的文章查询相关操作
func IndexController(group *gin.RouterGroup) {
	cachedGroup := performance.Group[any]{}
	postService := service.DefaultPostService
	categoryService := service.DefaultCategoryService
	tagService := service.DefaultTagService
	menuService := service.DefaultMenuService
	mediaService := service.DefaultMediaService
	userService := user.DefaultUserService
	searchService := service.DefaultGlobalSearchService

	indexGroup := group.Group(
		"/index",
		func(ctx *gin.Context) {
			// 禁止访问路径过长的请求
			if len(ctx.Request.RequestURI) >= 256 {
				err := fmt.Errorf("illegal request param, too long request uri: %d characters", len(ctx.Request.RequestURI))
				ctx.Error(util.CreateBizErr(
					"非法的请求参数",
					err,
				))
				ctx.Abort()
				return
			}
			ctx.Next()
		})

	// 查询公开文章的列表
	indexGroup.GET("/list", func(ctx *gin.Context) {
		cachedKey := ctx.Request.RequestURI

		resp, _, err := cachedGroup.Do(cachedKey, func() (value any, err error) {
			req := &dto.GetPostPageListReq{}
			if err = ctx.ShouldBindQuery(req); err != nil {
				err = util.CreateBizErr("参数错误: "+err.Error(), err)
				return
			}
			if resp, err := postService.GetVisiblePostsByPagination(req, ctx); err != nil {
				return nil, err
			} else {
				return resp, nil
			}
		}, time.Millisecond*500)

		if err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(resp, ctx)
		}
	})

	// 获取动态菜单列表
	indexGroup.GET("/menus", func(ctx *gin.Context) {
		cachedKey := ctx.Request.URL.Path

		resp, _, err := cachedGroup.Do(cachedKey, func() (value any, err error) {
			req := &dto.GetMenusReq{}
			if err = ctx.ShouldBindQuery(req); err != nil {
				err = util.CreateBizErr("参数错误: "+err.Error(), err)
				return
			}
			if resp, err := menuService.GetVisibleMenus(req, ctx); err != nil {
				return nil, err
			} else {
				return resp, nil
			}
		}, time.Millisecond*500)

		if err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(resp, ctx)
		}
	})

	// 获取所有允许公开展示的社交媒体标签
	indexGroup.GET("/pub-medias", func(ctx *gin.Context) {
		cachedKey := ctx.Request.URL.Path

		resp, _, err := cachedGroup.Do(cachedKey, func() (value any, err error) {
			req := &dto.GetMediaTagsReq{}
			if err = ctx.ShouldBindQuery(req); err != nil {
				err = util.CreateBizErr("参数错误: "+err.Error(), err)
				return
			}
			if resp, err := mediaService.GetVisibleMediaTags(req, ctx); err != nil {
				return nil, err
			} else {
				return resp, nil
			}
		}, time.Millisecond*500)

		if err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(resp, ctx)
		}
	})

	// 根据 ID 获取文章详情
	indexGroup.GET("/detail/:postID", func(ctx *gin.Context) {
		cachedKey := ctx.Request.URL.Path

		resp, _, err := cachedGroup.Do(cachedKey, func() (value any, err error) {
			req := &dto.GetPostDetailReq{}
			if err = ctx.ShouldBindUri(req); err != nil {
				err = util.CreateBizErr("参数错误: "+err.Error(), err)
				return
			}
			if resp, err := postService.GetVisiblePostById(req, ctx); err != nil {
				return nil, err
			} else {
				return resp, nil
			}
		}, time.Millisecond*500)
		if err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(resp, ctx)
		}
	})

	// 获取公开分类列表
	indexGroup.GET("/category/list", func(ctx *gin.Context) {
		cachedKey := ctx.Request.URL.Path
		resp, _, err := cachedGroup.Do(cachedKey, func() (value any, err error) {
			if resp, err := categoryService.GetAllVisibleCategories(&dto.GetCategoryListReq{}, ctx); err != nil {
				return nil, err
			} else {
				return resp, nil
			}
		}, time.Millisecond*500)
		if err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(resp, ctx)
		}
	})

	// 根据 ID 获取公开分类和其包含的文章列表
	indexGroup.GET("/category/subs/:catID", func(ctx *gin.Context) {
		cachedKey := ctx.Request.URL.Path
		resp, _, err := cachedGroup.Do(cachedKey, func() (value any, err error) {
			req := &dto.GetCategoryWithPostsReq{}
			if err = ctx.ShouldBindUri(req); err != nil {
				err = util.CreateBizErr("参数错误: "+err.Error(), err)
				return
			}
			if resp, err := categoryService.GetCategoryWithVisiblePosts(req, ctx); err != nil {
				return nil, err
			} else {
				return resp, nil
			}
		}, time.Millisecond*500)
		if err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(resp, ctx)
		}
	})

	// 获取公开的文章标签列表
	indexGroup.GET("/tag/list", func(ctx *gin.Context) {
		cachedKey := ctx.Request.RequestURI

		resp, _, err := cachedGroup.Do(cachedKey, func() (value any, err error) {
			req := &dto.GetTagPageListReq{}
			if err = ctx.ShouldBindQuery(req); err != nil {
				err = util.CreateBizErr("参数错误: "+err.Error(), err)
				return
			}

			if resp, err := tagService.GetVisibleTagPages(req, ctx); err != nil {
				return nil, err
			} else {
				return resp, nil
			}
		}, time.Millisecond*500)
		if err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(resp, ctx)
		}
	})

	// 获取已经登录的用户的基本信息(头像/主页链接等...)
	indexGroup.GET("/user/profile", func(ctx *gin.Context) {
		if resp, err := userService.GetLoginUserBasicProfile(ctx); err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(resp, ctx)
		}
	})

	// 根据标签名称查找公开的所有可见文章列表
	indexGroup.GET("/tag/relations", func(ctx *gin.Context) {
		cachedKey := ctx.Request.RequestURI
		resp, _, err := cachedGroup.Do(cachedKey, func() (value any, err error) {
			req := &dto.GetPostByTagNameReq{}
			if err = ctx.ShouldBindQuery(req); err != nil {
				err = util.CreateBizErr("参数错误: "+err.Error(), err)
				return

			}
			if resp, err := postService.GetVisiblePostsByTagName(req, ctx); err != nil {
				return nil, err
			} else {
				return resp, nil
			}
		}, time.Millisecond*500)
		if err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(resp, ctx)
		}
	})

	// 公开文章的全文搜索实现
	indexGroup.GET("/post/search", func(ctx *gin.Context) {
		cachedKey := ctx.Request.RequestURI
		resp, _, err := cachedGroup.Do(cachedKey, func() (value any, err error) {
			req := &dto.GlobalSearchReq{}
			if err = ctx.ShouldBindQuery(req); err != nil {
				err = util.CreateBizErr("参数错误: "+err.Error(), err)
				return
			}
			if resp, err := searchService.SearchKeywordPages(req, ctx); err != nil {
				return nil, err
			} else {
				return resp, nil
			}
		}, time.Millisecond*500)
		if err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(resp, ctx)
		}
	})
}
