package dto

import "space-domain/model"

type (
	BizOp struct {
		// 允许强制覆盖可能影响的业务(比如关联资源)
		ForceOverride bool `json:"forceOverride"`
	}

	UpdatePostReq struct {
		BizOp
		PostId       int64   `json:"postId"`
		AuthorId     int64   `json:"authorId"`
		Hide         byte    `json:"hide"`
		Title        string  `json:"title"`
		Content      string  `json:"content"`
		WordCount    int64   `json:"wordCount"`
		ReadTime     *int64  `json:"readTime"`
		PubTime      *int64  `json:"pubTime"`
		Category     *string `json:"category"`
		Tags         string  `json:"tags"`
		LastPubTime  *int64  `json:"lastPubTime"`
		Weight       *int    `json:"weight"`
		Views        *int64  `json:"views"`
		UpVote       *int64  `json:"upVote"`
		DownVote     *int64  `json:"downVote"`
		AllowComment *byte   `json:"allowComment"`
	}

	UpdatePostResp struct {
		model.Post
	}
)
