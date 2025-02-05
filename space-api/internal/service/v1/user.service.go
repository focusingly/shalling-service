package service

import (
	"encoding/json"
	"io"
	"space-api/conf"
	"space-api/util"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

var githubOauth2Config, googleOauth2Config *oauth2.Config

func init() {
	v := conf.GetProjectViper()
	githubOauth2Config = &oauth2.Config{
		ClientID:     v.GetString("oauth2Conf.github.clientId"),
		ClientSecret: v.GetString("oauth2Conf.github.clientSecret"),
		Endpoint:     github.Endpoint,
		RedirectURL:  v.GetString("oauth2Conf.github.redirectUrl"),
		Scopes:       v.GetStringSlice(v.GetString("oauth2Conf.github.scopes")),
	}

	googleOauth2Config = &oauth2.Config{
		ClientID:     v.GetString("oauth2Conf.google.clientId"),
		ClientSecret: v.GetString("oauth2Conf.google.clientSecret"),
		Endpoint:     google.Endpoint,
		RedirectURL:  v.GetString("oauth2Conf.google.redirectUrl"),
		Scopes:       v.GetStringSlice(v.GetString("oauth2Conf.google.scopes")),
	}
}

func GetGithubLoginURL(ctx *gin.Context) (val string, err error) {
	val = githubOauth2Config.AuthCodeURL("asdad")

	return
}

func GithubCallbackHandler(ctx *gin.Context) (val any, err error) {
	grantCode := ctx.DefaultQuery("code", "")
	// 获取授权码
	if grantCode == "" {
		err = &util.BizErr{
			Msg: "登录失败, 请重试",
		}
		return
	}
	oauthToken, err := githubOauth2Config.Exchange(ctx, grantCode)
	if err != nil {
		err = &util.BizErr{
			Msg: err.Error(),
		}
		return
	}
	client := githubOauth2Config.Client(ctx, oauthToken)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		err = &util.BizErr{
			Msg: err.Error(),
		}
		return
	}
	if resp != nil {
		defer resp.Body.Close()
	}

	var user struct {
		Login   string `json:"login"`
		Name    string `json:"name"`
		Email   string `json:"email"` // 公开邮箱
		Avatar  string `json:"avatar_url"`
		HtmlUrl string `json:"html_url"`
	}
	// 获取公开信息
	if err = json.NewDecoder(resp.Body).Decode(&user); err != nil {
		err = &util.BizErr{
			Msg: err.Error(),
		}
		return
	}
	// 获取用户私人电子邮件地址
	emailResp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		err = &util.BizErr{
			Msg: err.Error(),
		}
		return
	}

	defer emailResp.Body.Close()
	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}
	if err = json.NewDecoder(emailResp.Body).Decode(&emails); err != nil {
		err = &util.BizErr{
			Msg: err.Error(),
		}

		return
	}

	val = gin.H{
		"public": &user,
		"emails": &emails,
	}
	return
}

func GetGoogleLoginURL(ctx *gin.Context) (val string, err error) {
	val = googleOauth2Config.AuthCodeURL("asdad")

	return
}

func GoogleCallbackHandler(ctx *gin.Context) (val any, err error) {
	grantCode := ctx.DefaultQuery("code", "")
	// 获取授权码
	if grantCode == "" {
		ctx.Error(&util.BizErr{
			Msg: "登录失败, 请重试",
		})
		return
	}
	oauthToken, err := googleOauth2Config.Exchange(ctx, grantCode)
	if err != nil {
		err = &util.BizErr{
			Msg: err.Error(),
		}

		return
	}

	client := googleOauth2Config.Client(ctx, oauthToken)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		err = &util.BizErr{Msg: err.Error()}
		return
	}
	defer resp.Body.Close()
	bf, err := io.ReadAll(resp.Body)
	if err != nil {
		err = &util.BizErr{Msg: "序列化错误"}
		return
	}
	// 返回用户信息
	val = string(bf)

	return
}
