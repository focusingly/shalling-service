// 默认的 OSS3 文件直传服务
package service

import (
	"context"
	"fmt"
	"log"
	"path"
	"space-api/conf"
	"space-api/dto"
	"space-api/util"
	"space-api/util/performance"
	"space-domain/dao/biz"
	"space-domain/model"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type (
	_s3Service struct {
		*s3.Client
		presignClient *s3.PresignClient
		s3Conf        *conf.S3Conf
	}
)

var DefaultS3Service *_s3Service

func init() {
	s3Conf := conf.ProjectConf.GetS3Conf()
	if s3Conf == nil {
		return
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				s3Conf.AccessKeyID,
				s3Conf.AccessKeySecret,
				"",
			),
		),
		config.WithRegion("auto"),
	)

	if err != nil {
		log.Fatal("init s3 config error: ", err)
	}
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(s3Conf.EndPoint)
	})
	DefaultS3Service = &_s3Service{
		Client:        client,
		presignClient: s3.NewPresignClient(client),
		s3Conf:        s3Conf,
	}
}

func (s *_s3Service) CreateBucketByName(bucketName string, ctx context.Context) (err error) {
	_, err = s.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		err = util.CreateBizErr("创建存储桶失败", err)
		return
	}

	return
}

func (s *_s3Service) EnsureBucketExists(bucketName string, ctx context.Context) error {
	s3Op := biz.S3ObjectRecord
	_, err := s3Op.WithContext(ctx).
		Where(s3Op.BucketName.Eq(bucketName)).
		Take()
	if err != nil {
		_, e := s.HeadBucket(ctx, &s3.HeadBucketInput{
			Bucket: aws.String(bucketName),
		})
		if e != nil {
			return s.CreateBucketByName(bucketName, ctx)
		}
	}
	return nil
}

func (s *_s3Service) SyncToDatabase(req *dto.SyncS3RecordToDatabaseReq, ctx context.Context) (resp *dto.SyncToDatabaseResp, err error) {
	objInfo, err := s.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(req.BucketName),
		Key:    aws.String(req.ObjectKey),
	})
	if err != nil {
		err = util.CreateBizErr("对象不存在: "+err.Error(), err)
		return
	}

	err = biz.Q.Transaction(func(tx *biz.Query) error {
		s3Tx := tx.S3ObjectRecord

		condList := []gen.Condition{
			s3Tx.BucketName.Eq(req.BucketName),
			s3Tx.ObjectKey.Eq(req.ObjectKey),
		}

		if objInfo.ChecksumSHA256 != nil {
			condList = append(condList, s3Tx.Checksum.Eq(*objInfo.ChecksumSHA256))
		}

		record, e := s3Tx.WithContext(ctx).
			Where(condList...).
			Take()

		domain := s.s3Conf.LinkedDomain
		var visitURL = util.TernaryExpr(
			domain != "",
			fmt.Sprintf("https://%s/%s", domain, record.ObjectKey),
			fmt.Sprintf("%s/%s", s.s3Conf.EndPoint, record.ObjectKey),
		)
		// 不存在同步的数据
		if e != nil {
			e = s3Tx.WithContext(ctx).
				Create(&model.S3ObjectRecord{
					ObjectKey:    req.ObjectKey,
					FileName:     req.OriginFileName,
					Extension:    path.Ext(record.ObjectKey),
					FileSize:     *objInfo.ContentLength,
					BucketName:   req.BucketName,
					ChecksumType: "SHA256",
					Checksum:     req.Checksum,
					PubAvailable: 1,
				})
			if e != nil {
				return e
			}
			resp = &dto.SyncToDatabaseResp{
				FileSize:     *objInfo.ContentLength,
				LinkedDomain: s.s3Conf.LinkedDomain,
				ObjectKey:    req.ObjectKey,
				FullVisitURL: visitURL,
			}
		} else {
			resp = &dto.SyncToDatabaseResp{
				FileSize:     record.FileSize,
				LinkedDomain: s.s3Conf.LinkedDomain,
				ObjectKey:    record.ObjectKey,
				FullVisitURL: visitURL,
			}
		}
		return nil
	})

	if err != nil {
		err = util.CreateBizErr("同步到数据库失败", err)
		return
	}

	return

}

