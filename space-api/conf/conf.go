package conf

import (
	"bytes"
	_ "embed"

	"github.com/spf13/viper"
)

//go:embed config.yml
var confFile []byte

type (
	// AppConf 应用程序基本配置
	AppConf struct {
		Port           uint   `yaml:"port" json:"port"`
		ServerHint     string `yaml:"serverHint" json:"serverHint"`
		Salt           string `yaml:"salt" json:"salt"`
		MaxUserActive  int    `yaml:"maxUserActive" json:"maxUserActive"`
		ServerTimezone string `yaml:"serverTimezone" json:"serverTimezone"`
	}

	// Oauth2 认证配置
	Oauth2Conf struct {
		Endpoint     string   `yaml:"endPoint" json:"endpoint"`
		ClientId     string   `yaml:"clientId" json:"clientId"`
		ClientSecret string   `yaml:"clientSecret" json:"clientSecret"`
		Scopes       []string `yaml:"scopes" json:"scopes"`
	}

	// Jwt 配置
	JwtConf struct {
		Salt    string `yaml:"salt" json:"salt"`
		Expired struct {
			Unit  string `yaml:"unit" json:"unit"`
			Value int64  `yaml:"value" json:"value"`
		} `yaml:"expired" json:"expired"`
	}
)

var globalViper *viper.Viper

func init() {
	v := viper.New()
	v.SetConfigType("yaml")
	if err := v.ReadConfig(bytes.NewReader(confFile)); err != nil {
		panic(err)
	}
	globalViper = v
}

func GetProjectViper() *viper.Viper {
	if globalViper == nil {
		panic("viper not init")
	}
	return globalViper
}
