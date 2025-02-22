package comment

import (
	"fmt"
	"space-api/constants"
	"space-api/dto"
	"space-api/middleware/auth"
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

type commentService struct{}

var DefaultCommentService = &commentService{}

type detailTransferUser struct {
	Username    string
	Avatar      *string
	HomePageURL *string
}

func (servicePtr *commentService) VerifyAndCreateComment(req *dto.CreateCommentReq, ctx *gin.Context) (resp *dto.CreateCommentResp, err error) {
	loginSession, uErr := auth.GetCurrentLoginSession(ctx)

	if uErr != nil {
		return nil, uErr
	}

	postTx := biz.Post

	// 先找到文章
	post, e := postTx.WithContext(ctx).Where(postTx.ID.Eq(req.PostID)).Take()
	if e != nil {
		err = util.CreateBizErr("不存在文章", e)
		return
	}

	// 如果文章设为了不可见, 不可评论的情况, 分别判断情况
	if post.Hide != 0 || (post.AllowComment == 0) {
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

	commentTx := biz.Comment

	// 发表的评论信息如果是子评论, 那么要先检查根评论和他的父评论合法性在才能创建
	if req.ReplyToID != 0 || req.RootCommentID != 0 {
		// 查找根评论(只有根评论存在, 才有底下的评论)
		rootCmt, e := commentTx.WithContext(ctx).
			Where(
				commentTx.ID.Eq(req.RootCommentID), // 根评论的 ID
				commentTx.PostId.Eq(req.PostID),    // 文章存在
			).
			Take()
		if e != nil {
			// 不存在评论
			err = util.CreateBizErr("评论不可用", e)
			return
		}

		// 根评论为隐藏的条件下
		if rootCmt.Hide != 0 && loginSession.UserType != constants.Admin {
			err = util.CreateBizErr("评论不可用", fmt.Errorf("comment not available"))
			return
		}

		// 查找父评论,
		parentCmt, e := commentTx.WithContext(ctx).
			Where(
				commentTx.ID.Eq(req.ReplyToID), // 父评论的 ID(父评论本身可能就是根评论, 也可能是根评论下另一条较早发布的子评论)
				commentTx.PostId.Eq(req.PostID),
			).
			Take()
		if e != nil {
			// 不存在评论
			err = util.CreateBizErr("评论不可用", e)
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

func (servicePtr *commentService) UpdateComment(req *dto.UpdateCommentReq, ctx *gin.Context) (resp *dto.UpdateCommentResp, err error) {
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

// CreateCommentDirect 充当基础通用方法, 只创建, 不进行额外的关联关系校验
func (*commentService) createCommentDirect(req *dto.CreateCommentReq, ctx *gin.Context) (resp *dto.CreateCommentResp, err error) {
	err = biz.Q.Transaction(func(tx *biz.Query) error {
		commentTx := tx.Comment
		loginSession, uErr := auth.GetCurrentLoginSession(ctx)
		if uErr != nil {
			return uErr
		}

		var detailUserTmp *detailTransferUser
		switch loginSession.UserType {
		case constants.GoogleUser, constants.GithubUser: // oauth2 用户评论
			oauthUserTx := tx.OAuth2User
			oauthUser, e := oauthUserTx.WithContext(ctx).
				Select(oauthUserTx.Username, oauthUserTx.AvatarURL, oauthUserTx.HomepageLink).
				Where(
					oauthUserTx.ID.Eq(loginSession.UserID),
				).
				Take()
			if e != nil {
				return e
			} else {
				detailUserTmp = &detailTransferUser{
					Username:    oauthUser.Username,
					Avatar:      oauthUser.AvatarURL,
					HomePageURL: oauthUser.HomepageLink,
				}
			}

		case constants.LocalUser, constants.Admin: // 本地账户评论(包含管理员账户)
			localUserTx := tx.LocalUser
			localUser, e := localUserTx.WithContext(ctx).
				Select(localUserTx.DisplayName, localUserTx.AvatarURL, localUserTx.HomepageLink).
				Where(localUserTx.ID.Eq(loginSession.UserID)).
				Take()
			if e != nil {
				return e
			} else {
				detailUserTmp = &detailTransferUser{
					Username:    localUser.DisplayName,
					Avatar:      localUser.AvatarURL,
					HomePageURL: localUser.HomepageLink,
				}
			}

		default:
			return fmt.Errorf("un-support user type: %s", loginSession.UserType)
		}

		realIpAddr := inbound.GetRealIpWithContext(ctx)
		uaDetail := inbound.GetUserAgentFromContext(ctx)
		ipSource, _ := ip.GetIpSearcher().SearchByStr(realIpAddr)

		// 进行创建
		e := commentTx.WithContext(ctx).Create(
			&model.Comment{
				PostId:        req.PostID,
				UserId:        loginSession.ID,
				UserType:      loginSession.UserType,
				Avatar:        detailUserTmp.Avatar,
				HomePageURL:   detailUserTmp.HomePageURL,
				RootCommentId: req.RootCommentID,
				ReplyToId:     req.ReplyToID,
				Content:       strings.TrimSpace(req.Content),
				UpVote:        nil,
				DownVote:      nil,
				IpAddr:        realIpAddr,
				IpSource:      &ipSource,
				Useragent:     &uaDetail.Useragent,
				OS:            &uaDetail.OS,
				ClientName:    &uaDetail.ClientName,
				SubEmailReply: util.TernaryExpr(req.SubEmailNotify, 1, 0),
			},
		)
		if e != nil {
			return e
		}

		return nil
	})

	if err != nil {
		err = util.CreateBizErr("创建评论失败", err)
		return
	}

	resp = &dto.CreateCommentResp{}

	return
}

// GetVisibleRootCommentPages 获取根评论的分页
func (servicePtr *commentService) GetVisibleRootCommentPages(req *dto.GetRootCommentPagesReq, ctx *gin.Context) (resp *dto.GetRootCommentPagesResp, err error) {
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
func (*commentService) GetVisibleSubCommentPages(req *dto.GetSubCommentPagesReq, ctx *gin.Context) (resp *dto.GetSubCommentPagesResp, err error) {
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
func (servicePtr *commentService) GetAnyRootCommentPages(req *dto.GetRootCommentPagesReq, ctx *gin.Context) (resp *dto.GetRootCommentPagesResp, err error) {
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
func (*commentService) GetAnySubCommentPages(req *dto.GetSubCommentPagesReq, ctx *gin.Context) (resp *dto.GetSubCommentPagesResp, err error) {
	postCtx := biz.Post
	_, err = postCtx.WithContext(ctx).
		Where(postCtx.ID.Eq(req.PostID)).
		Take()
	if err != nil {
		err = util.CreateBizErr("当前文章不允许评论", err)
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
func (*commentService) DeleteSubComments(req *dto.DeleteSubCommentReq, ctx *gin.Context) (resp *dto.DeleteCategoryResp, err error) {
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
func (*commentService) DeleteRootComments(req *dto.DeleteRootCommentReq, ctx *gin.Context) (resp *dto.DeleteRootCommentResp, err error) {
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
