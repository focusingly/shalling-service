package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"space-api/conf"
	"space-api/constants"
	"space-api/dto"
	"space-api/middleware/inbound"
	"space-api/util"
	"space-api/util/arr"
	"space-api/util/encrypt"
	"space-api/util/id"
	"space-api/util/ip"
	"space-api/util/performance"
	"space-api/util/ptr"
	"space-api/util/verify"
	"space-domain/dao/biz"
	"space-domain/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/sync/errgroup"
)

type (
	GithubPub struct {
		Login             string    `json:"login"`
		ID                int64     `json:"id"`
		NodeID            string    `json:"node_id"`
		AvatarURL         string    `json:"avatar_url"`
		GravatarID        string    `json:"gravatar_id"`
		URL               string    `json:"url"`
		HtmlURL           string    `json:"html_url"`
		FollowersURL      string    `json:"followers_url"`
		FollowingURL      string    `json:"following_url"`
		GistsURL          string    `json:"gists_url"`
		StarredURL        string    `json:"starred_url"`
		SubscriptionsURL  string    `json:"subscriptions_url"`
		OrganizationsURL  string    `json:"organizations_url"`
		ReposURL          string    `json:"repos_url"`
		EventsURL         string    `json:"events_url"`
		ReceivedEventsURL string    `json:"received_events_url"`
		Type              string    `json:"type"`
		UserViewType      string    `json:"user_view_type"`
		SiteAdmin         bool      `json:"site_admin"`
		Name              string    `json:"name"`
		Blog              string    `json:"blog"`
		PublicRepos       int64     `json:"public_repos"`
		PublicGists       int64     `json:"public_gists"`
		Followers         int64     `json:"followers"`
		Following         int64     `json:"following"`
		CreatedAt         time.Time `json:"created_at"`
		UpdatedAt         time.Time `json:"updated_at"`
	}

	GithubEmailElement struct {
		Email      string  `json:"email"`
		Primary    bool    `json:"primary"`
		Verified   bool    `json:"verified"`
		Visibility *string `json:"visibility"`
	}

	GooglePub struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Name          string `json:"name"`
		GivenName     string `json:"given_name"`
		Picture       string `json:"picture"`
	}
)

var (
	githubOauth2Config, googleOauth2Config *oauth2.Config
	// 本地缓存空间
	authCache = inbound.GetMiddlewareRelativeAuthCache()
	_jwtConf  = conf.ProjectConf.GetJwtConf()
	_appConf  = conf.ProjectConf.GetAppConf()
)

type (
	IAuthService interface {
		AdminLogin(req *dto.AdminLoginReq, ctx *gin.Context) (resp *dto.AdminLoginResp, err error)
		HandleOauthLogin(req *dto.OauthLoginCallbackReq, ctx *gin.Context) (resp *dto.OauthLoginCallbackResp, err error)
		GetOauth2LoginGrantURL(req *dto.GetLoginURLReq, ctx *gin.Context) (resp dto.GetLoginURLResp, err error)
		ParseGithubCallback(req *dto.OauthLoginCallbackReq, ctx *gin.Context) (resp *model.OAuth2User, err error)
		GetGoogleLoginURL(ctx *gin.Context) (val string, err error)
		ParseGoogleCallback(req *dto.OauthLoginCallbackReq, ctx *gin.Context) (resp *model.OAuth2User, err error)
		GetRefreshOauth2Data(user *model.OAuth2User, ctx context.Context) (resp *model.OAuth2User, err error)
		GetGithubUserProfile(authClient *http.Client) (resp *OauthBasicUserProfile, err error)
		GetGoogleUserProfile(authClient *http.Client) (resp *OauthBasicUserProfile, err error)
		CurrentUserLogout(ctx *gin.Context) (resp *performance.Empty, err error)
	}
	authServiceImpl struct {
	}

	boData struct {
		UserID   int64
		UserType string
	}
)

var (
	_ IAuthService = (*authServiceImpl)(nil)

	DefaultAuthService IAuthService = &authServiceImpl{}
)

