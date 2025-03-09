package service

import (
	"fmt"
	"space-api/dto"
	"space-api/util"
	"space-api/util/arr"
	"space-api/util/id"
	"space-domain/dao/biz"
	"space-domain/model"

	"github.com/gin-gonic/gin"
)

type (
	ICategoryService interface {
		CreateOrUpdateCategory(req *dto.CreateOrUpdateCategoryReq, ctx *gin.Context) (resp *dto.CreateOrUpdateCategoryResp, err error)
		GetCategoryByID(id int64, ctx *gin.Context) (resp *model.Category, err error)
		GetCategoryWithAllPosts(req *dto.GetCategoryWithPostsReq, ctx *gin.Context) (resp *dto.GetCategoryWithPostsResp, err error)
		GetAllCategories(req *dto.GetCategoryListReq, ctx *gin.Context) (resp *dto.GetCategoryListResp, err error)
		GetAllVisibleCategories(req *dto.GetCategoryListReq, ctx *gin.Context) (resp *dto.GetCategoryListResp, err error)
		GetCategoryWithVisiblePosts(req *dto.GetCategoryWithPostsReq, ctx *gin.Context) (resp *dto.GetCategoryWithPostsResp, err error)
		DeleteCategoryByIDList(req *dto.DeleteCategoryReq, ctx *gin.Context) (resp *dto.DeleteCategoryResp, err error)
	}
	categoryServiceImpl struct{}
)

var (
	_ ICategoryService = (*categoryServiceImpl)(nil)

	DefaultCategoryService ICategoryService = &categoryServiceImpl{}
)

func (c *categoryServiceImpl) CreateOrUpdateCategory(req *dto.CreateOrUpdateCategoryReq, ctx *gin.Context) (resp *dto.CreateOrUpdateCategoryResp, err error) {
	var tagID int64
	err = biz.Q.Transaction(func(tx *biz.Query) error {
		catTx := tx.Category
		// 查找匹配分类
		findCat, e := catTx.WithContext(ctx).Take()
		// 不存在分类, 直接创建
		if e != nil {
			// 获取新的 ID
			tagID = id.GetSnowFlakeNode().Generate().Int64()
			e := catTx.WithContext(ctx).Create(
				&model.Category{
					BaseColumn: model.BaseColumn{
						ID: findCat.ID,
					},
					CategoryName: req.CategoryName,
					Color:        req.Color,
					IconUrl:      req.IconUrl,
				},
			)
			if e != nil {
				return e
			}
		} else {
			// 存在的情况下进行更新
			tagID = findCat.ID
			_, e := catTx.WithContext(ctx).
				Select(
					catTx.ID,
					catTx.CategoryName,
					catTx.Color,
					catTx.IconUrl,
				).
				Updates(findCat)
			if e != nil {
				return e
			}

			// 查找关联到分类的相关文章
			postTx := tx.Post
			postList, e := postTx.WithContext(ctx).
				Where(postTx.Category.Eq(findCat.CategoryName)).
				Find()
			if e != nil {
				return e
			}
			// 更新文章
			for _, p := range postList {
				_, e := postTx.WithContext(ctx).Updates(p)
				if e != nil {
					return e
				}
			}
		}

		return nil
	})

	if err != nil {
		err = util.CreateBizErr("创建/更新文章分类失败", err)
		return
	}

	catTx := biz.Category
	find, e := catTx.WithContext(ctx).Where(catTx.ID.Eq(tagID)).Take()
	if e != nil {
		err = util.CreateBizErr("创建/更新文章分类失败", e)
		return
	}
	resp = &dto.CreateOrUpdateCategoryResp{
		Category: find,
	}
	return
}

func (c *categoryServiceImpl) GetCategoryByID(id int64, ctx *gin.Context) (resp *model.Category, err error) {
	resp, err = biz.Category.WithContext(ctx).Where(biz.Category.ID.Eq(id)).Take()
	if err != nil {
		err = util.CreateBizErr("分类不存在", err)
	}
	return
}

