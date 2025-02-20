package common

import (
	"embed"
	"fmt"
	"hash/crc32"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// MapEmbedFS2Crc32 遍历 embed.FS 并返回文件路径到哈希值的映射
func MapEmbedFS2Crc32(embedFS *embed.FS) (map[string]string, error) {
	result := make(map[string]string)

	// 遍历文件系统
	err := fs.WalkDir(embedFS, ".", func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录
		if d.IsDir() {
			return nil
		}

		// 读取文件内容
		content, err := embedFS.ReadFile(filePath)
		if err != nil {
			return err
		}

		// 计算 CRC32 哈希
		checksum := crc32.ChecksumIEEE(content)
		hash := strconv.FormatUint(uint64(checksum), 16)

		// 规范化文件路径（使用正斜杠）
		normalizedPath := path.Clean(filePath)

		// 存储到映射中
		result[normalizedPath] = hash

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// CreateEmbedSpaAppHandler creates a handler function to serve an embedded Single Page Application (SPA).
//
// Parameters:
//   - urlPrefix: The default access path prefix.
//   - staticPath: The path to the static file system.
//   - skipApiPrefix: The prefix to skip for API routes.
//   - fsLike: The embedded file system.
//   - expireTime: The duration for which the cache is valid.
//
// Returns:
//   - gin.HandlerFunc: The handler function to serve the SPA.
func CreateEmbedSpaAppHandler(
	urlPrefix,
	staticPath,
	skipApiPrefix string,
	fsLike *embed.FS,
	expireTime time.Duration,
) gin.HandlerFunc {
	mapping, err := MapEmbedFS2Crc32(fsLike)
	if err != nil {
		log.Fatal("embed fs exists error: ", err)
	}
	skipApiPrefix = strings.TrimPrefix(skipApiPrefix, "/")
	pubCacheControlHeader := fmt.Sprintf("public, max-age=%d, immutable", expireTime/time.Second)
	modifiedTime := time.Now()

	return func(ctx *gin.Context) {
		// 去除 URL 前缀
		pathArg := strings.TrimPrefix(ctx.Request.URL.Path, urlPrefix)
		// 构建完整的文件路径
		fullFilePath := filepath.Join(staticPath, pathArg)
		// 对于前端的 HTML5 history 路由刷新采取回退路径
		fallbackPath := filepath.Join(staticPath, "index.html")

		// 检查文件是否存在于目录当中
		fileCrc32, ok := mapping[fullFilePath]
		if !ok {
			// 如果文件不存在，尝试检查是否为API路由
			if strings.HasPrefix(pathArg, skipApiPrefix) {
				ctx.Next()
				return
			}
			ctx.Header("Content-Type", gin.MIMEHTML)
			// 对于其他不存在的路径，返回 index.html，但不设置缓存
			http.ServeFileFS(
				ctx.Writer,
				ctx.Request,
				fsLike,
				fallbackPath,
			)
			return
		}

		if match := ctx.GetHeader("If-Modified-Since"); match != "" {
			// 客户端缓存未过期，返回 304 Not Modified
			clientModifiedTime, err := time.Parse(http.TimeFormat, match)
			if err == nil && modifiedTime.Before(clientModifiedTime.Add(time.Second)) {
				ctx.Status(http.StatusNotModified)
				return
			}
		}

		// 检查 Etag
		if match := ctx.GetHeader("If-None-Match"); match != "" && match == fileCrc32 {
			// 客户端缓存未过期，返回 304 Not Modified
			ctx.Status(http.StatusNotModified)
			return
		}

		// 设置 ETag 头
		ctx.Header("ETag", fileCrc32)
		// 设置旧平台/代理兼容性过期时间
		expiresTime := time.Now().Add(expireTime).Format(http.TimeFormat)
		ctx.Header("Expires", expiresTime)
		// 设置 Cache-Control 头，启用强缓存
		ctx.Header("Cache-Control", pubCacheControlHeader)
		// 设置 Last-Modified 头
		ctx.Header("Last-Modified", modifiedTime.Format(http.TimeFormat))

		// 设置文件标识
		if m := mime.TypeByExtension(path.Ext(fullFilePath)); m != "" {
			ctx.Header("Content-Type", m)
		}
		// 提供静态文件服务
		http.ServeFileFS(
			ctx.Writer,
			ctx.Request,
			fsLike,
			fullFilePath,
		)
	}
}