func (s *authServiceImpl) updateUserLoginSession(user *boData, bizTx *biz.Query, ctx *gin.Context) (resp *model.UserLoginSession, err error) {
	err = bizTx.Transaction(func(tx *biz.Query) error {
		loginSessionTx := tx.UserLoginSession
		// 获取用户所有已经登录的会话信息
		existsSessions, getSessionErr := loginSessionTx.WithContext(ctx).
			Where(loginSessionTx.UserID.Eq(user.UserID)).
			Order(loginSessionTx.ID.Desc()). // 进行排序, 比较新的数据放在前面
			Find()
		if getSessionErr != nil {
			return util.CreateBizErr("设置会话信息失败", getSessionErr)
		}

		nowEpochMill := time.Now().UnixMilli()
		// 更新列表
		updates := slices.DeleteFunc(
			slices.Clone(existsSessions),
			// 先轻量掉所有已经过期的时间戳
			func(s *model.UserLoginSession) bool {
				return nowEpochMill >= s.ExpiredAt
			},
		)
		if len(updates) >= _appConf.MaxUserActive-1 {
			// 淘汰最旧的数据(最大允许用户数来自系统配置)
			updates = updates[:_appConf.MaxUserActive-1]
		}

		ipAddr := inbound.GetRealIpWithContext(ctx)
		ua := inbound.GetUserAgentFromContext(ctx)
		to32Ip, _ := ip.Ipv4StringToU32(ipAddr)
		ipSource, _ := ip.GetIpSearcher().SearchByStr(ipAddr)
		// 新的用户会话信息
		newLoginSession := &model.UserLoginSession{
			BaseColumn: model.BaseColumn{
				ID: id.GetSnowFlakeNode().Generate().Int64(),
			},
			UserID:     user.UserID,
			IpU32Val:   &to32Ip,
			IpAddress:  &ipAddr,
			IpSource:   &ipSource,
			ExpiredAt:  time.Now().Add(_jwtConf.ParsedExpTime).UnixMilli(),
			UserType:   user.UserType,
			Useragent:  ua.Useragent,
			ClientName: ua.ClientName,
			OSName:     ua.OS,
		}

		token, getTokenErr := verify.CreateJwtToken(newLoginSession)
		if getTokenErr != nil {
			return getTokenErr
		}

		// 设置新的 token
		newLoginSession.Token = token
		// 存入新的会话
		updates = append(updates, newLoginSession)
		// 删除所有已存在的列表
		_, deleteSessionErr := loginSessionTx.WithContext(ctx).
			Where(loginSessionTx.ID.In(
				arr.MapSlice(
					existsSessions,
					func(_ int, s *model.UserLoginSession) int64 {
						return s.ID
					})...,
			)).
			Delete()
		if deleteSessionErr != nil {
			return deleteSessionErr
		}
		// 批量重新创建会话信息
		createSessionErr := loginSessionTx.WithContext(ctx).
			CreateInBatches(updates, appConf.MaxUserActive)
		if createSessionErr != nil {
			return createSessionErr
		}

		// set value
		resp = newLoginSession
		// 更新缓存中的数据
		cacheSpace := inbound.GetMiddlewareRelativeAuthCache()
		// 删除已经登录的会话信息
		for _, u := range existsSessions {
			cacheSpace.Delete(fmt.Sprintf("%d", u.ID))
		}
		now := time.Now()
		// 重新设置登录会话信息
		for _, u := range updates {
			cacheSpace.Set(fmt.Sprintf("%d", u.ID), u, time.UnixMilli(u.ExpiredAt).Sub(now))
		}

		return nil
	})

	if err != nil {
		err = util.CreateAuthErr("创建/更新登录会话失败", err)
		return
	}

	return
}

