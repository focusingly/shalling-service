package service

import (
	"space-api/dto"
	"space-api/util"
	"space-api/util/id"
	"space-domain/dao/biz"
	"space-domain/model"

	"github.com/gin-gonic/gin"
)

type (
	IMediaService interface {
		CreateOrUpdateMediaTag(req *dto.CreateOrUpdateSocialMediaReq, ctx *gin.Context) (resp *dto.CreateOrUpdateSocialMediaResp, err error)
		GetAnyMediaTags(req *dto.GetMediaTagsReq, ctx *gin.Context) (resp dto.GetMediaTagsResp, err error)
		GetVisibleMediaTags(req *dto.GetMediaTagsReq, ctx *gin.Context) (resp dto.GetMediaTagsResp, err error)
		DeleteMediaTagByIdList(req *dto.DeleteSocialMediaByIdListReq, ctx *gin.Context) (resp *dto.DeleteSocialMediaByIdListResp, err error)
	}
	mediaServiceImpl struct{}
)

var (
	_ IMediaService = (*mediaServiceImpl)(nil)

	DefaultMediaService IMediaService = &mediaServiceImpl{}
)

func (*mediaServiceImpl) CreateOrUpdateMediaTag(req *dto.CreateOrUpdateSocialMediaReq, ctx *gin.Context) (resp *dto.CreateOrUpdateSocialMediaResp, err error) {
	var mediaId int64 = 0

	err = biz.Q.Transaction(func(tx *biz.Query) error {
		mediaOp := tx.PubSocialMedia

		find, err := mediaOp.WithContext(ctx).
			Where(mediaOp.ID.Eq(req.Id)).
			Take()

		// 不存在, 那么进行创建
		if err != nil {
			// 赋予一个新的 ID
			mediaId = id.GetSnowFlakeNode().Generate().Int64()

			e := mediaOp.WithContext(ctx).Create(
				&model.PubSocialMedia{
					BaseColumn: model.BaseColumn{
						ID:   mediaId,
						Hide: req.Hide,
					},
					DisplayName: req.DisplayName,
					IconURL:     req.IconURL,
					OpenUrl:     req.OpenUrl,
				},
			)

			if e != nil {
				return e
			}
		} else {
			// 已经存在, 那么仅更新
			mediaId = find.ID
			_, err = mediaOp.WithContext(ctx).Updates(
				&model.PubSocialMedia{
					BaseColumn: model.BaseColumn{
						ID:   mediaId,
						Hide: req.Hide,
					},
					DisplayName: req.DisplayName,
					IconURL:     req.IconURL,
					OpenUrl:     req.OpenUrl,
				},
			)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		err = util.CreateBizErr("创建/更新公开媒体信息失败: "+err.Error(), err)
		return
	}

	find, err := biz.PubSocialMedia.WithContext(ctx).Where(biz.PubSocialMedia.ID.Eq(mediaId)).Take()

	if err != nil {
		err = util.CreateBizErr("创建/更新公开媒体信息失败: "+err.Error(), err)
		return
	}
	resp = &dto.CreateOrUpdateSocialMediaResp{
		PubSocialMedia: find,
	}

	return
}

func (*mediaServiceImpl) GetAnyMediaTags(req *dto.GetMediaTagsReq, ctx *gin.Context) (resp dto.GetMediaTagsResp, err error) {
	list, err := biz.PubSocialMedia.WithContext(ctx).
		Find()

	if err != nil {
		err = util.CreateBizErr("查找分页失败: "+err.Error(), err)
		return
	}
	resp = list

	return
}

// GetVisibleMediaTags 获取所有已经设置公开媒体展示标签
func (*mediaServiceImpl) GetVisibleMediaTags(req *dto.GetMediaTagsReq, ctx *gin.Context) (resp dto.GetMediaTagsResp, err error) {
	list, err := biz.PubSocialMedia.WithContext(ctx).
		Where(biz.PubSocialMedia.Hide.Eq(0)).
		Find()

	if err != nil {
		err = util.CreateBizErr("查找分页失败: "+err.Error(), err)
		return
	}

	resp = list
	return
}

func (*mediaServiceImpl) DeleteMediaTagByIdList(req *dto.DeleteSocialMediaByIdListReq, ctx *gin.Context) (resp *dto.DeleteSocialMediaByIdListResp, err error) {
	err = biz.Q.Transaction(func(tx *biz.Query) error {
		_, err := tx.PubSocialMedia.
			WithContext(ctx).
			Where(tx.PubSocialMedia.ID.In(req.IdList...)).
			Delete()
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		err = util.CreateBizErr("删除公开媒体信息: "+err.Error(), err)
	}
	resp = &dto.DeleteSocialMediaByIdListResp{}

	return
}
