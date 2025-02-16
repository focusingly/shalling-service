package dto

import (
	"space-api/dto/query"
	"space-domain/model"
)

type (
	GetBlockingPagesReq struct {
		BasePageParam
		WhereCondList []*query.WhereCond   `json:"whereCondList" yaml:"whereCondList" xml:"whereCondList" toml:"whereCondList"`
		SortCondList  []*query.OrderColumn `json:"sortCondList" yaml:"sortCondList" xml:"sortCondList" toml:"sortCondList"`
	}
	GetBlockingPagesResp struct {
		model.PageList[*model.BlockIPRecord]
	}
	DeleteBlockingRecordReq struct {
		WhereCondList []*query.WhereCond `json:"whereCondList" yaml:"whereCondList" xml:"whereCondList" toml:"whereCondList"`
	}
	DeleteBlockingRecordResp struct{}
)
