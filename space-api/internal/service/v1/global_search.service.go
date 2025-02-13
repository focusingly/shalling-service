package service

import (
	"bytes"
	"context"
	"space-api/db"
	"space-api/dto"
	"space-api/util"
	"space-api/util/arr"
	"space-api/util/ptr"
	"space-api/util/str"
	"space-domain/dao/biz"
	"space-domain/model"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yanyiwu/gojieba"
	"gorm.io/gorm"
)

type _searchService struct {
	cutter *gojieba.Jieba
	*gorm.DB
}

var DefaultGlobalSearchService = &_searchService{
	cutter: gojieba.NewJieba(),
	DB:     db.GetBizDB(),
}

func (s *_searchService) ExpireAllPostIndex(ctx context.Context) error {
	return s.
		Table("keyword_docs").
		WithContext(ctx).
		Exec("DELETE FROM keyword_docs").
		Error
}

// ReleaseExtension 销毁分词器(注意: 销毁之后无法继续继续使用, 仅且应当用于当服务需要关闭的资源释放方操作)
func (s *_searchService) ReleaseExtension() {
	s.cutter.Free()
}

// DeletePostSearchIndex 根据 ID 列表删除缓存
func (s *_searchService) DeletePostSearchIndex(postIDList []int64, ctx context.Context) error {
	return s.Table("keyword_docs").Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).
			Where( /* sql */ `post_id in ?`, postIDList).
			Delete(&model.Sqlite3KeywordDoc{}).
			Error
	})
}

func (s *_searchService) GetPostSearchIndexPages(req *dto.GetSearchIndexPagesReq, ctx context.Context) (resp *dto.GetSearchIndexPagesResp, err error) {
	var total int64
	// 计算记录数
	err = s.Table("keyword_docs").WithContext(ctx).Count(&total).Error
	if err != nil {
		err = util.CreateBizErr("查询文章索引记录数失败", err)
		return
	}
	docs := []*model.Sqlite3KeywordDoc{}
	offset, limit := req.Normalize()
	err = s.Table("keyword_docs").
		WithContext(ctx).
		Select( /* sql */ `post_id, weight, post_updated_at, record_created_at, record_updated_at`).
		Order( /* sql */ `record_updated_at desc`).
		Offset(offset).
		Limit(limit).
		Find(&docs).
		Error
	if err != nil {
		err = util.CreateBizErr("查询文章索引记录分页失败", err)
		return
	}

	resp = &dto.GetSearchIndexPagesResp{
		PageList: model.PageList[*model.Sqlite3KeywordDoc]{
			List:  docs,
			Page:  int64(*req.Page),
			Size:  int64(*req.Size),
			Total: total,
		},
	}
	return
}

func (s *_searchService) UpdatePostSearchIndex(post *model.Post) {
	const space = " "

	var titleBf bytes.Buffer
	for _, word := range s.cutter.CutForSearch(post.Title, true) {
		tWord := strings.TrimSpace(word)
		if len(tWord) > 1 {
			titleBf.WriteString(tWord)
			titleBf.WriteString(space)
		} else {
			if str.IsNotPunctuation(tWord) {
				titleBf.WriteString(tWord)
				titleBf.WriteString(space)
			}
		}
	}
	var contentBf bytes.Buffer
	for _, word := range s.cutter.CutForSearch(post.Content, true) {
		tWord := strings.TrimSpace(word)
		if len(tWord) > 1 {
			contentBf.WriteString(tWord)
			contentBf.WriteString(space)
		} else {
			if str.IsNotPunctuation(tWord) {
				contentBf.WriteString(tWord)
				contentBf.WriteString(space)
			}
		}
	}

	var postID int64
	ctx := context.Background()
	titleSp := strings.TrimSuffix(titleBf.String(), space)
	contentSp := strings.TrimPrefix(contentBf.String(), space)
	tmpQuery := &model.Sqlite3KeywordDoc{}

	s.Table("keyword_docs").Transaction(func(tx *gorm.DB) error {
		e := tx.
			Select("post_id").
			Where( /* sql */ `post_id = ?`, post.ID).
			Take(tmpQuery).Error

		now := time.Now().UnixMilli()
		// 创建新记录
		if e != nil {
			return tx.WithContext(ctx).
				Create(&model.Sqlite3KeywordDoc{
					PostID:          post.ID,
					TileSplit:       titleSp,
					ContentSplit:    contentSp,
					Weight:          util.TernaryExpr(ptr.IsNil(post.Weight), 0, *post.Weight),
					PostUpdatedAt:   post.UpdatedAt,
					RecordCreatedAt: now,
					RecordUpdatedAt: now,
				}).Error
		} else {
			// 仅更新记录
			postID = tmpQuery.PostID
			return tx.WithContext(ctx).
				Where( /* sql */ `post_id = ?`, postID).
				Updates(&model.Sqlite3KeywordDoc{
					PostID:          post.ID,
					TileSplit:       titleSp,
					ContentSplit:    contentSp,
					Weight:          util.TernaryExpr(ptr.IsNil(post.Weight), 0, *post.Weight),
					PostUpdatedAt:   post.UpdatedAt,
					RecordCreatedAt: tmpQuery.RecordCreatedAt,
					RecordUpdatedAt: now,
				}).Error

		}
	})
}

