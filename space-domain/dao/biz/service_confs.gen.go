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

func newServiceConf(db *gorm.DB, opts ...gen.DOOption) serviceConf {
	_serviceConf := serviceConf{}

	_serviceConf.serviceConfDo.UseDB(db, opts...)
	_serviceConf.serviceConfDo.UseModel(&model.ServiceConf{})

	tableName := _serviceConf.serviceConfDo.TableName()
	_serviceConf.ALL = field.NewAsterisk(tableName)
	_serviceConf.ID = field.NewInt64(tableName, "id")
	_serviceConf.CreatedAt = field.NewInt64(tableName, "created_at")
	_serviceConf.UpdatedAt = field.NewInt64(tableName, "updated_at")
	_serviceConf.Hide = field.NewInt(tableName, "hide")
	_serviceConf.ConfKey = field.NewString(tableName, "conf_key")
	_serviceConf.ConfVal = field.NewString(tableName, "conf_val")
	_serviceConf.Category = field.NewString(tableName, "category")

	_serviceConf.fillFieldMap()

	return _serviceConf
}

type serviceConf struct {
	serviceConfDo serviceConfDo

	ALL       field.Asterisk
	ID        field.Int64
	CreatedAt field.Int64
	UpdatedAt field.Int64
	Hide      field.Int
	ConfKey   field.String
	ConfVal   field.String
	Category  field.String

	fieldMap map[string]field.Expr
}

func (s serviceConf) Table(newTableName string) *serviceConf {
	s.serviceConfDo.UseTable(newTableName)
	return s.updateTableName(newTableName)
}

func (s serviceConf) As(alias string) *serviceConf {
	s.serviceConfDo.DO = *(s.serviceConfDo.As(alias).(*gen.DO))
	return s.updateTableName(alias)
}

func (s *serviceConf) updateTableName(table string) *serviceConf {
	s.ALL = field.NewAsterisk(table)
	s.ID = field.NewInt64(table, "id")
	s.CreatedAt = field.NewInt64(table, "created_at")
	s.UpdatedAt = field.NewInt64(table, "updated_at")
	s.Hide = field.NewInt(table, "hide")
	s.ConfKey = field.NewString(table, "conf_key")
	s.ConfVal = field.NewString(table, "conf_val")
	s.Category = field.NewString(table, "category")

	s.fillFieldMap()

	return s
}

func (s *serviceConf) WithContext(ctx context.Context) IServiceConfDo {
	return s.serviceConfDo.WithContext(ctx)
}

func (s serviceConf) TableName() string { return s.serviceConfDo.TableName() }

func (s serviceConf) Alias() string { return s.serviceConfDo.Alias() }

func (s serviceConf) Columns(cols ...field.Expr) gen.Columns { return s.serviceConfDo.Columns(cols...) }

func (s *serviceConf) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := s.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (s *serviceConf) fillFieldMap() {
	s.fieldMap = make(map[string]field.Expr, 7)
	s.fieldMap["id"] = s.ID
	s.fieldMap["created_at"] = s.CreatedAt
	s.fieldMap["updated_at"] = s.UpdatedAt
	s.fieldMap["hide"] = s.Hide
	s.fieldMap["conf_key"] = s.ConfKey
	s.fieldMap["conf_val"] = s.ConfVal
	s.fieldMap["category"] = s.Category
}

func (s serviceConf) clone(db *gorm.DB) serviceConf {
	s.serviceConfDo.ReplaceConnPool(db.Statement.ConnPool)
	return s
}

func (s serviceConf) replaceDB(db *gorm.DB) serviceConf {
	s.serviceConfDo.ReplaceDB(db)
	return s
}

type serviceConfDo struct{ gen.DO }

