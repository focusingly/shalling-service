package inbound

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var _ipMarked = "ip:" + uuid.NewString()

func getClientIPv4(c *gin.Context) string {
	// 优先尝试获取 Cloudflare 提供的真实客户端 IP
	if clientIP := c.Request.Header.Get("CF-Connecting-IP"); clientIP != "" {
		return clientIP
	}

	// 如果没有 CF-Connecting-IP，则检查 X-Forwarded-For
	if clientIP := c.Request.Header.Get("X-Forwarded-For"); clientIP != "" {
		ipList := strings.Split(clientIP, ",")
		return strings.TrimSpace(ipList[0]) // 获取第一个 IP 地址（客户端 IP）
	}

	// 默认返回 RemoteAddr（用于直接连接的客户端）
	return c.ClientIP()
}

func UseExtractIPv4Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(_ipMarked, getClientIPv4(ctx))

		ctx.Next()
	}
}

func TryGetRealIp(ctx *gin.Context) string {
	return ctx.GetString(_ipMarked)
}
