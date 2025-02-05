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
		Port       uint   `yaml:"port"`
		ServerHint string `yaml:"serverHint"`
	}

	// Oauth2 认证配置
	Oauth2Conf struct {
		Endpoint     string   `yaml:"endPoint"`
		ClientId     string   `yaml:"clientId"`
		ClientSecret string   `yaml:"clientSecret"`
		Scopes       []string `yaml:"scopes"`
	}

	// Jwt 配置
	JwtConf struct {
		Salt    string `yaml:"salt"`
		Expired struct {
			Unit  string `yaml:"unit"`
			Value int64  `yaml:"value"`
		} `yaml:"expired"`
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
