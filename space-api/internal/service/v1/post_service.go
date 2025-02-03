package service

import (
	"database/sql"
	"slices"
	"space-api/dto"
	"space-api/middleware"
	"space-api/util"
	"space-domain/dao/biz"
	"space-domain/model"
	"strings"

	"github.com/gin-gonic/gin"
)

// UpdateOrCreatePost 创建新的文章或者更新已有的文章
func UpdateOrCreatePost(req *dto.UpdatePostReq, ctx *gin.Context) (resp *dto.UpdatePostResp, err error) {
	// 被创建/更新的 文章的 ID
	var postId int64 = 0

	// 全部在事务内操作
	err = biz.Q.Transaction(func(tx *biz.Query) error {
		postOp := tx.Post
		tagOp := tx.Tag
		postTagRelationOp := tx.PostTagRelation

		// 查找已经存在的文章
		findPost, err := postOp.WithContext(ctx).Where(postOp.Id.Eq(req.PostId)).Take()
		// 如果当前不存在文章则直接创新新的文章
		if err != nil {
			// 同步更新 ID
			postId = util.GetSnowFlakeNode().Generate().Int64()
			// 获取当前登录的用户信息
			loginUser, err := middleware.GetCurrentLoginUser(ctx)
			if err != nil {
				return err
			}
			// 创建新的文章
			t := &model.Post{
				BaseColumn:   model.BaseColumn{Id: postId},
				Title:        req.Title,
				AuthorId:     loginUser.Id,
				Content:      req.Content,
				WordCount:    req.WordCount,
				ReadTime:     req.ReadTime,
				Category:     req.Category,
				Tags:         &req.Tags,
				LastPubTime:  req.LastPubTime,
				Weight:       req.Weight,
				Views:        req.Views,
				UpVote:       req.UpVote,
				DownVote:     req.DownVote,
				AllowComment: req.AllowComment,
			}
			if err := postOp.WithContext(ctx).Create(t); err != nil {
				return err
			}
		} else {
			// 表示文章存在, 操作为更新
			// 更改为文章的 ID
			postId = findPost.Id
			// 存在的情况下进行更新
			t := &model.Post{
				BaseColumn:   findPost.BaseColumn,
				Title:        req.Title,
				AuthorId:     findPost.AuthorId,
				Content:      req.Content,
				WordCount:    req.WordCount,
				ReadTime:     req.ReadTime,
				Category:     req.Category,
				Tags:         &req.Tags,
				LastPubTime:  req.LastPubTime,
				Weight:       req.Weight,
				Views:        req.Views,
				UpVote:       req.UpVote,
				DownVote:     req.DownVote,
				AllowComment: req.AllowComment,
			}

			if _, err := postOp.WithContext(ctx).Updates(t); err != nil {
				return err
			}
		}

		// 同步其它表的信息

		// 同步新的标签操作
		trimmedTags := []string{}
		for _, tag := range strings.Split(req.Tags, ",") {
			t := strings.TrimSpace(tag)
			// 过滤掉空格
			if t != "" {
				trimmedTags = append(trimmedTags, t)
			}
		}

		// 查找表里所有已经存在的标签
		distinctTags, err := tagOp.
			WithContext(ctx).
			Distinct(tagOp.TagName).
			Select(tagOp.TagName).
			Find()
		if err != nil {
			return err
		} else {
			// 通过过滤只留下需要新创建的标签
			filterTags := slices.DeleteFunc(slices.Clone(trimmedTags), func(tag string) bool {
				return slices.ContainsFunc(distinctTags, func(e *model.Tag) bool {
					return e.TagName == tag
				})
			})

			// 生成需要创建的新标签
			shouldCreateTags := []*model.Tag{}
			for _, tag := range filterTags {
				shouldCreateTags = append(shouldCreateTags, &model.Tag{
					TagName: tag,
				})
			}
			// 批量创建新标签
			if err := tagOp.WithContext(ctx).
				CreateInBatches(
					shouldCreateTags,
					util.TernaryExpr(len(shouldCreateTags) >= 64,
						64,
						len(shouldCreateTags),
					),
				); err != nil {
				return err
			}
		}
		// 先清空已经存在的所有 文章-标签映射关系, 然后重新创建
		// 删除所有文章 ID 为当前文章的 postTagRelation 记录
		if _, err := postTagRelationOp.WithContext(ctx).Where(postTagRelationOp.PostId.Eq(postId)).Delete(); err != nil {
			return err
		} else {
			// 查找所有在 post 里出现的 tag
			if findTags, err := tagOp.WithContext(ctx).Where(tagOp.TagName.In(trimmedTags...)).Find(); err != nil {
				return err
			} else {
				// 组装映射
				tagPostRelationsList := []*model.PostTagRelation{}
				for _, tag := range findTags {
					tagPostRelationsList = append(tagPostRelationsList, &model.PostTagRelation{
						TagId:  tag.Id,
						PostId: postId, // 当前这篇文章
					})
				}
				// 恢复映射关系
				if err := postTagRelationOp.WithContext(ctx).Create(tagPostRelationsList...); err != nil {
					return err
				}
			}
		}

		return nil
	}, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})

	// 判断前面的事务操作结果
	if err != nil {
		return
	}

	v, err := biz.Q.Post.WithContext(ctx).
		Where(biz.Q.Post.Id.Eq(postId)).
		Take()
	if err != nil {
		return nil, err
	} else {
		resp = &dto.UpdatePostResp{
			Post: *v,
		}
	}

	return
}

// GetPostList 获取文章分页的信息(不包括正文)
func GetPostList(req *dto.GetPostListReq, ctx *gin.Context) (resp *dto.GetPostListResp, err error) {

	return
}

// GetPostById 根据文章 ID 获取全量的文章信息
func GetPostById(req *dto.GetPostDetailReq, ctx *gin.Context) (resp *dto.GetPostDetailResp, err error) {
	val, err := biz.Post.WithContext(ctx).Where(biz.Post.Id.Eq(req.Id)).Take()
	if err != nil {
		return nil, &util.BizErr{
			Msg:    "查找文章失败",
			Reason: err,
		}
	}

	return &dto.GetPostDetailResp{
		Post: *val,
	}, nil
}

// GetPostById 根据文章 ID 获取全量的文章信息
func DeletePostByIdList(req *dto.DeletePostByIdListReq, ctx *gin.Context) (resp *dto.DeletePostByIdListResp, err error) {
	err = biz.Q.Transaction(func(tx *biz.Query) error {
		_, err = tx.Post.WithContext(ctx).Where(tx.Post.Id.In(req.IdList...)).Delete()
		return err
	})

	if err != nil {
		return nil, &util.BizErr{
			Msg:    "删除失败: " + err.Error(),
			Reason: err,
		}
	}

	return &dto.DeletePostByIdListResp{}, nil
}
