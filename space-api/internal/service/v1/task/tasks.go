package task

import (
	"bytes"
	"context"
	"fmt"
	"runtime/debug"
	"space-api/conf"
	"space-api/constants"
	"space-api/internal/service/v1"
	"space-api/internal/service/v1/monitor"
	"space-api/pack"
	"space-api/util/ptr"
	"space-domain/dao/biz"
	"space-domain/dao/extra"
	"space-domain/model"
	"sync"
	"text/template"
	"time"

	"github.com/robfig/cron/v3"
	"gopkg.in/gomail.v2"
)

type (
	// 自定义业务任务
	customTask struct {
		FuncName    string
		Task        func()
		Description string
	}

	taskRelation struct {
		entryID    cron.EntryID // 任务本身的 ID
		dbRecordID int64        // 任务在数据库中记录的 ID
	}

	taskMemRecord struct {
		sync.Mutex
		records []*taskRelation
	}

	taskWrapper struct {
		taskName string
		taskFunc func()
	}
)

var _ cron.Job = (*taskWrapper)(nil)

// Run implements cron.Job.
func (j *taskWrapper) Run() {
	now := time.Now().UnixMilli()
	defer func() {
		cost := time.Now().UnixMilli() - now
		var logInfo *model.LogInfo
		if fatalErr := recover(); fatalErr != nil {
			logInfo = &model.LogInfo{
				LogType:    string(constants.TaskExecute),
				Message:    j.taskName + ": 任务执行失败",
				Level:      string(constants.Fatal),
				CostTime:   cost,
				StackTrace: ptr.ToPtr(ptr.Bytes2String(debug.Stack())),
				CreatedAt:  now,
			}
		} else {
			logInfo = &model.LogInfo{
				LogType:   string(constants.TaskExecute),
				Message:   j.taskName + ": 任务执行成功",
				Level:     string(constants.Trace),
				CostTime:  cost,
				CreatedAt: now,
			}
		}
		extra.LogInfo.WithContext(context.Background()).Create(logInfo)
	}()

	j.taskFunc()
}

func newJobWrapper(taskName string, task func()) cron.Job {
	return &taskWrapper{
		taskName: taskName,
		taskFunc: task,
	}
}

func (s *taskMemRecord) appendJobRelation(relation *taskRelation) {
	s.Lock()
	defer s.Unlock()

	if relation == nil {
		panic("record can't nil")
	}
	s.records = append(s.records, &taskRelation{
		entryID:    relation.entryID,
		dbRecordID: relation.dbRecordID,
	})
}

