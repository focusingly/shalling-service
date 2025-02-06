package service

import (
	"encoding/json"
	"fmt"
	"space-api/conf"
	"space-api/constants"
	"space-api/util"
	"space-domain/model"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
	"golang.org/x/sync/errgroup"
)

type empty struct{}

type (
	GithubPub struct {
		Login             string    `json:"login"`
		ID                int64     `json:"id"`
		NodeID            string    `json:"node_id"`
		AvatarURL         string    `json:"avatar_url"`
		GravatarID        string    `json:"gravatar_id"`
		URL               string    `json:"url"`
		HTMLURL           string    `json:"html_url"`
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

var githubOauth2Config, googleOauth2Config *oauth2.Config

// 记录缓存
var authSpaceCache = util.DefaultJsonCache.Group("auth")

func init() {
	v := conf.GetProjectViper()
	githubOauth2Config = &oauth2.Config{
		ClientID:     v.GetString("oauth2Conf.github.clientId"),
		ClientSecret: v.GetString("oauth2Conf.github.clientSecret"),
		Endpoint:     github.Endpoint,
		RedirectURL:  v.GetString("oauth2Conf.github.redirectUrl"),
		Scopes:       v.GetStringSlice("oauth2Conf.github.scopes"),
	}

	googleOauth2Config = &oauth2.Config{
		ClientID:     v.GetString("oauth2Conf.google.clientId"),
		ClientSecret: v.GetString("oauth2Conf.google.clientSecret"),
		Endpoint:     google.Endpoint,
		RedirectURL:  v.GetString("oauth2Conf.google.redirectUrl"),
		Scopes:       v.GetStringSlice("oauth2Conf.google.scopes"),
	}
}

func GetGithubLoginURL(ctx *gin.Context) (url string, err error) {
	state := uuid.NewString()
	ttl := time.Minute * 5 / time.Second
	// 设置过期
	authSpaceCache.Set(state, new(empty), util.Second(ttl))
	url = githubOauth2Config.AuthCodeURL(state)

	return
}

func VerifyGithubCallback(ctx *gin.Context) (oauthUser *model.OAuthLogin, err error) {
	grantCode := ctx.DefaultQuery("code", "")
	state := ctx.DefaultQuery("state", "")
	// 判断授权码
	if grantCode == "" || state == "" {
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
	if err = authSpaceCache.GetAndDel(state, &empty{}); err != nil {
		return
	}

	// 使用授权码
	oauthToken, err := githubOauth2Config.Exchange(ctx, grantCode)
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
	githubPub := new(GithubPub)
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
			if err = json.NewDecoder(res.Body).Decode(githubPub); err != nil {
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
		oauthUser = &model.OAuthLogin{
			PlatformName:   constants.GithubUser,
			PlatformUserId: githubPub.ID,
			PrimaryEmail:   primaryEmail,
			AccessToken:    oauthToken.AccessToken,
			RefreshToken:   &oauthToken.RefreshToken,
			ExpiredAt:      &oauthToken.ExpiresIn,
			Scopes:         githubOauth2Config.Scopes,
		}
	}

	return
}

func GetGoogleLoginURL(ctx *gin.Context) (val string, err error) {
	state := uuid.NewString()
	if err = authSpaceCache.Set(state, &empty{}, util.Second(time.Minute*5/time.Second)); err != nil {
		return
	}
	val = googleOauth2Config.AuthCodeURL(state)
	return
}

func VerifyGoogleCallback(ctx *gin.Context) (resp *model.OAuthLogin, err error) {
	grantCode := ctx.DefaultQuery("code", "")
	state := ctx.DefaultQuery("state", "")
	// 基本校验
	if grantCode == "" || state == "" {
		err = &util.AuthErr{
			BizErr: util.BizErr{
				Msg:    "凭据校验失败",
				Reason: fmt.Errorf("the principal is illegal"),
			},
		}
		return
	}
	if err = authSpaceCache.GetAndDel(state, &empty{}); err != nil {
		err = &util.AuthErr{
			BizErr: util.BizErr{
				Msg:    "凭据校验失败",
				Reason: err,
			},
		}
		return
	}

	oauthToken, err := googleOauth2Config.Exchange(ctx, grantCode)
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
	googlePub := GooglePub{}
	if err = json.NewDecoder(res.Body).Decode(&googlePub); err != nil {
		err = &util.AuthErr{
			BizErr: util.BizErr{
				Reason: err,
				Msg:    "解析用户数据失败: " + err.Error(),
			},
		}

		return
	}

	id, err := strconv.ParseInt(googlePub.ID, 10, 64)
	if err != nil {
		err = &util.AuthErr{
			BizErr: util.BizErr{
				Msg:    "获取用户 ID 失败: " + err.Error(),
				Reason: err,
			},
		}
		return
	}

	resp = &model.OAuthLogin{
		PlatformName:   constants.GoogleUser,
		PlatformUserId: id,
		PrimaryEmail:   googlePub.Name,
		AccessToken:    oauthToken.AccessToken,
		RefreshToken:   &oauthToken.RefreshToken,
		ExpiredAt:      &oauthToken.ExpiresIn,
		Scopes:         googleOauth2Config.Scopes,
	}

	return
}
