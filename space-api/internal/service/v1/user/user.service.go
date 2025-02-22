package user

import (
	"fmt"
	"space-api/constants"
	"space-api/dto"
	"space-api/middleware/auth"
	"space-api/util"
	"space-api/util/arr"
	"space-api/util/encrypt"
	"space-domain/dao/biz"
	"space-domain/model"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"gorm.io/gen/field"
)

type userService struct{}

var DefaultUserService = &userService{}

// UpdateLocalUserProfile 更新本地用户处除了密码之外的配置信息
func (*userService) UpdateLocalUserProfile(req *dto.UpdateLocalUserBasicReq, ctx *gin.Context) (resp *dto.UpdateLocalUserResp, err error) {
	err = biz.Q.Transaction(func(tx *biz.Query) error {
		loginSession, e := auth.GetCurrentLoginSession(ctx)
		if e != nil {
			return e
		}

		userTx := tx.LocalUser
		findUser, e := userTx.WithContext(ctx).
			Where(userTx.ID.Eq(req.UserID)).
			Take()
		if e != nil {
			return e
		}

		// 如果要更改的配置不是当前自身登录的配置
		if loginSession.UserID != findUser.ID {
			// 只允许管理员修改直接修改其他账户的信息
			if loginSession.UserType != constants.Admin {
				return util.CreateBizErr(
					"权限不足",
					fmt.Errorf("permission required admin"),
				)
			}
		}

		_, e = userTx.WithContext(ctx).
			Select(
				userTx.Hide,
				userTx.Email,
				userTx.Username,
				userTx.DisplayName,
				userTx.Password,
				userTx.AvatarURL,
				userTx.HomepageLink,
				userTx.Phone,
				userTx.IsAdmin,
			).
			Where(userTx.ID.Eq(findUser.ID)).
			Updates(
				&model.LocalUser{
					BaseColumn:   findUser.BaseColumn,
					Email:        req.Email,
					Username:     req.Username,
					DisplayName:  req.DisplayName,
					Password:     findUser.Password,
					AvatarURL:    req.AvatarURL,
					HomepageLink: req.HomepageLink,
					Phone:        req.Phone,
					IsAdmin:      findUser.IsAdmin, // 管理员的权限不允许被重新分配, 只能始终保持
				},
			)

		if e != nil {
			return e
		}

		return nil
	})

	if err != nil {
		err = util.CreateBizErr("更新基本信息失败", err)
		return
	}

	return
}

func (*userService) UpdateLocalUserPassword(req *dto.UpdateLocalUserPassReq, ctx *gin.Context) (resp *dto.UpdateLocalUserPassResp, err error) {
	newPass := strings.TrimSpace(req.NewPassword)
	if utf8.RuneCountInString(newPass) < 8 {
		err = util.CreateBizErr("密码强度太弱, 请使用至少 8 位的密码", fmt.Errorf("new password strength too weak, must less has 8 character"))
	}

	err = biz.Q.Transaction(func(tx *biz.Query) error {
		loginSession, e := auth.GetCurrentLoginSession(ctx)
		if e != nil {
			return e
		}

		localTx := tx.LocalUser
		findUser, e := localTx.WithContext(ctx).
			Where(localTx.ID.Eq(req.UserID)).
			Take()
		if e != nil {
			return e
		}

		if loginSession.UserID != findUser.ID {
			// 判断当前登录用户是否为管理员
			currentUser, e := biz.LocalUser.WithContext(ctx).Where(biz.LocalUser.ID.Eq(loginSession.UserID)).Take()
			if e != nil {
				return e
			}
			if currentUser.IsAdmin == 0 {
				return util.CreateBizErr(
					"权限不足",
					fmt.Errorf("permission required admin"),
				)
			}

			// 如果要修改的用户也为管理员, 那么不允许
			if findUser.IsAdmin != 0 {
				return fmt.Errorf("修改的账户为管理员账户, 需要自身登录")
			}

			if !encrypt.ComparePassword(strings.TrimSpace(req.OldPassword), currentUser.Password) {
				return fmt.Errorf("密码错误")
			}
		}

		newHashedPass, e := encrypt.EncryptPasswordByBcrypt(strings.TrimSpace(req.NewPassword))
		if e != nil {
			return e
		}

		// 允许管理员用户修改自身和其他低权限用户的密码
		_, e = localTx.WithContext(ctx).
			Where(localTx.ID.Eq(findUser.ID)).
			Update(localTx.Password, newHashedPass)
		if e != nil {
			return e
		}

		return nil
	})

	return
}

