package dto

import (
	"space-domain/model"
)

type (
	UpdateOrCreatePostReq struct {
		WarningOverride `json:"warningOverride" yaml:"warningOverride" xml:"warningOverride" toml:"warningOverride"`
		PostId          int64    `json:"postId" yaml:"postId" xml:"postId" toml:"postId"`
		AuthorId        int64    `json:"authorId" yaml:"authorId" xml:"authorId" toml:"authorId"`
		Hide            byte     `json:"hide" yaml:"hide" xml:"hide" toml:"hide"`
		Title           string   `json:"title" yaml:"title" xml:"title" toml:"title"`
		Content         string   `json:"content" yaml:"content" xml:"content" toml:"content"`
		WordCount       int64    `json:"wordCount" yaml:"wordCount" xml:"wordCount" toml:"wordCount"`
		ReadTime        *int64   `json:"readTime" yaml:"readTime" xml:"readTime" toml:"readTime"`
		PubTime         *int64   `json:"pubTime" yaml:"pubTime" xml:"pubTime" toml:"pubTime"`
		Category        *string  `json:"category" yaml:"category" xml:"category" toml:"category"`
		Tags            []string `json:"tags" yaml:"tags" xml:"tags" toml:"tags"`
		Snippet         *string  `json:"snippet" yaml:"snippet" xml:"snippet" toml:"snippet"`
		LastPubTime     *int64   `json:"lastPubTime" yaml:"lastPubTime" xml:"lastPubTime" toml:"lastPubTime"`
		Weight          *int     `json:"weight" yaml:"weight" xml:"weight" toml:"weight"`
		Views           *int64   `json:"views" yaml:"views" xml:"views" toml:"views"`
		UpVote          *int64   `json:"upVote" yaml:"upVote" xml:"upVote" toml:"upVote"`
		DownVote        *int64   `json:"downVote" yaml:"downVote" xml:"downVote" toml:"downVote"`
		AllowComment    int      `json:"allowComment" yaml:"allowComment" xml:"allowComment" toml:"allowComment"`
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
		PostID int64 `json:"postID" yaml:"postID" xml:"postID" toml:"postID"`
	}
	GetPostDetailResp struct {
		model.Post
	}

	DeletePostByIdListReq struct {
		IdList []int64 `json:"idList"`
	}
	DeletePostByIdListResp struct {
	}

	GetPostByTagNameReq struct {
		TagName string `form:"tagName" json:"tagName" yaml:"tagName" xml:"tagName" toml:"tagName"`
	}
	GetPostByTagNameResp struct {
		Tag   *model.Tag    `json:"tag" yaml:"tag" xml:"tag" toml:"tag"`
		Posts []*model.Post `json:"posts" yaml:"posts" xml:"posts" toml:"posts"`
	}

	// 全文检索
	GlobalSearchReq struct {
		Keyword string `form:"keyword" json:"keyword" yaml:"keyword" xml:"keyword" toml:"keyword"`
	}

	GlobalSearchResp struct {
		model.PageList[*model.Post]
	}
)
