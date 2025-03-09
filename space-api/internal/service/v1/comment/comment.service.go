package comment

import (
	"fmt"
	"space-api/constants"
	"space-api/dto"
	"space-api/middleware/inbound"

	"space-api/util"
	"space-api/util/ip"
	"space-api/util/ptr"
	"space-domain/dao/biz"
	"space-domain/model"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gen"
)

type (
	ICommentService interface {
		CheckAndCreateComment(req *dto.CreateCommentReq, ctx *gin.Context) (resp *dto.CreateCommentResp, err error)
		UpdateComment(req *dto.UpdateCommentReq, ctx *gin.Context) (resp *dto.UpdateCommentResp, err error)
		GetVisibleRootCommentPages(req *dto.GetRootCommentPagesReq, ctx *gin.Context) (resp *dto.GetRootCommentPagesResp, err error)
		GetVisibleSubCommentPages(req *dto.GetSubCommentPagesReq, ctx *gin.Context) (resp *dto.GetSubCommentPagesResp, err error)
		GetAnyRootCommentPages(req *dto.GetRootCommentPagesReq, ctx *gin.Context) (resp *dto.GetRootCommentPagesResp, err error)
		GetAnySubCommentPages(req *dto.GetSubCommentPagesReq, ctx *gin.Context) (resp *dto.GetSubCommentPagesResp, err error)
		DeleteSubComments(req *dto.DeleteSubCommentReq, ctx *gin.Context) (resp *dto.DeleteCategoryResp, err error)
		DeleteRootComments(req *dto.DeleteRootCommentReq, ctx *gin.Context) (resp *dto.DeleteRootCommentResp, err error)
	}
	commentServiceImpl struct{}
	commentUserDetail  struct {
		Username    string
		Avatar      *string
		HomePageURL *string
	}
)

var (
	_ ICommentService = (*commentServiceImpl)(nil)

	DefaultCommentService ICommentService = &commentServiceImpl{}
)

// CheckAndCreateComment 检查创建评论的合法性并创建评论
func (servicePtr *commentServiceImpl) CheckAndCreateComment(req *dto.CreateCommentReq, ctx *gin.Context) (resp *dto.CreateCommentResp, err error) {
	// 当前登录会话的用户
	loginSession, getSessionErr := inbound.GetCurrentLoginSession(ctx)
	if getSessionErr != nil {
		return nil, getSessionErr
	}

	postOP := biz.Post

	// 先找到文章
	matchedPost, matchedPostErr := postOP.WithContext(ctx).Where(postOP.ID.Eq(req.PostID)).Take()
	if matchedPostErr != nil {
		err = util.CreateBizErr("不存在文章", matchedPostErr)
		return
	}

	switch {
	// 如果文章设为了不可见, 不可评论的情况, 分别判断情况
	case matchedPost.Hide != 0 || (matchedPost.AllowComment == 0):
		// 此时不允许 oauth2 登录的用户进行评论
		if loginSession.UserType == constants.GoogleUser || loginSession.UserType == constants.GithubUser {
			err = util.CreateBizErr("评论不可用", fmt.Errorf("comment not available"))
			return
		}
		// TODO 暂时只允许管理员用户进行评论
		if loginSession.UserType != constants.Admin {
			err = util.CreateBizErr("当前本地账户评论功能不可用", fmt.Errorf("comment not available"))
			return
		}
	}

	commentOP := biz.Comment
	// 发表的评论信息如果是子评论, 那么要先检查根评论和他的父评论合法性在才能创建
	if req.ReplyToID != 0 || req.RootCommentID != 0 {
		// 查找根评论(只有根评论存在, 才有底下的评论)
		rootCmt, rootCmdErr := commentOP.WithContext(ctx).
			Where(
				commentOP.ID.Eq(req.RootCommentID), // 根评论的 ID
				commentOP.PostId.Eq(req.PostID),    // 文章存在
			).
			Take()
		if rootCmdErr != nil {
			// 不存在评论
			err = util.CreateBizErr("评论不可用", rootCmdErr)
			return
		}

		// 根评论为隐藏的条件下只允许管理员评论
		if rootCmt.Hide != 0 && loginSession.UserType != constants.Admin {
			err = util.CreateBizErr("评论不可用", fmt.Errorf("comment not available"))
			return
		}

		// 查找父评论,
		parentCmt, parentCmdErr := commentOP.WithContext(ctx).
			Where(
				commentOP.ID.Eq(req.ReplyToID), // 父评论的 ID(父评论本身可能就是根评论, 也可能是根评论下另一条较早发布的子评论)
				commentOP.PostId.Eq(req.PostID),
			).
			Take()
		if parentCmdErr != nil {
			// 不存在评论
			err = util.CreateBizErr("评论不可用", parentCmdErr)
			return
		}

		// 父亲评论为隐藏的条件下
		if parentCmt.Hide != 0 && loginSession.UserType != constants.Admin {
			err = util.CreateBizErr("评论不可用", fmt.Errorf("comment not available"))
			return
		}
	}

	return servicePtr.createCommentDirect(req, ctx)
}

