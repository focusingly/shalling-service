package dto

import "space-domain/model"

type (
	CreateOrUpdateCategoryReq struct {
		CategoryName string  `json:"categoryName"`
		Color        *string `json:"color"`
		IconUrl      *string `json:"iconUrl"`
	}
	CreateOrUpdateCategoryResp struct {
		*model.Category
	}

	// TODO 暂不考虑分页, 分类列表实际不会太多
	GetCategoryListReq  struct{}
	GetCategoryListResp struct {
		List []*model.Category
	}

	// 获取分类和其文章关联的所有文章列表信息
	GetCategoryWithPostsReq struct {
		CatID int64 `uri:"catID" json:"catID"`
	}
	GetCategoryWithPostsResp struct {
		Category      *model.Category `json:"category"`
		RelationPosts []*model.Post   `json:"relationPosts"`
	}

	DeleteCategoryReq struct {
		IDList []int64
	}
	DeleteCategoryResp struct{}
)
