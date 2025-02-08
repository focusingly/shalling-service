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

func newUserLoginSession(db *gorm.DB, opts ...gen.DOOption) userLoginSession {
	_userLoginSession := userLoginSession{}

	_userLoginSession.userLoginSessionDo.UseDB(db, opts...)
	_userLoginSession.userLoginSessionDo.UseModel(&model.UserLoginSession{})

	tableName := _userLoginSession.userLoginSessionDo.TableName()
	_userLoginSession.ALL = field.NewAsterisk(tableName)
	_userLoginSession.ID = field.NewInt64(tableName, "id")
	_userLoginSession.CreatedAt = field.NewInt64(tableName, "created_at")
	_userLoginSession.UpdatedAt = field.NewInt64(tableName, "updated_at")
	_userLoginSession.Hide = field.NewInt(tableName, "hide")
	_userLoginSession.UserId = field.NewInt64(tableName, "user_id")
	_userLoginSession.UUID = field.NewString(tableName, "uuid")
	_userLoginSession.IpU32Val = field.NewUint32(tableName, "ip_u32_val")
	_userLoginSession.IpAddress = field.NewString(tableName, "ip_address")
	_userLoginSession.IpSource = field.NewString(tableName, "ip_source")
	_userLoginSession.ExpiredAt = field.NewInt64(tableName, "expired_at")
	_userLoginSession.UserType = field.NewString(tableName, "user_type")
	_userLoginSession.Token = field.NewString(tableName, "token")
	_userLoginSession.Useragent = field.NewString(tableName, "useragent")
	_userLoginSession.ClientName = field.NewString(tableName, "client_name")
	_userLoginSession.OsName = field.NewString(tableName, "os_name")

	_userLoginSession.fillFieldMap()

	return _userLoginSession
}

type userLoginSession struct {
	userLoginSessionDo userLoginSessionDo

	ALL        field.Asterisk
	ID         field.Int64
	CreatedAt  field.Int64
	UpdatedAt  field.Int64
	Hide       field.Int
	UserId     field.Int64
	UUID       field.String
	IpU32Val   field.Uint32
	IpAddress  field.String
	IpSource   field.String
	ExpiredAt  field.Int64
	UserType   field.String
	Token      field.String
	Useragent  field.String
	ClientName field.String
	OsName     field.String

	fieldMap map[string]field.Expr
}

func (u userLoginSession) Table(newTableName string) *userLoginSession {
	u.userLoginSessionDo.UseTable(newTableName)
	return u.updateTableName(newTableName)
}

func (u userLoginSession) As(alias string) *userLoginSession {
	u.userLoginSessionDo.DO = *(u.userLoginSessionDo.As(alias).(*gen.DO))
	return u.updateTableName(alias)
}

func (u *userLoginSession) updateTableName(table string) *userLoginSession {
	u.ALL = field.NewAsterisk(table)
	u.ID = field.NewInt64(table, "id")
	u.CreatedAt = field.NewInt64(table, "created_at")
	u.UpdatedAt = field.NewInt64(table, "updated_at")
	u.Hide = field.NewInt(table, "hide")
	u.UserId = field.NewInt64(table, "user_id")
	u.UUID = field.NewString(table, "uuid")
	u.IpU32Val = field.NewUint32(table, "ip_u32_val")
	u.IpAddress = field.NewString(table, "ip_address")
	u.IpSource = field.NewString(table, "ip_source")
	u.ExpiredAt = field.NewInt64(table, "expired_at")
	u.UserType = field.NewString(table, "user_type")
	u.Token = field.NewString(table, "token")
	u.Useragent = field.NewString(table, "useragent")
	u.ClientName = field.NewString(table, "client_name")
	u.OsName = field.NewString(table, "os_name")

	u.fillFieldMap()

	return u
}

func (u *userLoginSession) WithContext(ctx context.Context) IUserLoginSessionDo {
	return u.userLoginSessionDo.WithContext(ctx)
}

func (u userLoginSession) TableName() string { return u.userLoginSessionDo.TableName() }

func (u userLoginSession) Alias() string { return u.userLoginSessionDo.Alias() }

