package dto

import (
	"space-api/dto/query"
	"space-domain/model"
)

type (
	QueryUvCountReq struct {
		WhereCondList []*query.WhereCond `json:"whereCondList" yaml:"whereCondList" xml:"whereCondList" toml:"whereCondList"`
	}

	GetUvPagesReq struct {
		BasePageParam
		WhereCondList []*query.WhereCond   `json:"whereCondList" yaml:"whereCondList" xml:"whereCondList" toml:"whereCondList"`
		SortList      []*query.OrderColumn `json:"sortList" yaml:"sortList" xml:"sortList" toml:"sortList"`
	}
	GetUvPagesResp struct {
		model.PageList[*model.UVStatistic]
	}

	DeleteUVReq struct {
		WhereCondList []*query.WhereCond `json:"whereCondList" yaml:"whereCondList" xml:"whereCondList" toml:"whereCondList"`
	}
	DeleteUVResp struct {
	}

	GetDailyCountReq struct {
		Date string `form:"date" json:"date" yaml:"date" xml:"date" toml:"date"`
	}

	GetUVTrendReq struct {
		Days int `form:"days" json:"day" yaml:"day" xml:"day" toml:"day"`
	}
)
