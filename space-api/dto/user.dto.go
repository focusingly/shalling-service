package dto

import "space-domain/model"

type (
	GetLoginUserSessionsReq struct {
		BasePageParam `json:"basePageParam" yaml:"basePageParam" xml:"basePageParam" toml:"basePageParam"`
	}
	GetLoginUserSessionsResp struct {
		model.PageList[*model.UserLoginSession]
	}

	UpdateOauthUserReq struct {
		UserID    int64 `json:"useID" yaml:"useID" xml:"useID" toml:"useID"`
		Available bool  `json:"available" yaml:"available" xml:"available" toml:"available"`
	}
	UpdateOauthUserResp struct{}

	DeleteOauth2UserReq struct {
		IDList []int64
	}
	DeleteOauth2UserResp struct{}

	UpdateLocalUserBasicReq struct {
		UserID       int64   `json:"userID" yaml:"userID" xml:"userID" toml:"userID"`
		Email        *string `gorm:"type:varchar(255);null;comment:用户邮箱, 可用于找回密码" json:"email" yaml:"email" xml:"email" toml:"email"`
		Username     string  `gorm:"type:varchar(255);not null;unique;comment:登录的用户名称" json:"username" yaml:"username" xml:"username" toml:"username"`
		DisplayName  string  `gorm:"type:varchar(255);not null;comment:对外展示的用户名称" json:"displayName" yaml:"displayName" xml:"displayName" toml:"displayName"`
		Password     string  `gorm:"type:text;not null;comment:可用于找回账户的密码" json:"password" yaml:"password" xml:"password" toml:"password"`
		AvatarURL    *string `gorm:"type:text;null;comment:用户的头像链接" json:"avatarURL" yaml:"avatarURL" xml:"avatarURL" toml:"avatarURL"`
		HomepageLink *string `gorm:"type:text;null;comment:用户的主页链接" json:"homepageLink" yaml:"homepageLink" xml:"homepageLink" toml:"homepageLink"`
		Phone        *string `gorm:"type:varchar(255);null;comment:可用于找回账户的密码" json:"phone" yaml:"phone" xml:"phone" toml:"phone"`
		IsAdmin      int     `gorm:"type:smallint;default:0;comment:是否为超级管理员用户(大于 0 的都可以认为是)" json:"isAdmin" yaml:"isAdmin" xml:"isAdmin" toml:"isAdmin"`
	}
	UpdateLocalUserResp struct{}

	UpdateLocalUserPassReq struct {
		UserID      int64  `json:"userID" yaml:"userID" xml:"userID" toml:"userID"`
		OldPassword string `json:"oldPassword" yaml:"oldPassword" xml:"oldPassword" toml:"oldPassword"`
		NewPassword string `json:"newPassword" yaml:"newPassword" xml:"newPassword" toml:"newPassword"`
	}
	UpdateLocalUserPassResp struct{}

	ExpireUserLoginSessionReq struct {
		UUIDList []string `json:"uuidList" yaml:"uuidList" xml:"uuidList" toml:"uuidList"`
	}
	ExpireUserLoginSessionResp struct{}
)
