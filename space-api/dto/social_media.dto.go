package dto

import (
	"space-api/dto/query"
	"space-api/util/performance"
	"space-domain/model"
)

type (
	CreateOrUpdateSocialMediaReq struct {
		Id          int64  `json:"id,string"`
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
		IdList query.Int64Array `json:"idList"`
	}
	DeleteSocialMediaByIdListResp performance.Empty
)
