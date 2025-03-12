package dto

import "space-domain/model"

type (
	GetLocalFilesPaginationReq struct {
		BasePageParam
	}
	GetLocalFilesPaginationResp struct {
		model.PageList[*model.FileRecord]
	}
)
