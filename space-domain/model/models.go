package model

import (
	"space-api/util/id"

	"gorm.io/gorm"
)

// Biz Data Tables
type (
	// LocalUser 代表本地平台的登录用户的记录信息
	LocalUser struct {
		BaseColumn   `json:"baseColumn" yaml:"baseColumn" xml:"baseColumn" toml:"baseColumn"`
		Email        *string `gorm:"type:varchar(255);null;comment:用户邮箱, 可用于找回密码" json:"email" yaml:"email" xml:"email" toml:"email"`
		Username     string  `gorm:"type:varchar(255);not null;unique;comment:登录的用户名称" json:"username" yaml:"username" xml:"username" toml:"username"`
		DisplayName  string  `gorm:"type:varchar(255);not null;comment:对外展示的用户名称" json:"displayName" yaml:"displayName" xml:"displayName" toml:"displayName"`
		Password     string  `gorm:"type:text;not null;comment:可用于找回账户的密码" json:"password" yaml:"password" xml:"password" toml:"password"`
		AvatarURL    *string `gorm:"type:text;null;comment:用户的头像链接" json:"avatarURL" yaml:"avatarURL" xml:"avatarURL" toml:"avatarURL"`
		HomepageLink *string `gorm:"type:text;null;comment:用户的主页链接" json:"homepageLink" yaml:"homepageLink" xml:"homepageLink" toml:"homepageLink"`
		Phone        *string `gorm:"type:varchar(255);null;comment:可用于找回账户的密码" json:"phone" yaml:"phone" xml:"phone" toml:"phone"`
		IsAdmin      int     `gorm:"type:smallint;default:0;comment:是否为超级管理员用户(大于 0 的都可以认为是)" json:"isAdmin" yaml:"isAdmin" xml:"isAdmin" toml:"isAdmin"`
	}

	// OAuth2User Oauth2 用户的认证记录信息
	OAuth2User struct {
		BaseColumn     `json:"baseColumn" yaml:"baseColumn" xml:"baseColumn" toml:"baseColumn"`
		PlatformName   string   `gorm:"type:varchar(255);not null;comment:oauth2授权来源平台名称" json:"platformName" yaml:"platformName" xml:"platformName" toml:"platformName"`
		PlatformUserId string   `gorm:"type:varchar(255);not null;comment:oauth2 授权平台返回的用户标识 ID" json:"platformUserId" yaml:"platformUserId" xml:"platformUserId" toml:"platformUserId"`
		Username       string   `gorm:"type:varchar(255);not null;comment:oauth2 用户在平台的名称" json:"username" yaml:"username" xml:"username" toml:"username"`
		PrimaryEmail   string   `gorm:"type:varchar(255);not null;comment:用户主邮箱" json:"primaryEmail" yaml:"primaryEmail" xml:"primaryEmail" toml:"primaryEmail"`
		AccessToken    string   `gorm:"type:text;not null;comment:授权token" json:"accessToken" yaml:"accessToken" xml:"accessToken" toml:"accessToken"`
		RefreshToken   *string  `gorm:"type:text;null;comment:刷新 token(如果存在的话)" json:"refreshToken" yaml:"refreshToken" xml:"refreshToken" toml:"refreshToken"`
		ExpiredAt      *int64   `gorm:"type:bigint;null;comment:凭证 token 的有效截至时间(unix 毫秒时间戳)" json:"expiredAt" yaml:"expiredAt" xml:"expiredAt" toml:"expiredAt"`
		AvatarURL      *string  `gorm:"type:text;null;comment:用户的头像链接" json:"avatarURL" yaml:"avatarURL" xml:"avatarURL" toml:"avatarURL"`
		HomepageLink   *string  `gorm:"type:text;null;comment:用户的主页链接" json:"homepageLink" yaml:"homepageLink" xml:"homepageLink" toml:"homepageLink"`
		Scopes         []string `gorm:"type:text;null;serializer:json;comment:oauth2 申请的权限范围" json:"scopes" yaml:"scopes" xml:"scopes" toml:"scopes"`
		Enable         int      `gorm:"type:smallint;default:1;not null;comment:用户是否启用" json:"enable" yaml:"enable" xml:"enable" toml:"enable"`
	}

	// UserLoginSession 表示登录会话信息
	UserLoginSession struct {
		BaseColumn `json:"baseColumn" yaml:"baseColumn" xml:"baseColumn" toml:"baseColumn"`
		UserID     int64   `gorm:"type:bigint;not null;comment:用户的 ID" json:"userID" yaml:"userID" xml:"userID" toml:"userID"`
		UUID       string  `gorm:"type:varchar(255);not null;comment:在额外的缓存中用于标识 key, 也对应 token 中设置的 uuid" json:"uuid" yaml:"uuid" xml:"uuid" toml:"uuid"`
		IpU32Val   *uint32 `gorm:"type:int;null;comment:ipv4 地址的 uint32 表示值" json:"ipU32Val" yaml:"ipU32Val" xml:"ipU32Val" toml:"ipU32Val"`
		IpAddress  *string `gorm:"type:varchar(255);null;comment:ip 地址表示字符串" json:"ipAddress" yaml:"ipAddress" xml:"ipAddress" toml:"ipAddress"`
		IpSource   *string `gorm:"type:varchar(255);null;comment:ip 来源归属地" json:"ipSource" yaml:"ipSource" xml:"ipSource" toml:"ipSource"`
		ExpiredAt  int64   `gorm:"type:bigint;not null;comment:token 的过期时间,这个值为 unix 毫秒时间戳" json:"expiresAt" yaml:"expiredAt" xml:"expiredAt" toml:"expiredAt"`
		UserType   string  `gorm:"type:varchar(255);not null;comment:当前登录的用户类型标识" json:"userType" yaml:"userType" xml:"userType" toml:"userType"`
		Token      string  `gorm:"type:text;not null;comment:当前用户的凭据" json:"token" yaml:"token" xml:"token" toml:"token"`
		Useragent  string  `gorm:"type:varchar(255);comment:用户登录的平台标识" json:"useragent" yaml:"useragent" xml:"useragent" toml:"useragent"`
		ClientName string  `gorm:"type:varchar(255);comment:客户端名称" json:"clientName" yaml:"clientName" xml:"clientName" toml:"clientName"`
		OsName     string  `gorm:"type:varchar(255);comment:系统名称" json:"osName" yaml:"osName" xml:"osName" toml:"osName"`
	}

	// Post 文章
	Post struct {
		BaseColumn   `json:"baseColumn" yaml:"baseColumn" xml:"baseColumn" toml:"baseColumn"`
		Title        string   `gorm:"type:varchar(255);not null;comment:文章标题" json:"title" yaml:"title" xml:"title" toml:"title"`
		AuthorId     int64    `gorm:"type:bigint;not null;comment:文章作者的主键 ID" json:"authorId" yaml:"authorId" xml:"authorId" toml:"authorId"`
		Content      string   `gorm:"type:text;comment:文章内容" json:"content" yaml:"content" xml:"content" toml:"content"`
		WordCount    int64    `gorm:"type:bigint;not null;comment:字数统计" json:"wordCount" yaml:"wordCount" xml:"wordCount" toml:"wordCount"`
		ReadTime     *int64   `gorm:"type:bigint;null;comment:阅读时间 unix 毫秒时间戳" json:"readTime" yaml:"readTime" xml:"readTime" toml:"readTime"`
		Category     *string  `gorm:"type:varchar(255);null;comment:所属类别名称" json:"category" yaml:"category" xml:"category" toml:"category"`
		Tags         []string `gorm:"type:text;null;serializer:json;comment:包含的标签列表" json:"tags" yaml:"tags" xml:"tags" toml:"tags"`
		LastPubTime  *int64   `gorm:"type:bigint;null;comment:最后一次更新时间(允许手动选择或者不设置)" json:"lastPubTime" yaml:"lastPubTime" xml:"lastPubTime" toml:"lastPubTime"`
		Weight       *int     `gorm:"type:smallint;null;comment:可选的权重标识" json:"weight" yaml:"weight" xml:"weight" toml:"weight"`
		Views        *int64   `gorm:"type:bigint;null;comment:总浏览量" json:"views" yaml:"views" xml:"views" toml:"views"`
		UpVote       *int64   `gorm:"type:bigint;null;comment:赞成数" json:"upVote" yaml:"upVote" xml:"upVote" toml:"upVote"`
		DownVote     *int64   `gorm:"type:bigint;null;comment:否定数" json:"downVote" yaml:"downVote" xml:"downVote" toml:"downVote"`
		Snippet      *string  `gorm:"type:text;null;comment:文章摘要" json:"snippet" yaml:"snippet" xml:"snippet" toml:"snippet"`
		Lang         *string  `gorm:"type:varchar(255);null;comment:文章使用的语言" json:"lang" yaml:"lang" xml:"lang" toml:"lang"`
		AllowComment int      `gorm:"type:smallint;default:0;comment:是否允许评论, 默认为允许" json:"allowComment" yaml:"allowComment" xml:"allowComment" toml:"allowComment"`
	}

	// Tag 文章标签信息
	Tag struct {
		BaseColumn `json:"baseColumn" yaml:"baseColumn" xml:"baseColumn" toml:"baseColumn"`
		TagName    string  `gorm:"type:varchar(255);not null;unique;comment:标签名称" json:"tagName" yaml:"tagName" xml:"tagName" toml:"tagName"`
		Color      *string `gorm:"type:varchar(255);null;comment:标签颜色" json:"color" yaml:"color" xml:"color" toml:"color"`
		IconUrl    *string `gorm:"type:varchar(255);null;comment:标签图标" json:"iconUrl" yaml:"iconUrl" xml:"iconUrl" toml:"iconUrl"`
	}

	// PostTagRelation 文章和标签的关联关系
	PostTagRelation struct {
		BaseColumn `json:"baseColumn" yaml:"baseColumn" xml:"baseColumn" toml:"baseColumn"`
		PostId     int64 `gorm:"type:bigint;not null;comment:文章的主键 ID" json:"postId" yaml:"postId" xml:"postId" toml:"postId"`
		TagId      int64 `gorm:"type:bigint;not null;comment:标签的主键 ID" json:"tagId" yaml:"tagId" xml:"tagId" toml:"tagId"`
	}

	// Category 分类信息
	Category struct {
		BaseColumn   `json:"baseColumn" yaml:"baseColumn" xml:"baseColumn" toml:"baseColumn"`
		CategoryName string  `gorm:"type:varchar(255);not null;unique;comment:类型名称" json:"categoryName" yaml:"categoryName" xml:"categoryName" toml:"categoryName"`
		Color        *string `gorm:"type:varchar(255);null;comment:类型颜色;" json:"color" yaml:"color" xml:"color" toml:"color"`
		IconUrl      *string `gorm:"type:varchar(255);null;comment:标签图标" json:"iconUrl" yaml:"iconUrl" xml:"iconUrl" toml:"iconUrl"`
	}

	// 评论表
	Comment struct {
		BaseColumn    `json:"baseColumn" yaml:"baseColumn" xml:"baseColumn" toml:"baseColumn"`
		PostId        int64   `gorm:"type:bigint;not null;comment:文章的的 ID" json:"postId" yaml:"postId" xml:"postId" toml:"postId"`
		UserId        int64   `gorm:"type:bigint;not null;comment:评论用户 ID" json:"userId" yaml:"userId" xml:"userId" toml:"userId"`
		UserType      string  `gorm:"type:text;null;comment:用户类型" json:"userType" yaml:"userType" xml:"userType" toml:"userType"`
		Avatar        *string `gorm:"type:text;null;comment:用户头像的链接" json:"avatar" yaml:"avatar" xml:"avatar" toml:"avatar"`
		HomePageURL   *string `gorm:"type:text;null;comment:用户的公开主页" json:"homePageURL" yaml:"homePageURL" xml:"homePageURL" toml:"homePageURL"`
		RootCommentId int64   `gorm:"type:bigint;not null;default:0;comment:评论所属的根评论 ID, 用于查找评论下所有子评论的以及评论二级分页" json:"rootCommentId" yaml:"rootCommentId" xml:"rootCommentId" toml:"rootCommentId"`
		ReplyToId     int64   `gorm:"type:bigint;not null;default:0;comment:回复的上条评论 ID, 如果自身是根评论, 那么为 0" json:"replyToId" yaml:"replyToId" xml:"replyToId" toml:"replyToId"`
		Content       string  `gorm:"type:text;not null;comment:评论内容" json:"content" yaml:"content" xml:"content" toml:"content"`
		UpVote        *int64  `gorm:"type:bigint;null;comment:赞成数" json:"upVote" yaml:"upVote" xml:"upVote" toml:"upVote"`
		DownVote      *int64  `gorm:"type:bigint;null;comment:否定数" json:"downVote" yaml:"downVote" xml:"downVote" toml:"downVote"`
		IpAddr        string  `gorm:"type:varchar(255);null;comment:客户端的 IP 地址" json:"ipAddr" yaml:"ipAddr" xml:"ipAddr" toml:"ipAddr"`
		IpSource      *string `gorm:"type:varchar(255);null;comment:IP 归属地" json:"ipSource" yaml:"ipSource" xml:"ipSource" toml:"ipSource"`
		Useragent     *string `gorm:"type:text;null;comment:用户的客户端标识" json:"useragent" yaml:"useragent" xml:"useragent" toml:"useragent"`
		OS            *string `gorm:"type:text;null;comment:用户的操作系统平台" json:"os" yaml:"os" xml:"os" toml:"os"`
		SubEmailReply int     `gorm:"type:smallint;default:0;comment:是否订阅邮件的通知回复" json:"subEmailReply" yaml:"subEmailReply" xml:"subEmailReply" toml:"subEmailReply"`
		ClientName    *string `gorm:"type:text;null;comment:用户的访问平台名称" json:"clientName" yaml:"clientName" xml:"clientName" toml:"clientName"`
	}

	// 公开的社交媒体信息
	PubSocialMedia struct {
		BaseColumn  `json:"baseColumn" yaml:"baseColumn" xml:"baseColumn" toml:"baseColumn"`
		DisplayName string `gorm:"type:varchar(255);not null;comment:显示名称" json:"displayName" yaml:"displayName" xml:"displayName" toml:"displayName"`
		IconURL     string `gorm:"type:varchar(255);not null;comment:图标链接" json:"iconURL" yaml:"iconURL" xml:"iconURL" toml:"iconURL"`
		OpenUrl     string `gorm:"type:varchar(255);not null;comment:跳转链接" json:"openUrl" yaml:"openUrl" xml:"openUrl" toml:"openUrl"`
	}

	// 云函数(如 cloudflare 的 worker) 相关的代码片段
	CloudFn struct {
		BaseColumn `json:"baseColumn" yaml:"baseColumn" xml:"baseColumn" toml:"baseColumn"`
		FuncName   string `gorm:"not null;unique;comment:代码片段名称" json:"funcName" yaml:"funcName" xml:"funcName" toml:"funcName"`
		Code       string `gorm:"type:text;not null;comment:代码内容" json:"code" yaml:"code" xml:"code" toml:"code"`
		Lang       string `gorm:"type:varchar(255);not null;comment:语言类型" json:"lang" yaml:"lang" xml:"lang" toml:"lang"`
		Enable     int    `gorm:"type:smallint;default:0;not null;comment:是否启用,默认不启用(0)" json:"enable" yaml:"enable" xml:"enable" toml:"enable"`
	}

	// 本地的文件保存记录
	FileRecord struct {
		BaseColumn    `json:"baseColumn" yaml:"baseColumn" xml:"baseColumn" toml:"baseColumn"`
		FileName      string `gorm:"type:varchar(255);not null;comment:" json:"fileName" yaml:"fileName" xml:"fileName" toml:"fileName"`
		LocalLocation string `gorm:"varchar(255);not null;comment:本地存储路径" json:"localLocation" yaml:"localLocation" xml:"localLocation" toml:"localLocation"`
		Extension     string `gorm:"type:varchar(255);not null;comment:" json:"extension" yaml:"extension" xml:"extension" toml:"extension"`
		FileSize      int64  `gorm:"type:bigint;not null;comment:" json:"fileSize" yaml:"fileSize" xml:"fileSize" toml:"fileSize"`
		Category      string `gorm:"type:varchar(255);not null;comment:文件归类名称" json:"category" yaml:"category" xml:"category" toml:"category"`
		ChecksumType  string `gorm:"varchar(255);not null;comment:校验类型" json:"checksumType" yaml:"checksumType" xml:"checksumType" toml:"checksumType"`
		Checksum      string `gorm:"type:text;not null;comment:校验和" json:"checksum" yaml:"checksum" xml:"checksum" toml:"checksum"`
		PubAvailable  int    `gorm:"type:smallint;not null;default:0;comment:是否允许对外的公开访问(默认不允许: 0)" json:"pubAvailable" yaml:"pubAvailable" xml:"pubAvailable" toml:"pubAvailable"`
	}

	// 使用 OSS3 的相关服务存储记录
	S3ObjectRecord struct {
		BaseColumn   `json:"baseColumn" yaml:"baseColumn" xml:"baseColumn" toml:"baseColumn"`
		ObjectKey    string `gorm:"type:text;not null;comment:对象的唯一标识" json:"objectKey" yaml:"objectKey" xml:"objectKey" toml:"objectKey"`
		FileName     string `gorm:"type:varchar(255);not null;comment:" json:"fileName" yaml:"fileName" xml:"fileName" toml:"fileName"`
		Extension    string `gorm:"type:varchar(255);not null;comment:" json:"extension" yaml:"extension" xml:"extension" toml:"extension"`
		FileSize     int64  `gorm:"type:bigint;not null;comment:" json:"fileSize" yaml:"fileSize" xml:"fileSize" toml:"fileSize"`
		BucketName   string `gorm:"type:varchar(255);not null;comment:存储的桶名称" json:"category" yaml:"bucketName" xml:"bucketName" toml:"bucketName"`
		ChecksumType string `gorm:"varchar(255);not null;comment:校验类型" json:"checksumType" yaml:"checksumType" xml:"checksumType" toml:"checksumType"`
		Checksum     string `gorm:"type:text;not null;comment:校验和" json:"checksum" yaml:"checksum" xml:"checksum" toml:"checksum"`
		PubAvailable int    `gorm:"type:smallint;not null;default:0;comment:是否允许对外的公开访问(默认不允许: 0)" json:"pubAvailable" yaml:"pubAvailable" xml:"pubAvailable" toml:"pubAvailable"`
	}

	// 菜单组
	MenuGroup struct {
		BaseColumn `json:"baseColumn" yaml:"baseColumn" xml:"baseColumn" toml:"baseColumn"`
		MenuName   string `json:"menuName" yaml:"menuName" xml:"menuName" toml:"menuName"` // 菜单名称
		// ParentID  int64 // TODO 暂时不实现嵌套菜单
		RoutePath *string     `gorm:"type:varchar(255);null;comment:可选的前端路由地址" json:"routePath" yaml:"routePath" xml:"routePath" toml:"routePath"`
		PostLink  *int64      `gorm:"type:varchar(255);null;comment:可选的文章的 ID(如果是站内的文章的话)" json:"postLink" yaml:"postLink" xml:"postLink" toml:"postLink"`
		OpenWay   string      `gorm:"type:varchar(255);default:_self;not null;comment:链接的打开方式(如当前页面/打开新的页面)" json:"openWay" yaml:"openWay" xml:"openWay" toml:"openWay"`
		SubLinks  []*MenuLink `gorm:"type:text;null;comment:菜单包含的子链接列表" json:"subLinks" yaml:"subLinks" xml:"subLinks" toml:"subLinks"`
	}

	// 动态链接
	MenuLink struct {
		BaseColumn `json:"baseColumn" yaml:"baseColumn" xml:"baseColumn" toml:"baseColumn"`
		LinkName   string  `gorm:"type:varchar(255);not null;unique;comment:链接显示的名称" json:"linkName" yaml:"linkName" xml:"linkName" toml:"linkName"`
		RoutePath  *string `gorm:"type:varchar(255);null;comment:可选的前端的跳转路由地址" json:"routePath" yaml:"routePath" xml:"routePath" toml:"routePath"`
		PostLink   *string `gorm:"type:varchar(255);null;comment:可续的文章的 ID(如果是站内的文章的话)" json:"postLink" yaml:"postLink" xml:"postLink" toml:"postLink"`
		OpenWay    string  `gorm:"type:varchar(255);default:_self;not null;comment:链接的打开方式(如当前页面/打开新的页面)" json:"openWay" yaml:"openWay" xml:"openWay" toml:"openWay"`
	}

	// 友链
	FriendLink struct {
		BaseColumn  `json:"baseColumn" yaml:"baseColumn" xml:"baseColumn" toml:"baseColumn"`
		SiteURL     string  `gorm:"type:varchar(255);not null;comment:站点主页链接" json:"siteURL" yaml:"siteURL" xml:"siteURL" toml:"siteURL"`
		Owner       string  `gorm:"type:varchar(255);not null;comment:拥有者的名称" json:"owner" yaml:"owner" xml:"owner" toml:"owner"`
		ShortName   string  `gorm:"type:varchar(255);not null;comment:站点的简单描述信息" json:"shortName" yaml:"shortName" xml:"shortName" toml:"shortName"`
		Available   int     `gorm:"type:smallint;not null;default:0;comment:站点是否可用(表示是否可以正常访问, 默认为不可用:0)" json:"available" yaml:"available" xml:"available" toml:"available"`
		LogoURL     string  `gorm:"type:varchar(255);not null;comment:站点的页签或者站主的头像图片链接" json:"logoURL" yaml:"logoURL" xml:"logoURL" toml:"logoURL"`
		Description *string `gorm:"type:text;null;comment:可选的额外描述信息" json:"description" yaml:"description" xml:"description" toml:"description"`
		BgURL       *string `gorm:"type:varchar(255);null;comment:可选的展示卡片底图背景" json:"bgURL" yaml:"bgURL" xml:"bgURL" toml:"bgURL"`
	}

	// 定时任务管理
	CronJob struct {
		BaseColumn  `json:"baseColumn" yaml:"baseColumn" xml:"baseColumn" toml:"baseColumn"`
		JobFuncName string `gorm:"type:varchar(255);not null;comment:执行的任务函数名称" json:"jobFuncName" yaml:"jobFuncName" xml:"jobFuncName" toml:"jobFuncName"`
		CronExpr    string `gorm:"type:varchar(255);not null;comment:cron 表达式" json:"cronExpr" yaml:"cronExpr" xml:"cronExpr" toml:"cronExpr"`
		Status      string `gorm:"type:varchar(255);not null;comment:运行状态" json:"status" yaml:"status" xml:"status" toml:"status"`
		Enable      int    `gorm:"type:smallint;default:0;not null;comment:是否启用,默认不启用(0)" json:"enable" yaml:"enable" xml:"enable" toml:"enable"`
		Mark        string `gorm:"type:varchar(255);null;comment:可选的任务备注" json:"mark" yaml:"mark" xml:"mark" toml:"mark"`
	}

	// 被阻止的 IP 列表
	BlockIPRecord struct {
		BaseColumn  `json:"baseColumn" yaml:"baseColumn" xml:"baseColumn" toml:"baseColumn"`
		IPAddr      string `gorm:"type:varchar(255);unique;not null;comment:ip 地址" json:"ipAddr" yaml:"ipAddr" xml:"ipAddr" toml:"ipAddr"`
		IPSource    string `gorm:"type:varchar(255);null;comment:ip 来源" json:"ipSource" yaml:"ipSource" xml:"ipSource" toml:"ipSource"`
		UserAgent   string `gorm:"type:varchar(255);null;comment:用户代理标识" json:"userAgent" yaml:"userAgent" xml:"userAgent" toml:"userAgent"`
		LastRequest int64  `gorm:"type:bigint;comment:最后一次请求的时间" json:"lastRequest" yaml:"lastRequest" xml:"lastRequest" toml:"lastRequest"`
	}

	// ServiceConf 自定义配置
	ServiceConf struct {
		BaseColumn `json:"baseColumn" yaml:"baseColumn" xml:"baseColumn" toml:"baseColumn"`
		ConfKey    string `gorm:"type:varchar(255);not null;unique;comment:配置名称" json:"confKey" yaml:"confKey" xml:"confKey" toml:"confKey"`
		ConfVal    string `gorm:"type:text;comment;null;comment:配置值" json:"confVal" yaml:"confVal" xml:"confVal" toml:"confVal"`
		Category   string `gorm:"type:varchar(255);not null;comment:类型" json:"category" yaml:"category" xml:"category" toml:"category"`
	}

	// Sqlite3KeywordDoc sqlite3 虚拟表配置
	Sqlite3KeywordDoc struct {
		PostID          int64  `gorm:"column:post_id" json:"postID" yaml:"postID" xml:"postID" toml:"postID"`
		TileSplit       string `gorm:"column:title_split" json:"tileSplit" yaml:"tileSplit" xml:"tileSplit" toml:"tileSplit"`
		ContentSplit    string `gorm:"column:content_split" json:"contentSplit" yaml:"contentSplit" xml:"contentSplit" toml:"contentSplit"`
		Weight          int    `gorm:"column:weight" json:"weight" yaml:"weight" xml:"weight" toml:"weight"`
		PostUpdatedAt   int64  `gorm:"column:post_updated_at" json:"postUpdatedAt" yaml:"postUpdatedAt" xml:"postUpdatedAt" toml:"postUpdatedAt"`
		RecordCreatedAt int64  `gorm:"column:record_created_at" json:"recordCreatedAt" yaml:"recordCreatedAt" xml:"recordCreatedAt" toml:"recordCreatedAt"`
		RecordUpdatedAt int64  `gorm:"column:record_updated_at" json:"recordUpdatedAt" yaml:"recordUpdatedAt" xml:"recordUpdatedAt" toml:"recordUpdatedAt"`
	}
)