// ExpireAnyLoginSessions 删除登录会话信息
func (*userService) ExpireAnyLoginSessions(req *dto.ExpireUserLoginSessionReq, ctx *gin.Context) (resp *dto.ExpireUserLoginSessionResp, err error) {
	err = biz.Q.Transaction(func(tx *biz.Query) error {
		loginSessionTx := tx.UserLoginSession
		_, e := loginSessionTx.WithContext(ctx).
			Where(loginSessionTx.ID.In(req.IDList...)).
			Delete()
		if e != nil {
			return e
		}

		// 清理缓存空间
		cacheSpace := auth.GetMiddlewareRelativeAuthCache()
		for _, id := range req.IDList {
			cacheSpace.Delete(fmt.Sprintf("%d", id))
		}

		return nil
	})

	if err != nil {
		err = util.CreateBizErr("删除登录会话信息失败", err)
		return
	}
	resp = &dto.ExpireUserLoginSessionResp{}

	return
}

func (*userService) UpdateOauth2User(req *dto.UpdateOauthUserReq, ctx *gin.Context) (resp *dto.UpdateOauthUserResp, err error) {
	err = biz.Q.Transaction(func(tx *biz.Query) error {
		userTx := tx.OAuth2User
		_, e := userTx.WithContext(ctx).
			Where(userTx.ID.Eq(req.UserID)).
			Update(
				userTx.Enable,
				util.TernaryExpr(req.Available, 1, 0),
			)
		if e != nil {
			return e
		}

		// 如果用户被禁用
		if !req.Available {
			sessionTx := tx.UserLoginSession
			loginSessions, e := sessionTx.WithContext(ctx).
				Where(
					sessionTx.UserID.Eq(req.UserID),
					sessionTx.UserType.Neq(constants.Admin),
					sessionTx.UserType.Neq(constants.LocalUser),
				).
				Find()

			if e != nil {
				return e
			}

			_, e = sessionTx.WithContext(ctx).
				Where(sessionTx.ID.In(
					arr.MapSlice(loginSessions, func(_ int, t *model.UserLoginSession) int64 {
						return t.ID
					})...,
				)).
				Delete()

			if e != nil {
				return e
			}

			// 清空缓存空间
			cacheSpace := auth.GetMiddlewareRelativeAuthCache()
			for _, t := range loginSessions {
				cacheSpace.Delete(fmt.Sprintf("%d", t.ID))
			}
		}

		return nil
	})

	if err != nil {
		err = util.CreateBizErr("更新信息失败", err)
		return
	}

	resp = &dto.UpdateOauthUserResp{}
	return
}

// DeleteOauth2User 删除已经登录的 Oauth2 用户信息
func (*userService) DeleteOauth2User(req *dto.DeleteOauth2UserReq, ctx *gin.Context) (resp *dto.DeleteOauth2UserResp, err error) {
	err = biz.Q.Transaction(func(tx *biz.Query) error {
		userTx := tx.OAuth2User
		// 找到所有要被删除的用户

		removes, e := userTx.WithContext(ctx).
			Select(userTx.ID).
			Where(userTx.ID.In(req.IDList...)).
			Find()
		if e != nil {
			return e
		}

		removeUserIDList := arr.MapSlice(removes, func(_ int, t *model.OAuth2User) int64 {
			return t.ID
		})
		// 删除用户
		_, e = userTx.WithContext(ctx).
			Where(
				userTx.ID.In(
					removeUserIDList...,
				),
			).
			Delete()
		if e != nil {
			return e
		}

		// 查找匹配的登录会话
		loginSessions, e := tx.UserLoginSession.WithContext(ctx).
			Where(tx.UserLoginSession.UserID.In(removeUserIDList...)).
			Find()
		if e != nil {
			return e
		}

		// 删除会话记录
		_, e = tx.UserLoginSession.WithContext(ctx).
			Where(tx.UserLoginSession.UserID.In(
				arr.MapSlice(loginSessions, func(_ int, t *model.UserLoginSession) int64 {
					return t.ID
				})...,
			)).
			Delete()
		if e != nil {
			return e
		}
		// 清空缓存中的记录
		cacheSpace := auth.GetMiddlewareRelativeAuthCache()
		for _, s := range loginSessions {
			cacheSpace.Delete(fmt.Sprintf("%d", s.ID))
		}

		return nil
	})

	if err != nil {
		err = util.CreateBizErr("删除用户失败", err)
		return
	}

	resp = &dto.DeleteOauth2UserResp{}
	return
}

