package dto

import (
	"space-api/dto/query"
	"space-api/util/performance"
	"space-domain/model"
)

type (
	CreateOrUpdateMenuReq struct {
		ID        int64             `json:"id,string"`
		MenuName  string            `json:"menuName"`
		RoutePath *string           `json:"routePath"`
		PostLink  *int64            `json:"postLink"`
		OpenWay   string            `json:"openWay"`
		SubLinks  []*model.MenuLink `json:"subLinks"`
	}
	CreateOrUpdateMenuResp struct {
		*model.MenuGroup
	}

	GetMenusReq  struct{}
	GetMenusResp struct {
		Menus []*model.MenuGroup `json:"menus" yaml:"menus" xml:"menus" toml:"menus"`
	}

	DeleteMenuGroupsReq struct {
		IDList query.Int64Array `json:"idList" yaml:"idList" xml:"idList" toml:"idList"`
	}
	DeleteMenuGroupsResp = performance.Empty
)
