package dto

import "space-domain/model"

type (
	RegisteredJob struct {
		JobName     string `json:"jobName"`
		Description string `json:"description"`
	}
	CreateOrUpdateJobReq struct {
		DBRecordID  int64  `json:"dbRecordID"`
		JobFuncName string `json:"jobFuncName"`
		CronExpr    string `json:"cronExpr"`
		Enable      int    `json:"enable"`
		Mark        string `json:"mark"`
	}
	CreateOrUpdateJobResp struct{}

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
	DeleteRunningJobListResp struct {
	}

	RunJobReq struct {
		JobID int64 `json:"jobID" yaml:"jobID" xml:"jobID" toml:"jobID"`
	}
	RunJobResp struct{}
)