// GetCategoryWithAllPosts 获取所有的分类-文章关系, 不关心文章是否被隐藏/未公开
func (c *categoryServiceImpl) GetCategoryWithAllPosts(req *dto.GetCategoryWithPostsReq, ctx *gin.Context) (resp *dto.GetCategoryWithPostsResp, err error) {
	cat, err := c.GetCategoryByID(req.CatID, ctx)
	if err != nil {
		return
	}
	postTx := biz.Post
	postTx.Columns()
	postList, err := postTx.WithContext(ctx).
		Select(
			postTx.ID,
			postTx.CreatedAt,
			postTx.UpdatedAt,
			postTx.Hide,
			postTx.Title,
			postTx.AuthorId,
			// postTx.Content, // 减少数据携带量
			postTx.WordCount,
			postTx.Snippet,
			postTx.ReadTime,
			postTx.Category,
			postTx.Tags,
			postTx.LastPubTime,
			postTx.Weight,
			postTx.Views,
			postTx.UpVote,
			postTx.DownVote,
			postTx.AllowComment,
		).
		Where(postTx.Category.Eq(cat.CategoryName)).
		Find()
	if err != nil {
		err = util.CreateBizErr("查询分类的关联文章失败: "+err.Error(), err)
		return
	}

	resp = &dto.GetCategoryWithPostsResp{
		Category:      cat,
		RelationPosts: postList,
	}
	return
}

// GetAllCategories 获取所有的分类信息
func (c *categoryServiceImpl) GetAllCategories(req *dto.GetCategoryListReq, ctx *gin.Context) (resp *dto.GetCategoryListResp, err error) {
	catTx := biz.Category
	list, err := catTx.WithContext(ctx).Find()
	if err != nil {
		err = util.CreateBizErr("查询分类列表失败: "+err.Error(), err)
		return
	}
	resp = &dto.GetCategoryListResp{
		List: list,
	}

	return
}

// GetAllCategories 获取所有可见的分类信息
func (c *categoryServiceImpl) GetAllVisibleCategories(req *dto.GetCategoryListReq, ctx *gin.Context) (resp *dto.GetCategoryListResp, err error) {
	catTx := biz.Category
	list, err := catTx.WithContext(ctx).
		Where(catTx.Hide.Eq(0)).
		Find()
	if err != nil {
		err = util.CreateBizErr("查询分类列表失败: "+err.Error(), err)
		return
	}
	resp = &dto.GetCategoryListResp{
		List: list,
	}

	return
}

// GetCategoryWithVisiblePosts 获取 [分类-公开文章] 的数据
func (c *categoryServiceImpl) GetCategoryWithVisiblePosts(req *dto.GetCategoryWithPostsReq, ctx *gin.Context) (resp *dto.GetCategoryWithPostsResp, err error) {
	cat, err := c.GetCategoryByID(req.CatID, ctx)
	if err != nil {
		return
	}

	postTx := biz.Post
	postList, err := postTx.WithContext(ctx).
		Select(
			postTx.ID,
			postTx.CreatedAt,
			postTx.UpdatedAt,
			postTx.Hide,
			postTx.Title,
			postTx.AuthorId,
			// postTx.Content, // 减少数据携带量
			postTx.WordCount,
			postTx.ReadTime,
			postTx.Snippet,
			postTx.Category,
			postTx.Tags,
			postTx.LastPubTime,
			postTx.Weight,
			postTx.Views,
			postTx.UpVote,
			postTx.DownVote,
			postTx.AllowComment,
		).
		// 只查找公开的文章(非隐藏)
		Where(postTx.Category.Eq(cat.CategoryName), postTx.Hide.Eq(0)).
		Find()
	if err != nil {
		err = util.CreateBizErr("查询分类的关联文章失败: "+err.Error(), err)
		return
	}

	resp = &dto.GetCategoryWithPostsResp{
		Category:      cat,
		RelationPosts: postList,
	}

	return
}

func (c *categoryServiceImpl) DeleteCategoryByIDList(req *dto.DeleteCategoryReq, ctx *gin.Context) (resp *dto.DeleteCategoryResp, err error) {
	err = biz.Q.Transaction(func(tx *biz.Query) error {
		catTx := tx.Category
		catList, e := catTx.WithContext(ctx).
			Where(catTx.ID.In(req.IDList...)).
			Find()
		if e != nil {
			return e
		}
		postTx := tx.Post
		for _, cat := range catList {
			list, e := postTx.WithContext(ctx).
				Where(postTx.Category.Eq(cat.CategoryName)).
				Find()
			if e != nil {
				return e
			}
			// 存在相关文章的情况下, 不允许直接删除, 要先进行关系解绑
			if len(list) != 0 {
				return fmt.Errorf("存在关联的文章, 请先处理: %v",
					arr.MapSlice(list, func(_ int, p *model.Post) string {
						return p.Title
					}),
				)
			}
		}

		return nil
	})

	if err != nil {
		err = util.CreateBizErr("删除分类错误: "+err.Error(), err)
		return
	}

	resp = &dto.DeleteCategoryResp{}
	return
}
