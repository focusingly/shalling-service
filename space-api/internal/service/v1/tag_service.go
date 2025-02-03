package service

import (
	"fmt"
	"slices"
	"space-api/dto"
	"space-api/util"
	"space-api/util/ptr"
	"space-domain/dao/biz"
	"space-domain/model"
	"strings"

	"github.com/gin-gonic/gin"
)

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

func GetTagDetailById(req *dto.GetTagDetailReq, ctx *gin.Context) (resp *dto.GetTagDetailResp, err error) {
	f, e := biz.Tag.WithContext(ctx).Where(biz.Tag.Id.Eq(req.Id)).Take()
	if e != nil {
		err = &util.BizErr{
			Msg:    "为找打数据",
			Reason: err,
		}
		return
	}
	resp = &dto.GetTagDetailResp{
		Tag: *f,
	}

	return
}

func UpdateOrCreateTag(req *dto.CreateOrUpdateTagReq, ctx *gin.Context) (resp *dto.CreateOrUpdateTagResp, err error) {
	var tagId int64 = 0

	err = biz.Q.Transaction(func(tx *biz.Query) error {
		tagOp := tx.Tag
		req.TagName = strings.TrimSpace(req.TagName)
		findTag, err := tagOp.WithContext(ctx).Where(tagOp.Id.Eq(req.Id)).Take()
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
				err = &util.BizErr{
					Msg:    "创建新的标签失败: " + e.Error(),
					Reason: e,
				}
				return err
			}
		} else {
			// 找到相关记录, 更新本地和其它相关表数据
			tagId = findTag.Id

			postTagRelations, e := tx.PostTagRelation.WithContext(ctx).Where(tx.PostTagRelation.TagId.Eq(tagId)).Find()
			if e != nil {
				return &util.BizErr{
					Msg:    "查找标签关联数据失败: " + e.Error(),
					Reason: err,
				}
			}
			findPostIdList := []int64{}
			for _, relation := range postTagRelations {
				findPostIdList = append(findPostIdList, relation.PostId)
			}
			// 所有和当前要变更的标签相关联的文章列表
			relativePosts, e := tx.Post.
				WithContext(ctx).
				Where(tx.Post.Id.In(findPostIdList...)).
				Find()
			if e != nil {
				return &util.BizErr{
					Msg:    "更新相关数据失败: " + e.Error(),
					Reason: e,
				}
			}
			for _, relativePost := range relativePosts {
				if relativePost.Tags != nil {
					splits := slices.
						DeleteFunc(
							strings.Split(*relativePost.Tags, ","), func(tag string) bool {
								// 排除掉旧标签
								return tag == findTag.TagName
							},
						)
					// 添加新的请求携带的标签
					splits = append(splits, req.TagName)
					// 拼接字符串覆盖掉旧标签
					relativePost.Tags = ptr.ToPtr(strings.Join(splits, ","))
				}
			}
			// 批量更新文章的 tags
			_, e = tx.Post.WithContext(ctx).Updates(relativePosts)
			if e != nil {
				return &util.BizErr{
					Reason: e,
					Msg:    "更新关联数据失败: " + e.Error(),
				}
			}

			// 最后更新标签自身
			newTagVal := &model.Tag{
				BaseColumn: findTag.BaseColumn,
				TagName:    req.TagName,
				Color:      req.Color,
				IconUrl:    req.IconUrl,
			}
			_, e = tx.Tag.WithContext(ctx).Updates(newTagVal)
			if e != nil {
				return &util.BizErr{
					Reason: err,
					Msg:    "更新标签数据失败: " + e.Error(),
				}
			}
		}

		return nil
	})

	if err != nil {
		return
	}

	// 获取当前操作的标签最新值
	val, err := biz.Tag.WithContext(ctx).Where(biz.Tag.Id.Eq(tagId)).Take()
	if err != nil {
		return
	}

	resp = &dto.CreateOrUpdateTagResp{
		Tag: *val,
	}

	return
}
