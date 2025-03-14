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

func newSqlite3KeywordDoc(db *gorm.DB, opts ...gen.DOOption) sqlite3KeywordDoc {
	_sqlite3KeywordDoc := sqlite3KeywordDoc{}

	_sqlite3KeywordDoc.sqlite3KeywordDocDo.UseDB(db, opts...)
	_sqlite3KeywordDoc.sqlite3KeywordDocDo.UseModel(&model.Sqlite3KeywordDoc{})

	tableName := _sqlite3KeywordDoc.sqlite3KeywordDocDo.TableName()
	_sqlite3KeywordDoc.ALL = field.NewAsterisk(tableName)
	_sqlite3KeywordDoc.PostID = field.NewInt64(tableName, "post_id")
	_sqlite3KeywordDoc.TileSplit = field.NewString(tableName, "title_split")
	_sqlite3KeywordDoc.ContentSplit = field.NewString(tableName, "content_split")
	_sqlite3KeywordDoc.Weight = field.NewInt(tableName, "weight")
	_sqlite3KeywordDoc.PostUpdatedAt = field.NewInt64(tableName, "post_updated_at")
	_sqlite3KeywordDoc.RecordCreatedAt = field.NewInt64(tableName, "record_created_at")
	_sqlite3KeywordDoc.RecordUpdatedAt = field.NewInt64(tableName, "record_updated_at")

	_sqlite3KeywordDoc.fillFieldMap()

	return _sqlite3KeywordDoc
}

type sqlite3KeywordDoc struct {
	sqlite3KeywordDocDo sqlite3KeywordDocDo

	ALL             field.Asterisk
	PostID          field.Int64
	TileSplit       field.String
	ContentSplit    field.String
	Weight          field.Int
	PostUpdatedAt   field.Int64
	RecordCreatedAt field.Int64
	RecordUpdatedAt field.Int64

	fieldMap map[string]field.Expr
}

func (s sqlite3KeywordDoc) Table(newTableName string) *sqlite3KeywordDoc {
	s.sqlite3KeywordDocDo.UseTable(newTableName)
	return s.updateTableName(newTableName)
}

func (s sqlite3KeywordDoc) As(alias string) *sqlite3KeywordDoc {
	s.sqlite3KeywordDocDo.DO = *(s.sqlite3KeywordDocDo.As(alias).(*gen.DO))
	return s.updateTableName(alias)
}

func (s *sqlite3KeywordDoc) updateTableName(table string) *sqlite3KeywordDoc {
	s.ALL = field.NewAsterisk(table)
	s.PostID = field.NewInt64(table, "post_id")
	s.TileSplit = field.NewString(table, "title_split")
	s.ContentSplit = field.NewString(table, "content_split")
	s.Weight = field.NewInt(table, "weight")
	s.PostUpdatedAt = field.NewInt64(table, "post_updated_at")
	s.RecordCreatedAt = field.NewInt64(table, "record_created_at")
	s.RecordUpdatedAt = field.NewInt64(table, "record_updated_at")

	s.fillFieldMap()

	return s
}

func (s *sqlite3KeywordDoc) WithContext(ctx context.Context) ISqlite3KeywordDocDo {
	return s.sqlite3KeywordDocDo.WithContext(ctx)
}

func (s sqlite3KeywordDoc) TableName() string { return s.sqlite3KeywordDocDo.TableName() }

func (s sqlite3KeywordDoc) Alias() string { return s.sqlite3KeywordDocDo.Alias() }

func (s sqlite3KeywordDoc) Columns(cols ...field.Expr) gen.Columns {
	return s.sqlite3KeywordDocDo.Columns(cols...)
}

func (s *sqlite3KeywordDoc) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := s.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (s *sqlite3KeywordDoc) fillFieldMap() {
	s.fieldMap = make(map[string]field.Expr, 7)
	s.fieldMap["post_id"] = s.PostID
	s.fieldMap["title_split"] = s.TileSplit
	s.fieldMap["content_split"] = s.ContentSplit
	s.fieldMap["weight"] = s.Weight
	s.fieldMap["post_updated_at"] = s.PostUpdatedAt
	s.fieldMap["record_created_at"] = s.RecordCreatedAt
	s.fieldMap["record_updated_at"] = s.RecordUpdatedAt
}

func (s sqlite3KeywordDoc) clone(db *gorm.DB) sqlite3KeywordDoc {
	s.sqlite3KeywordDocDo.ReplaceConnPool(db.Statement.ConnPool)
	return s
}

func (s sqlite3KeywordDoc) replaceDB(db *gorm.DB) sqlite3KeywordDoc {
	s.sqlite3KeywordDocDo.ReplaceDB(db)
	return s
}

type sqlite3KeywordDocDo struct{ gen.DO }