// AdminLogin 后台管理员登录
func (s *authServiceImpl) AdminLogin(req *dto.AdminLoginReq, ctx *gin.Context) (resp *dto.AdminLoginResp, err error) {
	err = biz.Q.Transaction(func(tx *biz.Query) error {
		localUserTx := tx.LocalUser

		// 数据库中的本地账户
		findUser, e := localUserTx.WithContext(ctx).Where(
			localUserTx.Username.Eq(req.Username),
		).Take()
		if e != nil {
			return fmt.Errorf("不存在的用户或密码不正确")
		}
		if !encrypt.ComparePassword(req.Password, findUser.Password) {
			return fmt.Errorf("不存在的用户或密码不正确")
		}

		loginSession, e := s.updateUserLoginSession(&boData{
			UserID:   findUser.ID,
			UserType: constants.Admin,
		}, tx, ctx)
		if e != nil {
			return e
		}

		resp = &dto.AdminLoginResp{
			Token: loginSession.Token,
			UserBasicData: dto.UserBasicData{
				UserType:     loginSession.UserType,
				IsAdmin:      findUser.IsAdmin > 0,
				IconURL:      findUser.AvatarURL,
				HomePageLink: findUser.HomepageLink,
				DisplayName:  findUser.DisplayName, // 只展示对外公开的用户名, 降低攻击概率
				ExpiredAt:    loginSession.ExpiredAt,
			},
		}

		return nil
	})

	if err != nil {
		err = util.CreateAuthErr("登录失败: "+err.Error(), err)
		return
	}

	return
}

// HandleOauthLogin 处理 Oauth2 登录
func (s *authServiceImpl) HandleOauthLogin(req *dto.OauthLoginCallbackReq, ctx *gin.Context) (resp *dto.OauthLoginCallbackResp, err error) {
	var parsed *model.OAuth2User
	var e error
	switch req.Platform {
	case constants.GithubUser:
		parsed, e = s.ParseGoogleCallback(req, ctx)
	case constants.GoogleUser:
		parsed, e = s.ParseGoogleCallback(req, ctx)
	default:
		err = util.CreateAuthErr("不受支持的平台", fmt.Errorf("un-support auth platform: %s", req.Platform))
		return
	}
	if e != nil {
		err = util.CreateAuthErr("授权失败", err)
		return
	}

	// 添加/更新 三方账户
	err = biz.Q.Transaction(func(tx *biz.Query) error {
		var userID int64
		oauth2Tx := tx.OAuth2User
		find, e := oauth2Tx.WithContext(ctx).
			Where(
				oauth2Tx.PlatformUserId.Eq(parsed.PlatformUserId),
				oauth2Tx.PlatformName.Eq(parsed.PlatformName),
			).
			Take()
		// 第一次登录
		if e == nil {
			userID = id.GetSnowFlakeNode().Generate().Int64()

			e := oauth2Tx.WithContext(ctx).Create(
				&model.OAuth2User{
					BaseColumn: model.BaseColumn{
						ID: userID,
					},
					PlatformName:   parsed.PlatformName,
					PlatformUserId: parsed.PlatformUserId,
					Username:       parsed.Username,
					PrimaryEmail:   parsed.PrimaryEmail,
					AccessToken:    parsed.AccessToken,
					RefreshToken:   parsed.RefreshToken,
					ExpiredAt:      parsed.ExpiredAt,
					AvatarURL:      parsed.AvatarURL,
					HomepageLink:   parsed.HomepageLink,
					Scopes:         parsed.Scopes,
					Enable:         1,
				},
			)
			if e != nil {
				return e
			}
		} else {
			userID = find.ID
			// 阻止封禁的用户
			if find.Enable == 0 {
				e = fmt.Errorf("该账户封禁中")
				return e
			}
			_, e = oauth2Tx.WithContext(ctx).Updates(&model.OAuth2User{
				BaseColumn:     find.BaseColumn,
				PlatformName:   parsed.PlatformName,
				PlatformUserId: parsed.PlatformUserId,
				Username:       parsed.Username,
				PrimaryEmail:   parsed.PrimaryEmail,
				AccessToken:    parsed.AccessToken,
				RefreshToken:   parsed.RefreshToken,
				ExpiredAt:      parsed.ExpiredAt,
				AvatarURL:      parsed.HomepageLink,
				HomepageLink:   parsed.HomepageLink,
				Scopes:         parsed.Scopes,
			})
			if e != nil {
				return e
			}
		}

		u, e := oauth2Tx.WithContext(ctx).
			Where(oauth2Tx.ID.Eq(userID)).
			Take()
		if e != nil {
			return fmt.Errorf("同步数据失败")
		}

		loginSession, e := s.updateUserLoginSession(
			&boData{
				UserID:   userID,
				UserType: req.Platform,
			},
			tx,
			ctx,
		)
		if e != nil {
			return e
		}
		resp = &dto.OauthLoginCallbackResp{
			Token: loginSession.Token,
			UserBasicData: dto.UserBasicData{
				UserType:     loginSession.UserType,
				IsAdmin:      false,
				IconURL:      u.AvatarURL,
				HomePageLink: u.HomepageLink,
				DisplayName:  u.Username,
				ExpiredAt:    loginSession.ExpiredAt,
			},
		}

		return nil
	})

	if err != nil {
		err = util.CreateAuthErr("登录失败", err)
		return
	}

	return
}

