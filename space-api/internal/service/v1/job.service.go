package service

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"runtime/debug"
	"slices"
	"space-api/conf"
	"space-api/constants"
	"space-api/dto"
	"space-api/pack"
	"space-api/util"
	"space-api/util/arr"
	"space-api/util/id"
	"space-api/util/ptr"
	"space-domain/dao/biz"
	"space-domain/dao/extra"
	"space-domain/model"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"gopkg.in/gomail.v2"
)

type _taskService struct {
	mu       sync.Mutex
	initDone bool
}

var (
	DefaultTaskService = &_taskService{}
)

type (
	CustomTask struct {
		FuncName    string
		Task        func()
		Description string
	}

	_taskRelationRecord struct {
		entryID    cron.EntryID
		dbRecordID int64
	}

	_memSyncRecord struct {
		sync.Mutex
		records []*_taskRelationRecord
	}

	_jobWrapper struct {
		taskName string
		task     func()
	}
)

func createJobWrapper(taskName string, task func()) cron.Job {
	return &_jobWrapper{
		taskName: taskName,
		task:     task,
	}
}

func (s *_memSyncRecord) appendJobID(rc *_taskRelationRecord) {
	s.Lock()
	defer s.Unlock()
	if rc == nil {
		panic("record can't nil")
	}
	s.records = append(s.records, &_taskRelationRecord{
		entryID:    rc.entryID,
		dbRecordID: rc.dbRecordID,
	})
}

var _ cron.Job = (*_jobWrapper)(nil)

// Run implements cron.Job.
func (j *_jobWrapper) Run() {
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

	j.task()
}