func (u userLoginSession) Columns(cols ...field.Expr) gen.Columns {
	return u.userLoginSessionDo.Columns(cols...)
}

func (u *userLoginSession) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := u.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (u *userLoginSession) fillFieldMap() {
	u.fieldMap = make(map[string]field.Expr, 15)
	u.fieldMap["id"] = u.ID
	u.fieldMap["created_at"] = u.CreatedAt
	u.fieldMap["updated_at"] = u.UpdatedAt
	u.fieldMap["hide"] = u.Hide
	u.fieldMap["user_id"] = u.UserId
	u.fieldMap["uuid"] = u.UUID
	u.fieldMap["ip_u32_val"] = u.IpU32Val
	u.fieldMap["ip_address"] = u.IpAddress
	u.fieldMap["ip_source"] = u.IpSource
	u.fieldMap["expired_at"] = u.ExpiredAt
	u.fieldMap["user_type"] = u.UserType
	u.fieldMap["token"] = u.Token
	u.fieldMap["useragent"] = u.Useragent
	u.fieldMap["client_name"] = u.ClientName
	u.fieldMap["os_name"] = u.OsName
}

func (u userLoginSession) clone(db *gorm.DB) userLoginSession {
	u.userLoginSessionDo.ReplaceConnPool(db.Statement.ConnPool)
	return u
}

func (u userLoginSession) replaceDB(db *gorm.DB) userLoginSession {
	u.userLoginSessionDo.ReplaceDB(db)
	return u
}

type userLoginSessionDo struct{ gen.DO }

