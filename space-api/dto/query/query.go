package query

import (
	"database/sql/driver"
	"fmt"
	"space-api/util/arr"

	"github.com/huandu/xstrings"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type (
	Operation string

	// 查找条件
	WhereCond struct {
		Op     Operation `json:"op" yaml:"op" xml:"op" toml:"op"`
		Column string    `json:"column" yaml:"column" xml:"column" toml:"column"`
		Value  []any     `json:"value" yaml:"value" xml:"value" toml:"value"`
	}

	// 排序条件
	OrderColumn struct {
		Column string `json:"column" yaml:"column" xml:"column" toml:"column"`
		IsDesc bool   `json:"isDesc" yaml:"isDesc" xml:"isDesc" toml:"isDesc"`
	}
)

func (o *OrderColumn) ToOrderField(tableName string) field.Expr {
	f := field.NewField(tableName, xstrings.ToSnakeCase(o.Column))
	if o.IsDesc {
		return f.Desc()
	}

	return f
}

func (w *WhereCond) Parse() (err error) {
	if w.Value == nil {
		return fmt.Errorf("require condition value, but got nil")
	}
	switch w.Op {
	case Eq, Neq, Gt, Gte, Lt, Lte, Like:
		if len(w.Value) != 1 {
			err = fmt.Errorf("%s should only have one arg but got: %d", w.Op, len(w.Value))
		}
	case IsNull, NotNull:
	case In:
	default:
		err = fmt.Errorf("%s unknown operation for column", w.Op)
	}

	if err != nil {
		return
	}

	return
}

type baseValuer struct {
	val any
}

// Value implements driver.Valuer.
func (b *baseValuer) Value() (driver.Value, error) {
	return b.val, nil
}

var _ driver.Valuer = (*baseValuer)(nil)

func NewDriverValue(val any) driver.Valuer {
	return &baseValuer{val: val}
}

func ResolveCondList(tableName string, list []*WhereCond) (parsedCond []gen.Condition, err error) {
	for _, cond := range list {
		var expr field.Expr
		// check param
		if e := cond.Parse(); e != nil {
			return nil, e
		}

		f := field.NewField(tableName, xstrings.ToSnakeCase(cond.Column))
		switch cond.Op {
		case Eq:
			expr = f.Eq(NewDriverValue(cond.Value[0]))
		case Neq:
			expr = f.Neq(NewDriverValue(cond.Value[0]))
		case Gt:
			expr = f.Gt(NewDriverValue(cond.Value[0]))
		case Gte:
			expr = f.Gte(NewDriverValue(cond.Value[0]))
		case Lt:
			expr = f.Lt(NewDriverValue(cond.Value[0]))
		case Lte:
			expr = f.Lte(NewDriverValue(cond.Value[0]))
		case In:
			expr = f.In(arr.MapSlice(cond.Value, func(_ int, t any) driver.Valuer {
				return NewDriverValue(t)
			})...)
		case Like:
			expr = f.Like(NewDriverValue(cond.Value[0]))
		case IsNull:
			expr = f.IsNull()
		case NotNull:
			expr = f.IsNotNull()
		default:
			return nil, fmt.Errorf("un-support operator: %s", cond.Op)
		}

		parsedCond = append(parsedCond, expr)
	}

	return
}

func (o Operation) String() string {
	return string(o)
}

const (
	Eq      Operation = "="
	Neq     Operation = "!="
	Gt      Operation = ">"
	Gte     Operation = ">="
	Lt      Operation = "<"
	Lte     Operation = "<="
	In      Operation = "in"
	Like    Operation = "like"
	IsNull  Operation = "isNull"
	NotNull Operation = "notNull"
)
