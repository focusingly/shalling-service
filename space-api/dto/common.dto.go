package dto

import "space-api/util/ptr"

// BasePageParam 基本的分页参数
type BasePageParam struct {
	Page *int `json:"page" form:"page"`
	Size *int `json:"size" form:"size"`
}

type BizOp struct {
	// 允许强制覆盖可能影响的业务(比如关联资源)
	ForceOverride bool `json:"forceOverride"`
}

func (bp *BasePageParam) Resolve() (offset, limit int) {
	if bp.Page == nil || *bp.Page <= 0 {
		bp.Page = ptr.ToPtr(1)
	}

	if bp.Size == nil || *bp.Size < 0 {
		bp.Size = ptr.ToPtr(10)
	}

	if *bp.Size > 100 {
		bp.Size = ptr.ToPtr(100)
	}

	offset = (*bp.Page - 1) * (*bp.Size)
	limit = *bp.Size

	return
}