type IServiceConfDo interface {
	gen.SubQuery
	Debug() IServiceConfDo
	WithContext(ctx context.Context) IServiceConfDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IServiceConfDo
	WriteDB() IServiceConfDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IServiceConfDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IServiceConfDo
	Not(conds ...gen.Condition) IServiceConfDo
	Or(conds ...gen.Condition) IServiceConfDo
	Select(conds ...field.Expr) IServiceConfDo
	Where(conds ...gen.Condition) IServiceConfDo
	Order(conds ...field.Expr) IServiceConfDo
	Distinct(cols ...field.Expr) IServiceConfDo
	Omit(cols ...field.Expr) IServiceConfDo
	Join(table schema.Tabler, on ...field.Expr) IServiceConfDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IServiceConfDo
	RightJoin(table schema.Tabler, on ...field.Expr) IServiceConfDo
	Group(cols ...field.Expr) IServiceConfDo
	Having(conds ...gen.Condition) IServiceConfDo
	Limit(limit int) IServiceConfDo
	Offset(offset int) IServiceConfDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IServiceConfDo
	Unscoped() IServiceConfDo
	Create(values ...*model.ServiceConf) error
	CreateInBatches(values []*model.ServiceConf, batchSize int) error
	Save(values ...*model.ServiceConf) error
	First() (*model.ServiceConf, error)
	Take() (*model.ServiceConf, error)
	Last() (*model.ServiceConf, error)
	Find() ([]*model.ServiceConf, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.ServiceConf, err error)
	FindInBatches(result *[]*model.ServiceConf, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.ServiceConf) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IServiceConfDo
	Assign(attrs ...field.AssignExpr) IServiceConfDo
	Joins(fields ...field.RelationField) IServiceConfDo
	Preload(fields ...field.RelationField) IServiceConfDo
	FirstOrInit() (*model.ServiceConf, error)
	FirstOrCreate() (*model.ServiceConf, error)
	FindByPage(offset int, limit int) (result []*model.ServiceConf, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IServiceConfDo
	UnderlyingDB() *gorm.DB
	schema.Tabler

	SelectWithSorts(condList []model.WhereCond, sortList []model.SortColumn) (result []model.ServiceConf, err error)
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
func (s serviceConfDo) SelectWithSorts(condList []model.WhereCond, sortList []model.SortColumn) (result []model.ServiceConf, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("SELECT * FROM service_confs ")
	if len(condList) > 0 {
		generateSQL.WriteString("WHERE ")
		for index, cond := range condList {
			if index < len(condList)-1 {
				params = append(params, cond.Val)
				generateSQL.WriteString(s.Quote(cond.Column) + " = ? AND ")
			} else {
				params = append(params, cond.Val)
				generateSQL.WriteString(s.Quote(cond.Column) + " = ? ")
			}
		}
	}
	if len(sortList) > 0 {
		generateSQL.WriteString("ORDER BY ")
		for index, sort := range sortList {
			if index < len(sortList)-1 {
				if sort.Desc {
					generateSQL.WriteString(s.Quote(sort.Column) + " DESC, ")
				} else {
					generateSQL.WriteString(s.Quote(sort.Column) + " ASC, ")
				}
			} else {
				if sort.Desc {
					generateSQL.WriteString(s.Quote(sort.Column) + " DESC ")
				} else {
					generateSQL.WriteString(s.Quote(sort.Column) + " ASC ")
				}
			}
		}
	}

	var executeSQL *gorm.DB
	executeSQL = s.UnderlyingDB().Raw(generateSQL.String(), params...).Find(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

func (s serviceConfDo) Debug() IServiceConfDo {
	return s.withDO(s.DO.Debug())
}

func (s serviceConfDo) WithContext(ctx context.Context) IServiceConfDo {
	return s.withDO(s.DO.WithContext(ctx))
}

func (s serviceConfDo) ReadDB() IServiceConfDo {
	return s.Clauses(dbresolver.Read)
}

func (s serviceConfDo) WriteDB() IServiceConfDo {
	return s.Clauses(dbresolver.Write)
}

func (s serviceConfDo) Session(config *gorm.Session) IServiceConfDo {
	return s.withDO(s.DO.Session(config))
}

func (s serviceConfDo) Clauses(conds ...clause.Expression) IServiceConfDo {
	return s.withDO(s.DO.Clauses(conds...))
}

func (s serviceConfDo) Returning(value interface{}, columns ...string) IServiceConfDo {
	return s.withDO(s.DO.Returning(value, columns...))
}

func (s serviceConfDo) Not(conds ...gen.Condition) IServiceConfDo {
	return s.withDO(s.DO.Not(conds...))
}

func (s serviceConfDo) Or(conds ...gen.Condition) IServiceConfDo {
	return s.withDO(s.DO.Or(conds...))
}

func (s serviceConfDo) Select(conds ...field.Expr) IServiceConfDo {
	return s.withDO(s.DO.Select(conds...))
}

func (s serviceConfDo) Where(conds ...gen.Condition) IServiceConfDo {
	return s.withDO(s.DO.Where(conds...))
}

func (s serviceConfDo) Order(conds ...field.Expr) IServiceConfDo {
	return s.withDO(s.DO.Order(conds...))
}

func (s serviceConfDo) Distinct(cols ...field.Expr) IServiceConfDo {
	return s.withDO(s.DO.Distinct(cols...))
}

func (s serviceConfDo) Omit(cols ...field.Expr) IServiceConfDo {
	return s.withDO(s.DO.Omit(cols...))
}

func (s serviceConfDo) Join(table schema.Tabler, on ...field.Expr) IServiceConfDo {
	return s.withDO(s.DO.Join(table, on...))
}

func (s serviceConfDo) LeftJoin(table schema.Tabler, on ...field.Expr) IServiceConfDo {
	return s.withDO(s.DO.LeftJoin(table, on...))
}

func (s serviceConfDo) RightJoin(table schema.Tabler, on ...field.Expr) IServiceConfDo {
	return s.withDO(s.DO.RightJoin(table, on...))
}

func (s serviceConfDo) Group(cols ...field.Expr) IServiceConfDo {
	return s.withDO(s.DO.Group(cols...))
}

func (s serviceConfDo) Having(conds ...gen.Condition) IServiceConfDo {
	return s.withDO(s.DO.Having(conds...))
}

func (s serviceConfDo) Limit(limit int) IServiceConfDo {
	return s.withDO(s.DO.Limit(limit))
}

func (s serviceConfDo) Offset(offset int) IServiceConfDo {
	return s.withDO(s.DO.Offset(offset))
}

func (s serviceConfDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IServiceConfDo {
	return s.withDO(s.DO.Scopes(funcs...))
}

func (s serviceConfDo) Unscoped() IServiceConfDo {
	return s.withDO(s.DO.Unscoped())
}

func (s serviceConfDo) Create(values ...*model.ServiceConf) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Create(values)
}

func (s serviceConfDo) CreateInBatches(values []*model.ServiceConf, batchSize int) error {
	return s.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (s serviceConfDo) Save(values ...*model.ServiceConf) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Save(values)
}

func (s serviceConfDo) First() (*model.ServiceConf, error) {
	if result, err := s.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.ServiceConf), nil
	}
}

func (s serviceConfDo) Take() (*model.ServiceConf, error) {
	if result, err := s.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.ServiceConf), nil
	}
}

