package dto

import "space-domain/model"

type (
	CreateOrUpdateFriendLinkReq struct {
		ID          int64
		SiteURL     string
		Owner       string
		ShortName   string
		Available   int
		LogoURL     string
		Description *string
		BgURL       *string
	}
	CreateOrUpdateFriendLinkResp struct {
		*model.FriendLink
	}

	GetFriendLinksReq  struct{}
	GetFriendLinksResp struct {
		List []*model.FriendLink
	}

	DeleteFriendLinkReq struct {
		IDList []int64
	}
	DeleteFriendLinkResp struct{}
)