type IUserLoginSessionDo interface {
	gen.SubQuery
	Debug() IUserLoginSessionDo
	WithContext(ctx context.Context) IUserLoginSessionDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IUserLoginSessionDo
	WriteDB() IUserLoginSessionDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IUserLoginSessionDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IUserLoginSessionDo
	Not(conds ...gen.Condition) IUserLoginSessionDo
	Or(conds ...gen.Condition) IUserLoginSessionDo
	Select(conds ...field.Expr) IUserLoginSessionDo
	Where(conds ...gen.Condition) IUserLoginSessionDo
	Order(conds ...field.Expr) IUserLoginSessionDo
	Distinct(cols ...field.Expr) IUserLoginSessionDo
	Omit(cols ...field.Expr) IUserLoginSessionDo
	Join(table schema.Tabler, on ...field.Expr) IUserLoginSessionDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IUserLoginSessionDo
	RightJoin(table schema.Tabler, on ...field.Expr) IUserLoginSessionDo
	Group(cols ...field.Expr) IUserLoginSessionDo
	Having(conds ...gen.Condition) IUserLoginSessionDo
	Limit(limit int) IUserLoginSessionDo
	Offset(offset int) IUserLoginSessionDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IUserLoginSessionDo
	Unscoped() IUserLoginSessionDo
	Create(values ...*model.UserLoginSession) error
	CreateInBatches(values []*model.UserLoginSession, batchSize int) error
	Save(values ...*model.UserLoginSession) error
	First() (*model.UserLoginSession, error)
	Take() (*model.UserLoginSession, error)
	Last() (*model.UserLoginSession, error)
	Find() ([]*model.UserLoginSession, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.UserLoginSession, err error)
	FindInBatches(result *[]*model.UserLoginSession, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.UserLoginSession) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IUserLoginSessionDo
	Assign(attrs ...field.AssignExpr) IUserLoginSessionDo
	Joins(fields ...field.RelationField) IUserLoginSessionDo
	Preload(fields ...field.RelationField) IUserLoginSessionDo
	FirstOrInit() (*model.UserLoginSession, error)
	FirstOrCreate() (*model.UserLoginSession, error)
	FindByPage(offset int, limit int) (result []*model.UserLoginSession, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IUserLoginSessionDo
	UnderlyingDB() *gorm.DB
	schema.Tabler

	SelectWithSorts(condList []model.WhereCond, sortList []model.SortColumn) (result []model.UserLoginSession, err error)
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
func (u userLoginSessionDo) SelectWithSorts(condList []model.WhereCond, sortList []model.SortColumn) (result []model.UserLoginSession, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("SELECT * FROM user_login_sessions ")
	if len(condList) > 0 {
		generateSQL.WriteString("WHERE ")
		for index, cond := range condList {
			if index < len(condList)-1 {
				params = append(params, cond.Val)
				generateSQL.WriteString(u.Quote(cond.Column) + " = ? AND ")
			} else {
				params = append(params, cond.Val)
				generateSQL.WriteString(u.Quote(cond.Column) + " = ? ")
			}
		}
	}
	if len(sortList) > 0 {
		generateSQL.WriteString("ORDER BY ")
		for index, sort := range sortList {
			if index < len(sortList)-1 {
				if sort.Desc {
					generateSQL.WriteString(u.Quote(sort.Column) + " DESC, ")
				} else {
					generateSQL.WriteString(u.Quote(sort.Column) + " ASC, ")
				}
			} else {
				if sort.Desc {
					generateSQL.WriteString(u.Quote(sort.Column) + " DESC ")
				} else {
					generateSQL.WriteString(u.Quote(sort.Column) + " ASC ")
				}
			}
		}
	}

	var executeSQL *gorm.DB
	executeSQL = u.UnderlyingDB().Raw(generateSQL.String(), params...).Find(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

func (u userLoginSessionDo) Debug() IUserLoginSessionDo {
	return u.withDO(u.DO.Debug())
}

func (u userLoginSessionDo) WithContext(ctx context.Context) IUserLoginSessionDo {
	return u.withDO(u.DO.WithContext(ctx))
}

func (u userLoginSessionDo) ReadDB() IUserLoginSessionDo {
	return u.Clauses(dbresolver.Read)
}

func (u userLoginSessionDo) WriteDB() IUserLoginSessionDo {
	return u.Clauses(dbresolver.Write)
}

func (u userLoginSessionDo) Session(config *gorm.Session) IUserLoginSessionDo {
	return u.withDO(u.DO.Session(config))
}

func (u userLoginSessionDo) Clauses(conds ...clause.Expression) IUserLoginSessionDo {
	return u.withDO(u.DO.Clauses(conds...))
}

func (u userLoginSessionDo) Returning(value interface{}, columns ...string) IUserLoginSessionDo {
	return u.withDO(u.DO.Returning(value, columns...))
}

func (u userLoginSessionDo) Not(conds ...gen.Condition) IUserLoginSessionDo {
	return u.withDO(u.DO.Not(conds...))
}

func (u userLoginSessionDo) Or(conds ...gen.Condition) IUserLoginSessionDo {
	return u.withDO(u.DO.Or(conds...))
}

func (u userLoginSessionDo) Select(conds ...field.Expr) IUserLoginSessionDo {
	return u.withDO(u.DO.Select(conds...))
}

func (u userLoginSessionDo) Where(conds ...gen.Condition) IUserLoginSessionDo {
	return u.withDO(u.DO.Where(conds...))
}

func (u userLoginSessionDo) Order(conds ...field.Expr) IUserLoginSessionDo {
	return u.withDO(u.DO.Order(conds...))
}

func (u userLoginSessionDo) Distinct(cols ...field.Expr) IUserLoginSessionDo {
	return u.withDO(u.DO.Distinct(cols...))
}

func (u userLoginSessionDo) Omit(cols ...field.Expr) IUserLoginSessionDo {
	return u.withDO(u.DO.Omit(cols...))
}

func (u userLoginSessionDo) Join(table schema.Tabler, on ...field.Expr) IUserLoginSessionDo {
	return u.withDO(u.DO.Join(table, on...))
}

func (u userLoginSessionDo) LeftJoin(table schema.Tabler, on ...field.Expr) IUserLoginSessionDo {
	return u.withDO(u.DO.LeftJoin(table, on...))
}

func (u userLoginSessionDo) RightJoin(table schema.Tabler, on ...field.Expr) IUserLoginSessionDo {
	return u.withDO(u.DO.RightJoin(table, on...))
}

func (u userLoginSessionDo) Group(cols ...field.Expr) IUserLoginSessionDo {
	return u.withDO(u.DO.Group(cols...))
}

func (u userLoginSessionDo) Having(conds ...gen.Condition) IUserLoginSessionDo {
	return u.withDO(u.DO.Having(conds...))
}

func (u userLoginSessionDo) Limit(limit int) IUserLoginSessionDo {
	return u.withDO(u.DO.Limit(limit))
}

func (u userLoginSessionDo) Offset(offset int) IUserLoginSessionDo {
	return u.withDO(u.DO.Offset(offset))
}

func (u userLoginSessionDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IUserLoginSessionDo {
	return u.withDO(u.DO.Scopes(funcs...))
}

func (u userLoginSessionDo) Unscoped() IUserLoginSessionDo {
	return u.withDO(u.DO.Unscoped())
}

func (u userLoginSessionDo) Create(values ...*model.UserLoginSession) error {
	if len(values) == 0 {
		return nil
	}
	return u.DO.Create(values)
}

func (u userLoginSessionDo) CreateInBatches(values []*model.UserLoginSession, batchSize int) error {
	return u.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (u userLoginSessionDo) Save(values ...*model.UserLoginSession) error {
	if len(values) == 0 {
		return nil
	}
	return u.DO.Save(values)
}

func (u userLoginSessionDo) First() (*model.UserLoginSession, error) {
	if result, err := u.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserLoginSession), nil
	}
}

func (u userLoginSessionDo) Take() (*model.UserLoginSession, error) {
	if result, err := u.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserLoginSession), nil
	}
}

func (u userLoginSessionDo) Last() (*model.UserLoginSession, error) {
	if result, err := u.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserLoginSession), nil
	}
}

