package dto

import "space-api/util/ptr"

// BasePageParam 基本的分页参数
type BasePageParam struct {
	Page *int `json:"page" form:"page"`
	Size *int `json:"size" form:"size"`
}

// WarningOverride 对于涉及到可能覆盖的关联数据的强制操作
type WarningOverride struct {
	// 允许强制覆盖可能影响的业务(比如关联资源)
	ForceOverride bool `json:"forceOverride"`
}

// Normalize 规范化分页查询参数, 去除控制针, 控制每页最大查询数量
func (bp *BasePageParam) Normalize() (offset, limit int) {
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
