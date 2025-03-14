package admin

import (
	"space-api/dto"
	"space-api/internal/service/v1"
	"space-api/middleware/outbound"
	"space-api/util"

	"github.com/gin-gonic/gin"
)

func UseS3Controller(group *gin.RouterGroup) {
	s3Service := service.DefaultS3Service
	s3Group := group.Group("/s3")

	// 获取 S3 文件直传链接
	s3Group.POST("/upload-link", func(ctx *gin.Context) {
		req := &dto.GetUploadObjectURLReq{}
		if err := ctx.ShouldBindBodyWithJSON(req); err != nil {
			ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
			return
		}

		if resp, err := s3Service.GetClientDirectUploadURL(req, ctx); err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(resp, ctx)
		}
	})

	// 同步信息到数据库当中
	s3Group.POST("/sync", func(ctx *gin.Context) {
		req := &dto.SyncS3RecordToDatabaseReq{}
		if err := ctx.ShouldBindBodyWithJSON(req); err != nil {
			ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
			return
		}

		if resp, err := s3Service.SyncToDatabase(req, ctx); err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(resp, ctx)
		}
	})

	// 展示数据库中的缓存信息
	s3Group.POST("/list", func(ctx *gin.Context) {
		req := &dto.GetS3ObjectPagesReq{}
		if err := ctx.ShouldBindBodyWithJSON(req); err != nil {
			ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
			return
		}
		if resp, err := s3Service.GetBucketDetailPages(req, ctx); err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(resp, ctx)
		}
	})

	// 删除对象
	s3Group.POST("/delete", func(ctx *gin.Context) {
		req := &dto.DeleteS3ObjectPagesReq{}
		if err := ctx.ShouldBindBodyWithJSON(req); err != nil {
			ctx.Error(util.CreateBizErr("参数错误: "+err.Error(), err))
			return
		}
		if resp, err := s3Service.DeleteS3Object(req, ctx); err != nil {
			ctx.Error(err)
		} else {
			outbound.NotifyProduceResponse(resp, ctx)
		}
	})

}
