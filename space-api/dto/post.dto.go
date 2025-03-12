package dto

import (
	"space-api/dto/query"
	"space-domain/model"
)

type (
	UpdateOrCreatePostReq struct {
		WarningOverride `json:"warningOverride" yaml:"warningOverride" xml:"warningOverride" toml:"warningOverride"`
		*model.Post
	}

	UpdateOrCreatePostResp struct {
		*model.Post
	}

	GetPostPageListReq struct {
		BasePageParam
	}
	GetPostPageListResp struct {
		model.PageList[*model.Post]
	}

	GetPostDetailReq struct {
		PostID int64 `uri:"postID" binding:"required" json:"postID,string" yaml:"postID" xml:"postID" toml:"postID"`
	}
	GetPostDetailResp struct {
		model.Post
	}

	DeletePostByIdListReq struct {
		IdList query.Int64Array `json:"idList"`
	}
	DeletePostByIdListResp struct {
	}

	GetPostByTagNameReq struct {
		TagName string `form:"tagName" json:"tagName" yaml:"tagName" xml:"tagName" toml:"tagName"`
	}
	GetPostByTagNameResp struct {
		Tag   *model.Tag    `json:"tag" yaml:"tag" xml:"tag" toml:"tag"`
		Posts []*model.Post `json:"posts" yaml:"posts" xml:"posts" toml:"posts"`
	}
)