func (servicePtr *commentServiceImpl) UpdateComment(req *dto.UpdateCommentReq, ctx *gin.Context) (resp *dto.UpdateCommentResp, err error) {
	err = biz.Q.Transaction(func(tx *biz.Query) error {
		cmtTx := tx.Comment
		_, e := cmtTx.WithContext(ctx).
			Where(cmtTx.ID.Eq(req.ID)).
			Select(
				cmtTx.Content,
				cmtTx.UpVote,
				cmtTx.DownVote,
				cmtTx.Hide,
			).
			Updates(&model.Comment{
				BaseColumn: model.BaseColumn{
					Hide: util.TernaryExpr(req.Hide, 1, 0),
				},
				Content:  req.Content,
				UpVote:   req.UpVote,
				DownVote: req.DownVote,
			})

		return e
	})

	if err != nil {
		err = util.CreateBizErr("修改评论失败", err)
		return
	}
	return
}

// 充当基础通用方法, 只创建, 不进行额外的关联关系校验
func (*commentServiceImpl) createCommentDirect(req *dto.CreateCommentReq, ctx *gin.Context) (resp *dto.CreateCommentResp, err error) {
	err = biz.Q.Transaction(func(tx *biz.Query) error {
		commentTx := tx.Comment
		loginSession, sessionErr := inbound.GetCurrentLoginSession(ctx)
		if sessionErr != nil {
			return sessionErr
		}

		var currentUserDetail *commentUserDetail
		switch loginSession.UserType {
		// 来自 oauth2 登录用户的评论
		case constants.GoogleUser, constants.GithubUser:
			oauthUserTx := tx.OAuth2User
			oauthUser, oauthUserFindErr := oauthUserTx.WithContext(ctx).
				Select(oauthUserTx.Username, oauthUserTx.AvatarURL, oauthUserTx.HomepageLink).
				Where(
					oauthUserTx.ID.Eq(loginSession.UserID),
				).
				Take()
			if oauthUserFindErr != nil {
				return oauthUserFindErr
			}
			currentUserDetail = &commentUserDetail{
				Username:    oauthUser.Username,
				Avatar:      oauthUser.AvatarURL,
				HomePageURL: oauthUser.HomepageLink,
			}

		// 本地账户评论(包含管理员账户)
		case constants.LocalUser, constants.Admin:
			localUserTx := tx.LocalUser
			localUser, findLocalUserErr := localUserTx.WithContext(ctx).
				Select(localUserTx.DisplayName, localUserTx.AvatarURL, localUserTx.HomepageLink).
				Where(localUserTx.ID.Eq(loginSession.UserID)).
				Take()
			if findLocalUserErr != nil {
				return findLocalUserErr
			}
			currentUserDetail = &commentUserDetail{
				Username:    localUser.DisplayName, // 对于本地用户不直接使用登录的用户名称, 而是使用额外名称
				Avatar:      localUser.AvatarURL,
				HomePageURL: localUser.HomepageLink,
			}
		default:
			return fmt.Errorf("un-support user type: %s", loginSession.UserType)
		}

		realIpAddr := inbound.GetRealIpWithContext(ctx)
		uaDetail := inbound.GetUserAgentFromContext(ctx)
		ipSource, _ := ip.GetIpSearcher().SearchByStr(realIpAddr)

		// 进行创建
		createCommentErr := commentTx.WithContext(ctx).Create(
			&model.Comment{
				PostId:          req.PostID,
				UserId:          loginSession.ID,
				RootCommentId:   req.RootCommentID,
				ReplyToId:       req.ReplyToID,
				UserType:        loginSession.UserType,
				DisplayUsername: currentUserDetail.Username,
				Avatar:          currentUserDetail.Avatar,
				HomePageURL:     currentUserDetail.HomePageURL,
				Content:         strings.TrimSpace(req.Content),
				UpVote:          nil,
				DownVote:        nil,
				IpAddr:          realIpAddr,
				IpSource:        &ipSource,
				Useragent:       &uaDetail.Useragent,
				OS:              &uaDetail.OS,
				ClientName:      &uaDetail.ClientName,
				SubEmailReply:   util.TernaryExpr(req.SubEmailNotify, 1, 0),
			},
		)

		return createCommentErr
	})

	if err != nil {
		err = util.CreateBizErr("创建评论失败", err)
		return
	}

	resp = &dto.CreateCommentResp{}

	return
}