// SearchKeywordPages 全文搜索
func (s *_searchService) SearchKeywordPages(req *dto.GlobalSearchReq, ctx *gin.Context) (resp *dto.GlobalSearchResp, err error) {
	offset, limit := req.Normalize()
	keyword := strings.TrimSpace(req.Keyword)

	var total int64
	// 计算总条数
	err = s.Table("keyword_docs").
		WithContext(ctx).
		Model(&model.Sqlite3KeywordDoc{}).
		Where(`title_split MATCH ? or content_split MATCH ?`, keyword, keyword).
		Count(&total).
		Error
	if err != nil {
		err = util.CreateBizErr("查找数据失败", err)
		return
	}

	keywordMatches := []*model.Sqlite3KeywordDoc{}
	// 查找相关记录
	op := s.Table("keyword_docs").WithContext(ctx).
		Select(`post_id`). // 只选择相关 ID
		Where(`title_split MATCH ? or content_split MATCH ?`, keyword, keyword).
		Order("weight desc, post_updated_at desc").
		Offset(offset).
		Limit(limit).
		Find(&keywordMatches)
	err = op.Error
	if err != nil {
		err = util.CreateBizErr("查询数据失败", err)
		return
	}

	postOp := biz.Post
	// 回表寻找准确数据
	postMatches, err := postOp.WithContext(ctx).
		Where(postOp.ID.In(
			arr.MapSlice(keywordMatches, func(_ int, kw *model.Sqlite3KeywordDoc) int64 {
				return kw.PostID
			})...,
		)).
		Order(
			postOp.Weight.Desc(),
			postOp.CreatedAt.Desc(),
		).
		Find()
	if err != nil {
		err = util.CreateBizErr("查询数据数据失败", err)
		return
	}

	highlightList := arr.MapSlice(
		postMatches,
		func(_ int, post *model.Post) *dto.SearchHighlight {
			retVal := &dto.SearchHighlight{
				SubContent:               "",
				Title:                    post.Title,
				TitleHighLightIndex:      []int{},
				SubContentHighLightIndex: []int{},
				Category:                 post.Category,
				Tags:                     post.Tags,
				CreatedAt:                post.CreatedAt,
				PubAt:                    post.LastPubTime,
				Weight:                   post.Weight,
			}

			// 关键词在文章标题的高亮位置
			titleHlIdx := strings.Index(post.Title, keyword)
			// 转成 utf-8 位置
			if titleHlIdx != -1 {
				retVal.TitleHighLightIndex = []int{
					str.FindKeywordPositionRune(post.Title, keyword),
				}
			}

			contentKwIndex := strings.Index(post.Content, keyword)
			if contentKwIndex != -1 {
				sub, start := str.ExtractAroundKeyword(post.Content, keyword, 20)
				retVal.SubContent = sub
				retVal.SubContentHighLightIndex = []int{start}
			}

			return retVal
		},
	)

	resp = &dto.GlobalSearchResp{
		Keyword: keyword,
		List: model.PageList[*dto.SearchHighlight]{
			List:  highlightList,
			Page:  int64(*req.Page),
			Size:  int64(*req.Size),
			Total: total,
		},
	}
	return
}