func (s serviceConfDo) Last() (*model.ServiceConf, error) {
	if result, err := s.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.ServiceConf), nil
	}
}

func (s serviceConfDo) Find() ([]*model.ServiceConf, error) {
	result, err := s.DO.Find()
	return result.([]*model.ServiceConf), err
}

func (s serviceConfDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.ServiceConf, err error) {
	buf := make([]*model.ServiceConf, 0, batchSize)
	err = s.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (s serviceConfDo) FindInBatches(result *[]*model.ServiceConf, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return s.DO.FindInBatches(result, batchSize, fc)
}

func (s serviceConfDo) Attrs(attrs ...field.AssignExpr) IServiceConfDo {
	return s.withDO(s.DO.Attrs(attrs...))
}

func (s serviceConfDo) Assign(attrs ...field.AssignExpr) IServiceConfDo {
	return s.withDO(s.DO.Assign(attrs...))
}

func (s serviceConfDo) Joins(fields ...field.RelationField) IServiceConfDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Joins(_f))
	}
	return &s
}

func (s serviceConfDo) Preload(fields ...field.RelationField) IServiceConfDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Preload(_f))
	}
	return &s
}

func (s serviceConfDo) FirstOrInit() (*model.ServiceConf, error) {
	if result, err := s.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.ServiceConf), nil
	}
}

func (s serviceConfDo) FirstOrCreate() (*model.ServiceConf, error) {
	if result, err := s.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.ServiceConf), nil
	}
}

func (s serviceConfDo) FindByPage(offset int, limit int) (result []*model.ServiceConf, count int64, err error) {
	result, err = s.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = s.Offset(-1).Limit(-1).Count()
	return
}

func (s serviceConfDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = s.Count()
	if err != nil {
		return
	}

	err = s.Offset(offset).Limit(limit).Scan(result)
	return
}

func (s serviceConfDo) Scan(result interface{}) (err error) {
	return s.DO.Scan(result)
}

func (s serviceConfDo) Delete(models ...*model.ServiceConf) (result gen.ResultInfo, err error) {
	return s.DO.Delete(models)
}

func (s *serviceConfDo) withDO(do gen.Dao) *serviceConfDo {
	s.DO = *do.(*gen.DO)
	return s
}
