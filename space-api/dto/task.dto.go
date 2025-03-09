package dto

import (
	"space-api/util/performance"
	"space-domain/model"
)

type (
	RegisteredJob struct {
		JobName     string `json:"jobName"`
		Description string `json:"description"`
	}
	CreateOrUpdateJobReq struct {
		DBRecordID  int64  `json:"dbRecordID,string"`
		JobFuncName string `json:"jobFuncName"`
		CronExpr    string `json:"cronExpr"`
		Enable      int    `json:"enable"`
		Mark        string `json:"mark"`
	}
	CreateOrUpdateJobResp performance.Empty

	GetAvailableJobListReq  struct{}
	GetAvailableJobListResp struct {
		List []*RegisteredJob `json:"list" yaml:"list" xml:"list" toml:"list"`
	}

	GetRunningJobListReq  struct{}
	GetRunningJobListResp struct {
		List []*model.CronJob `json:"list" yaml:"list" xml:"list" toml:"list"`
	}

	DeleteRunningJobListReq struct {
		IDList []int64 `json:"idList"`
	}
	DeleteRunningJobListResp performance.Empty

	RunJobReq struct {
		JobID int64 `json:"jobID,string" yaml:"jobID" xml:"jobID" toml:"jobID"`
	}
	RunJobResp performance.Empty
)
