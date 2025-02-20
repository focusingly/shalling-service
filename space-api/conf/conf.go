package conf

import (
	"fmt"
	"log"
	"os"
	"path"
	"slices"
	"space-api/constants"
	"space-api/util"
	"space-api/util/arr"
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
		NotifyEmail      string                   `json:"notifyEmail" yaml:"notifyEmail" xml:"notifyEmail" toml:"notifyEmail"`
		ApiPrefix        string                   `json:"apiPrefix" yaml:"apiPrefix" xml:"apiPrefix" toml:"apiPrefix"`
		Certs            struct {
			Pem string `json:"pem" yaml:"pem" xml:"pem" toml:"pem"` // 证书配置路径
			Key string `json:"key" yaml:"key" xml:"key" toml:"key"` // 证书密钥路径
		} `json:"certs" yaml:"certs" xml:"certs" toml:"certs"`
	}

	// 邮件服务的描述配置
	MailSmtpConf struct {
		Host        string `json:"host" yaml:"host" xml:"host" toml:"host"`
		Port        int    `json:"port" yaml:"port" xml:"port" toml:"port"`
		Account     string `json:"account" yaml:"account" xml:"account" toml:"account"`
		Credential  string `json:"credential" yaml:"credential" xml:"credential" toml:"credential"`
		Primary     bool   `json:"primary" yaml:"primary" xml:"primary" toml:"primary"` // 是否被标记为首选邮箱
		Mark        string `json:"mark" yaml:"mark" xml:"mark" toml:"mark"`
		DefaultFrom string `json:"defaultFrom" yaml:"defaultFrom" xml:"defaultFrom" toml:"defaultFrom"`
		SpecificID  string `json:"specificID" yaml:"specificID" xml:"specificID" toml:"specificID"` // 配置项标识的唯一 ID
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
		LinkedDomain    string `json:"linkedDomain" yaml:"linkedDomain" xml:"linkedDomain" toml:"linkedDomain"`
	}
)

type projectRootConf struct {
	appConf          AppConf
	githubAuthConf   *Oauth2Conf
	googleAuthConf   *Oauth2Conf
	mailConfList     []*MailSmtpConf
	primaryEmailConf *MailSmtpConf
	cloudflareConf   *CloudflareConf
	s3Conf           *S3Conf
	jwtConf          JwtConf
	bizDBConf        DatabaseConf
	extraDBConf      DatabaseConf
}

var (
	defaultStore = path.Join(util.GetOrFallback(os.UserHomeDir, "./"), ".space-store")

	ProjectConf = &projectRootConf{
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
			ApiPrefix:        "/v1/api",
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
		mailConfList: make([]*MailSmtpConf, 0),
	}
)

func (c *projectRootConf) GetAppConf() *AppConf {
	return &c.appConf
}

// GetPrimaryMailConf 获取首选的邮箱配置
func (c *projectRootConf) GetPrimaryMailConf() *MailSmtpConf {
	if c.primaryEmailConf == nil {
		fmt.Printf("%s当前未配置任何主邮箱\n%s", constants.RED, constants.RESET)
	}

	return c.primaryEmailConf
}

func (c *projectRootConf) GetMailConfList() []*MailSmtpConf {
	return c.mailConfList
}

func (c *projectRootConf) GetCloudflareConf() *CloudflareConf {
	return c.cloudflareConf
}

func (c *projectRootConf) GetS3Conf() *S3Conf {
	return c.s3Conf
}

func (c *projectRootConf) GetJwtConf() *JwtConf {
	return &c.jwtConf
}

func (c *projectRootConf) GetGithubAuthConf() *Oauth2Conf {
	return c.githubAuthConf
}

func (c *projectRootConf) GetGoogleAuthConf() *Oauth2Conf {
	return c.googleAuthConf
}

func (c *projectRootConf) GetBizDBConf() *DatabaseConf {
	return &c.bizDBConf
}

func (c *projectRootConf) GetExtraDBConf() *DatabaseConf {
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

	if v.Get("emails") != nil {
		if e := v.UnmarshalKey("emails", &ProjectConf.mailConfList); e != nil {
			log.Fatal("set mail config err: ", e)
		} else {
			if len(ProjectConf.mailConfList) != 0 {
				idList := []string{}
				for _, cf := range ProjectConf.mailConfList {
					if cf.SpecificID == "" {
						log.Fatal("请提供一个具体的邮件的标识: ", cf)
					}
					if slices.Contains(idList, cf.SpecificID) {
						log.Fatal("重复的邮箱配置标识 ID")
					}
					idList = append(idList, cf.SpecificID)
				}
				primaries := arr.FilterSlice[*MailSmtpConf](ProjectConf.mailConfList, func(current *MailSmtpConf, index int) bool {
					return current.Primary
				})
				switch len(primaries) {
				case 0:
					log.Fatal("必须要提供一个主邮箱配置")
				case 1:
					// 设置主邮箱
					ProjectConf.primaryEmailConf = primaries[0]
				default:
					log.Fatal("只允许配置一个主邮箱配置, 但得到了多个: ", primaries)
				}
			}
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
				switch strings.ToLower(matched) {
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
