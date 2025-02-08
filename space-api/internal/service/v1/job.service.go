package service

import (
	"context"
	"fmt"
	"slices"
	"space-api/dto"
	"space-api/util"
	"space-api/util/arr"
	"space-domain/dao/extra"
	"space-domain/model"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
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

	jobRecord struct {
		entryID    cron.EntryID
		dbRecordID int64
	}

	memSyncRecord struct {
		sync.Mutex
		records []*jobRecord
	}
)

func (s *memSyncRecord) AppendJobID(rc *jobRecord) {
	s.Lock()
	defer s.Unlock()
	if rc == nil {
		panic("record can't nil")
	}
	s.records = append(s.records, &jobRecord{})
}

var (
	_singleScheduler = cron.New()
	_memJobRecords   = &memSyncRecord{}
	// 已经注册的所有任务
	// TODO 暂时只支持本地定义的方法
	RegisterJobs = []*CustomTask{
		{
			FuncName: "ClearOldLogs",
			Task: func() {
				fmt.Println("clear old logs...")
			},
			Description: "日志清理",
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
	jobs, err := extra.CronJob.WithContext(context.TODO()).Find()
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
			_memJobRecords.AppendJobID(&jobRecord{
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
	err = extra.Q.Transaction(func(tx *extra.Query) error {
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
			e := jobTx.WithContext(ctx).Create(
				&model.CronJob{
					BaseColumn:  model.BaseColumn{},
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
				entryID := _singleScheduler.Schedule(parsedCronExpr, cron.FuncJob(newTaskRecord.Task))
				// 同步到内存列表当中
				_memJobRecords.records = append(_memJobRecords.records, &jobRecord{
					entryID:    entryID,
					dbRecordID: findTask.ID,
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
			existsTaskIndex := slices.IndexFunc(_memJobRecords.records, func(rc *jobRecord) bool {
				return rc.dbRecordID == findTask.ID
			})
			if existsTaskIndex == -1 {
				return fmt.Errorf("未找到内存中的任务记录, 可能存在问题")
			}
			exitsTask := _memJobRecords.records[existsTaskIndex]
			// 取消掉旧任务
			_singleScheduler.Remove(exitsTask.entryID)
			// 构建新的任务列表
			updates := slices.Delete(_memJobRecords.records, existsTaskIndex, existsTaskIndex+1)
			if req.Enable != 0 {
				// 设置新的任务
				entryID := _singleScheduler.Schedule(parsedCronExpr, cron.FuncJob(newTaskRecord.Task))
				updates = append(updates, &jobRecord{
					entryID:    entryID,
					dbRecordID: findTask.ID,
				})
			}
			// 替换任务记录
			_memJobRecords.records = updates
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
	jobTx := extra.CronJob
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
	// 上锁
	_memJobRecords.Lock()
	defer _memJobRecords.Lock()

	// 查找是否在活动列表当中
	index := slices.IndexFunc(_memJobRecords.records, func(r *jobRecord) bool {
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
	list, err := extra.CronJob.WithContext(ctx).Find()
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

		jobTx := tx.CronJob
		_, e := jobTx.WithContext(ctx).
			Where(jobTx.ID.In(req.IDList...)).
			Delete()
		if e != nil {
			return e
		}
		shouldRemoves := arr.FilterSlice(
			_memJobRecords.records,
			func(memRecord *jobRecord, _ int) bool {
				return slices.ContainsFunc(req.IDList, func(dbID int64) bool {
					return memRecord.dbRecordID == dbID
				})
			})
		for _, d := range shouldRemoves {
			_singleScheduler.Remove(d.entryID)
		}

		replaces := arr.FilterSlice(_memJobRecords.records, func(memRecord *jobRecord, _ int) bool {
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