type ISqlite3KeywordDocDo interface {
	gen.SubQuery
	Debug() ISqlite3KeywordDocDo
	WithContext(ctx context.Context) ISqlite3KeywordDocDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() ISqlite3KeywordDocDo
	WriteDB() ISqlite3KeywordDocDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) ISqlite3KeywordDocDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) ISqlite3KeywordDocDo
	Not(conds ...gen.Condition) ISqlite3KeywordDocDo
	Or(conds ...gen.Condition) ISqlite3KeywordDocDo
	Select(conds ...field.Expr) ISqlite3KeywordDocDo
	Where(conds ...gen.Condition) ISqlite3KeywordDocDo
	Order(conds ...field.Expr) ISqlite3KeywordDocDo
	Distinct(cols ...field.Expr) ISqlite3KeywordDocDo
	Omit(cols ...field.Expr) ISqlite3KeywordDocDo
	Join(table schema.Tabler, on ...field.Expr) ISqlite3KeywordDocDo
	LeftJoin(table schema.Tabler, on ...field.Expr) ISqlite3KeywordDocDo
	RightJoin(table schema.Tabler, on ...field.Expr) ISqlite3KeywordDocDo
	Group(cols ...field.Expr) ISqlite3KeywordDocDo
	Having(conds ...gen.Condition) ISqlite3KeywordDocDo
	Limit(limit int) ISqlite3KeywordDocDo
	Offset(offset int) ISqlite3KeywordDocDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) ISqlite3KeywordDocDo
	Unscoped() ISqlite3KeywordDocDo
	Create(values ...*model.Sqlite3KeywordDoc) error
	CreateInBatches(values []*model.Sqlite3KeywordDoc, batchSize int) error
	Save(values ...*model.Sqlite3KeywordDoc) error
	First() (*model.Sqlite3KeywordDoc, error)
	Take() (*model.Sqlite3KeywordDoc, error)
	Last() (*model.Sqlite3KeywordDoc, error)
	Find() ([]*model.Sqlite3KeywordDoc, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Sqlite3KeywordDoc, err error)
	FindInBatches(result *[]*model.Sqlite3KeywordDoc, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.Sqlite3KeywordDoc) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) ISqlite3KeywordDocDo
	Assign(attrs ...field.AssignExpr) ISqlite3KeywordDocDo
	Joins(fields ...field.RelationField) ISqlite3KeywordDocDo
	Preload(fields ...field.RelationField) ISqlite3KeywordDocDo
	FirstOrInit() (*model.Sqlite3KeywordDoc, error)
	FirstOrCreate() (*model.Sqlite3KeywordDoc, error)
	FindByPage(offset int, limit int) (result []*model.Sqlite3KeywordDoc, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) ISqlite3KeywordDocDo
	UnderlyingDB() *gorm.DB
	schema.Tabler

	SelectWithSorts(condList []model.WhereCond, sortList []model.SortColumn) (result []model.Sqlite3KeywordDoc, err error)
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
func (s sqlite3KeywordDocDo) SelectWithSorts(condList []model.WhereCond, sortList []model.SortColumn) (result []model.Sqlite3KeywordDoc, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("SELECT * FROM sqlite3_keyword_docs ")
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

func (s sqlite3KeywordDocDo) Debug() ISqlite3KeywordDocDo {
	return s.withDO(s.DO.Debug())
}

func (s sqlite3KeywordDocDo) WithContext(ctx context.Context) ISqlite3KeywordDocDo {
	return s.withDO(s.DO.WithContext(ctx))
}

func (s sqlite3KeywordDocDo) ReadDB() ISqlite3KeywordDocDo {
	return s.Clauses(dbresolver.Read)
}

func (s sqlite3KeywordDocDo) WriteDB() ISqlite3KeywordDocDo {
	return s.Clauses(dbresolver.Write)
}

func (s sqlite3KeywordDocDo) Session(config *gorm.Session) ISqlite3KeywordDocDo {
	return s.withDO(s.DO.Session(config))
}

func (s sqlite3KeywordDocDo) Clauses(conds ...clause.Expression) ISqlite3KeywordDocDo {
	return s.withDO(s.DO.Clauses(conds...))
}

func (s sqlite3KeywordDocDo) Returning(value interface{}, columns ...string) ISqlite3KeywordDocDo {
	return s.withDO(s.DO.Returning(value, columns...))
}

func (s sqlite3KeywordDocDo) Not(conds ...gen.Condition) ISqlite3KeywordDocDo {
	return s.withDO(s.DO.Not(conds...))
}

func (s sqlite3KeywordDocDo) Or(conds ...gen.Condition) ISqlite3KeywordDocDo {
	return s.withDO(s.DO.Or(conds...))
}

func (s sqlite3KeywordDocDo) Select(conds ...field.Expr) ISqlite3KeywordDocDo {
	return s.withDO(s.DO.Select(conds...))
}

func (s sqlite3KeywordDocDo) Where(conds ...gen.Condition) ISqlite3KeywordDocDo {
	return s.withDO(s.DO.Where(conds...))
}

func (s sqlite3KeywordDocDo) Order(conds ...field.Expr) ISqlite3KeywordDocDo {
	return s.withDO(s.DO.Order(conds...))
}

func (s sqlite3KeywordDocDo) Distinct(cols ...field.Expr) ISqlite3KeywordDocDo {
	return s.withDO(s.DO.Distinct(cols...))
}

func (s sqlite3KeywordDocDo) Omit(cols ...field.Expr) ISqlite3KeywordDocDo {
	return s.withDO(s.DO.Omit(cols...))
}

func (s sqlite3KeywordDocDo) Join(table schema.Tabler, on ...field.Expr) ISqlite3KeywordDocDo {
	return s.withDO(s.DO.Join(table, on...))
}

func (s sqlite3KeywordDocDo) LeftJoin(table schema.Tabler, on ...field.Expr) ISqlite3KeywordDocDo {
	return s.withDO(s.DO.LeftJoin(table, on...))
}

func (s sqlite3KeywordDocDo) RightJoin(table schema.Tabler, on ...field.Expr) ISqlite3KeywordDocDo {
	return s.withDO(s.DO.RightJoin(table, on...))
}

func (s sqlite3KeywordDocDo) Group(cols ...field.Expr) ISqlite3KeywordDocDo {
	return s.withDO(s.DO.Group(cols...))
}

func (s sqlite3KeywordDocDo) Having(conds ...gen.Condition) ISqlite3KeywordDocDo {
	return s.withDO(s.DO.Having(conds...))
}

func (s sqlite3KeywordDocDo) Limit(limit int) ISqlite3KeywordDocDo {
	return s.withDO(s.DO.Limit(limit))
}

func (s sqlite3KeywordDocDo) Offset(offset int) ISqlite3KeywordDocDo {
	return s.withDO(s.DO.Offset(offset))
}

func (s sqlite3KeywordDocDo) Scopes(funcs ...func(gen.Dao) gen.Dao) ISqlite3KeywordDocDo {
	return s.withDO(s.DO.Scopes(funcs...))
}

func (s sqlite3KeywordDocDo) Unscoped() ISqlite3KeywordDocDo {
	return s.withDO(s.DO.Unscoped())
}

func (s sqlite3KeywordDocDo) Create(values ...*model.Sqlite3KeywordDoc) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Create(values)
}

