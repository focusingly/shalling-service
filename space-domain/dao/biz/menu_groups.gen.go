// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package biz

import (
	"context"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"space-domain/model"
)

func newMenuGroup(db *gorm.DB, opts ...gen.DOOption) menuGroup {
	_menuGroup := menuGroup{}

	_menuGroup.menuGroupDo.UseDB(db, opts...)
	_menuGroup.menuGroupDo.UseModel(&model.MenuGroup{})

	tableName := _menuGroup.menuGroupDo.TableName()
	_menuGroup.ALL = field.NewAsterisk(tableName)
	_menuGroup.ID = field.NewInt64(tableName, "id")
	_menuGroup.CreatedAt = field.NewInt64(tableName, "created_at")
	_menuGroup.UpdatedAt = field.NewInt64(tableName, "updated_at")
	_menuGroup.Hide = field.NewInt(tableName, "hide")
	_menuGroup.MenuName = field.NewString(tableName, "menu_name")
	_menuGroup.RoutePath = field.NewString(tableName, "route_path")
	_menuGroup.PostLink = field.NewInt64(tableName, "post_link")
	_menuGroup.OpenWay = field.NewString(tableName, "open_way")
	_menuGroup.SubLinks = field.NewField(tableName, "sub_links")

	_menuGroup.fillFieldMap()

	return _menuGroup
}

type menuGroup struct {
	menuGroupDo menuGroupDo

	ALL       field.Asterisk
	ID        field.Int64
	CreatedAt field.Int64
	UpdatedAt field.Int64
	Hide      field.Int
	MenuName  field.String
	RoutePath field.String
	PostLink  field.Int64
	OpenWay   field.String
	SubLinks  field.Field

	fieldMap map[string]field.Expr
}

func (m menuGroup) Table(newTableName string) *menuGroup {
	m.menuGroupDo.UseTable(newTableName)
	return m.updateTableName(newTableName)
}

func (m menuGroup) As(alias string) *menuGroup {
	m.menuGroupDo.DO = *(m.menuGroupDo.As(alias).(*gen.DO))
	return m.updateTableName(alias)
}

func (m *menuGroup) updateTableName(table string) *menuGroup {
	m.ALL = field.NewAsterisk(table)
	m.ID = field.NewInt64(table, "id")
	m.CreatedAt = field.NewInt64(table, "created_at")
	m.UpdatedAt = field.NewInt64(table, "updated_at")
	m.Hide = field.NewInt(table, "hide")
	m.MenuName = field.NewString(table, "menu_name")
	m.RoutePath = field.NewString(table, "route_path")
	m.PostLink = field.NewInt64(table, "post_link")
	m.OpenWay = field.NewString(table, "open_way")
	m.SubLinks = field.NewField(table, "sub_links")

	m.fillFieldMap()

	return m
}

func (m *menuGroup) WithContext(ctx context.Context) IMenuGroupDo {
	return m.menuGroupDo.WithContext(ctx)
}

func (m menuGroup) TableName() string { return m.menuGroupDo.TableName() }

func (m menuGroup) Alias() string { return m.menuGroupDo.Alias() }

func (m menuGroup) Columns(cols ...field.Expr) gen.Columns { return m.menuGroupDo.Columns(cols...) }

func (m *menuGroup) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := m.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (m *menuGroup) fillFieldMap() {
	m.fieldMap = make(map[string]field.Expr, 9)
	m.fieldMap["id"] = m.ID
	m.fieldMap["created_at"] = m.CreatedAt
	m.fieldMap["updated_at"] = m.UpdatedAt
	m.fieldMap["hide"] = m.Hide
	m.fieldMap["menu_name"] = m.MenuName
	m.fieldMap["route_path"] = m.RoutePath
	m.fieldMap["post_link"] = m.PostLink
	m.fieldMap["open_way"] = m.OpenWay
	m.fieldMap["sub_links"] = m.SubLinks
}

func (m menuGroup) clone(db *gorm.DB) menuGroup {
	m.menuGroupDo.ReplaceConnPool(db.Statement.ConnPool)
	return m
}

func (m menuGroup) replaceDB(db *gorm.DB) menuGroup {
	m.menuGroupDo.ReplaceDB(db)
	return m
}

type menuGroupDo struct{ gen.DO }