func (s *_s3Service) GetBucketDetailPages(req *dto.GetS3ObjectPagesReq, ctx context.Context) (resp *dto.GetS3ObjectPagesResp, err error) {
	s3Op := biz.S3ObjectRecord
	condList := []gen.Condition{}
	tableName := s3Op.TableName()

	if req.CondList != nil {
		for _, cond := range req.CondList {
			if expr, err := cond.ParseCond(tableName); err != nil {
				err = util.CreateBizErr("查询参数错误: "+err.Error(), err)
				return nil, err
			} else {
				condList = append(condList, expr)
			}
		}
	}
	orderList := []field.Expr{}
	if req.OrderList != nil {
		for _, o := range req.OrderList {
			orderList = append(orderList, o.ToOrderField(tableName))
		}
	}

	list, count, err := s3Op.WithContext(ctx).
		Where(condList...).
		Order(orderList...).
		FindByPage(req.Normalize())

	if err != nil {
		err = util.CreateBizErr("查询失败", err)
		return
	}

	resp = &dto.GetS3ObjectPagesResp{
		PageList: model.PageList[*model.S3ObjectRecord]{
			List:  list,
			Page:  int64(*req.Page),
			Size:  int64(*req.Size),
			Total: count,
		},
	}
	return
}

func (s *_s3Service) DeleteS3Object(req *dto.DeleteS3ObjectPagesReq, ctx context.Context) (resp *dto.DeleteS3ObjectPagesResp, err error) {
	removeList := []*model.S3ObjectRecord{}

	err = biz.Q.Transaction(func(tx *biz.Query) error {
		s3Tx := tx.S3ObjectRecord
		tableName := s3Tx.TableName()

		condList := []gen.Condition{}
		if req.CondList != nil {
			for _, cond := range req.CondList {
				if p, e := cond.ParseCond(tableName); e != nil {
					return e
				} else {
					condList = append(condList, p)
				}
			}
		}

		l, e := s3Tx.WithContext(ctx).
			Select(s3Tx.ObjectKey, s3Tx.BucketName).
			Where(condList...).
			Find()
		if e != nil {
			return e
		}
		removeList = l

		return nil
	})

	if err != nil {
		err = util.CreateBizErr("操作失败: "+err.Error(), err)
		return
	}

	// 异步通知删除
	performance.DefaultTaskRunner.Go(func() {
		for _, r := range removeList {
			s.Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
				Bucket: aws.String(r.BucketName),
				Key:    aws.String(r.ObjectKey),
			})
		}

	})

	resp = &dto.DeleteS3ObjectPagesResp{}
	return
}

// GetClientDirectUploadURL 获取客户端链的直传链接
func (s *_s3Service) GetClientDirectUploadURL(req *dto.GetUploadObjectURLReq, ctx context.Context) (resp *dto.GetUploadObjectURLResp, err error) {
	exp := time.Minute * 3

	err = s.EnsureBucketExists(req.Bucket, ctx)
	if err != nil {
		err = util.CreateBizErr("创建存储桶失败", err)
		return
	}

	newObjectKey := strings.ReplaceAll(uuid.NewString(), "-", "") + path.Ext(req.Filename)
	httpRequest, err := s.presignClient.PresignPutObject(
		ctx,
		&s3.PutObjectInput{
			Bucket:            aws.String(req.Bucket),
			Key:               aws.String(newObjectKey),
			ChecksumAlgorithm: types.ChecksumAlgorithmSha256,
			ChecksumSHA256:    req.SHA256,
			Metadata:          aws.ToStringMap(req.Metadata),
		},
		s3.WithPresignExpires(exp),
	)

	if err != nil {
		err = util.CreateBizErr("创建上传链接失败", err)
		return
	}

	resp = &dto.GetUploadObjectURLResp{
		ObjectKey:      newObjectKey,
		OriginFileName: req.Filename,
		Method:         httpRequest.Method,
		URL:            httpRequest.URL,
		ChecksumType:   "SHA256",
		Checksum:       util.TernaryExpr(req.SHA256 != nil, *req.SHA256, ""),
		ExpiredAt:      time.Now().Add(exp).UnixMilli(),
	}

	return
}
