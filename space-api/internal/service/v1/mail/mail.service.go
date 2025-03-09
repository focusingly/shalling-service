package mail

import (
	"fmt"
	"log"
	"space-api/conf"
	"space-api/dto"
	"space-api/util"

	"gopkg.in/gomail.v2"
)

type (
	IMailService interface {
		SendEmailByPrimary(req *dto.SendMailReq) (resp *dto.SendMailResp, err error)
		SendEmailBySelection(req *dto.SendMailWithSelectionReq) (resp *dto.SendMailResp, err error)
		GetConfList(req *dto.GetMailConfListReq) (resp *dto.GetMailConfListResp)
	}

	mailServiceImpl struct {
		primaryDialer *gomail.Dialer
		mapping       map[string]*confDialerMapping
	}
	confDialerMapping struct {
		conf   *conf.MailSmtpConf
		dialer *gomail.Dialer
	}
)

var (
	_                  IMailService = (*mailServiceImpl)(nil)
	DefaultMailService IMailService = newMailService()
)

func newMailService() IMailService {
	tmp := &mailServiceImpl{
		mapping: make(map[string]*confDialerMapping),
	}

	list := conf.ProjectConf.GetMailConfList()
	for _, mailConf := range list {
		t := &confDialerMapping{
			dialer: gomail.NewDialer(
				mailConf.Host,
				mailConf.Port,
				mailConf.Account,
				mailConf.Credential,
			),
			conf: mailConf,
		}
		tmp.mapping[mailConf.SpecificID] = t
		if mailConf.Primary {
			tmp.primaryDialer = t.dialer
		}
	}
	if tmp.primaryDialer == nil {
		log.Fatal("未设置主要的邮件发送配置参数")
	}

	return tmp
}

// SendEmailByPrimary 使用主邮件配置发送邮件
func (s *mailServiceImpl) SendEmailByPrimary(req *dto.SendMailReq) (resp *dto.SendMailResp, err error) {
	msg := gomail.NewMessage()
	msg.SetHeader("From", req.From)
	msg.SetHeader("To", req.To...)
	msg.SetHeader("Reply-To", req.ReplyTo)
	msg.SetHeader("Subject", req.Subject)
	msg.SetBody(req.BodyType, req.Body)
	if req.Tags != nil {
		msg.SetHeader("Tags", req.Tags...)
	}
	if req.Headers != nil {
		for h, v := range req.Headers {
			msg.SetHeader(h, v)
		}
	}
	err = s.primaryDialer.DialAndSend(msg)
	if err != nil {
		err = util.CreateBizErr("发送邮件失败: "+err.Error(), err)

		return
	}

	resp = &dto.SendMailResp{}
	return
}

// SendEmailByPrimary 使用主邮件配置发送邮件
func (s *mailServiceImpl) SendEmailBySelection(req *dto.SendMailWithSelectionReq) (resp *dto.SendMailResp, err error) {
	c, ok := s.mapping[req.SpecificID]
	if !ok {
		err = util.CreateBizErr("未找到匹配的邮件发送配置", fmt.Errorf("not found matched mail config by id: %s", req.SpecificID))
		return
	}

	content := req.Content
	msg := gomail.NewMessage()
	msg.SetHeader("From", content.From)
	msg.SetHeader("To", content.To...)
	msg.SetHeader("Reply-To", content.ReplyTo)
	msg.SetHeader("Subject", content.Subject)
	msg.SetBody(content.BodyType, content.Body)
	if content.Tags != nil {
		msg.SetHeader("Tags", content.Tags...)
	}
	if content.Headers != nil {
		for h, v := range content.Headers {
			msg.SetHeader(h, v)
		}
	}

	// 使用客户端请求使用的邮件发送者
	err = c.dialer.DialAndSend(msg)
	if err != nil {
		err = util.CreateBizErr("发送邮件失败: "+err.Error(), err)

		return
	}

	resp = &dto.SendMailResp{}
	return
}

func (s *mailServiceImpl) GetConfList(req *dto.GetMailConfListReq) (resp *dto.GetMailConfListResp) {
	list := []*dto.MailConf{}
	for _, m := range s.mapping {
		list = append(list, &dto.MailConf{
			Host:        m.conf.Host,
			Port:        m.conf.Port,
			Account:     m.conf.Account,
			Primary:     m.conf.Primary,
			Mark:        m.conf.Mark,
			DefaultFrom: m.conf.DefaultFrom,
			SpecificID:  m.conf.SpecificID,
		})
	}

	resp = &dto.GetMailConfListResp{
		List: list,
	}

	return
}
