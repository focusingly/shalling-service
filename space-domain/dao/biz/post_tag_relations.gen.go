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

func newPostTagRelation(db *gorm.DB, opts ...gen.DOOption) postTagRelation {
	_postTagRelation := postTagRelation{}

	_postTagRelation.postTagRelationDo.UseDB(db, opts...)
	_postTagRelation.postTagRelationDo.UseModel(&model.PostTagRelation{})

	tableName := _postTagRelation.postTagRelationDo.TableName()
	_postTagRelation.ALL = field.NewAsterisk(tableName)
	_postTagRelation.ID = field.NewInt64(tableName, "id")
	_postTagRelation.CreatedAt = field.NewInt64(tableName, "created_at")
	_postTagRelation.UpdatedAt = field.NewInt64(tableName, "updated_at")
	_postTagRelation.Hide = field.NewInt(tableName, "hide")
	_postTagRelation.PostId = field.NewInt64(tableName, "post_id")
	_postTagRelation.TagId = field.NewInt64(tableName, "tag_id")

	_postTagRelation.fillFieldMap()

	return _postTagRelation
}

type postTagRelation struct {
	postTagRelationDo postTagRelationDo

	ALL       field.Asterisk
	ID        field.Int64
	CreatedAt field.Int64
	UpdatedAt field.Int64
	Hide      field.Int
	PostId    field.Int64
	TagId     field.Int64

	fieldMap map[string]field.Expr
}

func (p postTagRelation) Table(newTableName string) *postTagRelation {
	p.postTagRelationDo.UseTable(newTableName)
	return p.updateTableName(newTableName)
}

func (p postTagRelation) As(alias string) *postTagRelation {
	p.postTagRelationDo.DO = *(p.postTagRelationDo.As(alias).(*gen.DO))
	return p.updateTableName(alias)
}

func (p *postTagRelation) updateTableName(table string) *postTagRelation {
	p.ALL = field.NewAsterisk(table)
	p.ID = field.NewInt64(table, "id")
	p.CreatedAt = field.NewInt64(table, "created_at")
	p.UpdatedAt = field.NewInt64(table, "updated_at")
	p.Hide = field.NewInt(table, "hide")
	p.PostId = field.NewInt64(table, "post_id")
	p.TagId = field.NewInt64(table, "tag_id")

	p.fillFieldMap()

	return p
}

func (p *postTagRelation) WithContext(ctx context.Context) IPostTagRelationDo {
	return p.postTagRelationDo.WithContext(ctx)
}

func (p postTagRelation) TableName() string { return p.postTagRelationDo.TableName() }

func (p postTagRelation) Alias() string { return p.postTagRelationDo.Alias() }

func (p postTagRelation) Columns(cols ...field.Expr) gen.Columns {
	return p.postTagRelationDo.Columns(cols...)
}

func (p *postTagRelation) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := p.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (p *postTagRelation) fillFieldMap() {
	p.fieldMap = make(map[string]field.Expr, 6)
	p.fieldMap["id"] = p.ID
	p.fieldMap["created_at"] = p.CreatedAt
	p.fieldMap["updated_at"] = p.UpdatedAt
	p.fieldMap["hide"] = p.Hide
	p.fieldMap["post_id"] = p.PostId
	p.fieldMap["tag_id"] = p.TagId
}

func (p postTagRelation) clone(db *gorm.DB) postTagRelation {
	p.postTagRelationDo.ReplaceConnPool(db.Statement.ConnPool)
	return p
}

func (p postTagRelation) replaceDB(db *gorm.DB) postTagRelation {
	p.postTagRelationDo.ReplaceDB(db)
	return p
}

type postTagRelationDo struct{ gen.DO }

