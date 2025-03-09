package dto

type (
	UserBasicData struct {
		UserType     string  `json:"userType" yaml:"userType" xml:"userType" toml:"userType"`
		IsAdmin      bool    `json:"isAdmin" yaml:"isAdmin" xml:"isAdmin" toml:"isAdmin"`
		IconURL      *string `json:"iconURL" yaml:"iconURL" xml:"iconURL" toml:"iconURL"`
		HomePageLink *string `json:"homePageLink" yaml:"homePageLink" xml:"homePageLink" toml:"homePageLink"`
		DisplayName  string  `json:"displayName" yaml:"displayName" xml:"displayName" toml:"displayName"`
		ExpiredAt    int64   `json:"expiredAt,string" yaml:"expiredAt" xml:"expiredAt" toml:"expiredAt"`
	}
)

type (
	GetLoginURLReq struct {
		OauthPlatform string `form:"oauthPlatform" json:"oauthPlatform"`
	}
	GetLoginURLResp = string

	OauthLoginCallbackReq struct {
		Platform  string `json:"platform" yaml:"platform" xml:"platform" toml:"platform"`     // 登录平台
		GrantCode string `json:"grantCode" yaml:"grantCode" xml:"grantCode" toml:"grantCode"` // 授权码
		State     string `json:"state" yaml:"state" xml:"state" toml:"state"`                 // 随机状态标识
	}
	OauthLoginCallbackResp struct {
		Token         string `json:"token" yaml:"token" xml:"token" toml:"token"`
		UserBasicData `json:"userBasicData" yaml:"userBasicData" xml:"userBasicData" toml:"userBasicData"`
	}

	AdminLoginReq struct {
		Username   string `json:"username"`
		Password   string `json:"password"`
		VerifyCode string `json:"verifyCode"`
	}
	AdminLoginResp struct {
		Token         string `json:"token" yaml:"token" xml:"token" toml:"token"`
		UserBasicData `json:"userBasicData" yaml:"userBasicData" xml:"userBasicData" toml:"userBasicData"`
	}
)
