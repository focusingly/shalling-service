package dto

import (
	"space-api/dto/query"
	"space-domain/model"
)

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
		List []*model.Category `json:"list" yaml:"list" xml:"list" toml:"list"`
	}

	// 获取分类和其文章关联的所有文章列表信息
	GetCategoryWithPostsReq struct {
		CatID int64 `uri:"catID" json:"catID,string"`
	}
	GetCategoryWithPostsResp struct {
		Category      *model.Category `json:"category" yaml:"category" xml:"category" toml:"category"`
		RelationPosts []*model.Post   `json:"relationPosts" yaml:"relationPosts" xml:"relationPosts" toml:"relationPosts"`
	}

	DeleteCategoryReq struct {
		IDList query.Int64Array `json:"idList" yaml:"idList" xml:"idList" toml:"idList"`
	}
	DeleteCategoryResp struct{}
)
