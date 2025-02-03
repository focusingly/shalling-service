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

func newMenuLink(db *gorm.DB, opts ...gen.DOOption) menuLink {
	_menuLink := menuLink{}

	_menuLink.menuLinkDo.UseDB(db, opts...)
	_menuLink.menuLinkDo.UseModel(&model.MenuLink{})

	tableName := _menuLink.menuLinkDo.TableName()
	_menuLink.ALL = field.NewAsterisk(tableName)
	_menuLink.Id = field.NewInt64(tableName, "id")
	_menuLink.CreatedAt = field.NewInt64(tableName, "created_at")
	_menuLink.UpdatedAt = field.NewInt64(tableName, "updated_at")
	_menuLink.Hide = field.NewUint8(tableName, "hide")
	_menuLink.DisplayName = field.NewString(tableName, "display_name")
	_menuLink.RoutePath = field.NewString(tableName, "route_path")
	_menuLink.LinkType = field.NewString(tableName, "link_type")
	_menuLink.OpenWay = field.NewString(tableName, "open_way")

	_menuLink.fillFieldMap()

	return _menuLink
}

type menuLink struct {
	menuLinkDo menuLinkDo

	ALL         field.Asterisk
	Id          field.Int64
	CreatedAt   field.Int64
	UpdatedAt   field.Int64
	Hide        field.Uint8
	DisplayName field.String
	RoutePath   field.String
	LinkType    field.String
	OpenWay     field.String

	fieldMap map[string]field.Expr
}

func (m menuLink) Table(newTableName string) *menuLink {
	m.menuLinkDo.UseTable(newTableName)
	return m.updateTableName(newTableName)
}

func (m menuLink) As(alias string) *menuLink {
	m.menuLinkDo.DO = *(m.menuLinkDo.As(alias).(*gen.DO))
	return m.updateTableName(alias)
}

func (m *menuLink) updateTableName(table string) *menuLink {
	m.ALL = field.NewAsterisk(table)
	m.Id = field.NewInt64(table, "id")
	m.CreatedAt = field.NewInt64(table, "created_at")
	m.UpdatedAt = field.NewInt64(table, "updated_at")
	m.Hide = field.NewUint8(table, "hide")
	m.DisplayName = field.NewString(table, "display_name")
	m.RoutePath = field.NewString(table, "route_path")
	m.LinkType = field.NewString(table, "link_type")
	m.OpenWay = field.NewString(table, "open_way")

	m.fillFieldMap()

	return m
}

func (m *menuLink) WithContext(ctx context.Context) IMenuLinkDo { return m.menuLinkDo.WithContext(ctx) }

func (m menuLink) TableName() string { return m.menuLinkDo.TableName() }

func (m menuLink) Alias() string { return m.menuLinkDo.Alias() }

func (m menuLink) Columns(cols ...field.Expr) gen.Columns { return m.menuLinkDo.Columns(cols...) }

func (m *menuLink) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := m.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (m *menuLink) fillFieldMap() {
	m.fieldMap = make(map[string]field.Expr, 8)
	m.fieldMap["id"] = m.Id
	m.fieldMap["created_at"] = m.CreatedAt
	m.fieldMap["updated_at"] = m.UpdatedAt
	m.fieldMap["hide"] = m.Hide
	m.fieldMap["display_name"] = m.DisplayName
	m.fieldMap["route_path"] = m.RoutePath
	m.fieldMap["link_type"] = m.LinkType
	m.fieldMap["open_way"] = m.OpenWay
}

func (m menuLink) clone(db *gorm.DB) menuLink {
	m.menuLinkDo.ReplaceConnPool(db.Statement.ConnPool)
	return m
}

func (m menuLink) replaceDB(db *gorm.DB) menuLink {
	m.menuLinkDo.ReplaceDB(db)
	return m
}

type menuLinkDo struct{ gen.DO }