// GetLocalUserLoginSessions 获取已登录用户的会话列表
func (*userService) GetLocalUserLoginSessions(req *dto.GetLoginUserSessionsReq, ctx *gin.Context) (resp *dto.GetLoginUserSessionsResp, err error) {
	sessionOp := biz.UserLoginSession
	currentLogin, e := auth.GetCurrentLoginSession(ctx)
	if e != nil {
		err = util.CreateAuthErr("无效的用户凭据", err)
		return
	}

	// 仅限管理员查看全部权限
	fieldList := []field.Expr{}
	if currentLogin.UserType == constants.Admin {
		fieldList = append(fieldList, sessionOp.ALL)
	} else {
		fieldList = append(fieldList,
			sessionOp.ID,
			sessionOp.CreatedAt,
			sessionOp.UpdatedAt,
			sessionOp.Hide,
			sessionOp.UserID,
			sessionOp.IpU32Val,
			sessionOp.IpAddress,
			sessionOp.IpSource,
			sessionOp.ExpiredAt,
			sessionOp.UserType,
			sessionOp.Useragent,
			sessionOp.ClientName,
			sessionOp.OSName,
		)
	}

	list, count, e := sessionOp.WithContext(ctx).
		Select(fieldList...).
		FindByPage(req.Normalize())

	if e != nil {
		err = util.CreateBizErr("查找数据失败: "+e.Error(), e)
	}

	resp = &dto.GetLoginUserSessionsResp{
		PageList: model.PageList[*model.UserLoginSession]{
			List:  list,
			Page:  int64(*req.Page),
			Size:  int64(*req.Size),
			Total: count,
		},
	}

	return
}

// GetLoginUserBasicProfile 获取已登录用户的基本信息(头像, 链接等)
func (*userService) GetLoginUserBasicProfile(ctx *gin.Context) (resp *dto.LoginUserBasicProfile, err error) {
	currentLogin, e := auth.GetCurrentLoginSession(ctx)

	if e != nil {
		err = util.CreateAuthErr("无效的用户凭据", err)
		return
	}

	switch currentLogin.UserType {
	case constants.Admin, constants.LocalUser:
		localOp := biz.LocalUser
		u, e := localOp.WithContext(ctx).
			Where(localOp.ID.Eq(currentLogin.UserID)).
			Take()
		if e != nil {
			err = util.CreateAuthErr("查找用户失败", err)
			return
		}
		resp = &dto.LoginUserBasicProfile{
			UserID:       u.ID,
			PlatformName: "本地用户",
			Username:     u.DisplayName,
			AvatarURL:    u.AvatarURL,
			HomepageLink: u.HomepageLink,
		}
	case constants.GithubUser, constants.GoogleUser:
		userOp := biz.OAuth2User
		u, e := userOp.WithContext(ctx).
			Where(userOp.ID.Eq(currentLogin.UserID)).
			Take()
		if e != nil {
			err = util.CreateAuthErr("查找用户失败", err)
			return
		}
		resp = &dto.LoginUserBasicProfile{
			UserID:       u.ID,
			PlatformName: u.PlatformName,
			Username:     u.Username,
			AvatarURL:    u.AvatarURL,
			HomepageLink: u.HomepageLink,
		}
	default:
		err = util.CreateBizErr("未注册的用户类型, 暂不支持", fmt.Errorf("un-registered user type"))
		return
	}

	return
}