// 已经注册的所有任务
// TODO 暂时只支持本地定义的方法
var RegisterJobs = []*customTask{
	// 清理冗余日志
	{
		FuncName: "clear_old_logs",
		Task: func() {
			logOp := extra.LogInfo
			sec := time.Now().AddDate(0, 0, -20).UnixMilli()
			logOp.WithContext(context.TODO()).
				Where(logOp.CreatedAt.Lte(sec)).
				Delete()
		},
		Description: "清空 20天 之前的冗余日志",
	},

	// 检查 cloudflare 账单并在产生扣费时发送邮件
	{
		FuncName: "check_cloudflare_billings",
		Task: func() {
			cfService := service.DefaultCloudflareService
			mailService := service.DefaultMailService
			appConf := conf.ProjectConf.GetAppConf()
			mailConf := conf.ProjectConf.GetMailConf()
			subs, err := cfService.GetExistsCost(context.TODO())

			// 获取账单信息失败
			if err != nil {
				t, err := template.New("cloud-flare-request-fault").Parse(string(pack.CheckBillingFaultTemplate))
				if err != nil {
					panic(err)
				}
				var bf = bytes.Buffer{}
				if e := t.Execute(
					&bf,
					map[string]any{
						"Link": "https://dash.cloudflare.com",
						"Time": "Asia/Shanghai " + time.Now().In(time.FixedZone("CST", 8*3600)).Format("2006-01-02 15:04:05"),
					}); e != nil {
					panic(e)
				}

				msg := gomail.NewMessage()
				msg.SetHeader("From", mailConf.Username)
				msg.SetHeader("To", appConf.NotifyEmail)
				msg.SetHeader("Subject", "获取 Cloudflare 账单信息失败, 请留意")
				msg.SetBody("text/html", bf.String())
				e := mailService.SendEmail(msg)
				if e != nil {
					panic(e)
				}

				return
			}

			if len(subs) != 0 {
				t, err := template.New("cloud-flare-billing").Parse(string(pack.BillingSubsCostTemplate))
				if err != nil {
					panic(err)
				}
				var bf = bytes.Buffer{}
				if e := t.Execute(
					&bf,
					map[string]any{
						"Subs": subs,
						"Time": "Asia/Shanghai " + time.Now().In(time.FixedZone("CST", 8*3600)).Format("2006-01-02 15:04:05"),
					}); e != nil {
					panic(e)
				}

				msg := gomail.NewMessage()
				msg.SetHeader("From", mailConf.Username)
				msg.SetHeader("To", appConf.NotifyEmail)
				msg.SetHeader("Subject", "Cloudflare 的订阅产生了费用")
				msg.SetBody("text/html", bf.String())
				e := mailService.SendEmail(msg)
				if e != nil {
					panic(e)
				}
			}

		},
		Description: "检查 cloudflare 是否产生了扣费行为(比如 r2 超出免费额度, 并发送邮件进行提醒)",
	},

	// 清理久远的访客日志
	{
		FuncName: "clear_old_uv_data",
		Task: func() {
			uvOp := biz.UVStatistic
			oldDate := time.Now().AddDate(0, 0, 30).UnixMilli()
			_, err := uvOp.WithContext(context.TODO()).
				Where(uvOp.CreatedAt.Lt(oldDate)).
				Delete()

			if err != nil {
				panic(err)
			}
		},
		Description: "清空前 30天 之前的 UV 记录(根据 unix 时间戳判断)",
	},

	// 同步缓存中的文章浏览量到数据库当中
	{
		FuncName: "sync_post_pv",
		Task: func() {
			postService := service.DefaultPostService
			err := postService.SyncAllPostViews(context.TODO())
			if err != nil {
				panic(err)
			}
		},
		Description: "同步文章的 PV 数据到数据库当中",
	},

	// 检查系统状态, 并在负载高时发送预警
	{
		FuncName: "check_system_load",
		Task: func() {
			mailService := service.DefaultMailService
			appConf := conf.ProjectConf.GetAppConf()
			mailConf := conf.ProjectConf.GetMailConf()
			monitorService := monitor.DefaultMonitorService
			pefStatus, err := monitorService.GetStatus()

			if err != nil {
				t, err := template.New("get-system-load-fault").Parse(string(pack.SystemLoadFaultTemplate))
				if err != nil {
					panic(err)
				}
				var bf = bytes.Buffer{}
				if e := t.Execute(
					&bf,
					map[string]any{
						"Time": "Asia/Shanghai " + time.Now().
							In(time.FixedZone("CST", 8*3600)).
							Format("2006-01-02 15:04:05"),
					}); e != nil {
					panic(e)
				}

				msg := gomail.NewMessage()
				msg.SetHeader("From", mailConf.Username)
				msg.SetHeader("To", appConf.NotifyEmail)
				msg.SetHeader("Subject", "获取系统负载信息失败, 请留意")
				msg.SetBody("text/html", bf.String())
				e := mailService.SendEmail(msg)
				if e != nil {
					panic(err)
				}
			} else {
				coreNum := len(pefStatus.CPUUsagePercent)
				var usage float64
				for _, p := range pefStatus.CPUUsagePercent {
					usage += p
				}

				if usage >= float64(coreNum*80) || pefStatus.Memory.UsedPercent >= 80 {
					fmt.Println(coreNum)
					t, err := template.New("high-system-load-alert").Parse(string(pack.SystemLoadAlertTemplate))
					if err != nil {
						panic(err)
					}
					var bf = bytes.Buffer{}

					infos := [][]string{
						{
							"CPU",
							fmt.Sprintf("%d 核, 总峰值: %d%%", coreNum, coreNum*100),
							fmt.Sprintf("%.2f%%", usage),
							fmt.Sprintf("%.2f%%", float64(coreNum*100)-usage),
							fmt.Sprintf("%.2f%%", usage/float64(coreNum*100)*100),
						},
						{
							"内存",
							fmt.Sprintf("总容量 %.2fmb", float64(pefStatus.Memory.Total)/1024/1024),
							fmt.Sprintf("%.2fmb", float64(pefStatus.Memory.Used)/1024/1024),
							fmt.Sprintf("%.2fmb", float64(pefStatus.Memory.Available)/1024/1024),
							fmt.Sprintf("%.2f%%", pefStatus.Memory.UsedPercent),
						},
						{
							"磁盘",
							fmt.Sprintf("总容量 %2fmb", float64(pefStatus.DiskUsage.Total)/1024/1024),
							fmt.Sprintf("%.2fmb", float64(pefStatus.DiskUsage.Used)/1024/1024),
							fmt.Sprintf("%.2fmb", float64(pefStatus.DiskUsage.Free)/1024/1024),
							fmt.Sprintf("%.2f%%", pefStatus.DiskUsage.UsedPercent),
						},
					}

					// 除了模板
					if e := t.Execute(
						&bf,
						map[string]any{
							"Time": "Asia/Shanghai " + time.Now().
								In(time.FixedZone("CST", 8*3600)).
								Format("2006-01-02 15:04:05"),
							"Infos": infos,
						}); e != nil {
						panic(e)
					}

					msg := gomail.NewMessage()
					msg.SetHeader("From", mailConf.Username)
					msg.SetHeader("To", appConf.NotifyEmail)
					msg.SetHeader("Subject", "当前系统负载较高, 请留意")
					msg.SetBody("text/html", bf.String())
					e := mailService.SendEmail(msg)

					if e != nil {
						panic(err)
					}
				}
			}
		},
		Description: "检查系统负载, 并在 CPU 整体负载h >= 80 % 或内存使用率 >=80 的情况下发送邮件预警",
	},
}
