package constants

type UserType = string

// 登录用户的类型
const (
	LocalUser  UserType = "local"
	GithubUser UserType = "github"
	GoogleUser UserType = "google"
)
