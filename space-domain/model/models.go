package model

import (
	"space-api/util/id"

	"gorm.io/gorm"
)

// Biz Data Tables
type (
	// LocalUser 代表本地平台的登录用户的记录信息
	LocalUser struct {
		BaseColumn   `json:"baseColumn"`
		Email        *string `gorm:"type:varchar(255);null;comment:用户邮箱, 可用于找回密码" json:"email"`
		Username     string  `gorm:"type:varchar(255);not null;unique;comment:登录的用户名称" json:"username"`
		DisplayName  string  `gorm:"type:varchar(255);not null;comment:对外展示的用户名称" json:"displayName"`
		Password     string  `gorm:"type:varchar(255);null;comment:可用于找回账户的密码" json:"password"`
		AvatarURL    *string `gorm:"type:text;null;comment:用户的头像链接"`
		HomepageLink *string `gorm:"type:text;null;comment:用户的主页链接"`
		Phone        *string `gorm:"type:varchar(255);null;comment:可用于找回账户的密码" json:"phone"`
		IsAdmin      int     `gorm:"type:smallint;default:0;comment:是否为超级管理员用户(大于 0 的都可以认为是)" json:"isAdmin"`
	}

	// OAuth2User Oauth2 用户的认证记录信息
	OAuth2User struct {
		BaseColumn     `json:"baseColumn"`
		PlatformName   string   `gorm:"type:varchar(255);not null;comment:oauth2授权来源平台名称" json:"platformName"`
		PlatformUserId string   `gorm:"type:varchar(255);not null;comment:oauth2 授权平台返回的用户标识 ID" json:"platformUserId"`
		Username       string   `gorm:"type:varchar(255);not null;comment:oauth2 用户在平台的名称" json:"username"`
		PrimaryEmail   string   `gorm:"type:varchar(255);not null;comment:用户主邮箱" json:"primaryEmail"`
		AccessToken    string   `gorm:"type:text;not null;comment:授权token" json:"accessToken"`
		RefreshToken   *string  `gorm:"type:text;null;comment:刷新 token(如果存在的话)" json:"refreshToken"`
		ExpiredAt      *int64   `gorm:"type:bigint;null;comment:凭证 token 的有效截至时间(unix 毫秒时间戳)" json:"expiredAt"`
		AvatarURL      *string  `gorm:"type:text;null;comment:用户的头像链接"`
		HomepageLink   *string  `gorm:"type:text;null;comment:用户的主页链接"`
		Scopes         []string `gorm:"type:text;null;serializer:json;comment:oauth2 申请的权限范围" json:"scopes"`
	}

	// UserLoginSession 表示登录会话信息
	UserLoginSession struct {
		BaseColumn `json:"baseColumn"`
		UserId     int64   `gorm:"type:bigint;not null;comment:用户的 ID" json:"userId"`
		UUID       string  `gorm:"type:varchar(255);not null;comment:在额外的缓存中用于标识 key, 也对应 token 中设置的 uuid" json:"uuid"`
		IpU32Val   *uint32 `gorm:"type:int;null;comment:ipv4 地址的 uint32 表示值" json:"ipU32Val"`
		IpAddress  *string `gorm:"type:varchar(255);null;comment:ip 地址表示字符串" json:"ipAddress"`
		IpSource   *string `gorm:"type:varchar(255);null;comment:ip 来源归属地" json:"ipSource"`
		ExpiredAt  int64   `gorm:"type:bigint;not null;comment:token 的过期时间,这个值为 unix 毫秒时间戳" json:"expiresAt"`
		UserType   string  `gorm:"type:varchar(255);not null;comment:当前登录的用户类型标识" json:"userType"`
		Token      string  `gorm:"type:text;not null;comment:当前用户的凭据" json:"token"`
		Useragent  string  `gorm:"type:varchar(255);comment:用户登录的平台标识" json:"useragent"`
		ClientName string  `gorm:"type:varchar(255);comment:客户端名称" json:"clientName"`
		OsName     string  `gorm:"type:varchar(255);comment:系统名称" json:"osName"`
	}

	// Post 文章
	Post struct {
		BaseColumn   `json:"baseColumn"`
		Title        string   `gorm:"type:varchar(255);not null;comment:文章标题" json:"title"`
		AuthorId     int64    `gorm:"type:bigint;not null;comment:文章作者的主键 ID" json:"authorId"`
		Content      string   `gorm:"type:text;comment:文章内容" json:"content"`
		WordCount    int64    `gorm:"type:bigint;not null;comment:字数统计" json:"wordCount"`
		ReadTime     *int64   `gorm:"type:bigint;null;comment:阅读时间 unix 毫秒时间戳" json:"readTime"`
		Category     *string  `gorm:"type:varchar(255);null;comment:所属类别名称" json:"category"`
		Tags         []string `gorm:"type:text;null;serializer:json;comment:包含的标签列表" json:"tags"`
		LastPubTime  *int64   `gorm:"type:bigint;null;comment:最后一次更新时间(允许手动选择或者不设置)" json:"lastPubTime"`
		Weight       *int     `gorm:"type:smallint;null;comment:可选的权重标识" json:"weight"`
		Views        *int64   `gorm:"type:bigint;null;comment:总浏览量" json:"views"`
		UpVote       *int64   `gorm:"type:bigint;null;comment:赞成数" json:"upVote"`
		DownVote     *int64   `gorm:"type:bigint;null;comment:否定数" json:"downVote"`
		AllowComment int      `gorm:"type:smallint;default:0;comment:是否允许评论, 默认为允许" json:"allowComment"`
	}

	// Tag 文章标签信息
	Tag struct {
		BaseColumn `json:"baseColumn"`
		TagName    string  `gorm:"type:varchar(255);not null;unique;comment:标签名称" json:"tagName"`
		Color      *string `gorm:"type:varchar(255);null;comment:标签颜色" json:"color"`
		IconUrl    *string `gorm:"type:varchar(255);null;comment:标签图标" json:"iconUrl"`
	}

	// PostTagRelation 文章和标签的关联关系
	PostTagRelation struct {
		BaseColumn `json:"baseColumn"`
		PostId     int64 `gorm:"type:bigint;not null;comment:文章的主键 ID" json:"postId"`
		TagId      int64 `gorm:"type:bigint;not null;comment:标签的主键 ID" json:"tagId"`
	}

	// Category 分类信息
	Category struct {
		BaseColumn   `json:"baseColumn"`
		CategoryName string  `gorm:"type:varchar(255);not null;unique;comment:类型名称" json:"categoryName"`
		Color        *string `gorm:"type:varchar(255);null;comment:类型颜色;" json:"color"`
		IconUrl      *string `gorm:"type:varchar(255);null;comment:标签图标" json:"iconUrl"`
	}

	// 评论表
	Comment struct {
		BaseColumn    `json:"baseColumn"`
		PostId        int64   `gorm:"type:bigint;not null;comment:文章的的 ID" json:"postId"`
		UserId        int64   `gorm:"type:bigint;not null;comment:评论用户 ID" json:"userId"`
		UserType      string  `gorm:"type:text;null;comment:用户类型" json:"userType"`
		Avatar        *string `gorm:"type:text;null;comment:用户头像的链接" json:"avatar"`
		HomePageURL   *string `gorm:"type:text;null;comment:用户的公开主页" json:"homePageURL"`
		RootCommentId int64   `gorm:"type:bigint;not null;default:0;comment:评论所属的根评论 ID, 用于查找评论下所有子评论的以及评论二级分页" json:"rootCommentId"`
		ReplyToId     int64   `gorm:"type:bigint;not null;default:0;comment:回复的上条评论 ID, 如果自身是根评论, 那么为 0" json:"replyToId"`
		Content       string  `gorm:"type:text;not null;comment:评论内容" json:"content"`
		UpVote        *int64  `gorm:"type:bigint;null;comment:赞成数" json:"upVote"`
		DownVote      *int64  `gorm:"type:bigint;null;comment:否定数" json:"downVote"`
		IpAddr        string  `gorm:"type:varchar(255);null;comment:客户端的 IP 地址" json:"ipAddr"`
		IpSource      *string `gorm:"type:varchar(255);null;comment:IP 归属地" json:"ipSource"`
		Useragent     *string `gorm:"type:text;null;comment:用户的客户端标识" json:"useragent"`
		OS            *string `gorm:"type:text;null;comment:用户的操作系统平台" json:"os"`
		SubEmailReply int     `gorm:"type:smallint;default:0;comment:是否订阅邮件的通知回复" json:"subEmailReply"`
		ClientName    *string `gorm:"type:text;null;comment:用户的访问平台名称" json:"clientName"`
	}

	// 公开的社交媒体信息
	PubSocialMedia struct {
		BaseColumn  `json:"baseColumn"`
		DisplayName string `gorm:"type:varchar(255);not null;comment:显示名称" json:"displayName"`
		IconURL     string `gorm:"type:varchar(255);not null;comment:图标链接" json:"iconURL"`
		OpenUrl     string `gorm:"type:varchar(255);not null;comment:跳转链接" json:"openUrl"`
	}

	// 云函数(如 cloudflare 的 worker) 相关的代码片段
	CloudFn struct {
		BaseColumn `json:"baseColumn"`
		FuncName   string `gorm:"not null;unique;comment:代码片段名称" json:"funcName"`
		Code       string `gorm:"type:text;not null;comment:代码内容" json:"code"`
		Lang       string `gorm:"type:varchar(255);not null;comment:语言类型" json:"lang"`
		Enable     int    `gorm:"type:smallint;default:0;not null;comment:是否启用,默认不启用(0)" json:"enable"`
	}

	// 本地的文件保存记录
	FileRecord struct {
		BaseColumn    `json:"baseColumn"`
		FileName      string `gorm:"type:varchar(255);not null;comment:" json:"fileName"`
		LocalLocation string `gorm:"varchar(255);not null;comment:本地存储路径" json:"localLocation"`
		Extension     string `gorm:"type:varchar(255);not null;comment:" json:"extension"`
		FileSize      int64  `gorm:"type:bigint;not null;comment:" json:"fileSize"`
		Category      string `gorm:"type:varchar(255);not null;comment:文件归类名称" json:"category"`
		ChecksumType  string `gorm:"varchar(255);not null;comment:校验类型" json:"checksumType"`
		Checksum      string `gorm:"type:text;not null;comment:校验和" json:"checksum"`
	}

	// 使用 OSS3 的相关服务存储记录
	OSS3Record struct {
		BaseColumn   `json:"baseColumn"`
		FileName     string `gorm:"type:varchar(255);not null;comment:" json:"fileName"`
		Extension    string `gorm:"type:varchar(255);not null;comment:" json:"extension"`
		FileSize     int64  `gorm:"type:bigint;not null;comment:" json:"fileSize"`
		BucketName   string `gorm:"type:varchar(255);not null;comment:存储的桶名称" json:"category"`
		VisitURL     string `gorm:"type:text;not null;comment:oss存储上的相对路径" json:"visitURL"`
		ChecksumType string `gorm:"varchar(255);not null;comment:校验类型" json:"checksumType"`
		Checksum     string `gorm:"type:text;not null;comment:校验和" json:"checksum"`
	}

	// 菜单组
	MenuGroup struct {
		BaseColumn
		MenuName string // 菜单名称
		// ParentID  int64 // TODO 暂时不实现嵌套菜单
		RoutePath *string     `gorm:"type:varchar(255);null;comment:可选的前端路由地址" json:"routePath"`
		PostLink  *int64      `gorm:"type:varchar(255);null;comment:可选的文章的 ID(如果是站内的文章的话)" json:"postLink"`
		OpenWay   string      `gorm:"type:varchar(255);default:_self;not null;comment:链接的打开方式(如当前页面/打开新的页面)" json:"openWay"`
		SubLinks  []*MenuLink `gorm:"type:varchar(255);default:_self;not null;comment:菜单包含的子链接列表" json:"subLinks"`
	}

	// 动态链接
	MenuLink struct {
		BaseColumn `json:"baseColumn"`
		LinkName   string  `gorm:"type:varchar(255);not null;unique;comment:链接显示的名称" json:"linkName"`
		RoutePath  *string `gorm:"type:varchar(255);null;comment:可选的前端的跳转路由地址" json:"routePath"`
		PostLink   *string `gorm:"type:varchar(255);null;comment:可续的文章的 ID(如果是站内的文章的话)" json:"postLink"`
		OpenWay    string  `gorm:"type:varchar(255);default:_self;not null;comment:链接的打开方式(如当前页面/打开新的页面)" json:"openWay"`
	}

	// 友链
	FriendLink struct {
		BaseColumn  `json:"baseColumn"`
		SiteURL     string  `gorm:"type:varchar(255);not null;comment:站点主页链接" json:"siteURL"`
		Owner       string  `gorm:"type:varchar(255);not null;comment:拥有者的名称" json:"owner"`
		ShortName   string  `gorm:"type:varchar(255);not null;comment:站点的简单描述信息" json:"shortName"`
		Available   int     `gorm:"type:smallint;not null;default:0;comment:站点是否可用(表示是否可以正常访问, 默认为不可用:0)" json:"available"`
		LogoURL     string  `gorm:"type:varchar(255);not null;comment:站点的页签或者站主的头像图片链接" json:"logoURL"`
		Description *string `gorm:"type:text;null;comment:可选的额外描述信息" json:"description"`
		BgURL       *string `gorm:"type:varchar(255);null;comment:可选的展示卡片底图背景" json:"bgURL"`
	}

	// ServiceConf 自定义配置
	ServiceConf struct {
		BaseColumn `json:"baseColumn"`
		ConfKey    string `gorm:"type:varchar(255);not null;unique;comment:配置名称" json:"confKey"`
		ConfVal    string `gorm:"type:text;comment;null;comment:配置值" json:"confVal"`
		Category   string `gorm:"type:varchar(255);not null;comment:类型" json:"category"`
	}
)

