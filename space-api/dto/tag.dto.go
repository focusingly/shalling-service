package dto

import "space-domain/model"

type (
	GetTagPageListReq struct {
		BasePageParam
	}

	GetTagPageListResp struct {
		model.PageList[*model.Tag]
	}

	GetTagDetailReq struct {
		Id int64 `uri:"id" json:"id"`
	}

	GetTagDetailResp struct {
		model.Tag
	}

	CreateOrUpdateTagReq struct {
		Id      int64   `json:"id"`
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
		IdList []int64 `json:"idList"`
	}

	DeleteTagByIdListResp struct{}
)
