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

func newCronJob(db *gorm.DB, opts ...gen.DOOption) cronJob {
	_cronJob := cronJob{}

	_cronJob.cronJobDo.UseDB(db, opts...)
	_cronJob.cronJobDo.UseModel(&model.CronJob{})

	tableName := _cronJob.cronJobDo.TableName()
	_cronJob.ALL = field.NewAsterisk(tableName)
	_cronJob.ID = field.NewInt64(tableName, "id")
	_cronJob.CreatedAt = field.NewInt64(tableName, "created_at")
	_cronJob.UpdatedAt = field.NewInt64(tableName, "updated_at")
	_cronJob.Hide = field.NewInt(tableName, "hide")
	_cronJob.JobFuncName = field.NewString(tableName, "job_func_name")
	_cronJob.CronExpr = field.NewString(tableName, "cron_expr")
	_cronJob.Status = field.NewString(tableName, "status")
	_cronJob.Enable = field.NewInt(tableName, "enable")
	_cronJob.Mark = field.NewString(tableName, "mark")

	_cronJob.fillFieldMap()

	return _cronJob
}

type cronJob struct {
	cronJobDo cronJobDo

	ALL         field.Asterisk
	ID          field.Int64
	CreatedAt   field.Int64
	UpdatedAt   field.Int64
	Hide        field.Int
	JobFuncName field.String
	CronExpr    field.String
	Status      field.String
	Enable      field.Int
	Mark        field.String

	fieldMap map[string]field.Expr
}

func (c cronJob) Table(newTableName string) *cronJob {
	c.cronJobDo.UseTable(newTableName)
	return c.updateTableName(newTableName)
}

func (c cronJob) As(alias string) *cronJob {
	c.cronJobDo.DO = *(c.cronJobDo.As(alias).(*gen.DO))
	return c.updateTableName(alias)
}

func (c *cronJob) updateTableName(table string) *cronJob {
	c.ALL = field.NewAsterisk(table)
	c.ID = field.NewInt64(table, "id")
	c.CreatedAt = field.NewInt64(table, "created_at")
	c.UpdatedAt = field.NewInt64(table, "updated_at")
	c.Hide = field.NewInt(table, "hide")
	c.JobFuncName = field.NewString(table, "job_func_name")
	c.CronExpr = field.NewString(table, "cron_expr")
	c.Status = field.NewString(table, "status")
	c.Enable = field.NewInt(table, "enable")
	c.Mark = field.NewString(table, "mark")

	c.fillFieldMap()

	return c
}

func (c *cronJob) WithContext(ctx context.Context) ICronJobDo { return c.cronJobDo.WithContext(ctx) }

func (c cronJob) TableName() string { return c.cronJobDo.TableName() }

func (c cronJob) Alias() string { return c.cronJobDo.Alias() }

func (c cronJob) Columns(cols ...field.Expr) gen.Columns { return c.cronJobDo.Columns(cols...) }

func (c *cronJob) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := c.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (c *cronJob) fillFieldMap() {
	c.fieldMap = make(map[string]field.Expr, 9)
	c.fieldMap["id"] = c.ID
	c.fieldMap["created_at"] = c.CreatedAt
	c.fieldMap["updated_at"] = c.UpdatedAt
	c.fieldMap["hide"] = c.Hide
	c.fieldMap["job_func_name"] = c.JobFuncName
	c.fieldMap["cron_expr"] = c.CronExpr
	c.fieldMap["status"] = c.Status
	c.fieldMap["enable"] = c.Enable
	c.fieldMap["mark"] = c.Mark
}

func (c cronJob) clone(db *gorm.DB) cronJob {
	c.cronJobDo.ReplaceConnPool(db.Statement.ConnPool)
	return c
}

func (c cronJob) replaceDB(db *gorm.DB) cronJob {
	c.cronJobDo.ReplaceDB(db)
	return c
}

type cronJobDo struct{ gen.DO }