// Extra tables
type (
	// 定时任务管理
	CronJob struct {
		BaseColumn  `json:"baseColumn"`
		JobFuncName string `gorm:"type:varchar(255);not null;comment:执行的任务函数名称" json:"jobFuncName"`
		CronExpr    string `gorm:"type:varchar(255);not null;comment:cron 表达式" json:"cronExpr"`
		Status      string `gorm:"type:varchar(255);not null;comment:运行状态" json:"status"`
		Enable      int    `gorm:"type:smallint;default:0;not null;comment:是否启用,默认不启用(0)" json:"enable"`
		Mark        string `gorm:"type:varchar(255);null;comment:可选的任务备注" json:"mark"`
	}

	LogRecord struct {
		ID        int64  `gorm:"primaryKey;autoIncrement:false;comment:日志 ID" json:"id"`
		Category  string `gorm:"type:varchar(255);not null;comment:日志类型" json:"category"`
		Content   string `gorm:"type:text;comment:日志内容" json:"content"`
		Source    string `gorm:"type:varchar(255);comment:来源信息" json:"source"`
		CreatedAt int64  `gorm:"autoCreateTime:milli" json:"createdAt"`
	}
)

func (logRecord *LogRecord) BeforeCreate(tx *gorm.DB) (err error) {
	logRecord.ID = id.GetSnowFlakeNode().Generate().Int64()

	return
}

func GetBizMigrateTables() []any {
	return []any{
		new(LocalUser),
		new(UserLoginSession),
		new(OAuth2User),
		new(Post),
		new(Tag),
		new(PostTagRelation),
		new(Category),
		new(Comment),
		new(PubSocialMedia),
		new(CloudFn),
		new(FileRecord),
		new(MenuGroup),
		new(MenuLink),
		new(FriendLink),
		new(ServiceConf),
	}
}

func GetExtraHelperMigrateTables() []any {
	return []any{
		new(CronJob),
		new(LogRecord),
	}
}
