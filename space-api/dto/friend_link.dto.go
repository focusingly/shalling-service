package dto

import (
	"space-api/dto/query"
	"space-domain/model"
)

type (
	CreateOrUpdateFriendLinkReq struct {
		ID          int64   `json:"id,string" yaml:"id" xml:"id" toml:"id"`
		SiteURL     string  `json:"siteURL" yaml:"siteURL" xml:"siteURL" toml:"siteURL"`
		Owner       string  `json:"owner" yaml:"owner" xml:"owner" toml:"owner"`
		ShortName   string  `json:"shortName" yaml:"shortName" xml:"shortName" toml:"shortName"`
		Available   int     `json:"available" yaml:"available" xml:"available" toml:"available"`
		LogoURL     string  `json:"logoURL" yaml:"logoURL" xml:"logoURL" toml:"logoURL"`
		Description *string `json:"description" yaml:"description" xml:"description" toml:"description"`
		BgURL       *string `json:"bgURL" yaml:"bgURL" xml:"bgURL" toml:"bgURL"`
	}
	CreateOrUpdateFriendLinkResp struct {
		*model.FriendLink
	}

	GetFriendLinksReq  struct{}
	GetFriendLinksResp struct {
		List []*model.FriendLink `json:"list" yaml:"list" xml:"list" toml:"list"`
	}

	DeleteFriendLinkReq struct {
		IDList query.Int64Array `json:"idList" yaml:"idList" xml:"idList" toml:"idList"`
	}
	DeleteFriendLinkResp struct{}
)
