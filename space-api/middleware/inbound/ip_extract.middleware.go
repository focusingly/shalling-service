package inbound

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var _ipMarked = "ip:" + uuid.NewString()

func getIpAddress(ctx *gin.Context) string {
	// 优先尝试获取 Cloudflare 提供的真实客户端 IP
	if clientIP := ctx.Request.Header.Get("CF-Connecting-IP"); clientIP != "" {
		return clientIP
	}

	// 如果没有 CF-Connecting-IP，则检查 X-Forwarded-For
	// 注: 如果使用反向代理等, 记得覆盖
	if clientIP := ctx.Request.Header.Get("X-Forwarded-For"); clientIP != "" {
		ipList := strings.Split(clientIP, ",")
		return strings.TrimSpace(ipList[0]) // 获取第一个 IP 地址（客户端 IP）
	}

	// 默认返回 RemoteAddr（用于直接连接的客户端）
	return ctx.ClientIP()
}

func UseExtractIPv4Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(_ipMarked, getIpAddress(ctx))
		ctx.Next()
	}
}

func GetRealIpWithContext(ctx *gin.Context) string {
	if ip, ok := ctx.Get(_ipMarked); !ok {
		return ctx.ClientIP()
	} else {
		return ip.(string)
	}
}
