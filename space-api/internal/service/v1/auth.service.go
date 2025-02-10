package service

import (
	"encoding/json"
	"fmt"
	"slices"
	"space-api/conf"
	"space-api/constants"
	"space-api/dto"
	"space-api/middleware/auth"
	"space-api/middleware/inbound"
	"space-api/util"
	"space-api/util/arr"
	"space-api/util/encrypt"
	"space-api/util/id"
	"space-api/util/ip"
	"space-api/util/performance"
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
	_authCache = auth.GetMiddlewareRelativeAuthCache()
	_jwtConf   = conf.ProjectConf.GetJwtConf()
	_appConf   = conf.ProjectConf.GetAppConf()
)

type _authService struct {
}

var DefaultAuthService = &_authService{}

type boData struct {
	UserID   int64
	UserType string
}

func (s *_authService) updateUserLoginSession(user *boData, bizTx *biz.Query, ctx *gin.Context) (resp *model.UserLoginSession, err error) {
	err = bizTx.Transaction(func(tx *biz.Query) error {
		loginSessionTx := tx.UserLoginSession
		// 获取用户所有已经登录的会话信息
		existsSessions, e := loginSessionTx.WithContext(ctx).
			Where(loginSessionTx.UserID.Eq(user.UserID)).
			Find()
		if e != nil {
			return util.CreateBizErr("设置会话信息失败", e)
		}

		nowEpochMill := time.Now().UnixMilli()
		// 更新列表
		updates := slices.DeleteFunc(
			slices.Clone(existsSessions),
			func(s *model.UserLoginSession) bool {
				return nowEpochMill >= s.ExpiredAt
			},
		)
		// 排序
		slices.SortFunc(
			updates,
			func(a, b *model.UserLoginSession) int {
				// 比较新的数据, 放在前面
				return int(b.ID - a.ID)
			},
		)
		if len(updates) >= _appConf.MaxUserActive-1 {
			// 淘汰末尾数据(最大允许用户数来自系统配置)
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
			UUID:       uuid.NewString(),
			IpU32Val:   &to32Ip,
			IpAddress:  &ipAddr,
			IpSource:   &ipSource,
			ExpiredAt:  time.Now().Add(_jwtConf.ParsedExpTime).UnixMilli(),
			UserType:   user.UserType,
			Useragent:  ua.Useragent,
			ClientName: ua.ClientName,
			OsName:     ua.OS,
		}

		token, e := verify.CreateJwtToken(newLoginSession)
		if e != nil {
			return e
		}

		// 设置新的 token
		newLoginSession.Token = token
		// 存入新的会话
		updates = append(updates, newLoginSession)
		// 删除所有已存在的列表
		_, e = loginSessionTx.WithContext(ctx).
			Where(loginSessionTx.ID.In(
				arr.MapSlice(
					existsSessions,
					func(_ int, s *model.UserLoginSession) int64 {
						return s.ID
					})...,
			)).
			Delete()
		if e != nil {
			return e
		}
		// 批量重新创建会话信息
		e = loginSessionTx.WithContext(ctx).
			CreateInBatches(updates, appConf.MaxUserActive)
		if e != nil {
			return e
		}

		// set value
		resp = newLoginSession
		// 更新缓存中的数据
		cacheSpace := auth.GetMiddlewareRelativeAuthCache()
		// 删除已经登录的会话信息
		for _, u := range existsSessions {
			cacheSpace.Delete(u.UUID)
		}
		nowUnixMilli := time.Now().UnixMilli()
		for _, u := range updates {
			cacheSpace.Set(u.UUID, u, performance.Second((u.ExpiredAt-nowUnixMilli)/1000))
		}

		return nil
	})
	if err != nil {
		err = util.CreateAuthErr("创建会话失败", err)
		return
	}

	return
}