type ICronJobDo interface {
	gen.SubQuery
	Debug() ICronJobDo
	WithContext(ctx context.Context) ICronJobDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() ICronJobDo
	WriteDB() ICronJobDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) ICronJobDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) ICronJobDo
	Not(conds ...gen.Condition) ICronJobDo
	Or(conds ...gen.Condition) ICronJobDo
	Select(conds ...field.Expr) ICronJobDo
	Where(conds ...gen.Condition) ICronJobDo
	Order(conds ...field.Expr) ICronJobDo
	Distinct(cols ...field.Expr) ICronJobDo
	Omit(cols ...field.Expr) ICronJobDo
	Join(table schema.Tabler, on ...field.Expr) ICronJobDo
	LeftJoin(table schema.Tabler, on ...field.Expr) ICronJobDo
	RightJoin(table schema.Tabler, on ...field.Expr) ICronJobDo
	Group(cols ...field.Expr) ICronJobDo
	Having(conds ...gen.Condition) ICronJobDo
	Limit(limit int) ICronJobDo
	Offset(offset int) ICronJobDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) ICronJobDo
	Unscoped() ICronJobDo
	Create(values ...*model.CronJob) error
	CreateInBatches(values []*model.CronJob, batchSize int) error
	Save(values ...*model.CronJob) error
	First() (*model.CronJob, error)
	Take() (*model.CronJob, error)
	Last() (*model.CronJob, error)
	Find() ([]*model.CronJob, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.CronJob, err error)
	FindInBatches(result *[]*model.CronJob, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.CronJob) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) ICronJobDo
	Assign(attrs ...field.AssignExpr) ICronJobDo
	Joins(fields ...field.RelationField) ICronJobDo
	Preload(fields ...field.RelationField) ICronJobDo
	FirstOrInit() (*model.CronJob, error)
	FirstOrCreate() (*model.CronJob, error)
	FindByPage(offset int, limit int) (result []*model.CronJob, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) ICronJobDo
	UnderlyingDB() *gorm.DB
	schema.Tabler

	SelectWithSorts(condList []model.WhereCond, sortList []model.SortColumn) (result []model.CronJob, err error)
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
func (c cronJobDo) SelectWithSorts(condList []model.WhereCond, sortList []model.SortColumn) (result []model.CronJob, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("SELECT * FROM cron_jobs ")
	if len(condList) > 0 {
		generateSQL.WriteString("WHERE ")
		for index, cond := range condList {
			if index < len(condList)-1 {
				params = append(params, cond.Val)
				generateSQL.WriteString(c.Quote(cond.Column) + " = ? AND ")
			} else {
				params = append(params, cond.Val)
				generateSQL.WriteString(c.Quote(cond.Column) + " = ? ")
			}
		}
	}
	if len(sortList) > 0 {
		generateSQL.WriteString("ORDER BY ")
		for index, sort := range sortList {
			if index < len(sortList)-1 {
				if sort.Desc {
					generateSQL.WriteString(c.Quote(sort.Column) + " DESC, ")
				} else {
					generateSQL.WriteString(c.Quote(sort.Column) + " ASC, ")
				}
			} else {
				if sort.Desc {
					generateSQL.WriteString(c.Quote(sort.Column) + " DESC ")
				} else {
					generateSQL.WriteString(c.Quote(sort.Column) + " ASC ")
				}
			}
		}
	}

	var executeSQL *gorm.DB
	executeSQL = c.UnderlyingDB().Raw(generateSQL.String(), params...).Find(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

func (c cronJobDo) Debug() ICronJobDo {
	return c.withDO(c.DO.Debug())
}

func (c cronJobDo) WithContext(ctx context.Context) ICronJobDo {
	return c.withDO(c.DO.WithContext(ctx))
}

func (c cronJobDo) ReadDB() ICronJobDo {
	return c.Clauses(dbresolver.Read)
}

func (c cronJobDo) WriteDB() ICronJobDo {
	return c.Clauses(dbresolver.Write)
}

func (c cronJobDo) Session(config *gorm.Session) ICronJobDo {
	return c.withDO(c.DO.Session(config))
}

func (c cronJobDo) Clauses(conds ...clause.Expression) ICronJobDo {
	return c.withDO(c.DO.Clauses(conds...))
}

func (c cronJobDo) Returning(value interface{}, columns ...string) ICronJobDo {
	return c.withDO(c.DO.Returning(value, columns...))
}

func (c cronJobDo) Not(conds ...gen.Condition) ICronJobDo {
	return c.withDO(c.DO.Not(conds...))
}

func (c cronJobDo) Or(conds ...gen.Condition) ICronJobDo {
	return c.withDO(c.DO.Or(conds...))
}

func (c cronJobDo) Select(conds ...field.Expr) ICronJobDo {
	return c.withDO(c.DO.Select(conds...))
}

func (c cronJobDo) Where(conds ...gen.Condition) ICronJobDo {
	return c.withDO(c.DO.Where(conds...))
}

func (c cronJobDo) Order(conds ...field.Expr) ICronJobDo {
	return c.withDO(c.DO.Order(conds...))
}

func (c cronJobDo) Distinct(cols ...field.Expr) ICronJobDo {
	return c.withDO(c.DO.Distinct(cols...))
}

func (c cronJobDo) Omit(cols ...field.Expr) ICronJobDo {
	return c.withDO(c.DO.Omit(cols...))
}

func (c cronJobDo) Join(table schema.Tabler, on ...field.Expr) ICronJobDo {
	return c.withDO(c.DO.Join(table, on...))
}

func (c cronJobDo) LeftJoin(table schema.Tabler, on ...field.Expr) ICronJobDo {
	return c.withDO(c.DO.LeftJoin(table, on...))
}

func (c cronJobDo) RightJoin(table schema.Tabler, on ...field.Expr) ICronJobDo {
	return c.withDO(c.DO.RightJoin(table, on...))
}

func (c cronJobDo) Group(cols ...field.Expr) ICronJobDo {
	return c.withDO(c.DO.Group(cols...))
}

func (c cronJobDo) Having(conds ...gen.Condition) ICronJobDo {
	return c.withDO(c.DO.Having(conds...))
}

func (c cronJobDo) Limit(limit int) ICronJobDo {
	return c.withDO(c.DO.Limit(limit))
}

func (c cronJobDo) Offset(offset int) ICronJobDo {
	return c.withDO(c.DO.Offset(offset))
}

func (c cronJobDo) Scopes(funcs ...func(gen.Dao) gen.Dao) ICronJobDo {
	return c.withDO(c.DO.Scopes(funcs...))
}

func (c cronJobDo) Unscoped() ICronJobDo {
	return c.withDO(c.DO.Unscoped())
}

func (c cronJobDo) Create(values ...*model.CronJob) error {
	if len(values) == 0 {
		return nil
	}
	return c.DO.Create(values)
}

func (c cronJobDo) CreateInBatches(values []*model.CronJob, batchSize int) error {
	return c.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (c cronJobDo) Save(values ...*model.CronJob) error {
	if len(values) == 0 {
		return nil
	}
	return c.DO.Save(values)
}

func (c cronJobDo) First() (*model.CronJob, error) {
	if result, err := c.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.CronJob), nil
	}
}

func (c cronJobDo) Take() (*model.CronJob, error) {
	if result, err := c.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.CronJob), nil
	}
}