type IMenuGroupDo interface {
	gen.SubQuery
	Debug() IMenuGroupDo
	WithContext(ctx context.Context) IMenuGroupDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IMenuGroupDo
	WriteDB() IMenuGroupDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IMenuGroupDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IMenuGroupDo
	Not(conds ...gen.Condition) IMenuGroupDo
	Or(conds ...gen.Condition) IMenuGroupDo
	Select(conds ...field.Expr) IMenuGroupDo
	Where(conds ...gen.Condition) IMenuGroupDo
	Order(conds ...field.Expr) IMenuGroupDo
	Distinct(cols ...field.Expr) IMenuGroupDo
	Omit(cols ...field.Expr) IMenuGroupDo
	Join(table schema.Tabler, on ...field.Expr) IMenuGroupDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IMenuGroupDo
	RightJoin(table schema.Tabler, on ...field.Expr) IMenuGroupDo
	Group(cols ...field.Expr) IMenuGroupDo
	Having(conds ...gen.Condition) IMenuGroupDo
	Limit(limit int) IMenuGroupDo
	Offset(offset int) IMenuGroupDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IMenuGroupDo
	Unscoped() IMenuGroupDo
	Create(values ...*model.MenuGroup) error
	CreateInBatches(values []*model.MenuGroup, batchSize int) error
	Save(values ...*model.MenuGroup) error
	First() (*model.MenuGroup, error)
	Take() (*model.MenuGroup, error)
	Last() (*model.MenuGroup, error)
	Find() ([]*model.MenuGroup, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.MenuGroup, err error)
	FindInBatches(result *[]*model.MenuGroup, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.MenuGroup) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IMenuGroupDo
	Assign(attrs ...field.AssignExpr) IMenuGroupDo
	Joins(fields ...field.RelationField) IMenuGroupDo
	Preload(fields ...field.RelationField) IMenuGroupDo
	FirstOrInit() (*model.MenuGroup, error)
	FirstOrCreate() (*model.MenuGroup, error)
	FindByPage(offset int, limit int) (result []*model.MenuGroup, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IMenuGroupDo
	UnderlyingDB() *gorm.DB
	schema.Tabler

	SelectWithSorts(condList []model.WhereCond, sortList []model.SortColumn) (result []model.MenuGroup, err error)
}

// SELECT * FROM @@table
//
//	{{if len(condList) > 0}}
//		WHERE
//		{{for index, cond := range condList}}
//			{{if index < len(condList)-1}}
//				@@cond.Column = @cond.Val AND
//			{{else}}
//				@@cond.Column = @cond.Val
//			{{end}}
//		{{end}}
//	{{end}}
//	{{if len(sortList) > 0}}
//		ORDER BY
//		{{for index, sort := range sortList}}
//			{{if index < len(sortList)-1}}
//				{{if sort.Desc}}
//					@@sort.Column DESC,
//				{{else}}
//					@@sort.Column ASC,
//				{{end}}
//			{{else}}
//				{{if sort.Desc}}
//					@@sort.Column DESC
//				{{else}}
//					@@sort.Column ASC
//				{{end}}
//			{{end}}
//		{{end}}
//	{{end}}
func (m menuGroupDo) SelectWithSorts(condList []model.WhereCond, sortList []model.SortColumn) (result []model.MenuGroup, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("SELECT * FROM menu_groups ")
	if len(condList) > 0 {
		generateSQL.WriteString("WHERE ")
		for index, cond := range condList {
			if index < len(condList)-1 {
				params = append(params, cond.Val)
				generateSQL.WriteString(m.Quote(cond.Column) + " = ? AND ")
			} else {
				params = append(params, cond.Val)
				generateSQL.WriteString(m.Quote(cond.Column) + " = ? ")
			}
		}
	}
	if len(sortList) > 0 {
		generateSQL.WriteString("ORDER BY ")
		for index, sort := range sortList {
			if index < len(sortList)-1 {
				if sort.Desc {
					generateSQL.WriteString(m.Quote(sort.Column) + " DESC, ")
				} else {
					generateSQL.WriteString(m.Quote(sort.Column) + " ASC, ")
				}
			} else {
				if sort.Desc {
					generateSQL.WriteString(m.Quote(sort.Column) + " DESC ")
				} else {
					generateSQL.WriteString(m.Quote(sort.Column) + " ASC ")
				}
			}
		}
	}

	var executeSQL *gorm.DB
	executeSQL = m.UnderlyingDB().Raw(generateSQL.String(), params...).Find(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

func (m menuGroupDo) Debug() IMenuGroupDo {
	return m.withDO(m.DO.Debug())
}

func (m menuGroupDo) WithContext(ctx context.Context) IMenuGroupDo {
	return m.withDO(m.DO.WithContext(ctx))
}

func (m menuGroupDo) ReadDB() IMenuGroupDo {
	return m.Clauses(dbresolver.Read)
}

func (m menuGroupDo) WriteDB() IMenuGroupDo {
	return m.Clauses(dbresolver.Write)
}

func (m menuGroupDo) Session(config *gorm.Session) IMenuGroupDo {
	return m.withDO(m.DO.Session(config))
}

func (m menuGroupDo) Clauses(conds ...clause.Expression) IMenuGroupDo {
	return m.withDO(m.DO.Clauses(conds...))
}

func (m menuGroupDo) Returning(value interface{}, columns ...string) IMenuGroupDo {
	return m.withDO(m.DO.Returning(value, columns...))
}

func (m menuGroupDo) Not(conds ...gen.Condition) IMenuGroupDo {
	return m.withDO(m.DO.Not(conds...))
}

func (m menuGroupDo) Or(conds ...gen.Condition) IMenuGroupDo {
	return m.withDO(m.DO.Or(conds...))
}

func (m menuGroupDo) Select(conds ...field.Expr) IMenuGroupDo {
	return m.withDO(m.DO.Select(conds...))
}

func (m menuGroupDo) Where(conds ...gen.Condition) IMenuGroupDo {
	return m.withDO(m.DO.Where(conds...))
}

func (m menuGroupDo) Order(conds ...field.Expr) IMenuGroupDo {
	return m.withDO(m.DO.Order(conds...))
}

func (m menuGroupDo) Distinct(cols ...field.Expr) IMenuGroupDo {
	return m.withDO(m.DO.Distinct(cols...))
}

func (m menuGroupDo) Omit(cols ...field.Expr) IMenuGroupDo {
	return m.withDO(m.DO.Omit(cols...))
}

func (m menuGroupDo) Join(table schema.Tabler, on ...field.Expr) IMenuGroupDo {
	return m.withDO(m.DO.Join(table, on...))
}

func (m menuGroupDo) LeftJoin(table schema.Tabler, on ...field.Expr) IMenuGroupDo {
	return m.withDO(m.DO.LeftJoin(table, on...))
}

func (m menuGroupDo) RightJoin(table schema.Tabler, on ...field.Expr) IMenuGroupDo {
	return m.withDO(m.DO.RightJoin(table, on...))
}

func (m menuGroupDo) Group(cols ...field.Expr) IMenuGroupDo {
	return m.withDO(m.DO.Group(cols...))
}

func (m menuGroupDo) Having(conds ...gen.Condition) IMenuGroupDo {
	return m.withDO(m.DO.Having(conds...))
}

func (m menuGroupDo) Limit(limit int) IMenuGroupDo {
	return m.withDO(m.DO.Limit(limit))
}

func (m menuGroupDo) Offset(offset int) IMenuGroupDo {
	return m.withDO(m.DO.Offset(offset))
}

func (m menuGroupDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IMenuGroupDo {
	return m.withDO(m.DO.Scopes(funcs...))
}

func (m menuGroupDo) Unscoped() IMenuGroupDo {
	return m.withDO(m.DO.Unscoped())
}

func (m menuGroupDo) Create(values ...*model.MenuGroup) error {
	if len(values) == 0 {
		return nil
	}
	return m.DO.Create(values)
}

func (m menuGroupDo) CreateInBatches(values []*model.MenuGroup, batchSize int) error {
	return m.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (m menuGroupDo) Save(values ...*model.MenuGroup) error {
	if len(values) == 0 {
		return nil
	}
	return m.DO.Save(values)
}

func (m menuGroupDo) First() (*model.MenuGroup, error) {
	if result, err := m.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.MenuGroup), nil
	}
}

func (m menuGroupDo) Take() (*model.MenuGroup, error) {
	if result, err := m.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.MenuGroup), nil
	}
}

func (m menuGroupDo) Last() (*model.MenuGroup, error) {
	if result, err := m.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.MenuGroup), nil
	}
}

