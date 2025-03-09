package service

import (
	"bytes"
	"context"
	"io"
	"space-api/db"
	"space-api/dto"
	"space-api/util"
	"space-api/util/arr"
	"space-api/util/str"
	"space-domain/dao/biz"
	"space-domain/model"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yanyiwu/gojieba"
	"gorm.io/gorm"
)

type (
	ISearchService interface {
		io.Closer
		ExpireAllPostIndex(ctx context.Context) error
		DeletePostSearchIndex(postIDList []int64, ctx context.Context) error
		GetPostSearchIndexPages(req *dto.GetSearchIndexPagesReq, ctx context.Context) (resp *dto.GetSearchIndexPagesResp, err error)
		UpdatePostSearchIndex(post *model.Post)
		SearchKeywordPages(req *dto.GlobalSearchReq, ctx *gin.Context) (resp *dto.GlobalSearchResp, err error)
	}

	searchServiceImpl struct {
		cutter *gojieba.Jieba
		*gorm.DB
	}
)

var (
	_ ISearchService = (*searchServiceImpl)(nil)

	DefaultGlobalSearchService = &searchServiceImpl{
		cutter: gojieba.NewJieba(),
		DB:     db.GetBizDB(),
	}
)

// Free 销毁分词器(注意: 销毁之后无法继续继续使用, 仅且应当用于当服务需要关闭的资源释放方操作)
func (s *searchServiceImpl) Close() error {
	s.cutter.Free()
	return nil
}

func (s *searchServiceImpl) ExpireAllPostIndex(ctx context.Context) error {
	return s.
		Table("keyword_docs").
		WithContext(ctx).
		Exec("DELETE FROM keyword_docs").
		Error
}

// DeletePostSearchIndex 根据 ID 列表删除缓存
func (s *searchServiceImpl) DeletePostSearchIndex(postIDList []int64, ctx context.Context) error {
	return biz.Q.Transaction(func(tx *biz.Query) error {
		docTx := tx.Sqlite3KeywordDoc
		_, e := docTx.WithContext(ctx).Where(docTx.PostID.In(postIDList...)).Delete()
		return e
	})
}

func (s *searchServiceImpl) GetPostSearchIndexPages(req *dto.GetSearchIndexPagesReq, ctx context.Context) (resp *dto.GetSearchIndexPagesResp, err error) {
	docOP := biz.Sqlite3KeywordDoc

	docs, count, err := docOP.WithContext(ctx).
		Select(
			docOP.PostID,
			docOP.Weight,
			docOP.PostUpdatedAt,
			docOP.RecordCreatedAt,
			docOP.RecordUpdatedAt,
		).
		FindByPage(req.Normalize())
	if err != nil {
		err = util.CreateBizErr("查询文章分词分页失败", err)
		return
	}

	resp = &dto.GetSearchIndexPagesResp{
		PageList: model.PageList[*model.Sqlite3KeywordDoc]{
			List:  docs,
			Page:  int64(*req.Page),
			Size:  int64(*req.Size),
			Total: count,
		},
	}
	return
}

func (s *searchServiceImpl) UpdatePostSearchIndex(post *model.Post) {
	const space = " "

	var titleBf bytes.Buffer
	titleCuts := s.cutter.CutForSearch(post.Title, true)
	titleCutsLen := len(titleCuts)
	for index, word := range titleCuts {
		tWord := strings.TrimSpace(word)
		if len(tWord) > 1 {
			titleBf.WriteString(tWord)
			if index != titleCutsLen-1 {
				titleBf.WriteString(space)
			}
		} else {
			if str.IsNotPunctuation(tWord) {
				titleBf.WriteString(tWord)
				if index != titleCutsLen-1 {
					titleBf.WriteString(space)
				}
			}
		}
	}

	var contentBf bytes.Buffer
	contentCuts := s.cutter.CutForSearch(post.Content, true)
	contentCutsLen := len(contentCuts)
	for index, word := range contentCuts {
		tWord := strings.TrimSpace(word)
		if len(tWord) > 1 {
			contentBf.WriteString(tWord)
			if index != contentCutsLen-1 {
				contentBf.WriteString(space)
			}
		} else {
			if str.IsNotPunctuation(tWord) {
				contentBf.WriteString(tWord)
				if index != contentCutsLen-1 {
					contentBf.WriteString(space)
				}
			}
		}
	}

	var postID int64
	ctx := context.TODO()

	biz.Q.Transaction(func(tx *biz.Query) error {
		docTx := tx.Sqlite3KeywordDoc
		tmpQuery, e := docTx.WithContext(ctx).
			Select(docTx.PostID).
			Where(docTx.PostID.Eq(post.ID)).
			Take()

		// 创建新记录
		if e != nil {
			return docTx.WithContext(ctx).
				Create(&model.Sqlite3KeywordDoc{
					PostID:       post.ID,
					TileSplit:    titleBf.String(),
					ContentSplit: contentBf.String(),
					Weight: util.TernaryExprWithProducer(
						post.Weight != nil,
						func() int {
							return *post.Weight
						},
						func() int {
							return 0
						},
					),
					PostUpdatedAt: post.UpdatedAt,
				})
		} else {
			// 仅更新记录
			postID = tmpQuery.PostID
			_, e := docTx.WithContext(ctx).
				Where(docTx.PostID.Eq(postID)).
				Updates(&model.Sqlite3KeywordDoc{
					PostID:       post.ID,
					TileSplit:    titleBf.String(),
					ContentSplit: contentBf.String(),
					Weight: util.TernaryExprWithProducer(
						post.Weight != nil,
						func() int {
							return *post.Weight
						},
						func() int {
							return 0
						},
					),
					PostUpdatedAt: post.UpdatedAt,
				})

			return e
		}

	})
}

// SearchKeywordPages 全文搜索
func (s *searchServiceImpl) SearchKeywordPages(req *dto.GlobalSearchReq, ctx *gin.Context) (resp *dto.GlobalSearchResp, err error) {
	offset, limit := req.Normalize()
	keyword := strings.TrimSpace(req.Keyword)
	docOP := biz.Sqlite3KeywordDoc

	var total int64
	// 计算总条数
	err = docOP.WithContext(ctx).
		UnderlyingDB().
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
	op := docOP.WithContext(ctx).
		UnderlyingDB().
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