// GetVisibleRootCommentPages 获取根评论的分页
func (servicePtr *commentServiceImpl) GetVisibleRootCommentPages(req *dto.GetRootCommentPagesReq, ctx *gin.Context) (resp *dto.GetRootCommentPagesResp, err error) {
	// 必须的文章开启评和可见论才允许获取评论
	postCtx := biz.Post
	_, err = postCtx.WithContext(ctx).
		Where(postCtx.ID.Eq(req.PostID), postCtx.Hide.Eq(0), postCtx.AllowComment.Neq(0)).
		Take()
	if err != nil {
		err = util.CreateBizErr("当前文章不允许评论", err)
		return
	}

	commentCtx := biz.Comment
	rootList, count, err := commentCtx.WithContext(ctx).
		Select(
			commentCtx.ID,
			commentCtx.CreatedAt,
			commentCtx.UpdatedAt,
			commentCtx.PostId,
			commentCtx.UserId,
			commentCtx.UserType,
			commentCtx.DisplayUsername,
			commentCtx.Avatar,
			commentCtx.HomePageURL,
			commentCtx.RootCommentId,
			commentCtx.ReplyToId,
			commentCtx.Content,
			commentCtx.UpVote,
			commentCtx.DownVote,
			commentCtx.IpSource,
			commentCtx.ClientName,
		).
		Order(commentCtx.UpVote.Desc(), commentCtx.CreatedAt.Desc()).
		FindByPage(req.Normalize())

	if err != nil {
		err = util.CreateBizErr("查询评论失败", err)
		return
	}

	nestedList := []*dto.NestedComments{}
	for _, cmt := range rootList {
		// 规范化显示的 IP 地址
		if cmt.IpSource != nil {
			splits := strings.Split(*cmt.IpSource, "|")
			if len(splits) > 0 {
				t := ""
				if splits[0] != "0" {
					t = splits[0]
				}
				if splits[len(splits)-1] != "0" && splits[len(splits)-1] != "内网IP" {
					t += splits[len(splits)-1]
				}
				cmt.IpSource = &t
			} else {
				cmt.IpSource = ptr.ToPtr("")
			}
		}

		// 获取子评论
		subs, e := servicePtr.GetVisibleSubCommentPages(&dto.GetSubCommentPagesReq{
			BasePageParam: dto.BasePageParam{Page: ptr.ToPtr(1), Size: ptr.ToPtr(10)},
			PostID:        req.PostID,
			RootCommentID: cmt.ID,
		}, ctx)
		if e != nil {
			return
		}
		// 填充子列表
		nestedList = append(nestedList, &dto.NestedComments{
			RootComment: cmt,
			Subs:        subs,
		})
	}

	resp = &dto.GetRootCommentPagesResp{
		PageList: model.PageList[*dto.NestedComments]{
			List:  nestedList,
			Page:  int64(*req.Page),
			Size:  int64(*req.Size),
			Total: count,
		},
	}

	return
}

