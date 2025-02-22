package service

import (
	"context"
	"fmt"
	"slices"
	"space-api/constants"
	"space-api/dto"
	"space-api/middleware/auth"
	"space-api/util"
	"space-api/util/arr"
	"space-api/util/id"
	"space-api/util/performance"
	"space-api/util/ptr"
	"space-domain/dao/biz"
	"space-domain/model"
	"strings"

	"github.com/bytedance/gopkg/util/gopool"
	"github.com/gin-gonic/gin"
)

type postService struct {
	searchService *_searchService
	executor      gopool.Pool
	visitCache    *performance.JsonCache
}

var DefaultPostService = &postService{
	searchService: DefaultGlobalSearchService,
	executor:      performance.DefaultTaskRunner,
	visitCache:    performance.NewCache(constants.MB * 4),
}

// CreateOrUpdatePost 创建/更新文章, 取决于是否存在已有的文章
func (s *postService) CreateOrUpdatePost(req *dto.UpdateOrCreatePostReq, ctx *gin.Context) (resp *dto.UpdateOrCreatePostResp, err error) {
	// 被创建/更新的 文章的 ID
	var postId int64 = 0

	// 全部在事务内操作
	err = biz.Q.Transaction(func(tx *biz.Query) error {
		postOp := tx.Post
		tagOp := tx.Tag
		postTagRelationOp := tx.PostTagRelation

		// 标准化标签(去除首尾空格和移除纯空白字符串)
		if req.Tags != nil {
			filters := []string{}
			for _, tag := range req.Tags {
				t := strings.TrimSpace(tag)
				if t != "" {
					filters = append(filters, t)
				}
			}
			if len(filters) == 0 {
				req.Tags = nil
			} else {
				req.Tags = filters
			}
		}

		// 查找已经存在的文章
		findPost, err := postOp.WithContext(ctx).Where(postOp.ID.Eq(req.PostId)).Take()
		// 如果当前不存在文章则直接创新新的文章
		if err != nil {
			// 同步更新 ID
			postId = id.GetSnowFlakeNode().Generate().Int64()
			// 获取当前登录的用户信息
			loginUser, err := auth.GetCurrentLoginSession(ctx)
			if err != nil {
				return err
			}

			// 创建新的文章
			t := &model.Post{
				BaseColumn:   model.BaseColumn{ID: postId},
				Title:        req.Title,
				AuthorId:     loginUser.ID,
				Content:      req.Content,
				WordCount:    req.WordCount,
				ReadTime:     req.ReadTime,
				Category:     req.Category,
				Tags:         req.Tags,
				Snippet:      req.Snippet,
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
			postId = findPost.ID
			// 存在的情况下进行更新
			t := &model.Post{
				BaseColumn:   findPost.BaseColumn,
				Title:        req.Title,
				AuthorId:     findPost.AuthorId,
				Content:      req.Content,
				WordCount:    req.WordCount,
				ReadTime:     req.ReadTime,
				Category:     req.Category,
				Tags:         req.Tags,
				LastPubTime:  req.LastPubTime,
				Weight:       req.Weight,
				Views:        req.Views,
				Snippet:      req.Snippet,
				UpVote:       req.UpVote,
				DownVote:     req.DownVote,
				AllowComment: req.AllowComment,
			}

			_, err := postOp.WithContext(ctx).
				Where(postOp.ID.Eq(postId)).
				Select(
					postOp.ID,
					postOp.Hide,
					postOp.Snippet,
					postOp.Title,
					postOp.AuthorId,
					postOp.Content,
					postOp.WordCount,
					postOp.ReadTime,
					postOp.Category,
					postOp.Tags,
					postOp.LastPubTime,
					postOp.Weight,
					postOp.Views,
					postOp.Snippet,
					postOp.UpVote,
					postOp.DownVote,
					postOp.AllowComment,
				).
				Updates(t)
			if err != nil {
				return err
			}
		}

		// 同步其它表的信息
		// 同步新的标签操作
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
			filterTags := slices.DeleteFunc(
				slices.Clone(util.TernaryExpr(
					req.Tags != nil,
					req.Tags,
					[]string{},
				)),
				func(tag string) bool {
					return slices.ContainsFunc(distinctTags, func(e *model.Tag) bool {
						// 去掉所有已经存在的标签, 避免重复创建
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
			err := tagOp.WithContext(ctx).
				CreateInBatches(
					shouldCreateTags,
					64,
				)
			if err != nil {
				return err
			}
		}

		// 先清空已经存在的所有 文章-标签映射关系, 然后重新创建
		// 删除所有文章 ID 为当前文章的 postTagRelation 记录
		_, err = postTagRelationOp.WithContext(ctx).Where(postTagRelationOp.PostId.Eq(postId)).Delete()
		// 删除失败
		if err != nil {
			return err
		} else {
			// 删除成功, 需要重新恢复映射关系
			// 查找所有在 post 里出现的 tag
			findRequireTags, err := tagOp.WithContext(ctx).
				Distinct(tagOp.TagName).
				Select(tagOp.ID).
				Where(tagOp.TagName.In(req.Tags...)).
				Find()

			if err != nil {
				return err
			}
			tagPostRelationsList := arr.MapSlice(findRequireTags, func(i int, tag *model.Tag) *model.PostTagRelation {
				return &model.PostTagRelation{
					TagId:  tag.ID,
					PostId: postId, // 当前的这篇文章
				}
			})

			// 恢复映射关系
			if err := postTagRelationOp.WithContext(ctx).CreateInBatches(tagPostRelationsList, 64); err != nil {
				return err
			}
		}

		// 更新分类的关系
		if req.Category != nil {
			catTX := tx.Category
			_, e := catTX.WithContext(ctx).Take()
			// 如果分类还不存在, 那么创建对应的分类信息
			if e != nil {
				if e := catTX.WithContext(ctx).Create(&model.Category{
					CategoryName: *req.Category,
				}); e != nil {
					return e
				}
			}
		}

		return nil
	})

	// 判断前面的事务操作结果
	if err != nil {
		err = &util.BizErr{
			Reason: err,
			Msg:    "更新/创建文章失败: " + err.Error(),
		}

		return
	}

	post, err := biz.Q.Post.WithContext(ctx).
		Where(biz.Q.Post.ID.Eq(postId)).
		Take()
	if err != nil {
		err = util.CreateBizErr("更新/创建文章失败", err)
		return
	}
	resp = &dto.UpdateOrCreatePostResp{
		Post: *post,
	}

	s.executor.Go(func() {
		if resp.Hide != 0 {
			s.searchService.DeletePostSearchIndex([]int64{postId}, context.Background())
		} else {
			s.searchService.UpdatePostSearchIndex(post)
		}
	})

	return
}

// GetAllPostList 获取所有文章分页的信息(不包括正文内容)
func (s *postService) GetAllPostList(req *dto.GetPostPageListReq, ctx *gin.Context) (resp *dto.GetPostPageListResp, err error) {
	postOp := biz.Post

	result, count, err := postOp.
		WithContext(ctx).
		Select(postOp.ID,
			postOp.CreatedAt,
			postOp.UpdatedAt,
			postOp.Hide,
			postOp.Title,
			postOp.AuthorId,
			postOp.WordCount,
			postOp.Snippet,
			postOp.ReadTime,
			postOp.Category,
			postOp.Tags,
			postOp.LastPubTime,
			postOp.Weight,
			postOp.Views,
			postOp.UpVote,
			postOp.DownVote,
			postOp.AllowComment,
		).
		FindByPage(req.Normalize())

	if err != nil {
		return nil, &util.BizErr{
			Msg:    "查询错误",
			Reason: err,
		}
	}

	for _, post := range result {
		post.Views = ptr.ToPtr(s.GetCachedViewCountOrFallback(post, false))
	}

	return &dto.GetPostPageListResp{
		PageList: model.PageList[*model.Post]{
			List:  result,
			Page:  int64(*req.Page),
			Size:  int64(*req.Size),
			Total: count,
		},
	}, nil
}

// GetCachedViewCountOrFallback 获取文章在缓存中的技术值, 如果为公开访问, 那么还递增对应的缓存计数
func (s *postService) GetCachedViewCountOrFallback(post *model.Post, isPubMode bool) int64 {
	key := fmt.Sprintf("%d", post.ID)
	shouldIncr := util.TernaryExpr[int64](isPubMode, 1, 0)

	var fallbackViewCount = util.TernaryExprWithProducer(
		post.Views != nil,
		func() int64 {
			return shouldIncr + (*post.Views)
		},
		func() int64 {
			return shouldIncr
		},
	)

	var getViews int64 = 0
	var err error = nil

	if isPubMode {
		getViews, err = s.visitCache.IncAndGet(key, 1, 0)
	} else {
		getViews, err = s.visitCache.GetInt64(key)
	}

	// 处理访问量合法性
	switch {
	case
		err != nil,                   // 缓存中不存在
		fallbackViewCount > getViews: // 缓存中的计数值比数据库记录值低, 那么需要使用数据库提供的

		s.visitCache.Set(key, fallbackViewCount)
		getViews = fallbackViewCount
	}

	return getViews
}

// GetVisiblePostList 获取可见文章分页的信息(不包括正文内容)
func (s *postService) GetVisiblePostList(req *dto.GetPostPageListReq, ctx *gin.Context) (resp *dto.GetPostPageListResp, err error) {
	postOp := biz.Post

	// 只允许可见的文章
	result, count, err := postOp.
		WithContext(ctx).
		Where(postOp.Hide.Eq(0)).
		Select(postOp.ID,
			postOp.CreatedAt,
			postOp.UpdatedAt,
			postOp.Hide,
			postOp.Title,
			postOp.AuthorId,
			// postOp.Content,
			postOp.Snippet,
			postOp.WordCount,
			postOp.ReadTime,
			postOp.Category,
			postOp.Tags,
			postOp.LastPubTime,
			postOp.Weight,
			postOp.Views,
			postOp.UpVote,
			postOp.DownVote,
			postOp.AllowComment,
		).
		Order(postOp.CreatedAt.Distinct()).
		FindByPage(req.Normalize())

	if err != nil {
		return nil, &util.BizErr{
			Msg:    "查询错误",
			Reason: err,
		}
	}

	for _, post := range result {
		post.Views = ptr.ToPtr(s.GetCachedViewCountOrFallback(post, false))
	}

	return &dto.GetPostPageListResp{
		PageList: model.PageList[*model.Post]{
			List:  result,
			Page:  int64(*req.Page),
			Size:  int64(*req.Size),
			Total: count,
		},
	}, nil
}

// GetPostById 根据文章 ID 获取全量的文章信息
func (s *postService) GetPostById(req *dto.GetPostDetailReq, ctx *gin.Context) (resp *dto.GetPostDetailResp, err error) {
	val, err := biz.Post.WithContext(ctx).Where(biz.Post.ID.Eq(req.PostID)).Take()
	if err != nil {
		return nil, &util.BizErr{
			Msg:    "查找文章失败",
			Reason: err,
		}
	}

	val.Views = ptr.ToPtr(s.GetCachedViewCountOrFallback(val, false))

	return &dto.GetPostDetailResp{
		Post: *val,
	}, nil
}

func (s *postService) SyncAllPostViews(ctx context.Context) (err error) {
	err = biz.Q.Transaction(func(tx *biz.Query) error {
		bizTx := tx.Post

		// 找到所有公开的文章, 和缓存的访问数量进行对比/同步
		list, e := bizTx.WithContext(ctx).
			Select(bizTx.ID, bizTx.Views).
			Find()
		if e != nil {
			return e
		}

		// 逐条更新文章的访问量
		for _, post := range list {
			key := fmt.Sprintf("%d", post.ID)
			fallbackViews := util.TernaryExprWithProducer(
				post.Views != nil,
				func() int64 {
					return *post.Views
				}, func() int64 {
					return 0
				})

			cachedCount, cachedErr := s.visitCache.GetInt64(key)
			switch {
			case cachedErr != nil:
				s.visitCache.Set(key, fallbackViews, 0)
			case fallbackViews > cachedCount:
				s.visitCache.Set(key, fallbackViews, 0)
			default:
				// 更新缓存访问数到数据库当中
				if _, e = bizTx.WithContext(ctx).
					Where(bizTx.ID.Eq(post.ID)).
					Update(bizTx.Views, cachedCount); e != nil {
					return e
				}
			}
		}

		return nil
	})
	if err != nil {
		err = util.CreateBizErr("同步页面访问数失败: "+err.Error(), err)
	}

	return
}

func (s *postService) ClearPostsViewsCache() {
	s.visitCache.ClearAll()
}

func (s *postService) ExpirePubViewsCacheByID(postID int64) {
	s.visitCache.Delete(fmt.Sprintf("%d", postID))
}

// GetVisiblePostById 根据文章 ID 获取可见的全量的文章信息
func (s *postService) GetVisiblePostById(req *dto.GetPostDetailReq, ctx *gin.Context) (resp *dto.GetPostDetailResp, err error) {
	pubPost, err := biz.Post.WithContext(ctx).
		Where(
			biz.Post.ID.Eq(req.PostID),
			biz.Post.Hide.Eq(0),
		).
		Take()
	if err != nil {
		return nil, &util.BizErr{
			Msg:    "查找文章失败",
			Reason: err,
		}
	}

	pubPost.Views = ptr.ToPtr(s.GetCachedViewCountOrFallback(pubPost, true))

	return &dto.GetPostDetailResp{
		Post: *pubPost,
	}, nil
}

// GetPostById 根据文章 ID 获取全量的文章信息
func (s *postService) DeletePostByIdList(req *dto.DeletePostByIdListReq, ctx *gin.Context) (resp *dto.DeletePostByIdListResp, err error) {
	err = biz.Q.Transaction(func(tx *biz.Query) error {
		_, err = tx.Post.WithContext(ctx).Where(tx.Post.ID.In(req.IdList...)).Delete()
		return err
	})

	if err != nil {
		return nil, &util.BizErr{
			Msg:    "删除失败: " + err.Error(),
			Reason: err,
		}
	}

	// 清空相关的全文搜索索引
	s.executor.Go(func() {
		for _, id := range req.IdList {
			s.ExpirePubViewsCacheByID(id)
		}
		s.searchService.DeletePostSearchIndex(req.IdList, context.Background())
	})

	return &dto.DeletePostByIdListResp{}, nil
}

func (s *postService) GetVisiblePostsByTagName(req *dto.GetPostByTagNameReq, ctx *gin.Context) (resp *dto.GetPostByTagNameResp, err error) {
	tagOp := biz.Tag
	// 找到匹配的标签
	tag, err := tagOp.WithContext(ctx).
		Where(tagOp.TagName.Eq(req.TagName), tagOp.Hide.Eq(0)).
		Take()
	if err != nil {
		err = util.CreateBizErr("没有相关的标签", err)
		return
	}

	tagPostOp := biz.PostTagRelation
	postOp := biz.Post

	// select post.* from postTagRelation left join post on post.id.eq postTagRelation.post_id
	// where post.hide = 0
	postsList := []*model.Post{}
	err = tagPostOp.WithContext(ctx).
		Select(
			postOp.ID,
			postOp.CreatedAt,
			postOp.UpdatedAt,
			postOp.Hide,
			postOp.Title,
			postOp.AuthorId,
			// postOp.Content,省略文章内容, 减少传输
			postOp.WordCount,
			postOp.ReadTime,
			postOp.Snippet,
			postOp.Category,
			postOp.Tags,
			postOp.LastPubTime,
			postOp.Weight,
			postOp.Views,
			postOp.UpVote,
			postOp.DownVote,
			postOp.Lang,
			postOp.AllowComment,
		).
		LeftJoin(
			postOp,
			postOp.ID.EqCol(tagPostOp.PostId),
		).
		Where(postOp.Hide.Eq(0)).
		Scan(&postsList)
	if err != nil {
		err = util.CreateBizErr("查找数据失败", err)
		return
	}

	for _, post := range postsList {
		post.Views = ptr.ToPtr(s.GetCachedViewCountOrFallback(post, false))
	}
	resp = &dto.GetPostByTagNameResp{
		Tag:   tag,
		Posts: postsList,
	}

	return
}
