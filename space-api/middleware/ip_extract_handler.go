package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func getClientIP(c *gin.Context) string {
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

func UseExtractIpHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("extract-real-ip", getClientIP(ctx))
	}
}

func TryGetRealIp(ctx *gin.Context) string {
	return ctx.GetString("extract-real-ip")
}