// GetRootCommentPages 获取子评论的分页
func (*commentServiceImpl) GetVisibleSubCommentPages(req *dto.GetSubCommentPagesReq, ctx *gin.Context) (resp *dto.GetSubCommentPagesResp, err error) {
	// 必须的文章开启评论才允许获取评论
	postCtx := biz.Post
	_, err = postCtx.WithContext(ctx).
		Where(postCtx.ID.Eq(req.PostID), postCtx.Hide.Eq(0), postCtx.AllowComment.Neq(0)).
		Take()
	if err != nil {
		err = util.CreateBizErr("当前文章不允许评论", err)
		return
	}

	commentCtx := biz.Comment
	list, count, err := commentCtx.WithContext(ctx).
		Select(
			commentCtx.ID,
			commentCtx.CreatedAt,
			commentCtx.UpdatedAt,
			commentCtx.PostId,
			commentCtx.DisplayUsername,
			commentCtx.UserId,
			commentCtx.UserType,
			commentCtx.Avatar,
			commentCtx.HomePageURL,
			commentCtx.RootCommentId,
			commentCtx.ReplyToId,
			commentCtx.Content,
			commentCtx.UpVote,
			commentCtx.DownVote,
			commentCtx.IpSource,
			commentCtx.ClientName,
		).
		Where(
			// 匹配文章
			commentCtx.RootCommentId.Eq(req.RootCommentID),
			commentCtx.ReplyToId.Neq(0),
		).
		// 设置排序, 根据点赞数量, 创建时间逐一进行降序排列
		Order(commentCtx.UpVote.Desc(), commentCtx.CreatedAt.Desc()).
		FindByPage(req.BasePageParam.Normalize())

	if err != nil {
		err = util.CreateBizErr("查询评论失败", err)
		return
	}

	// 规范化显示的 IP 地址
	for _, cmt := range list {
		if cmt.IpSource != nil {
			splits := strings.Split(*cmt.IpSource, "|")
			if len(splits) > 0 {
				t := ""
				if splits[0] != "0" {
					t = splits[0]
				}
				if splits[len(splits)-1] != "0" && splits[len(splits)-1] != "内网IP" {
					t += splits[len(splits)-1]
				}
				cmt.IpSource = &t
			} else {
				cmt.IpSource = ptr.ToPtr("")
			}
		}
	}

	resp = &model.PageList[*model.Comment]{
		List:  list,
		Page:  int64(*req.Page),
		Size:  int64(*req.Size),
		Total: count,
	}

	return
}

