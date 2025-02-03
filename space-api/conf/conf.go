package conf

import (
	"bytes"
	_ "embed"

	"github.com/spf13/viper"
)

//go:embed config.yml
var confFile []byte

type (
	AppConf struct {
		Port       uint   `yaml:"port"`
		ServerHint string `yaml:"serverHint"`
	}

	Oauth2Conf struct {
		Endpoint     string   `yaml:"endPoint"`
		ClientId     string   `yaml:"clientId"`
		ClientSecret string   `yaml:"clientSecret"`
		Scopes       []string `yaml:"scopes"`
	}

	JwtConf struct {
		Salt    string `yaml:"salt"`
		Expired struct {
			Unit  string `yaml:"unit"`
			Setup uint64 `yaml:"setup"`
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
