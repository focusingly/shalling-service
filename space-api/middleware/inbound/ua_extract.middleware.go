package inbound

import (
	"github.com/dineshgowda24/browser"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var _uaInjectMark = uuid.NewString()

type UADetail struct {
	ClientName string `json:"clientName"` // 用户使用的平台名称, 如: Chrome, Edge, Postman...
	Version    string `json:"version"`    // 平台版本
	Useragent  string `json:"useragent"`  // 原始的 UA 标识
	IsBot      bool   `json:"isBot"`      // 是否未搜索引擎爬虫
	OS         string `json:"os"`         // 系统类型
	OSVersion  string `json:"osVersion"`  // 系统版本
	IsMobile   bool   `json:"isMobile"`
}

func UseUseragentParserMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var uaDetail *UADetail

		userAgent := ctx.GetHeader("User-Agent")
		parsedUA, err := browser.NewBrowser(userAgent)
		if err != nil {
			uaDetail = &UADetail{
				ClientName: "",
				Version:    "",
				Useragent:  userAgent,
				IsBot:      false,
				OS:         "",
				OSVersion:  "",
				IsMobile:   false,
			}
		} else {
			uaDetail = &UADetail{
				ClientName: parsedUA.Name(),
				Version:    parsedUA.Version(),
				Useragent:  userAgent,
				IsBot:      parsedUA.Bot().IsBot(),
				OS:         parsedUA.Platform().Name(),
				OSVersion:  parsedUA.Platform().Version(),
				IsMobile:   parsedUA.Device().IsMobile(),
			}
		}

		ctx.Set(_uaInjectMark, uaDetail)
		ctx.Next()
	}
}

func GetUserAgentFromContext(ctx *gin.Context) *UADetail {
	if val, ok := ctx.Get(_uaInjectMark); !ok {
		return &UADetail{
			ClientName: "",
			Version:    "",
			Useragent:  ctx.Request.Header.Get("User-Agent"),
			IsBot:      false,
			OS:         "",
			OSVersion:  "",
		}
	} else {
		return val.(*UADetail)
	}
}