func (s sqlite3KeywordDocDo) CreateInBatches(values []*model.Sqlite3KeywordDoc, batchSize int) error {
	return s.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (s sqlite3KeywordDocDo) Save(values ...*model.Sqlite3KeywordDoc) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Save(values)
}

func (s sqlite3KeywordDocDo) First() (*model.Sqlite3KeywordDoc, error) {
	if result, err := s.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Sqlite3KeywordDoc), nil
	}
}

func (s sqlite3KeywordDocDo) Take() (*model.Sqlite3KeywordDoc, error) {
	if result, err := s.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Sqlite3KeywordDoc), nil
	}
}

func (s sqlite3KeywordDocDo) Last() (*model.Sqlite3KeywordDoc, error) {
	if result, err := s.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Sqlite3KeywordDoc), nil
	}
}

func (s sqlite3KeywordDocDo) Find() ([]*model.Sqlite3KeywordDoc, error) {
	result, err := s.DO.Find()
	return result.([]*model.Sqlite3KeywordDoc), err
}

func (s sqlite3KeywordDocDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Sqlite3KeywordDoc, err error) {
	buf := make([]*model.Sqlite3KeywordDoc, 0, batchSize)
	err = s.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (s sqlite3KeywordDocDo) FindInBatches(result *[]*model.Sqlite3KeywordDoc, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return s.DO.FindInBatches(result, batchSize, fc)
}

func (s sqlite3KeywordDocDo) Attrs(attrs ...field.AssignExpr) ISqlite3KeywordDocDo {
	return s.withDO(s.DO.Attrs(attrs...))
}

func (s sqlite3KeywordDocDo) Assign(attrs ...field.AssignExpr) ISqlite3KeywordDocDo {
	return s.withDO(s.DO.Assign(attrs...))
}

func (s sqlite3KeywordDocDo) Joins(fields ...field.RelationField) ISqlite3KeywordDocDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Joins(_f))
	}
	return &s
}

func (s sqlite3KeywordDocDo) Preload(fields ...field.RelationField) ISqlite3KeywordDocDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Preload(_f))
	}
	return &s
}

func (s sqlite3KeywordDocDo) FirstOrInit() (*model.Sqlite3KeywordDoc, error) {
	if result, err := s.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Sqlite3KeywordDoc), nil
	}
}

func (s sqlite3KeywordDocDo) FirstOrCreate() (*model.Sqlite3KeywordDoc, error) {
	if result, err := s.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Sqlite3KeywordDoc), nil
	}
}

func (s sqlite3KeywordDocDo) FindByPage(offset int, limit int) (result []*model.Sqlite3KeywordDoc, count int64, err error) {
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

func (s sqlite3KeywordDocDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = s.Count()
	if err != nil {
		return
	}

	err = s.Offset(offset).Limit(limit).Scan(result)
	return
}

func (s sqlite3KeywordDocDo) Scan(result interface{}) (err error) {
	return s.DO.Scan(result)
}

func (s sqlite3KeywordDocDo) Delete(models ...*model.Sqlite3KeywordDoc) (result gen.ResultInfo, err error) {
	return s.DO.Delete(models)
}

func (s *sqlite3KeywordDocDo) withDO(do gen.Dao) *sqlite3KeywordDocDo {
	s.DO = *do.(*gen.DO)
	return s
}