func (*authServiceImpl) GetOauth2LoginGrantURL(req *dto.GetLoginURLReq, ctx *gin.Context) (resp dto.GetLoginURLResp, err error) {
	state := uuid.NewString()

	// 设置过期时间
	err = authCache.Set(state, new(performance.Empty), time.Minute*5)
	if err != nil {
		err = util.CreateBizErr("设置校验状态失败: "+err.Error(), err)
		return
	}
	switch req.OauthPlatform {
	case "github":
		resp = githubOauth2Config.AuthCodeURL(state)
	case "google":
		resp = googleOauth2Config.AuthCodeURL(state)
	default:
		err = util.CreateBizErr(
			"暂不支持的验证平台: "+req.OauthPlatform,
			fmt.Errorf("un-support oauth2 platform: %s", req.OauthPlatform),
		)
		return
	}

	return
}

func (s *authServiceImpl) ParseGithubCallback(req *dto.OauthLoginCallbackReq, ctx *gin.Context) (resp *model.OAuth2User, err error) {
	// 判断授权码
	if req.GrantCode == "" || req.State == "" {
		err = util.CreateAuthErr(
			"获取授权信息失败, 请重试",
			fmt.Errorf("grant code not exits"),
		)
		return
	}

	// 判断 state
	// 判断缓存里的情况
	if err = authCache.GetAndDel(req.State, &performance.Empty{}); err != nil {
		return
	}
	// 使用授权码
	oauthToken, err := githubOauth2Config.Exchange(ctx, req.GrantCode)
	if err != nil {
		err = util.CreateAuthErr(
			"获取授权信息失败, 请重试",
			err,
		)

		return
	}
	client := githubOauth2Config.Client(ctx, oauthToken)
	if githubUserProfile, e := s.GetGithubUserProfile(client); e != nil {
		err = util.CreateAuthErr("获取用户信息失败, 请重试", e)
		return
	} else {
		resp = &model.OAuth2User{
			PlatformName:   constants.GithubUser,
			PlatformUserId: githubUserProfile.ID,
			Username:       githubUserProfile.Username,
			PrimaryEmail:   githubUserProfile.PrimaryEmail,
			AccessToken:    oauthToken.AccessToken,
			RefreshToken:   oauthToken.RefreshToken,
			ExpiredAt:      &oauthToken.ExpiresIn,
			AvatarURL:      githubUserProfile.AvatarURL,
			HomepageLink:   githubUserProfile.HomepageLink,
			Scopes:         githubOauth2Config.Scopes,
		}
	}

	return
}

func (*authServiceImpl) GetGoogleLoginURL(ctx *gin.Context) (val string, err error) {
	state := uuid.NewString()
	if err = authCache.Set(state, &performance.Empty{}, time.Minute*5); err != nil {
		return
	}
	val = googleOauth2Config.AuthCodeURL(state)
	return
}

