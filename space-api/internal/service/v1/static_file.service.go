package service

import (
	"context"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"mime"
	"net/http"
	"os"
	"path"
	"space-api/conf"
	"space-api/constants"
	"space-api/util/performance"
	"space-domain/dao/biz"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type _staticFileService struct {
	cache *performance.JsonCache
}

const cacheControlExpired = time.Hour * 24 * 15 / time.Second

var (
	DefaultStaticFileService = &_staticFileService{
		cache: performance.NewCache(constants.MB * 4),
	}
	// 设置静态资源缓存时间为 15 天
	cacheControlHeader = fmt.Sprintf("public, max-age=%d, immutable", cacheControlExpired)
	appConf            = conf.ProjectConf.GetAppConf()
	staticDirPrefix    = path.Clean(appConf.StaticDir)
)

func (s *_staticFileService) ExposeInnerCacher() *performance.JsonCache {
	return s.cache
}

func (s *_staticFileService) IsPubVisible(locationName string) bool {
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
		s.cache.Set(locationName, &performance.Empty{}, performance.Second(cacheControlExpired))
		return true
	}
}

func (s *_staticFileService) InvalidateAllVisibleCache(locationName string) {
	s.cache.ClearAll()
}

// HandlePubVisit 处理公共的访问的文件服务, 并添加缓存标识, 限制未公开文件访问
func (s *_staticFileService) HandlePubVisit(rawFileParam string, ctx *gin.Context) {
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
	expiresTime := time.Now().Add(time.Hour * 24 * 15).Format(http.TimeFormat)
	ctx.Header("Expires", expiresTime)

	// 设置 Cache-Control 头，启用强缓存
	ctx.Header("Cache-Control", cacheControlHeader)

	// 设置 Last-Modified 头
	ctx.Header("Last-Modified", modifiedTime.Format(http.TimeFormat))

	// 设置文件标识
	if m := mime.TypeByExtension(path.Ext(rawFileParam)); m != "" {
		ctx.Header("Content-Type", m)
	}
	// 返回文件
	ctx.File(fullPath)
}

func (s *_staticFileService) InvalidateCache(locationName string) {
	s.cache.Delete(locationName)
}

// HandleAllVisit 处理所有的静态资源, 并且不设置缓存策略(包括未公开的, 适合于管理员使用)
func (s *_staticFileService) HandleAllVisit(rawFileParam string, ctx *gin.Context) {
	// 获取请求的文件路径
	// 使用 path.Clean 清理路径, 防止路径跳跃比如 ../)
	cleanPath := path.Clean(rawFileParam)
	// 确保路径仍然在静态目录内
	fullPath := path.Join(appConf.StaticDir, cleanPath)
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
	// 设置文件标识
	if m := mime.TypeByExtension(path.Ext(rawFileParam)); m != "" {
		ctx.Header("Content-Type", m)
	}
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
	for _, allowedReferer := range allowedRefererList {
		if referer == allowedReferer {
			return true
		}
	}
	return false
}
