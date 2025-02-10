package service

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

	"github.com/gin-gonic/gin"
	"gorm.io/gen/field"
)

type _userService struct{}

var DefaultUserService = &_userService{}

func (*_userService) UpdateLocalUserBasic(req *dto.UpdateLocalUserBasicReq, ctx *gin.Context) (resp *dto.UpdateLocalUserResp, err error) {
	err = biz.Q.Transaction(func(tx *biz.Query) error {
		currentLogin, e := auth.GetCurrentLoginSession(ctx)
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
		if currentLogin.UserID != findUser.ID {
			// 判断当前登录用户是否为管理员
			f, e := biz.LocalUser.WithContext(ctx).Where(biz.LocalUser.ID.Eq(currentLogin.UserID)).Take()
			if e != nil || (f.IsAdmin == 0) {
				return &util.BizErr{
					Msg:    "权限不足",
					Reason: fmt.Errorf("permission required admin"),
				}
			}
		}
		_, e = localTx.WithContext(ctx).
			Select(
				localTx.Hide,
				localTx.Email,
				localTx.Username,
				localTx.DisplayName,
				localTx.Password,
				localTx.AvatarURL,
				localTx.HomepageLink,
				localTx.Phone,
				localTx.IsAdmin,
			).
			Where(localTx.ID.Eq(findUser.ID)).
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
					IsAdmin:      findUser.IsAdmin,
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

func (*_userService) UpdateLocalUserPassword(req *dto.UpdateLocalUserPassReq, ctx *gin.Context) (resp *dto.UpdateLocalUserPassResp, err error) {
	if len(strings.TrimSpace(req.NewPassword)) < 8 {
		err = util.CreateBizErr("密码强度太弱, 请使用至少 8 位的密码", fmt.Errorf("new password strength too weak, must less has 8 character"))
	}
	err = biz.Q.Transaction(func(tx *biz.Query) error {
		currentLogin, e := auth.GetCurrentLoginSession(ctx)
		if e != nil {
			return e
		}
		localTx := tx.LocalUser

		findLocalUser, e := localTx.WithContext(ctx).
			Where(localTx.ID.Eq(req.UserID)).
			Take()
		if e != nil {
			return e
		}
		if currentLogin.UserID != findLocalUser.ID {
			// 判断当前登录用户是否为管理员
			currentUser, e := biz.LocalUser.WithContext(ctx).Where(biz.LocalUser.ID.Eq(currentLogin.UserID)).Take()
			if e != nil || (currentUser.IsAdmin == 0) {
				return &util.BizErr{
					Msg:    "权限不足",
					Reason: fmt.Errorf("permission required admin"),
				}
			}
			// 如果要修改的用户也为管理员, 那么不允许
			if findLocalUser.IsAdmin != 0 {
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

		// 允许管理员用户修改其他低一级权限的所有密码
		_, e = localTx.WithContext(ctx).
			Where(localTx.ID.Eq(findLocalUser.ID)).
			Update(localTx.Password, newHashedPass)
		if e != nil {
			return e
		}

		return nil
	})

	return
}

func (*_userService) ExpireLoginSessions(req *dto.ExpireUserLoginSessionReq, ctx *gin.Context) (resp *dto.ExpireUserLoginSessionResp, err error) {
	cacheSpace := auth.GetMiddlewareRelativeAuthCache()
	err = biz.Q.Transaction(func(tx *biz.Query) error {
		loginSessionTx := tx.UserLoginSession
		_, e := loginSessionTx.WithContext(ctx).
			Where(loginSessionTx.UUID.In(req.UUIDList...)).
			Delete()
		if e != nil {
			return e
		}
		// 清理缓存空间
		for _, uid := range req.UUIDList {
			cacheSpace.Delete(uid)
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

func (*_userService) UpdateOauth2User(req *dto.UpdateOauthUserReq, ctx *gin.Context) (resp *dto.UpdateOauthUserResp, err error) {
	cacheSpace := auth.GetMiddlewareRelativeAuthCache()
	err = biz.Q.Transaction(func(tx *biz.Query) error {
		userTx := tx.OAuth2User
		_, e := userTx.WithContext(ctx).
			Where(userTx.ID.Eq(req.UserID)).
			Update(userTx.Enable, util.TernaryExpr(req.Available, 1, 0))

		if e != nil {
			return e
		}

		// 如果用户被禁用
		if !req.Available {
			sessionTx := tx.UserLoginSession
			l, e := sessionTx.WithContext(ctx).
				Where(
					sessionTx.UserID.Eq(req.UserID),
					sessionTx.UserType.Neq(constants.LocalUser),
				).
				Find()
			if e != nil {
				return e
			}
			_, e = sessionTx.WithContext(ctx).
				Where(sessionTx.ID.In(
					arr.MapSlice(l, func(_ int, t *model.UserLoginSession) int64 {
						return t.ID
					})...,
				)).
				Delete()

			if e != nil {
				return e
			}

			// 清空缓存空间
			for _, t := range l {
				cacheSpace.Delete(t.UUID)
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

func (*_userService) DeleteOauth2User(req *dto.UpdateOauthUserReq, ctx *gin.Context) (resp *dto.UpdateOauthUserResp, err error) {
	err = biz.Q.Transaction(func(tx *biz.Query) error {
		userTx := tx.OAuth2User

		l, e := userTx.WithContext(ctx).Find()
		if e != nil {
			return e
		}
		// 删除用户
		_, e = userTx.WithContext(ctx).
			Where(
				userTx.ID.In(
					arr.MapSlice(l, func(_ int, t *model.OAuth2User) int64 {
						return t.ID
					})...,
				),
			).
			Delete()
		if e != nil {
			return e
		}
		// 查找登录会话
		sessions, e := tx.UserLoginSession.WithContext(ctx).Find()
		if e != nil {
			return e
		}
		// 删除会话
		_, e = tx.UserLoginSession.WithContext(ctx).
			Where(tx.UserLoginSession.UserID.In(
				arr.MapSlice(sessions, func(_ int, t *model.UserLoginSession) int64 {
					return t.ID
				})...,
			)).
			Delete()
		if e != nil {
			return e
		}
		cacheSpace := auth.GetMiddlewareRelativeAuthCache()
		// 清空凭据
		for _, s := range sessions {
			cacheSpace.Delete(s.UUID)
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

func (*_userService) GetLoginSessions(req *dto.GetLoginUserSessionsReq, ctx *gin.Context) (resp *dto.GetLoginUserSessionsResp, err error) {
	sessionOp := biz.UserLoginSession

	currentLogin, e := auth.GetCurrentLoginSession(ctx)
	if e != nil {
		return
	}

	localOp := biz.LocalUser
	f, e := localOp.WithContext(ctx).
		Where(localOp.ID.Eq(currentLogin.UserID)).
		Take()
	if e != nil {
		err = util.CreateBizErr("查询失败", err)
	}

	// 仅限管理员查看全部权限
	fieldList := []field.Expr{}
	if f.IsAdmin == 0 {
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
			sessionOp.OsName,
		)
	}

	l, count, e := sessionOp.WithContext(ctx).
		Select(fieldList...).
		FindByPage(req.Normalize())
	if e != nil {
		err = util.CreateBizErr("查找数据失败: "+e.Error(), e)
	}

	resp = &dto.GetLoginUserSessionsResp{
		PageList: model.PageList[*model.UserLoginSession]{
			List:  l,
			Page:  int64(*req.Page),
			Size:  int64(*req.Size),
			Total: count,
		},
	}

	return
}