// AdminLogin 后台管理员登录
func (s *_authService) AdminLogin(req *dto.AdminLoginReq, ctx *gin.Context) (resp *dto.AdminLoginResp, err error) {
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
			UserType: constants.LocalUser,
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

func (s *_authService) HandleOauthLogin(req *dto.OauthLoginCallbackReq, ctx *gin.Context) (resp *dto.OauthLoginCallbackResp, err error) {
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

func (*_authService) GetOauth2LoginGrantURL(req *dto.GetLoginURLReq, ctx *gin.Context) (resp dto.GetLoginURLResp, err error) {
	state := uuid.NewString()
	ttl := time.Minute * 5 / time.Second
	// 设置过期时间
	err = _authCache.Set(state, new(performance.Empty), performance.Second(ttl))
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

func (*_authService) ParseGithubCallback(req *dto.OauthLoginCallbackReq, ctx *gin.Context) (resp *model.OAuth2User, err error) {
	// 判断授权码
	if req.GrantCode == "" || req.State == "" {
		err = &util.AuthErr{
			BizErr: util.BizErr{
				Msg:    "获取授权信息失败, 请重试",
				Reason: fmt.Errorf("grant code not exits"),
			},
		}
		return
	}

	// 判断 state
	// 判断缓存里的情况
	if err = _authCache.GetAndDel(req.State, &performance.Empty{}); err != nil {
		return
	}
	// 使用授权码
	oauthToken, err := githubOauth2Config.Exchange(ctx, req.GrantCode)
	if err != nil {
		err = &util.AuthErr{
			BizErr: util.BizErr{
				Msg: err.Error(),
			},
		}

		return
	}

	client := githubOauth2Config.Client(ctx, oauthToken)
	var primaryEmail string
	githubPubDetail := new(GithubPub)
	emailList := []GithubEmailElement{}
	var group errgroup.Group

	// 读取公开信息
	group.Go(func() error {
		res, err := client.Get("https://api.github.com/user")
		if err != nil {
			err = &util.BizErr{
				Msg: err.Error(),
			}
			return err
		}
		if res != nil {
			defer res.Body.Close()
			// 获取公开信息
			if err = json.NewDecoder(res.Body).Decode(githubPubDetail); err != nil {
				err = &util.AuthErr{
					BizErr: util.BizErr{
						Msg:    "解码错误: " + err.Error(),
						Reason: err,
					},
				}
				return err
			}
		}
		return nil
	})
	// 获取主邮箱
	group.Go(func() error {
		// 获取用户私人电子邮件地址
		emailResp, e := client.Get("https://api.github.com/user/emails")
		if e != nil {
			e = &util.BizErr{
				Msg: e.Error(),
			}
			return e
		}
		if emailResp != nil {
			defer emailResp.Body.Close()
			if err := json.NewDecoder(emailResp.Body).Decode(&emailList); err != nil {
				return err
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
	if err = group.Wait(); err != nil {
		err = &util.AuthErr{
			BizErr: util.BizErr{
				Msg:    "获取用户信息失败, 请重试",
				Reason: err,
			},
		}

		return
	} else {
		resp = &model.OAuth2User{
			PlatformName:   constants.GithubUser,
			PlatformUserId: fmt.Sprintf("%d", githubPubDetail.ID),
			Username:       githubPubDetail.Login,
			PrimaryEmail:   primaryEmail,
			AccessToken:    oauthToken.AccessToken,
			RefreshToken:   &oauthToken.RefreshToken,
			ExpiredAt:      &oauthToken.ExpiresIn,
			AvatarURL:      &githubPubDetail.AvatarURL,
			HomepageLink:   &githubPubDetail.HtmlURL,
			Scopes:         githubOauth2Config.Scopes,
		}
	}

	return
}

func (*_authService) GetGoogleLoginURL(ctx *gin.Context) (val string, err error) {
	state := uuid.NewString()
	if err = _authCache.Set(state, &performance.Empty{}, performance.Second(time.Minute*5/time.Second)); err != nil {
		return
	}
	val = googleOauth2Config.AuthCodeURL(state)
	return
}

func (*_authService) ParseGoogleCallback(req *dto.OauthLoginCallbackReq, ctx *gin.Context) (resp *model.OAuth2User, err error) {
	// 基本校验
	if req.GrantCode == "" || req.State == "" {
		err = &util.AuthErr{
			BizErr: util.BizErr{
				Msg:    "凭据校验失败",
				Reason: fmt.Errorf("the principal is illegal"),
			},
		}
		return
	}
	if err = _authCache.GetAndDel(req.State, &performance.Empty{}); err != nil {
		err = &util.AuthErr{
			BizErr: util.BizErr{
				Msg:    "凭据校验失败",
				Reason: err,
			},
		}
		return
	}

	oauthToken, err := googleOauth2Config.Exchange(ctx, req.GrantCode)
	if err != nil {
		err = &util.AuthErr{
			BizErr: util.BizErr{
				Msg:    "获取凭证失败" + err.Error(),
				Reason: err,
			},
		}

		return
	}

	client := googleOauth2Config.Client(ctx, oauthToken)
	res, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil || res.Body == nil {
		err = &util.AuthErr{
			BizErr: util.BizErr{
				Reason: err,
				Msg:    "获取用户数据失败: " + err.Error(),
			},
		}

		return
	}

	defer res.Body.Close()

	googlePubDetail := GooglePub{}
	if err = json.NewDecoder(res.Body).Decode(&googlePubDetail); err != nil {
		err = &util.AuthErr{
			BizErr: util.BizErr{
				Reason: err,
				Msg:    "解析用户数据失败: " + err.Error(),
			},
		}

		return
	}

	resp = &model.OAuth2User{
		PlatformName:   constants.GoogleUser,
		PlatformUserId: googlePubDetail.ID,
		Username:       googlePubDetail.Name,
		PrimaryEmail:   googlePubDetail.Email,
		AccessToken:    oauthToken.AccessToken,
		RefreshToken:   &oauthToken.RefreshToken,
		ExpiredAt:      &oauthToken.ExpiresIn,
		AvatarURL:      &googlePubDetail.Picture,
		HomepageLink:   new(string),
		Scopes:         googleOauth2Config.Scopes,
	}

	return
}

func (*_authService) Logout(ctx *gin.Context) (resp *performance.Empty, err error) {
	exitsSession, e := auth.GetCurrentLoginSession(ctx)
	if e != nil {
		err = util.CreateAuthErr("无效凭据", e)
		return
	}

	err = biz.Q.Transaction(func(tx *biz.Query) error {
		loginSessionTx := tx.UserLoginSession
		_, e := loginSessionTx.WithContext(ctx).Where(loginSessionTx.UUID.Eq(exitsSession.UUID)).Delete()
		if e != nil {
			return e
		}

		// expire user
		auth.GetMiddlewareRelativeAuthCache().Delete(exitsSession.UUID)
		return nil
	})

	if err != nil {
		err = util.CreateBizErr("退出登录失败", err)
		return
	}

	return
}
