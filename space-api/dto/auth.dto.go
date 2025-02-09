package dto

type (
	UserBasicData struct {
		UserType     string  `json:"userType" yaml:"userType" xml:"userType" toml:"userType"`
		IsAdmin      bool    `json:"isAdmin" yaml:"isAdmin" xml:"isAdmin" toml:"isAdmin"`
		IconURL      *string `json:"iconURL" yaml:"iconURL" xml:"iconURL" toml:"iconURL"`
		HomePageLink *string `json:"homePageLink" yaml:"homePageLink" xml:"homePageLink" toml:"homePageLink"`
		DisplayName  string  `json:"displayName" yaml:"displayName" xml:"displayName" toml:"displayName"`
		ExpiredAt    int64   `json:"expiredAt" yaml:"expiredAt" xml:"expiredAt" toml:"expiredAt"`
	}
)

type (
	GetLoginURLReq struct {
		OauthPlatform string `form:"oauthPlatform" json:"oauthPlatform"`
	}
	GetLoginURLResp string

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
