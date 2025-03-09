package dto

import (
	"space-api/dto/query"
	"space-domain/model"
)

type (
	GetTagPageListReq struct {
		BasePageParam
	}
	GetTagPageListResp struct {
		model.PageList[*model.Tag]
	}

	GetTagDetailReq struct {
		Id int64 `uri:"id" json:"id,string"`
	}
	GetTagDetailResp struct {
		model.Tag
	}

	CreateOrUpdateTagReq struct {
		Id      int64   `json:"id,string"`
		Hide    int     `json:"hide"`
		TagName string  `json:"tagName"`
		Color   *string `json:"color"`
		IconUrl *string `json:"iconUrl"`
	}
	CreateOrUpdateTagResp struct {
		model.Tag
	}

	DeleteTagByIdListReq struct {
		WarningOverride
		IdList query.Int64Array `json:"idList"`
	}
	DeleteTagByIdListResp struct{}
)