type IPostTagRelationDo interface {
	gen.SubQuery
	Debug() IPostTagRelationDo
	WithContext(ctx context.Context) IPostTagRelationDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IPostTagRelationDo
	WriteDB() IPostTagRelationDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IPostTagRelationDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IPostTagRelationDo
	Not(conds ...gen.Condition) IPostTagRelationDo
	Or(conds ...gen.Condition) IPostTagRelationDo
	Select(conds ...field.Expr) IPostTagRelationDo
	Where(conds ...gen.Condition) IPostTagRelationDo
	Order(conds ...field.Expr) IPostTagRelationDo
	Distinct(cols ...field.Expr) IPostTagRelationDo
	Omit(cols ...field.Expr) IPostTagRelationDo
	Join(table schema.Tabler, on ...field.Expr) IPostTagRelationDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IPostTagRelationDo
	RightJoin(table schema.Tabler, on ...field.Expr) IPostTagRelationDo
	Group(cols ...field.Expr) IPostTagRelationDo
	Having(conds ...gen.Condition) IPostTagRelationDo
	Limit(limit int) IPostTagRelationDo
	Offset(offset int) IPostTagRelationDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IPostTagRelationDo
	Unscoped() IPostTagRelationDo
	Create(values ...*model.PostTagRelation) error
	CreateInBatches(values []*model.PostTagRelation, batchSize int) error
	Save(values ...*model.PostTagRelation) error
	First() (*model.PostTagRelation, error)
	Take() (*model.PostTagRelation, error)
	Last() (*model.PostTagRelation, error)
	Find() ([]*model.PostTagRelation, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.PostTagRelation, err error)
	FindInBatches(result *[]*model.PostTagRelation, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.PostTagRelation) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IPostTagRelationDo
	Assign(attrs ...field.AssignExpr) IPostTagRelationDo
	Joins(fields ...field.RelationField) IPostTagRelationDo
	Preload(fields ...field.RelationField) IPostTagRelationDo
	FirstOrInit() (*model.PostTagRelation, error)
	FirstOrCreate() (*model.PostTagRelation, error)
	FindByPage(offset int, limit int) (result []*model.PostTagRelation, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IPostTagRelationDo
	UnderlyingDB() *gorm.DB
	schema.Tabler

	SelectWithSorts(condList []model.WhereCond, sortList []model.SortColumn) (result []model.PostTagRelation, err error)
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
func (p postTagRelationDo) SelectWithSorts(condList []model.WhereCond, sortList []model.SortColumn) (result []model.PostTagRelation, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("SELECT * FROM post_tag_relations ")
	if len(condList) > 0 {
		generateSQL.WriteString("WHERE ")
		for index, cond := range condList {
			if index < len(condList)-1 {
				params = append(params, cond.Val)
				generateSQL.WriteString(p.Quote(cond.Column) + " = ? AND ")
			} else {
				params = append(params, cond.Val)
				generateSQL.WriteString(p.Quote(cond.Column) + " = ? ")
			}
		}
	}
	if len(sortList) > 0 {
		generateSQL.WriteString("ORDER BY ")
		for index, sort := range sortList {
			if index < len(sortList)-1 {
				if sort.Desc {
					generateSQL.WriteString(p.Quote(sort.Column) + " DESC, ")
				} else {
					generateSQL.WriteString(p.Quote(sort.Column) + " ASC, ")
				}
			} else {
				if sort.Desc {
					generateSQL.WriteString(p.Quote(sort.Column) + " DESC ")
				} else {
					generateSQL.WriteString(p.Quote(sort.Column) + " ASC ")
				}
			}
		}
	}

	var executeSQL *gorm.DB
	executeSQL = p.UnderlyingDB().Raw(generateSQL.String(), params...).Find(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

func (p postTagRelationDo) Debug() IPostTagRelationDo {
	return p.withDO(p.DO.Debug())
}

func (p postTagRelationDo) WithContext(ctx context.Context) IPostTagRelationDo {
	return p.withDO(p.DO.WithContext(ctx))
}

func (p postTagRelationDo) ReadDB() IPostTagRelationDo {
	return p.Clauses(dbresolver.Read)
}

func (p postTagRelationDo) WriteDB() IPostTagRelationDo {
	return p.Clauses(dbresolver.Write)
}

func (p postTagRelationDo) Session(config *gorm.Session) IPostTagRelationDo {
	return p.withDO(p.DO.Session(config))
}

func (p postTagRelationDo) Clauses(conds ...clause.Expression) IPostTagRelationDo {
	return p.withDO(p.DO.Clauses(conds...))
}

func (p postTagRelationDo) Returning(value interface{}, columns ...string) IPostTagRelationDo {
	return p.withDO(p.DO.Returning(value, columns...))
}

func (p postTagRelationDo) Not(conds ...gen.Condition) IPostTagRelationDo {
	return p.withDO(p.DO.Not(conds...))
}

func (p postTagRelationDo) Or(conds ...gen.Condition) IPostTagRelationDo {
	return p.withDO(p.DO.Or(conds...))
}

func (p postTagRelationDo) Select(conds ...field.Expr) IPostTagRelationDo {
	return p.withDO(p.DO.Select(conds...))
}

func (p postTagRelationDo) Where(conds ...gen.Condition) IPostTagRelationDo {
	return p.withDO(p.DO.Where(conds...))
}

func (p postTagRelationDo) Order(conds ...field.Expr) IPostTagRelationDo {
	return p.withDO(p.DO.Order(conds...))
}

func (p postTagRelationDo) Distinct(cols ...field.Expr) IPostTagRelationDo {
	return p.withDO(p.DO.Distinct(cols...))
}

func (p postTagRelationDo) Omit(cols ...field.Expr) IPostTagRelationDo {
	return p.withDO(p.DO.Omit(cols...))
}

func (p postTagRelationDo) Join(table schema.Tabler, on ...field.Expr) IPostTagRelationDo {
	return p.withDO(p.DO.Join(table, on...))
}

func (p postTagRelationDo) LeftJoin(table schema.Tabler, on ...field.Expr) IPostTagRelationDo {
	return p.withDO(p.DO.LeftJoin(table, on...))
}

func (p postTagRelationDo) RightJoin(table schema.Tabler, on ...field.Expr) IPostTagRelationDo {
	return p.withDO(p.DO.RightJoin(table, on...))
}

func (p postTagRelationDo) Group(cols ...field.Expr) IPostTagRelationDo {
	return p.withDO(p.DO.Group(cols...))
}

func (p postTagRelationDo) Having(conds ...gen.Condition) IPostTagRelationDo {
	return p.withDO(p.DO.Having(conds...))
}

func (p postTagRelationDo) Limit(limit int) IPostTagRelationDo {
	return p.withDO(p.DO.Limit(limit))
}

func (p postTagRelationDo) Offset(offset int) IPostTagRelationDo {
	return p.withDO(p.DO.Offset(offset))
}

func (p postTagRelationDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IPostTagRelationDo {
	return p.withDO(p.DO.Scopes(funcs...))
}

func (p postTagRelationDo) Unscoped() IPostTagRelationDo {
	return p.withDO(p.DO.Unscoped())
}

func (p postTagRelationDo) Create(values ...*model.PostTagRelation) error {
	if len(values) == 0 {
		return nil
	}
	return p.DO.Create(values)
}

func (p postTagRelationDo) CreateInBatches(values []*model.PostTagRelation, batchSize int) error {
	return p.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (p postTagRelationDo) Save(values ...*model.PostTagRelation) error {
	if len(values) == 0 {
		return nil
	}
	return p.DO.Save(values)
}

func (p postTagRelationDo) First() (*model.PostTagRelation, error) {
	if result, err := p.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.PostTagRelation), nil
	}
}

func (p postTagRelationDo) Take() (*model.PostTagRelation, error) {
	if result, err := p.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.PostTagRelation), nil
	}
}

func (p postTagRelationDo) Last() (*model.PostTagRelation, error) {
	if result, err := p.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.PostTagRelation), nil
	}
}

