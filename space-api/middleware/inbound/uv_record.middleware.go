package inbound

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"space-api/util"
	"space-api/util/ip"
	"space-api/util/performance"
	"space-domain/dao/biz"
	"space-domain/model"
	"strings"
	"sync"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// UVManager 负责管理UV统计的核心逻辑
type UVManager struct {
	ExcludePaths []string
	CookieMaxAge int
	mu           sync.Mutex
	dailyCache   map[string]map[string]bool // date -> visitorHash -> exists
}

// NewUVManager 创建一个新的UV管理器
func NewUVManager(excludePath ...string) *UVManager {
	return &UVManager{
		ExcludePaths: excludePath,
		CookieMaxAge: 86400, // 24小时
		dailyCache:   make(map[string]map[string]bool),
	}
}

// CreateUVMiddleware 创建一个Gin中间件来统计UV
func (m *UVManager) CreateUVMiddleware() gin.HandlerFunc {
	var searcher = ip.GetIpSearcher()

	return func(ctx *gin.Context) {
		// 跳过排除路径
		path := ctx.Request.URL.Path
		for _, excludePath := range m.ExcludePaths {
			if strings.HasPrefix(path, excludePath) {
				ctx.Next()
				return
			}
		}

		// 获取或创建会话
		session := sessions.Default(ctx)
		sessionID := session.Get("session_id")

		if sessionID == nil {
			// 生成新的会话 ID
			sessionID = m.generateSessionID(ctx.Request)
			session.Set("session_id", sessionID)
			// 设置过期时间
			session.Options(sessions.Options{
				MaxAge:   m.CookieMaxAge,
				Secure:   true,
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
			})
			session.Save()
		}

		// 生成访客哈希，结合IP、UserAgent和SessionID
		realIP := GetRealIpWithContext(ctx)
		ipSource, _ := searcher.SearchByStr(realIP)
		visitorHash := m.generateVisitorHash(realIP, ctx.Request.UserAgent(), sessionID.(string))

		ua := GetUserAgentFromContext(ctx)
		// 获取当前时间
		now := time.Now()
		timestamp := now.UnixMilli()
		today := now.Format("2006-01-02")
		// 检查今天是否已记录该访客
		if !m.isVisitorRecordedToday(visitorHash, today) {
			// 记录新访客
			m.addNewRecord(
				visitorHash,
				realIP,
				ipSource,
				ua,
				timestamp,
				today,
			)
		}

		ctx.Next()
	}
}

// 检查访客今天是否已被记录
func (m *UVManager) isVisitorRecordedToday(visitorHash, today string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 检查内存缓存
	if visitors, exists := m.dailyCache[today]; exists {
		if _, visited := visitors[visitorHash]; visited {
			return true
		}
	} else {
		// 初始化当天的缓存
		m.dailyCache[today] = make(map[string]bool)

		// 从数据库加载当天已记录的访客
		var hashes []string
		uvOp := biz.UVStatistic
		uvOp.
			WithContext(context.TODO()).
			Where(uvOp.VisitDate.Eq(today)).
			Pluck(uvOp.VisitorHash, &hashes)

		for _, hash := range hashes {
			m.dailyCache[today][hash] = true
		}

		// 清理旧缓存（保留最近7天）
		m.cleanOldCache(today)

		if _, visited := m.dailyCache[today][visitorHash]; visited {
			return true
		}
	}

	return false
}

// 记录新的访客
func (m *UVManager) addNewRecord(visitorHash, ip, ipSource string, ua *UADetail, timestamp int64, date string) {
	func() {
		// 更新内存缓存
		m.mu.Lock()
		defer m.mu.Unlock()

		if _, exists := m.dailyCache[date]; !exists {
			m.dailyCache[date] = make(map[string]bool)
		}
		m.dailyCache[date][visitorHash] = true
	}()

	// 存储到数据库
	uvStat := model.UVStatistic{
		VisitorHash: visitorHash,
		IP:          ip,
		IPSource:    ipSource,
		UserAgent:   ua.Useragent,
		ClientName:  ua.ClientName,
		IsMobile:    util.TernaryExpr(ua.IsMobile, 1, 0),
		LikeBot:     util.TernaryExpr(ua.IsBot, 1, 0),
		OS:          ua.OS,
		VisitTime:   timestamp,
		VisitDate:   date,
	}

	// 使用 gopool 去异步执行
	performance.DefaultTaskRunner.Go(func() {
		biz.UVStatistic.
			WithContext(context.TODO()).
			Create(&uvStat)
	})
}

// 清理旧缓存
func (m *UVManager) cleanOldCache(today string) {
	todayDate, _ := time.Parse("2006-01-02", today)

	for dateStr := range m.dailyCache {
		cacheDate, _ := time.Parse("2006-01-02", dateStr)
		// 清理 7 天前的缓存
		if todayDate.Sub(cacheDate) > 7*24*time.Hour {
			delete(m.dailyCache, dateStr)
		}
	}
}

// 生成会话ID
func (m *UVManager) generateSessionID(req *http.Request) string {
	timestamp := time.Now().Unix()
	data := fmt.Sprintf("%s-%s-%s-%d",
		req.RemoteAddr,
		req.UserAgent(),
		req.Header.Get("Accept-Language"),
		timestamp,
	)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// 生成访客哈希
func (m *UVManager) generateVisitorHash(ip, userAgent, sessionID string) string {
	data := fmt.Sprintf("%s-%s-%s", ip, userAgent, sessionID)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
