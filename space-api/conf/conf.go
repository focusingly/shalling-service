package conf

import (
	"fmt"
	"log"
	"os"
	"path"
	"slices"
	"space-api/constants"
	"space-api/util"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type (
	// AppConf 应用程序基本配置
	AppConf struct {
		Port             uint                     `yaml:"port" json:"port" xml:"port" toml:"port"`
		ServerHint       string                   `yaml:"serverHint" json:"serverHint" xml:"serverHint" toml:"serverHint"`
		Salt             string                   `yaml:"salt" json:"salt" xml:"salt" toml:"salt"`
		MaxUserActive    int                      `yaml:"maxUserActive" json:"maxUserActive" xml:"maxUserActive" toml:"maxUserActive"`
		ServerTimezone   string                   `yaml:"serverTimezone" json:"serverTimezone" xml:"serverTimezone" toml:"serverTimezone"`
		NodeID           int64                    `json:"nodeID" yaml:"nodeID" xml:"nodeID" toml:"nodeID"`
		StaticDir        string                   `json:"staticDir" yaml:"staticDir" xml:"staticDir" toml:"staticDir"`
		GlobalUploadSize string                   `json:"globalUploadSize" yaml:"globalUploadSize" xml:"globalUploadSize" toml:"globalUploadSize"`
		ParsedUploadSize constants.MemoryByteSize `json:"parsedUploadSize" yaml:"parsedUploadSize" xml:"parsedUploadSize" toml:"parsedUploadSize"`
	}

	MailConf struct {
		Host     string `json:"host" yaml:"host" xml:"host" toml:"host"`
		Username string `json:"username" yaml:"username" xml:"username" toml:"username"`
		Password string `json:"password" yaml:"password" xml:"password" toml:"password"`
		Port     int    `json:"port" yaml:"port" xml:"port" toml:"port"`
	}

	// Oauth2 认证配置
	Oauth2Conf struct {
		EndPoint     string   `json:"endPoint" yaml:"endPoint" xml:"endPoint" toml:"endPoint"`
		ClientId     string   `yaml:"clientId" json:"clientId" xml:"clientId" toml:"clientId"`
		ClientSecret string   `yaml:"clientSecret" json:"clientSecret" xml:"clientSecret" toml:"clientSecret"`
		Scopes       []string `yaml:"scopes" json:"scopes" xml:"scopes" toml:"scopes"`
	}

	// Jwt 配置
	JwtConf struct {
		Salt          string        `yaml:"salt" json:"salt" xml:"salt" toml:"salt"`
		Expired       string        `json:"expired" yaml:"expired" xml:"expired" toml:"expired"`
		ParsedExpTime time.Duration `json:"parsedExpTime" yaml:"parsedExpTime" xml:"parsedExpTime" toml:"parsedExpTime"`
	}

	// 数据库配置
	DatabaseConf struct {
		DBName string `json:"dbName" yaml:"dbName" xml:"dbName" toml:"dbName"`
		DBType string `json:"dbType" yaml:"dbType" xml:"dbType" toml:"dbType"`
		Dsn    string `json:"dsn" yaml:"dsn" xml:"dsn" toml:"dsn"`
		Mark   string `json:"mark" yaml:"mark" xml:"mark" toml:"mark"`
	}

	CloudflareConf struct {
		AccountID string `json:"accountID" yaml:"accountID" xml:"accountID" toml:"accountID"`
		ApiKey    string `json:"apiKey" yaml:"apiKey" xml:"apiKey" toml:"apiKey"`
		Email     string `json:"email" yaml:"email" xml:"email" toml:"email"`
	}

	S3Conf struct {
		AccountID       string `json:"accountID" yaml:"accountID" xml:"accountID" toml:"accountID"`
		AccessKeyID     string `json:"accessKeyID" yaml:"accessKeyID" xml:"accessKeyID" toml:"accessKeyID"`
		AccessKeySecret string `json:"accessKeySecret" yaml:"accessKeySecret" xml:"accessKeySecret" toml:"accessKeySecret"`
		Token           string `json:"token" yaml:"token" xml:"token" toml:"token"`
		BucketName      string `json:"bucketName" yaml:"bucketName" xml:"bucketName" toml:"bucketName"`
		EndPoint        string `json:"endPoint" yaml:"endPoint" xml:"endPoint" toml:"endPoint"`
	}
)

type _confScr struct {
	appConf        AppConf
	githubAuthConf *Oauth2Conf
	googleAuthConf *Oauth2Conf
	mailConf       *MailConf
	cloudflareConf *CloudflareConf
	s3Conf         *S3Conf
	jwtConf        JwtConf
	bizDBConf      DatabaseConf
	extraDBConf    DatabaseConf
}

var (
	_defaultStore = path.Join(util.GetOrFallback(os.UserHomeDir, "./"), ".space-store")
	ProjectConf   = &_confScr{
		appConf: AppConf{
			NodeID:         1,
			Port:           uint(8088),
			ServerHint:     "Shalling Space",
			Salt:           uuid.NewString(),
			MaxUserActive:  3,
			ServerTimezone: "",
			// 设置默认路径
			StaticDir: path.Join(util.GetOrFallback(func() (string, error) {
				p, e := os.UserHomeDir()
				if e != nil {
					t := fmt.Sprintf(
						"%sGet Default User Home Dir Fail, Use A fallback value %s replaced%s",
						constants.RED,
						constants.BG_CYAN,
						constants.RESET,
					)
					fmt.Println(t)
				}
				return p, e
			}, "./"), ".space-store", "files"),
			GlobalUploadSize: "32m",
			ParsedUploadSize: constants.MB * 32, // 全局的最大本地文件上传大小, 32 MB
		},
		jwtConf: JwtConf{
			Salt:          uuid.NewString(),
			Expired:       "15d",
			ParsedExpTime: time.Hour * 24 * 15,
		},
		bizDBConf: DatabaseConf{
			DBName: "bizDB",
			DBType: "sqlite",
			Dsn: path.Join(util.GetOrFallback(func() (string, error) {
				p, e := os.UserHomeDir()
				if e != nil {
					t := fmt.Sprintf(
						"%sGet Default User Home Dir Fail, Use A fallback value %s replaced%s",
						constants.RED,
						constants.BG_CYAN,
						constants.RESET,
					)
					fmt.Println(t)
				}
				return p, e
			}, "./"), ".space-store", "db", "biz-db.sqlite"),
			Mark: "biz db",
		},
		extraDBConf: DatabaseConf{
			DBName: "extraDB",
			DBType: "sqlite",
			Dsn: path.Join(util.GetOrFallback(func() (string, error) {
				p, e := os.UserHomeDir()
				if e != nil {
					t := fmt.Sprintf(
						"%sGet Default User Home Dir Fail, Use A fallback value %s replaced%s",
						constants.RED,
						constants.BG_CYAN,
						constants.RESET,
					)
					fmt.Println(t)
				}
				return p, e
			}, "./"), ".space-store", "db", "extra-db.sqlite"),
			Mark: "extra db(for log, config...)",
		},
	}
)

func (c *_confScr) GetAppConf() *AppConf {
	return &c.appConf
}

func (c *_confScr) GetMailConf() *MailConf {
	return c.mailConf
}

func (c *_confScr) GetCloudflareConf() *CloudflareConf {
	return c.cloudflareConf
}

func (c *_confScr) GetS3Conf() *S3Conf {
	return c.s3Conf
}

func (c *_confScr) GetJwtConf() *JwtConf {
	return &c.jwtConf
}

func (c *_confScr) GetGithubAuthConf() *Oauth2Conf {
	return c.githubAuthConf
}

func (c *_confScr) GetGoogleAuthConf() *Oauth2Conf {
	return c.googleAuthConf
}

func (c *_confScr) GetBizDBConf() *DatabaseConf {
	return &c.bizDBConf
}

func (c *_confScr) GetExtraDBConf() *DatabaseConf {
	return &c.extraDBConf
}

func init() {
	cfLoc, _ := GetParsedArgs()
	// 直接使用默认配置
	if cfLoc == "" {
		return
	}

	v := viper.New()
	ext := path.Ext(cfLoc)
	if len(ext) < 2 || !strings.HasPrefix(ext, ".") {
		log.Fatal("un-known extension")
	}
	v.SetConfigType(ext[1:])
	baseName := path.Base(cfLoc)
	v.SetConfigName(baseName[:len(baseName)-len(ext)])
	v.AddConfigPath(path.Dir(cfLoc))

	if e := v.ReadInConfig(); e != nil {
		log.Fatal("read config error: ", e)
	}

	if v.Get("email") != nil {
		if e := v.UnmarshalKey("email", &ProjectConf.mailConf); e != nil {
			log.Fatal("set mail config err: ", e)
		}
	}

	if v.Get("cloudflare") != nil {
		if e := v.UnmarshalKey("cloudflare", &ProjectConf.cloudflareConf); e != nil {
			log.Fatal("set cloudflare config err: ", e)
		}
	}

	if v.Get("s3") != nil {
		if e := v.UnmarshalKey("s3", &ProjectConf.s3Conf); e != nil {
			log.Fatal("set s3 storage config err: ", e)
		}
	}

	if e := v.UnmarshalKey("app", &ProjectConf.appConf); e != nil {
		log.Fatal("set config error: ", e)
	} else {
		cf := &ProjectConf.appConf
		units := []string{"byte", "kb", "mb", "gb"}
		matched := "byte"
		if !slices.ContainsFunc(units, func(u string) bool {
			t := strings.HasSuffix(
				strings.ToLower(cf.GlobalUploadSize),
				u,
			)
			if t {
				matched = u
			}
			return t
		}) {
			log.Fatal("un-support file unit size: ", cf.GlobalUploadSize)
		} else {
			sub := cf.GlobalUploadSize[:len(cf.GlobalUploadSize)-len(matched)]
			if u, e := strconv.ParseInt(sub, 10, 64); e != nil || u <= 0 {
				log.Fatal("illegal file size: ", cf.GlobalUploadSize)
			} else {
				switch matched {
				case "byte":
					cf.ParsedUploadSize = constants.Byte * constants.MemoryByteSize(u)
				case "kb":
					cf.ParsedUploadSize = constants.KB * constants.MemoryByteSize(u)
				case "mb":
					cf.ParsedUploadSize = constants.MB * constants.MemoryByteSize(u)
				case "gb":
					cf.ParsedUploadSize = constants.GB * constants.MemoryByteSize(u)
				}
			}
		}
	}

	if e := v.UnmarshalKey("dataSource.db.bizDB", &ProjectConf.bizDBConf); e != nil {
		log.Fatal("set config error: ", e)
	}

	if e := v.UnmarshalKey("dataSource.db.extraDB", &ProjectConf.extraDBConf); e != nil {
		log.Fatal("set config error: ", e)
	}

	if v.Get("oauth2Conf.github") != nil {
		if e := v.UnmarshalKey("oauth2Conf.github", &ProjectConf.githubAuthConf); e != nil {
			log.Fatal("set config error: ", e)
		}
	}

	if v.Get("oauth2Conf.google") != nil {
		if e := v.UnmarshalKey("oauth2Conf.google", &ProjectConf.googleAuthConf); e != nil {
			log.Fatal("set config error: ", e)
		}
	}

	if e := v.UnmarshalKey("jwtConf", &ProjectConf.jwtConf); e != nil {
		log.Fatal("set config error: ", e)
	} else {
		exp := ProjectConf.jwtConf.Expired
		if !(len(exp) > 1) {
			log.Fatal("un support expired time config: ", ProjectConf.jwtConf.Expired)
		}
		if !slices.ContainsFunc([]string{"s", "m", "h", "d"}, func(str string) bool {
			return strings.HasSuffix(exp, str)
		}) {
			log.Fatal("un support expired time config: ", ProjectConf.jwtConf.Expired)
		}
		if val, err := strconv.Atoi(exp[:len(exp)-1]); err != nil {
			log.Fatal("un support expired time config: ", ProjectConf.jwtConf.Expired)
		} else {
			var d time.Duration
			switch exp[len(exp)-1:] {
			case "s":
				d = time.Second * time.Duration(val)
			case "m":
				d = time.Minute * time.Duration(val)
			case "h":
				d = time.Hour * time.Duration(val)
			case "d":
				d = time.Hour * 24 * time.Duration(val)
			}
			if d <= 0 {
				log.Fatal("require a positive expired time, but got: ", exp)
			}
			ProjectConf.jwtConf.ParsedExpTime = d
		}
	}

}
