package service

import (
	"space-api/conf"
	"space-api/util"

	"gopkg.in/gomail.v2"
)

type _mailService struct {
	dialer *gomail.Dialer
}

var DefaultMailService *_mailService

func init() {
	mailConf := conf.ProjectConf.GetMailConf()
	if mailConf != nil {
		DefaultMailService = &_mailService{
			dialer: gomail.NewDialer(
				mailConf.Host,
				mailConf.Port,
				mailConf.Username,
				mailConf.Password,
			),
		}
	}
}

func (s *_mailService) SendEmail(msg *gomail.Message) error {
	err := s.dialer.DialAndSend(msg)

	if err != nil {
		return util.CreateBizErr("发送邮件失败: "+err.Error(), err)
	}

	return nil
}