func (s *authServiceImpl) ParseGoogleCallback(req *dto.OauthLoginCallbackReq, ctx *gin.Context) (resp *model.OAuth2User, err error) {
	// 基本校验
	if req.GrantCode == "" || req.State == "" {
		err = util.CreateAuthErr(
			"凭据校验失败",
			fmt.Errorf("the principal is illegal"),
		)
		return
	}
	if err = authCache.GetAndDel(req.State, &performance.Empty{}); err != nil {
		err = util.CreateAuthErr(
			"凭据校验失败",
			err,
		)
		return
	}

	oauthToken, err := googleOauth2Config.Exchange(ctx, req.GrantCode)
	if err != nil {
		err = util.CreateAuthErr(
			"获取凭证失败"+err.Error(),
			err,
		)

		return
	}

	client := googleOauth2Config.Client(ctx, oauthToken)
	googleUserProfile, err := s.GetGoogleUserProfile(client)

	if err != nil {
		return
	}

	resp = &model.OAuth2User{
		PlatformName:   constants.GoogleUser,
		PlatformUserId: googleUserProfile.ID,
		Username:       googleUserProfile.Username,
		PrimaryEmail:   googleUserProfile.PrimaryEmail,
		AccessToken:    oauthToken.AccessToken,
		RefreshToken:   oauthToken.RefreshToken,
		ExpiredAt:      &oauthToken.ExpiresIn,
		AvatarURL:      googleUserProfile.AvatarURL,
		HomepageLink:   googleUserProfile.HomepageLink,
		Scopes:         googleOauth2Config.Scopes,
	}

	return
}

// GetRefreshOauth2Data 获取新的用户凭据
func (s *authServiceImpl) GetRefreshOauth2Data(user *model.OAuth2User, ctx context.Context) (resp *model.OAuth2User, err error) {
	var userProfile *OauthBasicUserProfile
	var newToken oauth2.Token

	switch user.PlatformName {
	case constants.GithubUser:
		// github 目前返回的 token 没有设置过期时间, 除非被用户取消授权
		newToken = oauth2.Token{
			AccessToken:  user.AccessToken,
			RefreshToken: user.RefreshToken,
		}
		userProfile, err = s.GetGithubUserProfile(githubOauth2Config.Client(ctx, &oauth2.Token{
			AccessToken: user.AccessToken,
			TokenType:   "Bearer",
		}))

		if err != nil {
			return
		}
	case constants.GoogleUser:
		t, e := googleOauth2Config.TokenSource(ctx, &oauth2.Token{
			RefreshToken: user.RefreshToken,
		}).Token()
		if e != nil {
			err = util.CreateAuthErr("获取新的凭据失败", e)
			return
		}
		newToken = oauth2.Token{
			AccessToken: t.AccessToken,
			RefreshToken: util.TernaryExpr(
				t.RefreshToken != "",
				t.RefreshToken,
				user.RefreshToken,
			),
			ExpiresIn: t.ExpiresIn,
		}
		// 获取新的配置文件
		userProfile, err = s.GetGoogleUserProfile(githubOauth2Config.Client(ctx, t))
		if err != nil {
			return
		}
	default:
		return nil, fmt.Errorf("un-support auth platform: %s", user.PlatformName)
	}

	resp = &model.OAuth2User{
		BaseColumn:     user.BaseColumn,
		PlatformName:   user.PlatformName,
		PlatformUserId: user.PlatformUserId,
		Username:       userProfile.Username,
		PrimaryEmail:   user.PrimaryEmail,
		AccessToken:    newToken.AccessToken,
		RefreshToken:   newToken.RefreshToken,
		ExpiredAt:      ptr.ToPtr(newToken.ExpiresIn),
		AvatarURL:      userProfile.AvatarURL,
		HomepageLink:   userProfile.HomepageLink,
		Scopes:         user.Scopes,
		Enable:         user.Enable,
	}

	return
}

type OauthBasicUserProfile struct {
	ID           string
	Username     string
	PrimaryEmail string
	HomepageLink *string
	AvatarURL    *string
}

