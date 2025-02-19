package dto

type (
	SendMailReq struct {
		From     string            `json:"from" yaml:"from" xml:"from" toml:"from"`             // 邮件发送者的地址
		To       []string          `json:"to" yaml:"to" xml:"to" toml:"to"`                     // 邮件接收者地址列表
		ReplyTo  string            `json:"replyTo" yaml:"replyTo" xml:"replyTo" toml:"replyTo"` // 提供给对方的邮件回复地址
		Tags     []string          `json:"tags" yaml:"tags" xml:"tags" toml:"tags"`             // 标签
		Subject  string            `json:"subject" yaml:"subject" xml:"subject" toml:"subject"`
		Body     string            `json:"body" yaml:"body" xml:"body" toml:"body"` // 邮件正文内容, 可以是 html 有可以是纯文本内容
		BodyType string            `json:"bodyType" yaml:"bodyType" xml:"bodyType" toml:"bodyType"`
		Headers  map[string]string `json:"headers" yaml:"headers" xml:"headers" toml:"headers"` // 自定义头部
	}
	SendMailWithSelectionReq struct {
		SpecificID string      `json:"specificID" yaml:"specificID" xml:"specificID" toml:"specificID"`
		Content    SendMailReq `json:"content" yaml:"content" xml:"content" toml:"content"`
	}
	SendMailResp struct {
	}

	MailConf struct {
		Host        string `json:"host" yaml:"host" xml:"host" toml:"host"`
		Port        int    `json:"port" yaml:"port" xml:"port" toml:"port"`
		Account     string `json:"account" yaml:"account" xml:"account" toml:"account"`
		Primary     bool   `json:"primary" yaml:"primary" xml:"primary" toml:"primary"` // 是否被标记为首选邮箱
		Mark        string `json:"mark" yaml:"mark" xml:"mark" toml:"mark"`
		DefaultFrom string `json:"defaultFrom" yaml:"defaultFrom" xml:"defaultFrom" toml:"defaultFrom"`
		SpecificID  string `json:"specificID" yaml:"specificID" xml:"specificID" toml:"specificID"` // 配置项标识的唯一 ID
	}
	GetMailConfListReq  struct{}
	GetMailConfListResp struct {
		List []*MailConf `json:"list" yaml:"list" xml:"list" toml:"list"`
	}
)
