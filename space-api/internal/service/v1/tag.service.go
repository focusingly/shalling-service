package service

import (
	"fmt"
	"space-api/dto"
	"space-api/util"
	"space-api/util/arr"
	"space-domain/dao/biz"
	"space-domain/model"
	"strings"

	"github.com/gin-gonic/gin"
)

// CreateOrUpdateTag 创建/更新 标签
func CreateOrUpdateTag(req *dto.CreateOrUpdateTagReq, ctx *gin.Context) (resp *dto.CreateOrUpdateTagResp, err error) {
	// tag 的 ID
	var tagId int64 = 0

	err = biz.Q.Transaction(func(tx *biz.Query) error {
		tagOp := tx.Tag
		req.TagName = strings.TrimSpace(req.TagName)

		// 找到需要被更新的 tag
		oldTag, err := tagOp.WithContext(ctx).Where(tagOp.Id.Eq(req.Id)).Take()

		// 未找到相关的的记录, 插入新的记录
		if err != nil {
			// 更新 id 值
			tagId = util.GetSnowFlakeNode().Generate().Int64()
			e := tagOp.WithContext(ctx).Create(&model.Tag{
				BaseColumn: model.BaseColumn{Id: tagId},
				TagName:    req.TagName,
				Color:      req.Color,
				IconUrl:    req.IconUrl,
			})
			if e != nil {
				return err
			}
		} else {
			// 找到相关记录, 更新本地和其它相关表数据
			tagId = oldTag.Id

			// 找到文章-标签 关联记录
			postTagRelations, e := tx.PostTagRelation.
				WithContext(ctx).
				Where(tx.PostTagRelation.TagId.Eq(tagId)).
				Find()
			if e != nil {
				return e
			}

			// 需要进行更新的文章 ID 列表
			relativePostIdList := arr.MapSlice(
				postTagRelations,
				func(_ int, relation *model.PostTagRelation) int64 {
					return relation.PostId
				},
			)

			// 找到所有要更新的文章
			shouldUpdatePosts, e := tx.Post.
				WithContext(ctx).
				Where(tx.Post.Id.In(relativePostIdList...)).
				Find()
			if e != nil {
				return e
			}

			// 更新 文章-标签 之间的关系
			for _, tmpPost := range shouldUpdatePosts {
				if tmpPost.Tags != nil {
					// 查找就地更新掉旧的 Tag
					for index, tg := range tmpPost.Tags {
						if tg == oldTag.TagName {
							tmpPost.Tags[index] = req.TagName
						}
					}
				}
			}

			// 更新文章自身的 tags 字段
			for _, newPost := range shouldUpdatePosts {
				_, e = tx.Post.WithContext(ctx).Updates(newPost)
				if e != nil {
					return e
				}
			}

			// 最后更新标签自身
			newTagVal := &model.Tag{
				BaseColumn: model.BaseColumn{
					Id:   oldTag.Id,
					Hide: req.Hide,
				},
				TagName: req.TagName,
				Color:   req.Color,
				IconUrl: req.IconUrl,
			}
			_, e = tx.Tag.WithContext(ctx).Select(
				tx.Tag.Id,
				tx.Tag.Hide,
				tx.Tag.TagName,
				tx.Tag.Color,
				tx.Tag.IconUrl,
			).Updates(newTagVal)

			if e != nil {
				return e
			}
		}

		return nil
	})

	if err != nil {
		err = util.CreateBizErr("更新失败: "+err.Error(), err)
		return
	}

	// 获取当前操作的标签最新值
	val, err := biz.Tag.WithContext(ctx).Where(biz.Tag.Id.Eq(tagId)).Take()
	if err != nil {
		err = util.CreateBizErr("更新失败: "+err.Error(), err)
		return
	}

	resp = &dto.CreateOrUpdateTagResp{
		Tag: *val,
	}

	return
}

func GetTagDetailById(req *dto.GetTagDetailReq, ctx *gin.Context) (resp *dto.GetTagDetailResp, err error) {
	f, e := biz.Tag.WithContext(ctx).Where(biz.Tag.Id.Eq(req.Id)).Take()
	if e != nil {
		err = &util.BizErr{
			Msg:    "未找到数据",
			Reason: err,
		}
		return
	}
	resp = &dto.GetTagDetailResp{
		Tag: *f,
	}

	return
}

func GetTagPageList(req *dto.GetTagPageListReq, ctx *gin.Context) (resp *dto.GetTagPageListResp, err error) {
	tagOp := biz.Tag
	list, count, err := tagOp.WithContext(ctx).FindByPage(req.Resolve())
	if err != nil {
		return nil, &util.BizErr{
			Reason: err,
			Msg:    "查询失败" + err.Error(),
		}
	}

	return &dto.GetTagPageListResp{
		PageList: model.PageList[*model.Tag]{
			List:  list,
			Page:  int64(*req.Page),
			Size:  int64(*req.Size),
			Total: count,
		},
	}, nil
}

func DeleteTagByIdList(req *dto.DeleteTagByIdListReq, ctx *gin.Context) (resp *dto.DeleteTagByIdListResp, err error) {
	query := biz.Q
	err = query.Transaction(func(tx *biz.Query) error {
		tagOp := tx.Tag

		list, e := tx.PostTagRelation.WithContext(ctx).Where(tx.PostTagRelation.TagId.In(req.IdList...)).Select(tx.PostTagRelation.Id).Find()
		if err != nil {
			err = &util.BizErr{
				Reason: err,
				Msg:    "查询关联数据失败: " + e.Error(),
			}
		}
		if len(list) != 0 && !req.ForceOverride {
			return &util.BizErr{
				Msg:    fmt.Sprintf("存在相关联的数据: %d 条", len(list)),
				Reason: fmt.Errorf("%v", list),
			}
		}

		_, e = tx.Tag.WithContext(ctx).Where(tagOp.Id.In(req.IdList...)).Delete()
		if err != nil {
			err = &util.BizErr{
				Reason: err,
				Msg:    "删除错误: " + e.Error(),
			}
			return err
		}
		return nil
	})

	if err != nil {
		return
	}

	return &dto.DeleteTagByIdListResp{}, nil
}
