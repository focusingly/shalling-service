// 默认的 OSS3 文件直传服务
package service

type (
	_ossService struct{}
)

var DefaultOSSService = &_ossService{}

func (*_ossService) Upload() (resp string, err error) {
	return
}