// Extra tables
type (
	LogInfo struct {
		ID            int64   `gorm:"primaryKey;autoIncrement:false;comment:日志 ID" json:"id" yaml:"id" xml:"id" toml:"id"`
		LogType       string  `gorm:"type:varchar(255);not null;comment:日志类型" json:"category" yaml:"logType" xml:"logType" toml:"logType"`
		Message       string  `gorm:"type:varchar(255);not null;comment:消息" json:"message" yaml:"message" xml:"message" toml:"message"`
		Level         string  `gorm:"type:varchar(255);not null;comment:" json:"level" yaml:"level" xml:"level" toml:"level"`
		CostTime      int64   `gorm:"type:bigint;not null;comment:执行耗费时间" json:"costTime" yaml:"costTime" xml:"costTime" toml:"costTime"`
		RequestMethod *string `gorm:"type:varchar(255);comment:请求方式(如果是 web 请求的话)" json:"requestMethod" yaml:"requestMethod" xml:"requestMethod" toml:"requestMethod"`
		RequestURI    *string `gorm:"type:text;null;comment:访问的URI(如果是 web 请求的话)" json:"requestURI" yaml:"requestURI" xml:"requestURI" toml:"requestURI"`
		StackTrace    *string `gorm:"type:text;null;comment:错误堆栈(如果发生了严重错误的话)" json:"stackTrace" yaml:"stackTrace" xml:"stackTrace" toml:"stackTrace"`
		IPAddr        *string `gorm:"type:varchar(255);null;comment:ip 地址" json:"ipAddr" yaml:"ipAddr" xml:"ipAddr" toml:"ipAddr"`
		IPSource      *string `gorm:"type:varchar(255);null;comment:ip 归属地" json:"ipSource" yaml:"ipSource" xml:"ipSource" toml:"ipSource"`
		Useragent     *string `gorm:"type:varchar(255);null;comment:客户端标识" json:"useragent" yaml:"useragent" xml:"useragent" toml:"useragent"`
		CreatedAt     int64   `gorm:"type:bigint;autoCreateTime:milli;comment:创建时间" json:"createdAt" yaml:"createdAt" xml:"createdAt" toml:"createdAt"`
	}
)

func (logRecord *LogInfo) BeforeCreate(tx *gorm.DB) (err error) {
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
		new(S3ObjectRecord),
		new(Category),
		new(Comment),
		new(PubSocialMedia),
		new(CloudFn),
		new(FileRecord),
		new(MenuGroup),
		new(MenuLink),
		new(FriendLink),
		new(BlockIPRecord),
		new(ServiceConf),
		new(CronJob),
	}
}

func GetExtraHelperMigrateTables() []any {
	return []any{
		new(LogInfo),
	}
}