// GetAnyRootCommentPages 获取任何的根评论分页
func (servicePtr *commentServiceImpl) GetAnyRootCommentPages(req *dto.GetRootCommentPagesReq, ctx *gin.Context) (resp *dto.GetRootCommentPagesResp, err error) {
	// 必须的文章开启评和可见论才允许获取评论
	postCtx := biz.Post
	_, err = postCtx.WithContext(ctx).
		Where(postCtx.ID.Eq(req.PostID)).
		Take()
	if err != nil {
		err = util.CreateBizErr("获取失败", err)
		return
	}

	commentCtx := biz.Comment
	rootList, count, err := commentCtx.WithContext(ctx).
		Order(commentCtx.UpVote.Desc(), commentCtx.CreatedAt.Desc()).
		FindByPage(req.Normalize())

	if err != nil {
		err = util.CreateBizErr("查询评论失败", err)
		return
	}

	nestedList := []*dto.NestedComments{}
	for _, cmt := range rootList {
		// 获取子评论
		subs, e := servicePtr.GetAnySubCommentPages(&dto.GetSubCommentPagesReq{
			BasePageParam: dto.BasePageParam{Page: ptr.ToPtr(1), Size: ptr.ToPtr(10)},
			PostID:        req.PostID,
			RootCommentID: cmt.ID,
		}, ctx)
		if e != nil {
			return
		}
		// 填充子列表
		nestedList = append(nestedList, &dto.NestedComments{
			RootComment: cmt,
			Subs:        subs,
		})
	}

	resp = &dto.GetRootCommentPagesResp{
		PageList: model.PageList[*dto.NestedComments]{
			List:  nestedList,
			Page:  int64(*req.Page),
			Size:  int64(*req.Size),
			Total: count,
		},
	}

	return
}

// GetAnySubCommentPages 获取任何的子评论分页数据
func (*commentServiceImpl) GetAnySubCommentPages(req *dto.GetSubCommentPagesReq, ctx *gin.Context) (resp *dto.GetSubCommentPagesResp, err error) {
	postCtx := biz.Post
	_, findPostErr := postCtx.WithContext(ctx).
		Where(postCtx.ID.Eq(req.PostID)).
		Take()
	if findPostErr != nil {
		err = util.CreateBizErr("当前文章不存在", findPostErr)
		return
	}

	commentCtx := biz.Comment
	list, count, err := commentCtx.WithContext(ctx).
		Where(
			// 匹配文章
			commentCtx.RootCommentId.Eq(req.RootCommentID),
			commentCtx.ReplyToId.Neq(0),
		).
		// 设置排序, 根据点赞数量, 创建时间逐一进行降序排列
		Order(commentCtx.UpVote.Desc(), commentCtx.CreatedAt.Desc()).
		FindByPage(req.BasePageParam.Normalize())

	if err != nil {
		err = util.CreateBizErr("查询评论失败", err)
		return
	}

	resp = &model.PageList[*model.Comment]{
		List:  list,
		Page:  int64(*req.Page),
		Size:  int64(*req.Size),
		Total: count,
	}

	return
}

// DeleteSubComments 删除评论
func (*commentServiceImpl) DeleteSubComments(req *dto.DeleteSubCommentReq, ctx *gin.Context) (resp *dto.DeleteCategoryResp, err error) {
	err = biz.Q.Transaction(func(tx *biz.Query) error {
		cmtTx := tx.Comment
		tableName := cmtTx.TableName()

		condList := []gen.Condition{}
		if req.CondList != nil {
			for _, cond := range req.CondList {
				p, e := cond.ParseCond(tableName)
				if e != nil {
					return e
				}
				condList = append(condList, p)
			}
		}

		// 确保不会直接删到根评论
		condList = append(condList, cmtTx.RootCommentId.Neq(0))

		_, e := cmtTx.WithContext(ctx).
			Where(condList...).
			Delete()

		return e
	})

	if err != nil {
		err = util.CreateBizErr("删除评论失败", err)
		return
	}
	resp = &dto.DeleteCategoryResp{}
	return
}

// DeleteSubComments 删除评论
func (*commentServiceImpl) DeleteRootComments(req *dto.DeleteRootCommentReq, ctx *gin.Context) (resp *dto.DeleteRootCommentResp, err error) {
	err = biz.Q.Transaction(func(tx *biz.Query) error {
		cmtTx := tx.Comment
		_, e := cmtTx.WithContext(ctx).
			Where(cmtTx.ID.In(req.IDList...)).
			Or(cmtTx.RootCommentId.In(req.IDList...)).
			Delete()

		return e
	})

	if err != nil {
		err = util.CreateBizErr("删除根评论失败", err)
		return
	}
	resp = &dto.DeleteRootCommentResp{}
	return
}
