package dto

import "space-domain/model"

type (
	GetSocialMediaPageListReq struct {
		BasePageParam
	}

	GetSocialMediaPageListResp struct {
		model.PageList[*model.PubSocialMedia]
	}

	GetSocialMediaDetailReq struct {
		Id int64 `uri:"id" json:"id"`
	}

	GetSocialMediaDetailResp struct {
		model.PubSocialMedia
	}

	CreateOrUpdateSocialMediaReq struct {
		Id          int64  `json:"id"`
		Hide        byte   `json:"hide"`
		DisplayName string `gorm:"type:varchar(255);not null;comment:显示名称" json:"displayName"`
		IconURL     string `gorm:"type:varchar(255);not null;comment:图标链接" json:"iconURL"`
		OpenUrl     string `gorm:"type:varchar(255);not null;comment:跳转链接" json:"openUrl"`
	}

	CreateOrUpdateSocialMediaResp struct {
		model.PubSocialMedia
	}

	DeleteSocialMediaByIdListReq struct {
		WarningOverride
		IdList []int64 `json:"idList"`
	}

	DeleteSocialMediaByIdListResp struct{}
)