type IMenuLinkDo interface {
	gen.SubQuery
	Debug() IMenuLinkDo
	WithContext(ctx context.Context) IMenuLinkDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IMenuLinkDo
	WriteDB() IMenuLinkDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IMenuLinkDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IMenuLinkDo
	Not(conds ...gen.Condition) IMenuLinkDo
	Or(conds ...gen.Condition) IMenuLinkDo
	Select(conds ...field.Expr) IMenuLinkDo
	Where(conds ...gen.Condition) IMenuLinkDo
	Order(conds ...field.Expr) IMenuLinkDo
	Distinct(cols ...field.Expr) IMenuLinkDo
	Omit(cols ...field.Expr) IMenuLinkDo
	Join(table schema.Tabler, on ...field.Expr) IMenuLinkDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IMenuLinkDo
	RightJoin(table schema.Tabler, on ...field.Expr) IMenuLinkDo
	Group(cols ...field.Expr) IMenuLinkDo
	Having(conds ...gen.Condition) IMenuLinkDo
	Limit(limit int) IMenuLinkDo
	Offset(offset int) IMenuLinkDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IMenuLinkDo
	Unscoped() IMenuLinkDo
	Create(values ...*model.MenuLink) error
	CreateInBatches(values []*model.MenuLink, batchSize int) error
	Save(values ...*model.MenuLink) error
	First() (*model.MenuLink, error)
	Take() (*model.MenuLink, error)
	Last() (*model.MenuLink, error)
	Find() ([]*model.MenuLink, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.MenuLink, err error)
	FindInBatches(result *[]*model.MenuLink, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.MenuLink) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IMenuLinkDo
	Assign(attrs ...field.AssignExpr) IMenuLinkDo
	Joins(fields ...field.RelationField) IMenuLinkDo
	Preload(fields ...field.RelationField) IMenuLinkDo
	FirstOrInit() (*model.MenuLink, error)
	FirstOrCreate() (*model.MenuLink, error)
	FindByPage(offset int, limit int) (result []*model.MenuLink, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IMenuLinkDo
	UnderlyingDB() *gorm.DB
	schema.Tabler

	SelectWithSorts(condList []model.WhereCond, sortList []model.SortColumn) (result []model.MenuLink, err error)
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
func (m menuLinkDo) SelectWithSorts(condList []model.WhereCond, sortList []model.SortColumn) (result []model.MenuLink, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("SELECT * FROM menu_links ")
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

func (m menuLinkDo) Debug() IMenuLinkDo {
	return m.withDO(m.DO.Debug())
}

func (m menuLinkDo) WithContext(ctx context.Context) IMenuLinkDo {
	return m.withDO(m.DO.WithContext(ctx))
}

func (m menuLinkDo) ReadDB() IMenuLinkDo {
	return m.Clauses(dbresolver.Read)
}

func (m menuLinkDo) WriteDB() IMenuLinkDo {
	return m.Clauses(dbresolver.Write)
}

func (m menuLinkDo) Session(config *gorm.Session) IMenuLinkDo {
	return m.withDO(m.DO.Session(config))
}

func (m menuLinkDo) Clauses(conds ...clause.Expression) IMenuLinkDo {
	return m.withDO(m.DO.Clauses(conds...))
}

func (m menuLinkDo) Returning(value interface{}, columns ...string) IMenuLinkDo {
	return m.withDO(m.DO.Returning(value, columns...))
}

func (m menuLinkDo) Not(conds ...gen.Condition) IMenuLinkDo {
	return m.withDO(m.DO.Not(conds...))
}

func (m menuLinkDo) Or(conds ...gen.Condition) IMenuLinkDo {
	return m.withDO(m.DO.Or(conds...))
}

func (m menuLinkDo) Select(conds ...field.Expr) IMenuLinkDo {
	return m.withDO(m.DO.Select(conds...))
}

func (m menuLinkDo) Where(conds ...gen.Condition) IMenuLinkDo {
	return m.withDO(m.DO.Where(conds...))
}

func (m menuLinkDo) Order(conds ...field.Expr) IMenuLinkDo {
	return m.withDO(m.DO.Order(conds...))
}

func (m menuLinkDo) Distinct(cols ...field.Expr) IMenuLinkDo {
	return m.withDO(m.DO.Distinct(cols...))
}

func (m menuLinkDo) Omit(cols ...field.Expr) IMenuLinkDo {
	return m.withDO(m.DO.Omit(cols...))
}

func (m menuLinkDo) Join(table schema.Tabler, on ...field.Expr) IMenuLinkDo {
	return m.withDO(m.DO.Join(table, on...))
}

func (m menuLinkDo) LeftJoin(table schema.Tabler, on ...field.Expr) IMenuLinkDo {
	return m.withDO(m.DO.LeftJoin(table, on...))
}

func (m menuLinkDo) RightJoin(table schema.Tabler, on ...field.Expr) IMenuLinkDo {
	return m.withDO(m.DO.RightJoin(table, on...))
}

func (m menuLinkDo) Group(cols ...field.Expr) IMenuLinkDo {
	return m.withDO(m.DO.Group(cols...))
}

func (m menuLinkDo) Having(conds ...gen.Condition) IMenuLinkDo {
	return m.withDO(m.DO.Having(conds...))
}

func (m menuLinkDo) Limit(limit int) IMenuLinkDo {
	return m.withDO(m.DO.Limit(limit))
}

func (m menuLinkDo) Offset(offset int) IMenuLinkDo {
	return m.withDO(m.DO.Offset(offset))
}

func (m menuLinkDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IMenuLinkDo {
	return m.withDO(m.DO.Scopes(funcs...))
}

func (m menuLinkDo) Unscoped() IMenuLinkDo {
	return m.withDO(m.DO.Unscoped())
}

func (m menuLinkDo) Create(values ...*model.MenuLink) error {
	if len(values) == 0 {
		return nil
	}
	return m.DO.Create(values)
}

func (m menuLinkDo) CreateInBatches(values []*model.MenuLink, batchSize int) error {
	return m.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (m menuLinkDo) Save(values ...*model.MenuLink) error {
	if len(values) == 0 {
		return nil
	}
	return m.DO.Save(values)
}

func (m menuLinkDo) First() (*model.MenuLink, error) {
	if result, err := m.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.MenuLink), nil
	}
}

func (m menuLinkDo) Take() (*model.MenuLink, error) {
	if result, err := m.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.MenuLink), nil
	}
}

func (m menuLinkDo) Last() (*model.MenuLink, error) {
	if result, err := m.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.MenuLink), nil
	}
}

