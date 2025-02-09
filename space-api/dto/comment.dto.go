package dto

import "space-domain/model"

type (
	CreateCommentReq struct {
		PostID         int64  `json:"postID"`         // 文章的 ID
		RootCommentID  int64  `json:"rootCommentID"`  // 根评论的 ID, 如果为 0, 表示自身就是根评论
		ReplyToID      int64  `json:"replyToID"`      // 回复的上一条评论的 ID, 如果为 0, 表示没回复任何人, 即当前为根评论
		Content        string `json:"content"`        // 回复内容
		SubEmailNotify bool   `json:"subEmailNotify"` // 是否订阅邮件回复通知
	}
	CreateCommentResp struct{}

	GetSubCommentPagesReq struct {
		BasePageParam `json:"basePageParam"`
		PostID        int64 `form:"postID" json:"postId"`
		RootCommentID int64 `form:"rootCommentID" json:"rootCommentID"`
	}
	GetSubCommentPagesResp struct {
		model.PageList[*model.Comment]
	}

	NestedComments struct {
		RootComment *model.Comment `json:"rootComment"`
		Subs        model.PageList[*model.Comment]
	}
	GetRootCommentPagesReq struct {
		BasePageParam `json:"basePageParam"`
		PostID        int64 `form:"postID" json:"postId"`
	}
	GetRootCommentPagesResp struct {
		model.PageList[*NestedComments]
	}
)
