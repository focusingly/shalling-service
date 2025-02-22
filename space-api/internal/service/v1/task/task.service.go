package task

import (
	"context"
	"fmt"
	"slices"
	"space-api/dto"
	"space-api/util"
	"space-api/util/arr"
	"space-api/util/id"
	"space-domain/dao/biz"
	"space-domain/dao/extra"
	"space-domain/model"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

type taskService struct {
	mu       sync.Mutex
	initDone bool
}

var DefaultTaskService = &taskService{}

var (
	singleScheduler = cron.New()
	memTaskRecords  = &taskMemRecord{}
)

// ResumeTasksFromPersistData 从数据库中恢复执行, 如果已经执行过了, 那么会直接返回
func (s *taskService) ResumeTasksFromPersistData() (err error) {
	if s.initDone {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// 从数据库中恢复任务记录
	jobs, err := biz.CronJob.WithContext(context.TODO()).Find()
	if err != nil {
		return
	}

	for _, job := range jobs {
		index := slices.IndexFunc(RegisterJobs, func(j *customTask) bool {
			return j.FuncName == job.JobFuncName
		})
		if index == -1 {
			return fmt.Errorf("非注册的定时任务: %s", job.JobFuncName)
		}
		if job.Enable != 0 {
			entryID, err := singleScheduler.AddFunc(job.CronExpr, RegisterJobs[index].Task)
			if err != nil {
				return err
			}
			memTaskRecords.appendJobRelation(&taskRelation{
				entryID:    entryID,
				dbRecordID: job.ID,
			})
		}
	}

	singleScheduler.Start()
	s.initDone = true

	return
}

func (*taskService) GetAvailableJobList(req *dto.GetAvailableJobListReq, ctx *gin.Context) (resp *dto.GetAvailableJobListResp, err error) {
	return &dto.GetAvailableJobListResp{
		List: arr.MapSlice(RegisterJobs, func(_ int, job *customTask) *dto.RegisteredJob {
			return &dto.RegisteredJob{
				JobName:     job.FuncName,
				Description: job.Description,
			}
		}),
	}, nil
}

func (*taskService) CreateOrUpdateNewJob(req *dto.CreateOrUpdateJobReq, ctx *gin.Context) (resp *dto.CreateOrUpdateJobResp, err error) {
	err = biz.Q.Transaction(func(tx *biz.Query) error {
		// 加锁操作
		memTaskRecords.Lock()
		defer memTaskRecords.Unlock()

		jobTx := tx.CronJob
		// 前置检查参数
		// TODO 暂时只允许注册任务

		// 任务索引所在的位置
		newTaskIndex := slices.IndexFunc(RegisterJobs, func(j *customTask) bool {
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
				entryID := singleScheduler.Schedule(
					parsedCronExpr,
					newJobWrapper(newTaskRecord.FuncName, newTaskRecord.Task),
				)
				// 同步到内存列表当中
				memTaskRecords.records = append(memTaskRecords.records, &taskRelation{
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
			existsTaskIndex := slices.IndexFunc(memTaskRecords.records, func(rc *taskRelation) bool {
				return rc.dbRecordID == findTask.ID
			})

			// 此前任务就已经在运行的状态
			if existsTaskIndex != -1 {
				// 存在运行的任务
				exitsTask := memTaskRecords.records[existsTaskIndex]

				// 取消掉旧任务
				singleScheduler.Remove(exitsTask.entryID)

				// 构建新的任务列表 -> 移除掉旧任务
				updates := slices.Delete(memTaskRecords.records, existsTaskIndex, existsTaskIndex+1)

				// 需要启动新的任务情况
				if req.Enable != 0 {
					// 设置新的任务
					entryID := singleScheduler.Schedule(parsedCronExpr, newJobWrapper(newTaskRecord.FuncName, newTaskRecord.Task))
					// 追加任务
					updates = append(updates, &taskRelation{
						entryID:    entryID,
						dbRecordID: findTask.ID,
					})
				}

				// 整体替换任务记录
				memTaskRecords.records = updates
			} else {
				// 此前并无已经存在运行的任务

				if req.Enable != 0 {
					// 设置新的任务
					entryID := singleScheduler.Schedule(parsedCronExpr, newJobWrapper(newTaskRecord.FuncName, newTaskRecord.Task))

					// 追加任务记录
					memTaskRecords.records = append(memTaskRecords.records, &taskRelation{
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
func (*taskService) RunJobImmediately(req *dto.RunJobReq, ctx *gin.Context) (resp *dto.RunJobResp, err error) {
	// 上锁
	memTaskRecords.Lock()
	defer memTaskRecords.Unlock()

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
	index := slices.IndexFunc(memTaskRecords.records, func(r *taskRelation) bool {
		return r.dbRecordID == job.ID
	})
	if index == -1 {
		err = util.CreateBizErr("未找到活动任务, 可能存在问题", fmt.Errorf("could not find ant matched job in memory"))
		return
	}
	entry := singleScheduler.Entry(memTaskRecords.records[index].entryID)
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
func (*taskService) GetRunningJobs(req *dto.GetRunningJobListReq, ctx *gin.Context) (resp *dto.GetRunningJobListResp, err error) {
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
func (*taskService) DeleteRunningJobs(req *dto.DeleteRunningJobListReq, ctx *gin.Context) (resp *dto.DeleteRunningJobListResp, err error) {
	err = extra.Q.Transaction(func(tx *extra.Query) error {
		// 加锁操作
		memTaskRecords.Lock()
		defer memTaskRecords.Unlock()

		jobTx := biz.CronJob
		_, e := jobTx.WithContext(ctx).
			Where(jobTx.ID.In(req.IDList...)).
			Delete()
		if e != nil {
			return e
		}
		shouldRemoves := arr.FilterSlice(
			memTaskRecords.records,
			func(memRecord *taskRelation, _ int) bool {
				return slices.ContainsFunc(req.IDList, func(dbID int64) bool {
					return memRecord.dbRecordID == dbID
				})
			})
		for _, d := range shouldRemoves {
			singleScheduler.Remove(d.entryID)
		}

		replaces := arr.FilterSlice(memTaskRecords.records, func(memRecord *taskRelation, _ int) bool {
			return !slices.ContainsFunc(
				req.IDList,
				func(dbID int64) bool {
					return memRecord.dbRecordID == dbID
				},
			)
		})
		memTaskRecords.records = replaces

		return nil
	})

	if err != nil {
		err = util.CreateBizErr("删除任务失败", err)
		return
	}
	resp = &dto.DeleteRunningJobListResp{}
	return
}