func (u userLoginSessionDo) Find() ([]*model.UserLoginSession, error) {
	result, err := u.DO.Find()
	return result.([]*model.UserLoginSession), err
}

func (u userLoginSessionDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.UserLoginSession, err error) {
	buf := make([]*model.UserLoginSession, 0, batchSize)
	err = u.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (u userLoginSessionDo) FindInBatches(result *[]*model.UserLoginSession, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return u.DO.FindInBatches(result, batchSize, fc)
}

func (u userLoginSessionDo) Attrs(attrs ...field.AssignExpr) IUserLoginSessionDo {
	return u.withDO(u.DO.Attrs(attrs...))
}

func (u userLoginSessionDo) Assign(attrs ...field.AssignExpr) IUserLoginSessionDo {
	return u.withDO(u.DO.Assign(attrs...))
}

func (u userLoginSessionDo) Joins(fields ...field.RelationField) IUserLoginSessionDo {
	for _, _f := range fields {
		u = *u.withDO(u.DO.Joins(_f))
	}
	return &u
}

func (u userLoginSessionDo) Preload(fields ...field.RelationField) IUserLoginSessionDo {
	for _, _f := range fields {
		u = *u.withDO(u.DO.Preload(_f))
	}
	return &u
}

func (u userLoginSessionDo) FirstOrInit() (*model.UserLoginSession, error) {
	if result, err := u.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserLoginSession), nil
	}
}

func (u userLoginSessionDo) FirstOrCreate() (*model.UserLoginSession, error) {
	if result, err := u.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserLoginSession), nil
	}
}

func (u userLoginSessionDo) FindByPage(offset int, limit int) (result []*model.UserLoginSession, count int64, err error) {
	result, err = u.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = u.Offset(-1).Limit(-1).Count()
	return
}

func (u userLoginSessionDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = u.Count()
	if err != nil {
		return
	}

	err = u.Offset(offset).Limit(limit).Scan(result)
	return
}

func (u userLoginSessionDo) Scan(result interface{}) (err error) {
	return u.DO.Scan(result)
}

func (u userLoginSessionDo) Delete(models ...*model.UserLoginSession) (result gen.ResultInfo, err error) {
	return u.DO.Delete(models)
}

func (u *userLoginSessionDo) withDO(do gen.Dao) *userLoginSessionDo {
	u.DO = *do.(*gen.DO)
	return u
}