func (s *authServiceImpl) GetGithubUserProfile(authClient *http.Client) (resp *OauthBasicUserProfile, err error) {
	var group errgroup.Group
	var primaryEmail string
	githubPubDetail := &GithubPub{}
	emailList := []*GithubEmailElement{}

	// 读取公开信息
	group.Go(func() error {
		res, err := authClient.Get("https://api.github.com/user")
		if err != nil {
			err = util.CreateBizErr(
				"获取用户数据失败: "+err.Error(),
				err,
			)
			return err
		}

		defer res.Body.Close()
		// 获取公开信息
		if err = json.NewDecoder(res.Body).Decode(githubPubDetail); err != nil {
			err = util.CreateBizErr(
				"获取用户数据失败: "+err.Error(),
				err,
			)
			return err
		}
		return nil
	})

	// 获取主邮箱
	group.Go(func() error {
		// 获取用户私人电子邮件地址
		emailResp, e := authClient.Get("https://api.github.com/user/emails")
		if e != nil {
			e = &util.BizErr{
				Msg: e.Error(),
			}
			return e
		}
		if emailResp != nil {
			defer emailResp.Body.Close()
			if e := json.NewDecoder(emailResp.Body).Decode(&emailList); e != nil {
				// 获取公开信息
				e = util.CreateBizErr(
					"获取用户数据失败: "+e.Error(),
					e,
				)
				return e
			}
			if len(emailList) == 0 {
				return fmt.Errorf("can't get primary email")
			}
			for _, el := range emailList {
				if el.Primary {
					primaryEmail = el.Email
					return nil
				}
			}
			return fmt.Errorf("can't get primary email")
		}

		return nil
	})

	err = group.Wait()
	if err == nil {
		resp = &OauthBasicUserProfile{
			ID:           fmt.Sprintf("%d", githubPubDetail.ID),
			Username:     githubPubDetail.Login,
			PrimaryEmail: primaryEmail,
			HomepageLink: &githubPubDetail.HtmlURL,
			AvatarURL:    &githubPubDetail.AvatarURL,
		}
	}

	return

}

func (s *authServiceImpl) GetGoogleUserProfile(authClient *http.Client) (resp *OauthBasicUserProfile, err error) {
	res, err := authClient.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil || res.Body == nil {
		err = util.CreateBizErr(
			"获取用户数据失败: "+err.Error(),
			err,
		)
		return
	}
	defer res.Body.Close()
	googlePubDetail := GooglePub{}
	if err = json.NewDecoder(res.Body).Decode(&googlePubDetail); err != nil {
		err = util.CreateBizErr(
			"解析用户数据失败: "+err.Error(),
			err,
		)
		return
	}

	resp = &OauthBasicUserProfile{
		ID:           googlePubDetail.ID,
		Username:     googlePubDetail.Name,
		PrimaryEmail: googlePubDetail.Email,
		HomepageLink: nil,
		AvatarURL:    &googlePubDetail.Picture,
	}

	return
}

// CurrentUserLogout 当前已登录的用户退出登录(根据请求头中携带的 Bearer Token 判断)
func (*authServiceImpl) CurrentUserLogout(ctx *gin.Context) (resp *performance.Empty, err error) {
	exitsSession, e := inbound.GetCurrentLoginSession(ctx)
	if e != nil {
		err = util.CreateAuthErr("无效凭据", e)
		return
	}

	err = biz.Q.Transaction(func(tx *biz.Query) error {
		loginSessionTx := tx.UserLoginSession
		_, e := loginSessionTx.WithContext(ctx).Where(loginSessionTx.ID.Eq(exitsSession.ID)).Delete()
		if e != nil {
			return e
		}

		// expire user
		inbound.GetMiddlewareRelativeAuthCache().Delete(fmt.Sprintf("%d", exitsSession.ID))
		return nil
	})

	if err != nil {
		err = util.CreateBizErr("退出登录失败", err)
		return
	}

	return
}
