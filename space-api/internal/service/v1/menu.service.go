package service

import (
	"space-api/dto"
	"space-api/util"
	"space-api/util/arr"
	"space-api/util/id"
	"space-domain/dao/biz"
	"space-domain/model"

	"github.com/gin-gonic/gin"
)

type _menuService struct{}

var DefaultMenuService = &_menuService{}

func (*_menuService) CreateOrUpdateMenu(req *dto.CreateOrUpdateMenuReq, ctx *gin.Context) (resp *dto.CreateOrUpdateMenuResp, err error) {
	var menuID int64
	err = biz.Q.Transaction(func(tx *biz.Query) error {
		menuTx := tx.MenuGroup
		exits, e := menuTx.WithContext(ctx).Where(menuTx.ID.Eq(menuID)).Take()
		// 不存在, 直接添加即可
		if e != nil {
			menuID = id.GetSnowFlakeNode().Generate().Int64()
			e := menuTx.WithContext(ctx).Create(&model.MenuGroup{
				BaseColumn: model.BaseColumn{ID: menuID},
				MenuName:   req.MenuName,
				RoutePath:  req.RoutePath,
				PostLink:   req.PostLink,
				OpenWay:    req.OpenWay,
				SubLinks: util.TernaryExpr(
					req.SubLinks == nil || (len(req.SubLinks) == 0),
					nil,
					req.SubLinks,
				),
			})
			if e != nil {
				return e
			}
		} else {
			// 存在的情况下进行更新
			menuID = exits.ID
			_, e := menuTx.WithContext(ctx).
				Select(
					menuTx.Hide,
					menuTx.MenuName,
					menuTx.RoutePath,
					menuTx.PostLink,
					menuTx.OpenWay,
					menuTx.SubLinks,
				).
				Where(menuTx.ID.Eq(menuID)).
				Updates(&model.MenuGroup{
					BaseColumn: exits.BaseColumn,
					MenuName:   req.MenuName,
					RoutePath:  req.RoutePath,
					PostLink:   req.PostLink,
					OpenWay:    req.OpenWay,
					SubLinks: util.TernaryExpr(
						req.SubLinks == nil || (len(req.SubLinks) == 0),
						nil,
						req.SubLinks,
					),
				})
			if e != nil {
				return e
			}
		}

		return nil
	})

	if err != nil {
		err = util.CreateBizErr("创建/更新菜单错误: "+err.Error(), err)
		return
	}

	menuTx := biz.MenuGroup
	find, e := menuTx.WithContext(ctx).Where(menuTx.ID.Eq(menuID)).Take()
	if e != nil {
		err = util.CreateBizErr("创建/更新菜单错误: "+e.Error(), e)
		return
	}
	resp = &dto.CreateOrUpdateMenuResp{
		MenuGroup: find,
	}

	return
}

func (*_menuService) GetVisibleMenus(req *dto.GetMenusReq, ctx *gin.Context) (resp *dto.GetMenusResp, err error) {
	menuCtx := biz.MenuGroup
	list, err := menuCtx.WithContext(ctx).Where(menuCtx.Hide.Neq(0)).Find()
	if err != nil {
		err = util.CreateBizErr("查找菜单列表错误", err)
		return
	}

	// 过滤掉隐藏菜单
	for _, menuGroup := range list {
		if menuGroup.SubLinks != nil {
			menuGroup.SubLinks = arr.FilterSlice(menuGroup.SubLinks, func(current *model.MenuLink, index int) bool {
				return current.Hide != 0
			})
		}

	}
	resp = &dto.GetMenusResp{
		Menus: list,
	}
	return
}

func (*_menuService) GetAllMenus(req *dto.GetMenusReq, ctx *gin.Context) (resp *dto.GetMenusResp, err error) {
	menuTx := biz.MenuGroup
	list, err := menuTx.WithContext(ctx).Find()
	if err != nil {
		err = util.CreateBizErr("查找菜单列表错误", err)
		return
	}

	resp = &dto.GetMenusResp{
		Menus: list,
	}

	return
}

func (*_menuService) DeleteMenuGroupByIDList(req *dto.DeleteMenuGroupsReq, ctx *gin.Context) (resp *dto.DeleteMenuGroupsResp, err error) {
	err = biz.Q.Transaction(func(tx *biz.Query) error {
		menuTx := tx.MenuGroup
		_, e := menuTx.WithContext(ctx).Where(menuTx.ID.In(req.IDList...)).Delete()
		if e != nil {
			return e
		}

		return nil
	})

	if err != nil {
		err = util.CreateBizErr("删除菜单失败: "+err.Error(), err)
		return
	}

	resp = &dto.DeleteMenuGroupsResp{}
	return
}