func (m menuLinkDo) Find() ([]*model.MenuLink, error) {
	result, err := m.DO.Find()
	return result.([]*model.MenuLink), err
}

func (m menuLinkDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.MenuLink, err error) {
	buf := make([]*model.MenuLink, 0, batchSize)
	err = m.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (m menuLinkDo) FindInBatches(result *[]*model.MenuLink, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return m.DO.FindInBatches(result, batchSize, fc)
}

func (m menuLinkDo) Attrs(attrs ...field.AssignExpr) IMenuLinkDo {
	return m.withDO(m.DO.Attrs(attrs...))
}

func (m menuLinkDo) Assign(attrs ...field.AssignExpr) IMenuLinkDo {
	return m.withDO(m.DO.Assign(attrs...))
}

func (m menuLinkDo) Joins(fields ...field.RelationField) IMenuLinkDo {
	for _, _f := range fields {
		m = *m.withDO(m.DO.Joins(_f))
	}
	return &m
}

func (m menuLinkDo) Preload(fields ...field.RelationField) IMenuLinkDo {
	for _, _f := range fields {
		m = *m.withDO(m.DO.Preload(_f))
	}
	return &m
}

func (m menuLinkDo) FirstOrInit() (*model.MenuLink, error) {
	if result, err := m.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.MenuLink), nil
	}
}

func (m menuLinkDo) FirstOrCreate() (*model.MenuLink, error) {
	if result, err := m.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.MenuLink), nil
	}
}

func (m menuLinkDo) FindByPage(offset int, limit int) (result []*model.MenuLink, count int64, err error) {
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

func (m menuLinkDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = m.Count()
	if err != nil {
		return
	}

	err = m.Offset(offset).Limit(limit).Scan(result)
	return
}

func (m menuLinkDo) Scan(result interface{}) (err error) {
	return m.DO.Scan(result)
}

func (m menuLinkDo) Delete(models ...*model.MenuLink) (result gen.ResultInfo, err error) {
	return m.DO.Delete(models)
}

func (m *menuLinkDo) withDO(do gen.Dao) *menuLinkDo {
	m.DO = *do.(*gen.DO)
	return m
}
