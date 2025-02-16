package dto

import (
	"space-api/dto/query"
	"space-api/util/performance"
	"space-domain/model"
)

type (
	GetLogPagesReq struct {
		BasePageParam
		Conditions   []*query.WhereCond   `json:"conditions" yaml:"conditions" xml:"conditions" toml:"conditions"`
		OrderColumns []*query.OrderColumn `json:"orderColumns" yaml:"orderColumns" xml:"orderColumns" toml:"orderColumns"`
	}
	GetLogPagesResp = model.PageList[*model.LogInfo]

	DeleteLogReq struct {
		Conditions []*query.WhereCond `json:"conditions" yaml:"conditions" xml:"conditions" toml:"conditions"`
	}
	DeleteLogResp = performance.Empty

	DumpLogReq struct {
		Format string `form:"format" uri:"format" json:"format" yaml:"format" xml:"format" toml:"format"` // 导出格式
	}
)
