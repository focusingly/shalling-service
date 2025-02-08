package service

import (
	"space-api/dto"
	"space-api/util"
	"space-api/util/id"
	"space-domain/dao/biz"
	"space-domain/model"

	"github.com/gin-gonic/gin"
)

type mediaService struct{}

var DefaultMediaService *mediaService = &mediaService{}

func (*mediaService) CreateOrUpdateMediaTag(req *dto.CreateOrUpdateSocialMediaReq, ctx *gin.Context) (resp *dto.CreateOrUpdateSocialMediaResp, err error) {
	var mediaId int64 = 0

	err = biz.Q.Transaction(func(tx *biz.Query) error {
		mediaOp := tx.PubSocialMedia

		find, err := mediaOp.WithContext(ctx).
			Where(mediaOp.Id.Eq(req.Id)).
			Take()

		// 不存在, 那么进行创建
		if err != nil {
			// 赋予一个新的 ID
			mediaId = id.GetSnowFlakeNode().Generate().Int64()

			e := mediaOp.WithContext(ctx).Create(
				&model.PubSocialMedia{
					BaseColumn: model.BaseColumn{
						Id:   mediaId,
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
			mediaId = find.Id
			_, err = mediaOp.WithContext(ctx).Updates(
				&model.PubSocialMedia{
					BaseColumn: model.BaseColumn{
						Id:   mediaId,
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

	find, err := biz.PubSocialMedia.WithContext(ctx).Where(biz.PubSocialMedia.Id.Eq(mediaId)).Take()

	if err != nil {
		err = util.CreateBizErr("创建/更新公开媒体信息失败: "+err.Error(), err)
		return
	}
	resp = &dto.CreateOrUpdateSocialMediaResp{
		PubSocialMedia: *find,
	}

	return
}

func (*mediaService) GetMediaTagDetailById(req *dto.GetSocialMediaDetailReq, ctx *gin.Context) (resp *dto.GetSocialMediaDetailResp, err error) {
	find, err := biz.PubSocialMedia.
		WithContext(ctx).
		Where(biz.PubSocialMedia.Id.Eq(req.Id)).
		Take()

	if err != nil {
		err = util.CreateBizErr("查找记录失败: "+err.Error(), err)
		return
	}
	resp = &dto.GetSocialMediaDetailResp{
		PubSocialMedia: *find,
	}

	return
}

func (*mediaService) GetMediaTagPages(req *dto.GetSocialMediaPageListReq, ctx *gin.Context) (resp *dto.GetSocialMediaPageListResp, err error) {
	list, count, err := biz.PubSocialMedia.WithContext(ctx).FindByPage(req.BasePageParam.Normalize())
	if err != nil {
		err = util.CreateBizErr("查找分页失败: "+err.Error(), err)
		return
	}
	resp = &dto.GetSocialMediaPageListResp{
		PageList: model.PageList[*model.PubSocialMedia]{
			List:  list,
			Page:  int64(*req.Page),
			Size:  int64(*req.Size),
			Total: count,
		},
	}

	return
}

func (*mediaService) DeleteMediaTagByIdList(req *dto.DeleteSocialMediaByIdListReq, ctx *gin.Context) (resp *dto.DeleteSocialMediaByIdListResp, err error) {
	err = biz.Q.Transaction(func(tx *biz.Query) error {
		_, err := tx.PubSocialMedia.
			WithContext(ctx).
			Where(tx.PubSocialMedia.Id.In(req.IdList...)).
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
