package service

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
)

type commentService struct{}

var DefaultCommentService = &commentService{}

type detailTransferUser struct {
	Username    string
	Avatar      *string
	HomePageURL *string
}

// CreateCommentDirect 充当基础通用方法, 只创建, 不进行额外的关联关系校验
func (*commentService) createCommentDirect(req *dto.CreateCommentReq, ctx *gin.Context) (resp *dto.CreateCommentResp, err error) {
	err = biz.Q.Transaction(func(tx *biz.Query) error {
		commentTx := tx.Comment
		loginSession, _ := auth.GetCurrentLoginSession(ctx)
		var detailUserTmp *detailTransferUser
		switch loginSession.UserType {
		case constants.GoogleUser, constants.GithubUser:
			oauthUserTx := tx.OAuth2User
			oauthUser, e := oauthUserTx.WithContext(ctx).
				Select(oauthUserTx.Username, oauthUserTx.AvatarURL, oauthUserTx.HomepageLink).
				Where(
					oauthUserTx.ID.Eq(loginSession.UserId),
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
		case constants.LocalUser:
			localUserTx := tx.LocalUser
			localUser, e := localUserTx.WithContext(ctx).
				Select(localUserTx.DisplayName, localUserTx.AvatarURL, localUserTx.HomepageLink).
				Where(localUserTx.ID.Eq(loginSession.UserId)).
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
		ua := inbound.GetUserAgentFromContext(ctx)
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
				UpVote:        new(int64),
				DownVote:      new(int64),
				IpAddr:        realIpAddr,
				IpSource:      &ipSource,
				Useragent:     &ua.Useragent,
				OS:            &ua.OS,
				ClientName:    &ua.ClientName,
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

func (servicePtr *commentService) SimpleVerifyAndCreateComment(req *dto.CreateCommentReq, ctx *gin.Context) (resp *dto.CreateCommentResp, err error) {
	loginSession, _ := auth.GetCurrentLoginSession(ctx)
	postTx := biz.Post
	// 找到文章
	post, e := postTx.WithContext(ctx).Where(postTx.ID.Eq(req.PostID)).Take()
	if e != nil {
		err = util.CreateBizErr("不存在文章", e)
		return
	}

	// 如果文章设为了不可见, 不可评论的情况, 分别判断情况
	if post.Hide != 0 || (post.AllowComment == 0) {
		// 此时不允许 oauth2 登录的用户进行评论
		if loginSession.UserType != constants.LocalUser {
			err = util.CreateBizErr("评论不可用", fmt.Errorf("comment not available"))
			return
		}

		// 如果是本地用户, 只允许管理员进行操作
		localUserCtx := biz.LocalUser
		_, e := localUserCtx.WithContext(ctx).
			Where(localUserCtx.IsAdmin.Neq(0), localUserCtx.ID.Eq(loginSession.ID)).
			Take()
		if e != nil {
			err = util.CreateBizErr("当前本地账户评论功能不可用", fmt.Errorf("comment not available"))
			return
		}
	}

	commentTx := biz.Comment
	// 表示发表的 [不是] 根评论信息(即为 '次级/子' 评论)
	if req.ReplyToID != 0 || req.RootCommentID != 0 {
		// 先尝试查找根评论
		rootCmt, e := commentTx.WithContext(ctx).
			Where(commentTx.ID.Eq(req.RootCommentID), commentTx.PostId.Eq(req.PostID)).
			Take()
		if e != nil {
			err = util.CreateBizErr("评论不可用", e)
			return
		}

		// 评论为隐藏的条件下
		if rootCmt.Hide != 0 {
			if loginSession.UserType != constants.LocalUser {
				err = util.CreateBizErr("评论不可用", fmt.Errorf("comment not available"))
				return
			}

			// 如果是本地用户, 只允许管理员进行操作
			localUserCtx := biz.LocalUser
			_, e := localUserCtx.WithContext(ctx).
				Where(localUserCtx.IsAdmin.Neq(0), localUserCtx.ID.Eq(loginSession.ID)).
				Take()
			if e != nil {
				err = util.CreateBizErr("当前本地账户评论功能不可用", fmt.Errorf("comment not available"))
				return
			}
		}
	}
	// 如果是发表子评论信息, 文章允许评论那么都可以通过

	return servicePtr.createCommentDirect(req, ctx)
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
				if splits[len(splits)-1] != "0" {
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
			Subs:        subs.PageList,
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
				if splits[len(splits)-1] != "0" {
					t += splits[len(splits)-1]
				}
				cmt.IpSource = &t
			} else {
				cmt.IpSource = ptr.ToPtr("")
			}
		}
	}

	resp = &dto.GetSubCommentPagesResp{
		PageList: model.PageList[*model.Comment]{
			List:  list,
			Page:  int64(*req.Page),
			Size:  int64(*req.Size),
			Total: count,
		},
	}

	return
}
