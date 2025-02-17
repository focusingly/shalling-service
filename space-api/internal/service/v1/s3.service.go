// 默认的 OSS3 文件直传服务
package service

import (
	"context"
	"log"
	"space-api/conf"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type (
	_s3Service struct {
		*s3.Client
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

	DefaultS3Service = &_s3Service{
		s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(s3Conf.EndPoint)
		}),
	}

}

func (s *_s3Service) Upload() (resp string, err error) {

	return
}