var (
	_singleScheduler = cron.New()
	_memJobRecords   = &_memSyncRecord{}
	// 已经注册的所有任务
	// TODO 暂时只支持本地定义的方法
	RegisterJobs = []*CustomTask{
		{
			FuncName: "clear_old_logs",
			Task: func() {
				logOp := extra.LogInfo
				sec := time.Now().AddDate(0, 0, -19).UnixMilli()
				logOp.WithContext(context.TODO()).
					Where(logOp.CreatedAt.Lte(sec)).
					Delete()
			},
			Description: "清空 10天 之前的日志",
		},
		{
			FuncName: "check_cloudflare_billings",
			Task: func() {
				cfService := DefaultCloudflareService
				mailService := DefaultMailService
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
	}
)

// ResumeFromDatabase 从数据库中恢复执行, 如果已经执行过了, 那么会直接返回
func (servicePtr *_taskService) ResumeFromDatabase() (err error) {
	if servicePtr.initDone {
		return
	}

	servicePtr.mu.Lock()
	defer servicePtr.mu.Unlock()

	// 从数据库中恢复任务记录
	jobs, err := biz.CronJob.WithContext(context.TODO()).Find()
	if err != nil {
		return
	}

	for _, job := range jobs {
		index := slices.IndexFunc(RegisterJobs, func(j *CustomTask) bool {
			return j.FuncName == job.JobFuncName
		})
		if index == -1 {
			return fmt.Errorf("非注册的定时任务: %s", job.JobFuncName)
		}
		if job.Enable != 0 {
			entryID, err := _singleScheduler.AddFunc(job.CronExpr, RegisterJobs[index].Task)
			if err != nil {
				return err
			}
			_memJobRecords.appendJobID(&_taskRelationRecord{
				entryID:    entryID,
				dbRecordID: job.ID,
			})
		}
	}

	_singleScheduler.Start()
	servicePtr.initDone = true

	return
}

func (*_taskService) GetAvailableJobList(req *dto.GetAvailableJobListReq, ctx *gin.Context) (resp *dto.GetAvailableJobListResp, err error) {
	return &dto.GetAvailableJobListResp{
		List: arr.MapSlice(RegisterJobs, func(_ int, job *CustomTask) *dto.RegisteredJob {
			return &dto.RegisteredJob{
				JobName:     job.FuncName,
				Description: job.Description,
			}
		}),
	}, nil
}

func (*_taskService) CreateOrUpdateNewJob(req *dto.CreateOrUpdateJobReq, ctx *gin.Context) (resp *dto.CreateOrUpdateJobResp, err error) {
	err = biz.Q.Transaction(func(tx *biz.Query) error {
		// 加锁操作
		_memJobRecords.Lock()
		defer _memJobRecords.Unlock()

		jobTx := tx.CronJob
		// 前置检查参数
		// TODO 暂时只允许注册任务

		// 任务索引所在的位置
		newTaskIndex := slices.IndexFunc(RegisterJobs, func(j *CustomTask) bool {
			return j.FuncName == req.JobFuncName
		})
		if newTaskIndex == -1 {
			return fmt.Errorf("未注册的任务: %s", req.JobFuncName)
		}
		newTaskRecord := RegisterJobs[newTaskIndex]
		// 检查 cron 表达式是否受到支持
		parsedCronExpr, err := cron.ParseStandard(req.CronExpr)
		if err != nil {
			return fmt.Errorf("cron 表达式解析错误: "+err.Error(), err)
		}

		// 查找已经存在的任务
		findTask, e := jobTx.WithContext(ctx).
			Where(jobTx.ID.Eq(req.DBRecordID)).
			Take()
		// 不存在任务, 进行新建
		if e != nil {
			newJobId := id.GetSnowFlakeNode().Generate().Int64()
			e := jobTx.WithContext(ctx).Create(
				&model.CronJob{
					BaseColumn: model.BaseColumn{
						ID: newJobId,
					},
					JobFuncName: req.JobFuncName,
					CronExpr:    req.CronExpr,
					Status:      util.TernaryExpr(req.Enable == 0, "inactive", "active"),
					Enable:      req.Enable,
					Mark:        req.Mark,
				},
			)
			if e != nil {
				return e
			}
			if req.Enable != 0 {
				entryID := _singleScheduler.Schedule(
					parsedCronExpr,
					createJobWrapper(newTaskRecord.FuncName, newTaskRecord.Task),
				)
				// 同步到内存列表当中
				_memJobRecords.records = append(_memJobRecords.records, &_taskRelationRecord{
					entryID:    entryID,
					dbRecordID: newJobId,
				})
			}
		} else {
			// 任务存在, 进行更新
			update := &model.CronJob{
				BaseColumn:  findTask.BaseColumn,
				JobFuncName: req.JobFuncName,
				CronExpr:    req.CronExpr,
				Status:      util.TernaryExpr(req.Enable == 0, "inactive", "active"),
				Enable:      req.Enable,
				Mark:        req.Mark,
			}
			// 更新数据库中的记录
			_, e := jobTx.WithContext(ctx).
				Select(
					jobTx.Hide,
					jobTx.JobFuncName,
					jobTx.CronExpr,
					jobTx.Status,
					jobTx.Enable,
					jobTx.Mark,
				).
				Where(jobTx.ID.Eq(findTask.ID)).
				Updates(update)
			if e != nil {
				return e
			}

			// 处理内存中的记录
			existsTaskIndex := slices.IndexFunc(_memJobRecords.records, func(rc *_taskRelationRecord) bool {
				return rc.dbRecordID == findTask.ID
			})

			// 此前任务就已经在运行的状态
			if existsTaskIndex != -1 {
				// 存在运行的任务
				exitsTask := _memJobRecords.records[existsTaskIndex]

				// 取消掉旧任务
				_singleScheduler.Remove(exitsTask.entryID)

				// 构建新的任务列表 -> 移除掉旧任务
				updates := slices.Delete(_memJobRecords.records, existsTaskIndex, existsTaskIndex+1)

				// 需要启动新的任务情况
				if req.Enable != 0 {
					// 设置新的任务
					entryID := _singleScheduler.Schedule(parsedCronExpr, createJobWrapper(newTaskRecord.FuncName, newTaskRecord.Task))
					// 追加任务
					updates = append(updates, &_taskRelationRecord{
						entryID:    entryID,
						dbRecordID: findTask.ID,
					})
				}

				// 整体替换任务记录
				_memJobRecords.records = updates
			} else {
				// 此前并无已经存在运行的任务

				if req.Enable != 0 {
					// 设置新的任务
					entryID := _singleScheduler.Schedule(parsedCronExpr, createJobWrapper(newTaskRecord.FuncName, newTaskRecord.Task))

					// 追加任务记录
					_memJobRecords.records = append(_memJobRecords.records, &_taskRelationRecord{
						entryID:    entryID,
						dbRecordID: findTask.ID,
					})
				}
			}

		}
		return nil
	})

	if err != nil {
		err = util.CreateBizErr("创建/更新任务失败: "+err.Error(), err)
		return
	}

	resp = &dto.CreateOrUpdateJobResp{}
	return
}

// 立马执行任务
func (*_taskService) RunJobImmediately(req *dto.RunJobReq, ctx *gin.Context) (resp *dto.RunJobResp, err error) {
	// 上锁
	_memJobRecords.Lock()
	defer _memJobRecords.Unlock()

	jobTx := biz.CronJob
	job, err := jobTx.WithContext(ctx).
		Where(jobTx.ID.Eq(req.JobID)).
		Take()
	if err != nil {
		err = util.CreateBizErr("任务不存在", err)
		return
	}
	// 只允许已经启动的任务
	if job.Enable == 0 {
		err = util.CreateBizErr("当前任务未启用", fmt.Errorf("current job not enable"))
		return
	}

	// 查找是否在活动列表当中
	index := slices.IndexFunc(_memJobRecords.records, func(r *_taskRelationRecord) bool {
		return r.dbRecordID == job.ID
	})
	if index == -1 {
		err = util.CreateBizErr("未找到活动任务, 可能存在问题", fmt.Errorf("could not find ant matched job in memory"))
		return
	}
	entry := _singleScheduler.Entry(_memJobRecords.records[index].entryID)
	if entry.Job == nil {
		err = util.CreateBizErr("未找到活动任务, 可能存在问题", fmt.Errorf("could not find ant matched job in memory"))
		return
	}
	// 直接获取并执行
	entry.Job.Run()

	resp = &dto.RunJobResp{}
	return
}

// 获取已经添加到数据库中的任务列表
func (*_taskService) GetRunningJobs(req *dto.GetRunningJobListReq, ctx *gin.Context) (resp *dto.GetRunningJobListResp, err error) {
	list, err := biz.CronJob.WithContext(ctx).Find()
	if err != nil {
		err = util.CreateBizErr("查询任务列表失败", err)
		return
	}
	resp = &dto.GetRunningJobListResp{
		List: list,
	}
	return
}

// 获取已经添加到数据库中的任务列表
func (*_taskService) DeleteRunningJobs(req *dto.DeleteRunningJobListReq, ctx *gin.Context) (resp *dto.DeleteRunningJobListResp, err error) {
	err = extra.Q.Transaction(func(tx *extra.Query) error {
		// 加锁操作
		_memJobRecords.Lock()
		defer _memJobRecords.Unlock()

		jobTx := biz.CronJob
		_, e := jobTx.WithContext(ctx).
			Where(jobTx.ID.In(req.IDList...)).
			Delete()
		if e != nil {
			return e
		}
		shouldRemoves := arr.FilterSlice(
			_memJobRecords.records,
			func(memRecord *_taskRelationRecord, _ int) bool {
				return slices.ContainsFunc(req.IDList, func(dbID int64) bool {
					return memRecord.dbRecordID == dbID
				})
			})
		for _, d := range shouldRemoves {
			_singleScheduler.Remove(d.entryID)
		}

		replaces := arr.FilterSlice(_memJobRecords.records, func(memRecord *_taskRelationRecord, _ int) bool {
			return !slices.ContainsFunc(
				req.IDList,
				func(dbID int64) bool {
					return memRecord.dbRecordID == dbID
				},
			)
		})
		_memJobRecords.records = replaces

		return nil
	})

	if err != nil {
		err = util.CreateBizErr("删除任务失败", err)
		return
	}
	resp = &dto.DeleteRunningJobListResp{}
	return
}
