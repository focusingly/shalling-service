package service

import (
	"slices"
	"space-api/dto"
	"space-api/middleware"
	"space-api/util"
	"space-domain/dao/biz"
	"space-domain/model"
	"strings"

	"github.com/gin-gonic/gin"
)

func UpdateOrCreatePost(req *dto.UpdatePostReq, ctx *gin.Context) (resp *dto.UpdatePostResp, err error) {
	// 被创建/更新的 文章的 ID
	var postId int64 = 0

	err = biz.Q.Transaction(func(tx *biz.Query) error {
		postOp := tx.Post
		tagOp := tx.Tag
		postTagRelationOp := tx.PostTagRelation

		// 查找已经存在的文章
		if findPost, err := postOp.
			WithContext(ctx).
			Where(postOp.Id.Eq(req.PostId)).
			First(); err != nil {
			// 如果当前不存在文章则直接创新新的文章
			// 同步更新 ID
			postId = util.GetSnowFlakeNode().Generate().Int64()
			if u, err := middleware.GetCurrentLoginUser(ctx); err != nil {
				return err
			} else {
				// 创建新的文章
				if err := postOp.WithContext(ctx).Create(&model.Post{
					BaseColumn:   model.BaseColumn{Id: postId},
					Title:        req.Title,
					AuthorId:     u.Id,
					Content:      *req.Category,
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
				}); err != nil {
					return err
				}
			}
		} else {
			// 更新 ID
			postId = findPost.Id
			// 存在的情况下进行更新
			if _, err := postOp.WithContext(ctx).Updates(findPost); err != nil {
				return err
			}
		}

		// 后台管理设置的新标签
		trimmedTags := []string{}
		for _, tag := range strings.Split(req.Tags, ",") {
			trimmedTags = append(trimmedTags, strings.TrimSpace(tag))
		}

		// 去重查找所有的标签
		if distinctTags, err := tagOp.WithContext(ctx).
			Distinct(tagOp.TagName).
			Select(tagOp.TagName).
			Find(); err != nil {
			return err
		} else {
			// 通过过滤只留下需要新创建的标签
			filterTags := slices.DeleteFunc(trimmedTags, func(t string) bool {
				return slices.ContainsFunc(distinctTags, func(e *model.Tag) bool {
					return e.TagName == t
				})
			})
			// 生成需要创建的新标签
			shouldCreateTags := []*model.Tag{}
			for _, t := range filterTags {
				shouldCreateTags = append(shouldCreateTags, &model.Tag{
					TagName: t,
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
		// 清空已经存在的所有 文章-标签映射关系, 并重新创建
		if _, err := postTagRelationOp.WithContext(ctx).Where(postTagRelationOp.PostId.Eq(postId)).Delete(); err != nil {
			return err
		}
		// 重新创建映射关系
		if secs, err := tagOp.WithContext(ctx).Where(tagOp.TagName.In(trimmedTags...)).Find(); err != nil {
			return err
		} else {
			// 组装映射
			rl := []*model.PostTagRelation{}
			for _, t := range secs {
				rl = append(rl, &model.PostTagRelation{
					TagId:  t.Id,
					PostId: postId,
				})
			}
			// 创建关系
			if err := postTagRelationOp.WithContext(ctx).Create(rl...); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return
	}

	if v, err := biz.Post.WithContext(ctx).
		Where(biz.Post.Id.Eq(postId)).
		First(); err != nil {
		return nil, err
	} else {
		resp = &dto.UpdatePostResp{
			Post: *v,
		}
	}

	return
}
