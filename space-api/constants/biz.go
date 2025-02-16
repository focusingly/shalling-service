package constants

// UserType 登录用户的类型
type UserType = string

const (
	LocalUser  UserType = "local"
	Admin      UserType = "admin"
	GithubUser UserType = "github"
	GoogleUser UserType = "google"
)

// LogLevel 自定义收集的日志级别
type LogLevel = string

const (
	Trace UserType = "trace"
	Info  UserType = "info"
	Warn  UserType = "warn"
	Error UserType = "error"
	Fatal UserType = "fatal"
)

// LogType 自定义的日志类型
type LogType = string

const (
	APIRequest  LogType = "apiRequest"
	TaskExecute LogType = "taskExecute"
)
