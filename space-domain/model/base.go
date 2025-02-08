package model

import (
	"space-api/util/id"

	"gorm.io/gorm"
)

// BaseColumn 表的基础字段信息
type BaseColumn struct {
	ID        int64 `gorm:"type:bigint;primaryKey;autoIncrement:false;comment:主键" json:"id"`
	CreatedAt int64 `gorm:"type:bigint;autoCreateTime:milli;not null;comment:创建时间, unix 毫秒" json:"createdAt"`
	UpdatedAt int64 `gorm:"type:bigint;autoUpdateTime:milli;not null;comment:更新时间, unix 毫秒时间戳" json:"updatedAt"`
	Hide      int   `gorm:"type:smallint;not null;default:0;comment:是否隐藏, 默认为 0 不隐藏" json:"hide"`
}

// BeforeCreate 设置自定义 ID 插入
func (base *BaseColumn) BeforeCreate(tx *gorm.DB) (err error) {
	// 如果 ID 为 0, 那么手动变更
	if base.ID == 0 {
		base.ID = id.GetSnowFlakeNode().Generate().Int64()
	}

	return
}

type (
	// Where 简单赛选条件
	WhereCond struct {
		Column string `json:"column"` // 字段
		Val    any    `json:"val"`    // 参考值
	}

	// PageQuery 分页查询结果包装
	PageList[T any] struct {
		List  []T   `json:"list"`  // 查询的数据
		Page  int64 `json:"page"`  // 当前页数
		Size  int64 `json:"size"`  // 每页数量
		Total int64 `json:"total"` // 总记录数
	}

	// SortColumn 排序方式
	SortColumn struct {
		Column string `json:"column"` // 字段名称
		Desc   bool   `json:"desc"`   // 排序方式
	}

	// 分页查询以及排序的依据
	PageQueryCond struct {
		Page  int          `json:"page"`  // 页数
		Size  int          `json:"size"`  // 每页条数
		Sorts []SortColumn `json:"sorts"` // 字段排序规则
	}
)