func (c cronJobDo) Last() (*model.CronJob, error) {
	if result, err := c.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.CronJob), nil
	}
}

func (c cronJobDo) Find() ([]*model.CronJob, error) {
	result, err := c.DO.Find()
	return result.([]*model.CronJob), err
}

func (c cronJobDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.CronJob, err error) {
	buf := make([]*model.CronJob, 0, batchSize)
	err = c.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (c cronJobDo) FindInBatches(result *[]*model.CronJob, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return c.DO.FindInBatches(result, batchSize, fc)
}

func (c cronJobDo) Attrs(attrs ...field.AssignExpr) ICronJobDo {
	return c.withDO(c.DO.Attrs(attrs...))
}

func (c cronJobDo) Assign(attrs ...field.AssignExpr) ICronJobDo {
	return c.withDO(c.DO.Assign(attrs...))
}

func (c cronJobDo) Joins(fields ...field.RelationField) ICronJobDo {
	for _, _f := range fields {
		c = *c.withDO(c.DO.Joins(_f))
	}
	return &c
}

func (c cronJobDo) Preload(fields ...field.RelationField) ICronJobDo {
	for _, _f := range fields {
		c = *c.withDO(c.DO.Preload(_f))
	}
	return &c
}

func (c cronJobDo) FirstOrInit() (*model.CronJob, error) {
	if result, err := c.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.CronJob), nil
	}
}

func (c cronJobDo) FirstOrCreate() (*model.CronJob, error) {
	if result, err := c.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.CronJob), nil
	}
}

func (c cronJobDo) FindByPage(offset int, limit int) (result []*model.CronJob, count int64, err error) {
	result, err = c.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = c.Offset(-1).Limit(-1).Count()
	return
}

func (c cronJobDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = c.Count()
	if err != nil {
		return
	}

	err = c.Offset(offset).Limit(limit).Scan(result)
	return
}

func (c cronJobDo) Scan(result interface{}) (err error) {
	return c.DO.Scan(result)
}

func (c cronJobDo) Delete(models ...*model.CronJob) (result gen.ResultInfo, err error) {
	return c.DO.Delete(models)
}

func (c *cronJobDo) withDO(do gen.Dao) *cronJobDo {
	c.DO = *do.(*gen.DO)
	return c
}
