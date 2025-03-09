package service

import (
	"space-api/dto"
	"space-api/util"
	"space-api/util/id"
	"space-domain/dao/biz"
	"space-domain/model"

	"github.com/gin-gonic/gin"
)

type (
	IFriendLinkService interface {
		CreateOrUpdateFriendLink(req *dto.CreateOrUpdateFriendLinkReq, ctx *gin.Context) (resp *dto.CreateOrUpdateFriendLinkResp, err error)
		GetVisibleFriendLinks(req *dto.GetFriendLinksReq, ctx *gin.Context) (resp *dto.GetFriendLinksResp, err error)
		GetAllFriendLinks(req *dto.GetFriendLinksReq, ctx *gin.Context) (resp *dto.GetFriendLinksResp, err error)
		DeleteFriendLinkByIDList(req *dto.DeleteFriendLinkReq, ctx *gin.Context) (resp *dto.DeleteFriendLinkResp, err error)
	}
	friendLinkServiceImpl struct{}
)

var (
	_ IFriendLinkService = (*friendLinkServiceImpl)(nil)

	DefaultFriendLinkService IFriendLinkService = &friendLinkServiceImpl{}
)

func (*friendLinkServiceImpl) CreateOrUpdateFriendLink(req *dto.CreateOrUpdateFriendLinkReq, ctx *gin.Context) (resp *dto.CreateOrUpdateFriendLinkResp, err error) {
	var friendLinkID int64
	err = biz.Q.Transaction(func(tx *biz.Query) error {
		friendLinkTx := tx.FriendLink
		exits, e := friendLinkTx.WithContext(ctx).Where(friendLinkTx.ID.Eq(friendLinkID)).Take()
		// 不存在, 直接添加即可
		if e != nil {
			friendLinkID = id.GetSnowFlakeNode().Generate().Int64()
			e := friendLinkTx.WithContext(ctx).Create(
				&model.FriendLink{
					BaseColumn: model.BaseColumn{
						ID: friendLinkID,
					},
					SiteURL:     req.SiteURL,
					Owner:       req.Owner,
					ShortName:   req.ShortName,
					Available:   req.Available,
					LogoURL:     req.LogoURL,
					Description: req.Description,
					BgURL:       req.BgURL,
				},
			)
			if e != nil {
				return e
			}
		} else {
			// 存在的情况下进行更新
			friendLinkID = exits.ID
			_, e := friendLinkTx.WithContext(ctx).
				Select().
				Where(friendLinkTx.ID.Eq(friendLinkID)).
				Updates(&model.FriendLink{
					BaseColumn:  exits.BaseColumn,
					SiteURL:     req.SiteURL,
					Owner:       req.Owner,
					ShortName:   req.ShortName,
					Available:   req.Available,
					LogoURL:     req.LogoURL,
					Description: req.Description,
					BgURL:       req.BgURL,
				})
			if e != nil {
				return e
			}
		}

		return nil
	})

	if err != nil {
		err = util.CreateBizErr("创建/更新友链错误: "+err.Error(), err)
		return
	}

	friendLinkTx := biz.FriendLink
	find, e := friendLinkTx.WithContext(ctx).
		Where(friendLinkTx.ID.Eq(friendLinkID)).
		Take()
	if e != nil {
		err = util.CreateBizErr("创建/更新友链错误: "+e.Error(), e)
		return
	}
	resp = &dto.CreateOrUpdateFriendLinkResp{
		FriendLink: find,
	}

	return
}

func (*friendLinkServiceImpl) GetVisibleFriendLinks(req *dto.GetFriendLinksReq, ctx *gin.Context) (resp *dto.GetFriendLinksResp, err error) {
	friendLinkTx := biz.FriendLink
	list, err := friendLinkTx.
		WithContext(ctx).
		Select(
			friendLinkTx.CreatedAt,
			friendLinkTx.SiteURL,
			friendLinkTx.Owner,
			friendLinkTx.ShortName,
			friendLinkTx.Available,
			friendLinkTx.LogoURL,
			friendLinkTx.Description,
			friendLinkTx.BgURL,
		).
		Where(
			friendLinkTx.Hide.Eq(0),
		).
		Find()
	if err != nil {
		err = util.CreateBizErr("查找友链列表错误", err)
		return
	}

	resp = &dto.GetFriendLinksResp{
		List: list,
	}
	return
}

func (*friendLinkServiceImpl) GetAllFriendLinks(req *dto.GetFriendLinksReq, ctx *gin.Context) (resp *dto.GetFriendLinksResp, err error) {
	friendLinkTx := biz.FriendLink
	list, err := friendLinkTx.
		WithContext(ctx).
		Find()
	if err != nil {
		err = util.CreateBizErr("查找友链列表错误", err)
		return
	}

	resp = &dto.GetFriendLinksResp{
		List: list,
	}
	return
}

func (*friendLinkServiceImpl) DeleteFriendLinkByIDList(req *dto.DeleteFriendLinkReq, ctx *gin.Context) (resp *dto.DeleteFriendLinkResp, err error) {
	err = biz.Q.Transaction(func(tx *biz.Query) error {
		friendLinkTx := tx.FriendLink
		_, e := friendLinkTx.WithContext(ctx).Where(friendLinkTx.ID.In(req.IDList...)).Delete()
		if e != nil {
			return e
		}

		return nil
	})

	if err != nil {
		err = util.CreateBizErr("删除友链失败: "+err.Error(), err)
		return
	}

	resp = &dto.DeleteFriendLinkResp{}
	return
}
