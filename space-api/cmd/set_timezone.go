package cmd

import (
	"log"
	"space-api/conf"
	"time"
)

func setTimeZone() {
	appConf := conf.ProjectConf.GetAppConf()
	// TODO 时区设置, 暂不设置, 统一全部直接使用 unix 时间戳; 数据格式化由客户端自己解析
	// 定时任务的时区直接遵循服务器所设置的时区
	if appConf.ServerTimezone != "" {
		if tz, err := time.LoadLocation(appConf.ServerTimezone); err != nil {
			log.Fatal("获取时区失败: ", err)
		} else {
			time.Local = tz
		}
	}
}