func (m menuGroupDo) Find() ([]*model.MenuGroup, error) {
	result, err := m.DO.Find()
	return result.([]*model.MenuGroup), err
}

func (m menuGroupDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.MenuGroup, err error) {
	buf := make([]*model.MenuGroup, 0, batchSize)
	err = m.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (m menuGroupDo) FindInBatches(result *[]*model.MenuGroup, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return m.DO.FindInBatches(result, batchSize, fc)
}

func (m menuGroupDo) Attrs(attrs ...field.AssignExpr) IMenuGroupDo {
	return m.withDO(m.DO.Attrs(attrs...))
}

func (m menuGroupDo) Assign(attrs ...field.AssignExpr) IMenuGroupDo {
	return m.withDO(m.DO.Assign(attrs...))
}

func (m menuGroupDo) Joins(fields ...field.RelationField) IMenuGroupDo {
	for _, _f := range fields {
		m = *m.withDO(m.DO.Joins(_f))
	}
	return &m
}

func (m menuGroupDo) Preload(fields ...field.RelationField) IMenuGroupDo {
	for _, _f := range fields {
		m = *m.withDO(m.DO.Preload(_f))
	}
	return &m
}

func (m menuGroupDo) FirstOrInit() (*model.MenuGroup, error) {
	if result, err := m.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.MenuGroup), nil
	}
}

func (m menuGroupDo) FirstOrCreate() (*model.MenuGroup, error) {
	if result, err := m.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.MenuGroup), nil
	}
}

func (m menuGroupDo) FindByPage(offset int, limit int) (result []*model.MenuGroup, count int64, err error) {
	result, err = m.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = m.Offset(-1).Limit(-1).Count()
	return
}

func (m menuGroupDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = m.Count()
	if err != nil {
		return
	}

	err = m.Offset(offset).Limit(limit).Scan(result)
	return
}

func (m menuGroupDo) Scan(result interface{}) (err error) {
	return m.DO.Scan(result)
}

func (m menuGroupDo) Delete(models ...*model.MenuGroup) (result gen.ResultInfo, err error) {
	return m.DO.Delete(models)
}

func (m *menuGroupDo) withDO(do gen.Dao) *menuGroupDo {
	m.DO = *do.(*gen.DO)
	return m
}
