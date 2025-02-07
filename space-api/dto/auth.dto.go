package dto

type (
	UserBasicData struct {
		UserType     string  `json:"userType"`
		IsAdmin      bool    `json:"isAdmin"`
		IconURL      *string `json:"iconURL"`
		HomePageLink *string `json:"homepageLink"`
		DisplayName  string  `json:"displayName"`
		ExpiredAt    int64   `json:"expiredAt"`
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
		Token         string `json:"token"`
		UserBasicData `json:"userBasicData"`
	}
)
