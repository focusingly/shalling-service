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

	AddBlockingIPReq struct {
		IPAddr      string `json:"ipAddr" yaml:"ipAddr" xml:"ipAddr" toml:"ipAddr"`
		IPSource    string `json:"ipSource" yaml:"ipSource" xml:"ipSource" toml:"ipSource"`
		UserAgent   string `json:"userAgent" yaml:"userAgent" xml:"userAgent" toml:"userAgent"`
		LastRequest int64  `json:"lastRequest" yaml:"lastRequest" xml:"lastRequest" toml:"lastRequest"`
	}
	AddBlockingIPResp struct{}

	DeleteBlockingRecordReq struct {
		WhereCondList []*query.WhereCond `json:"whereCondList" yaml:"whereCondList" xml:"whereCondList" toml:"whereCondList"`
	}
	DeleteBlockingRecordResp struct{}
)