func (p postTagRelationDo) Find() ([]*model.PostTagRelation, error) {
	result, err := p.DO.Find()
	return result.([]*model.PostTagRelation), err
}

func (p postTagRelationDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.PostTagRelation, err error) {
	buf := make([]*model.PostTagRelation, 0, batchSize)
	err = p.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (p postTagRelationDo) FindInBatches(result *[]*model.PostTagRelation, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return p.DO.FindInBatches(result, batchSize, fc)
}

func (p postTagRelationDo) Attrs(attrs ...field.AssignExpr) IPostTagRelationDo {
	return p.withDO(p.DO.Attrs(attrs...))
}

func (p postTagRelationDo) Assign(attrs ...field.AssignExpr) IPostTagRelationDo {
	return p.withDO(p.DO.Assign(attrs...))
}

func (p postTagRelationDo) Joins(fields ...field.RelationField) IPostTagRelationDo {
	for _, _f := range fields {
		p = *p.withDO(p.DO.Joins(_f))
	}
	return &p
}

func (p postTagRelationDo) Preload(fields ...field.RelationField) IPostTagRelationDo {
	for _, _f := range fields {
		p = *p.withDO(p.DO.Preload(_f))
	}
	return &p
}

func (p postTagRelationDo) FirstOrInit() (*model.PostTagRelation, error) {
	if result, err := p.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.PostTagRelation), nil
	}
}

func (p postTagRelationDo) FirstOrCreate() (*model.PostTagRelation, error) {
	if result, err := p.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.PostTagRelation), nil
	}
}

func (p postTagRelationDo) FindByPage(offset int, limit int) (result []*model.PostTagRelation, count int64, err error) {
	result, err = p.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = p.Offset(-1).Limit(-1).Count()
	return
}

func (p postTagRelationDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = p.Count()
	if err != nil {
		return
	}

	err = p.Offset(offset).Limit(limit).Scan(result)
	return
}

func (p postTagRelationDo) Scan(result interface{}) (err error) {
	return p.DO.Scan(result)
}

func (p postTagRelationDo) Delete(models ...*model.PostTagRelation) (result gen.ResultInfo, err error) {
	return p.DO.Delete(models)
}

func (p *postTagRelationDo) withDO(do gen.Dao) *postTagRelationDo {
	p.DO = *do.(*gen.DO)
	return p
}
