package dto

import (
	"space-domain/model"
)

type (
	UpdateOrCreatePostReq struct {
		WarningOverride
		PostId       int64    `json:"postId"`
		AuthorId     int64    `json:"authorId"`
		Hide         byte     `json:"hide"`
		Title        string   `json:"title"`
		Content      string   `json:"content"`
		WordCount    int64    `json:"wordCount"`
		ReadTime     *int64   `json:"readTime"`
		PubTime      *int64   `json:"pubTime"`
		Category     *string  `json:"category"`
		Tags         []string `json:"tags"`
		LastPubTime  *int64   `json:"lastPubTime"`
		Weight       *int     `json:"weight"`
		Views        *int64   `json:"views"`
		UpVote       *int64   `json:"upVote"`
		DownVote     *int64   `json:"downVote"`
		AllowComment int      `json:"allowComment"`
	}

	UpdateOrCreatePostResp struct {
		model.Post
	}

	GetPostPageListReq struct {
		BasePageParam
	}

	GetPostPageListResp struct {
		model.PageList[*model.Post]
	}

	GetPostDetailReq struct {
		Id int64 `uri:"id" json:"id"`
	}

	GetPostDetailResp struct {
		model.Post
	}

	DeletePostByIdListReq struct {
		IdList []int64 `json:"idList"`
	}

	DeletePostByIdListResp struct {
	}
)
