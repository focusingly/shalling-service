package dto

import "space-domain/model"

type (
	CreateOrUpdateSocialMediaReq struct {
		Id          int64  `json:"id"`
		Hide        int    `json:"hide"`
		DisplayName string `gorm:"type:varchar(255);not null;comment:显示名称" json:"displayName"`
		IconURL     string `gorm:"type:varchar(255);not null;comment:图标链接" json:"iconURL"`
		OpenUrl     string `gorm:"type:varchar(255);not null;comment:跳转链接" json:"openUrl"`
	}
	CreateOrUpdateSocialMediaResp struct {
		*model.PubSocialMedia
	}

	GetMediaTagsReq  struct{}
	GetMediaTagsResp = []*model.PubSocialMedia

	DeleteSocialMediaByIdListReq struct {
		WarningOverride
		IdList []int64 `json:"idList"`
	}
	DeleteSocialMediaByIdListResp struct{}
)
