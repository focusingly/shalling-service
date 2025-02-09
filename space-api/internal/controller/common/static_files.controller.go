package common

import (
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"net/http"
	"os"
	"path"
	"space-api/conf"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func UseStaticFilesController(group *gin.RouterGroup) {
	fileGroup := group.Group("/static")
	const cacheControlExpired = time.Hour * 24 * 15 / time.Second
	// 设置静态资源缓存时间为 15 天
	var cacheControlHeader = fmt.Sprintf("public, max-age=%d, immutable", cacheControlExpired)

	appConf := conf.ProjectConf.GetAppConf()

	var staticPrefix = path.Clean(appConf.StaticDir)
	fileGroup.GET("*file", func(ctx *gin.Context) {
		// 获取请求的文件路径
		rawFileParam := ctx.Param("file")
		// 使用 path.Clean 清理路径, 防止路径跳跃比如 ../)
		cleanPath := path.Clean(rawFileParam)
		// 确保路径仍然在静态目录内
		fullPath := path.Join(appConf.StaticDir, cleanPath)
		if !strings.HasPrefix(fullPath, staticPrefix) {
			ctx.Status(http.StatusNotFound)
			return
		}

		// 获取文件信息
		fileInfo, err := os.Stat(fullPath)
		if err != nil {
			ctx.Status(http.StatusNotFound)
			return
		}

		// 协商缓存：设置 ETag 和 Last-Modified
		modifiedTime := fileInfo.ModTime()

		if match := ctx.GetHeader("If-Modified-Since"); match != "" {
			// 客户端缓存未过期，返回 304 Not Modified
			clientModifiedTime, err := time.Parse(http.TimeFormat, match)
			if err == nil && modifiedTime.Before(clientModifiedTime.Add(time.Second)) {
				ctx.Status(http.StatusNotModified)
				return
			}
		}

		// 生成本地 Etag
		etag := generateETag(fullPath)
		// 检查 Etag
		if match := ctx.GetHeader("If-None-Match"); match != "" && match == etag {
			// 客户端缓存未过期，返回 304 Not Modified
			ctx.Status(http.StatusNotModified)
			return
		}

		// 设置 ETag 头
		ctx.Header("ETag", etag)

		// 设置旧平台/代理兼容性过期时间
		expiresTime := time.Now().Add(time.Hour * 24 * 15).Format(http.TimeFormat)
		ctx.Header("Expires", expiresTime)

		// 设置 Cache-Control 头，启用强缓存
		ctx.Header("Cache-Control", cacheControlHeader)

		// 设置 Last-Modified 头
		ctx.Header("Last-Modified", modifiedTime.Format(http.TimeFormat))

		// 返回文件
		ctx.File(fullPath)
	})

}

// 生成 ETag，使用 CRC32 进行计算
func generateETag(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		return ""
	}
	defer file.Close()
	hasher := crc32.NewIEEE()
	_, err = file.WriteTo(hasher)
	if err != nil {
		return ""
	}
	etag := hex.EncodeToString(hasher.Sum(nil))
	return etag
}

// 检查 Referer 是否有效
func IsValidReferer(referer string, allowedRefererList []string) bool {
	for _, allowedReferer := range allowedRefererList {
		if referer == allowedReferer {
			return true
		}
	}
	return false
}
