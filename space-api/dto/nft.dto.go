package dto

import "space-api/util/performance"

type (
	AddNftBanIPReq struct {
		IPList []string `json:"ipList" yaml:"ipList" xml:"ipList" toml:"ipList"`
	}
	AddNftBanIPResp performance.Empty

	UnbindNftIPReq struct {
		IPList []string `json:"ipList" yaml:"ipList" xml:"ipList" toml:"ipList"`
	}
	UnbindNftIPResp performance.Empty
)
