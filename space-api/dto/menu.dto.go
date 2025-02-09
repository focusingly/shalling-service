package dto

import "space-domain/model"

type (
	CreateOrUpdateMenuReq struct {
		ID        int64             `json:"id"`
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
		IDList []int64 `json:"idList" yaml:"idList" xml:"idList" toml:"idList"`
	}
	DeleteMenuGroupsResp struct {
	}
)
