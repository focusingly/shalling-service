package dto

type (
	CreateCommentReq struct {
		PostID         int64  `json:"postID"`         // 文章的 ID
		RootCommentID  int64  `json:"rootCommentID"`  // 根评论的 ID, 如果为 0, 表示自身就是根评论
		ReplyToID      int64  `json:"replyToID"`      // 回复的上一条评论的 ID, 如果为 0, 表示没回复任何人, 即当前为根评论
		Content        string `json:"content"`        // 回复内容
		SubEmailNotify bool   `json:"subEmailNotify"` // 是否订阅邮件回复通知
	}
	CreateCommentResp struct{}

	GetRootCommentPagesReq  struct{}
	GetRootCommentPagesResp struct{}

	GetSubCommentPagesReq  struct{}
	GetSubCommentPagesResp struct{}
)
