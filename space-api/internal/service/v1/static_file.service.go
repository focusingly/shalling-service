package service

import (
	"context"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"net/http"
	"os"
	"path"
	"space-api/conf"
	"space-api/util/performance"
	"space-domain/dao/biz"
	"strings"
	"time"

	"slices"

	"github.com/gin-gonic/gin"
)

type (
	IStaticFileService interface {
		IsPubVisible(locationName string) bool
		InvalidateAllVisibleCache(locationName string)
		InvalidateCache(locationName string)
		HandlePubVisit(rawFileParam string, ctx *gin.Context)
		HandleAnyVisit(rawFileParam string, ctx *gin.Context)
	}

	staticFileServiceImpl struct {
		cache performance.CacheGroupInf
	}
)

const (
	pubCacheControlExpired   = time.Hour * 24 * 15
	adminCacheControlExpired = time.Minute * 1
)

var (
	_ IStaticFileService = (*staticFileServiceImpl)(nil)

	DefaultStaticFileService IStaticFileService = &staticFileServiceImpl{
		cache: performance.DefaultJsonCache.Group("static-file"),
	}
	// 公开静态资源缓存时间为 15 天
	pubCacheControlHeader = fmt.Sprintf("public, max-age=%d, immutable", pubCacheControlExpired/time.Second)
	// 管理员静态资源缓存时间为 2 分钟
	adminCacheControlHeader = fmt.Sprintf("public, max-age=%d, immutable", adminCacheControlExpired/time.Second)
	appConf                 = conf.ProjectConf.GetAppConf()
	staticDirPrefix         = path.Clean(appConf.StaticDir)
)

func (s *staticFileServiceImpl) ExposeInnerCacher() performance.CacheGroupInf {
	return s.cache
}

func (s *staticFileServiceImpl) IsPubVisible(locationName string) bool {
	if e := s.cache.Get(locationName, &performance.Empty{}); e == nil {
		return true
	}
	fileOp := biz.FileRecord
	// 匹配是否存在可用公开文件
	_, err := fileOp.WithContext(context.Background()).
		Where(
			fileOp.LocalLocation.Eq(locationName),
			fileOp.PubAvailable.Neq(0),
			fileOp.Hide.Eq(0),
		).
		Take()
	if err != nil {
		return false
	} else {
		s.cache.Set(locationName, &performance.Empty{}, pubCacheControlExpired)
		return true
	}
}

func (s *staticFileServiceImpl) InvalidateAllVisibleCache(locationName string) {
	s.cache.ClearAll()
}

// HandlePubVisit 处理公共的访问的文件服务, 并添加缓存标识, 限制未公开文件访问
func (s *staticFileServiceImpl) HandlePubVisit(rawFileParam string, ctx *gin.Context) {
	// 获取请求的文件路径
	// 使用 path.Clean 清理路径, 防止路径跳跃比如 ../)
	cleanFileName := path.Clean(rawFileParam)
	// 确保路径仍然在静态目录内
	fullPath := path.Join(appConf.StaticDir, cleanFileName)
	if !strings.HasPrefix(fullPath, staticDirPrefix) {
		ctx.Status(http.StatusNotFound)
		return
	}

	// 限制为只允许访问公共资源
	if !s.IsPubVisible(cleanFileName[1:]) {
		ctx.Status(http.StatusNotFound)
		return
	}

	// 获取文件信息
	fileInfo, err := os.Stat(fullPath)
	if err != nil || fileInfo.IsDir() {
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
	etag := generateCRC32ETag(fullPath)
	// 检查 Etag
	if match := ctx.GetHeader("If-None-Match"); match != "" && match == etag {
		// 客户端缓存未过期，返回 304 Not Modified
		ctx.Status(http.StatusNotModified)
		return
	}

	// 设置 ETag 头
	ctx.Header("ETag", etag)

	// 设置旧平台/代理兼容性过期时间
	expiresTime := time.Now().Add(pubCacheControlExpired).Format(http.TimeFormat)
	ctx.Header("Expires", expiresTime)

	// 设置 Cache-Control 头，启用强缓存
	ctx.Header("Cache-Control", pubCacheControlHeader)

	// 设置 Last-Modified 头
	ctx.Header("Last-Modified", modifiedTime.Format(http.TimeFormat))

	ctx.File(fullPath)
}

func (s *staticFileServiceImpl) InvalidateCache(locationName string) {
	s.cache.Delete(locationName)
}

// HandleAnyVisit 处理所有的静态资源(包括未公开的, 适合于管理员使用), 并且设置较短的缓存策略
func (s *staticFileServiceImpl) HandleAnyVisit(rawFileParam string, ctx *gin.Context) {
	// 获取请求的文件路径
	// 使用 path.Clean 清理路径, 防止路径跳跃比如 ../)
	cleanFileName := path.Clean(rawFileParam)
	// 确保路径仍然在静态目录内
	fullPath := path.Join(appConf.StaticDir, cleanFileName)
	if !strings.HasPrefix(fullPath, staticDirPrefix) {
		ctx.Status(http.StatusNotFound)
		return
	}

	// 获取文件信息
	fileInfo, err := os.Stat(fullPath)
	if err != nil || fileInfo.IsDir() {
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
	etag := generateCRC32ETag(fullPath)
	// 检查 Etag
	if match := ctx.GetHeader("If-None-Match"); match != "" && match == etag {
		// 客户端缓存未过期，返回 304 Not Modified
		ctx.Status(http.StatusNotModified)
		return
	}

	// 设置 ETag 头
	ctx.Header("ETag", etag)

	// 设置旧平台/代理兼容性过期时间
	expiresTime := time.Now().Add(adminCacheControlExpired).Format(http.TimeFormat)
	ctx.Header("Expires", expiresTime)

	// 设置 Cache-Control 头，启用强缓存
	ctx.Header("Cache-Control", adminCacheControlHeader)

	// 设置 Last-Modified 头
	ctx.Header("Last-Modified", modifiedTime.Format(http.TimeFormat))

	// 返回文件
	ctx.File(fullPath)
}

// 生成 ETag，使用 CRC32 进行计算
func generateCRC32ETag(filePath string) string {
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
	return slices.Contains(allowedRefererList, referer)
}
