package dto

import (
	"space-api/dto/query"
	"space-domain/model"
)

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
	GetSubCommentPagesResp = model.PageList[*model.Comment]

	NestedComments struct {
		RootComment *model.Comment `json:"rootComment"`
		Subs        *model.PageList[*model.Comment]
	}
	GetRootCommentPagesReq struct {
		BasePageParam `json:"basePageParam"`
		PostID        int64 `form:"postID" json:"postId"`
	}
	GetRootCommentPagesResp struct {
		model.PageList[*NestedComments]
	}

	DeleteSubCommentReq struct {
		CondList []*query.WhereCond `json:"condList" yaml:"condList" xml:"condList" toml:"condList"`
	}
	DeleteSubCommentResp struct{}

	DeleteRootCommentReq struct {
		IDList []int64 `json:"idList" yaml:"idList" xml:"idList" toml:"idList"`
	}
	DeleteRootCommentResp struct{}

	UpdateCommentReq struct {
		ID       int64  `json:"id" yaml:"id" xml:"id" toml:"id"`
		Content  string `json:"content" yaml:"content" xml:"content" toml:"content"`
		UpVote   *int64 `json:"upVote" yaml:"upVote" xml:"upVote" toml:"upVote"`
		DownVote *int64 `json:"downVote" yaml:"downVote" xml:"downVote" toml:"downVote"`
		Hide     bool   `json:"hide" yaml:"hide" xml:"hide" toml:"hide"`
	}
	UpdateCommentResp struct{}
)
