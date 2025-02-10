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

func newOAuth2User(db *gorm.DB, opts ...gen.DOOption) oAuth2User {
	_oAuth2User := oAuth2User{}

	_oAuth2User.oAuth2UserDo.UseDB(db, opts...)
	_oAuth2User.oAuth2UserDo.UseModel(&model.OAuth2User{})

	tableName := _oAuth2User.oAuth2UserDo.TableName()
	_oAuth2User.ALL = field.NewAsterisk(tableName)
	_oAuth2User.ID = field.NewInt64(tableName, "id")
	_oAuth2User.CreatedAt = field.NewInt64(tableName, "created_at")
	_oAuth2User.UpdatedAt = field.NewInt64(tableName, "updated_at")
	_oAuth2User.Hide = field.NewInt(tableName, "hide")
	_oAuth2User.PlatformName = field.NewString(tableName, "platform_name")
	_oAuth2User.PlatformUserId = field.NewString(tableName, "platform_user_id")
	_oAuth2User.Username = field.NewString(tableName, "username")
	_oAuth2User.PrimaryEmail = field.NewString(tableName, "primary_email")
	_oAuth2User.AccessToken = field.NewString(tableName, "access_token")
	_oAuth2User.RefreshToken = field.NewString(tableName, "refresh_token")
	_oAuth2User.ExpiredAt = field.NewInt64(tableName, "expired_at")
	_oAuth2User.AvatarURL = field.NewString(tableName, "avatar_url")
	_oAuth2User.HomepageLink = field.NewString(tableName, "homepage_link")
	_oAuth2User.Scopes = field.NewField(tableName, "scopes")
	_oAuth2User.Enable = field.NewInt(tableName, "enable")

	_oAuth2User.fillFieldMap()

	return _oAuth2User
}

type oAuth2User struct {
	oAuth2UserDo oAuth2UserDo

	ALL            field.Asterisk
	ID             field.Int64
	CreatedAt      field.Int64
	UpdatedAt      field.Int64
	Hide           field.Int
	PlatformName   field.String
	PlatformUserId field.String
	Username       field.String
	PrimaryEmail   field.String
	AccessToken    field.String
	RefreshToken   field.String
	ExpiredAt      field.Int64
	AvatarURL      field.String
	HomepageLink   field.String
	Scopes         field.Field
	Enable         field.Int

	fieldMap map[string]field.Expr
}

func (o oAuth2User) Table(newTableName string) *oAuth2User {
	o.oAuth2UserDo.UseTable(newTableName)
	return o.updateTableName(newTableName)
}

func (o oAuth2User) As(alias string) *oAuth2User {
	o.oAuth2UserDo.DO = *(o.oAuth2UserDo.As(alias).(*gen.DO))
	return o.updateTableName(alias)
}

func (o *oAuth2User) updateTableName(table string) *oAuth2User {
	o.ALL = field.NewAsterisk(table)
	o.ID = field.NewInt64(table, "id")
	o.CreatedAt = field.NewInt64(table, "created_at")
	o.UpdatedAt = field.NewInt64(table, "updated_at")
	o.Hide = field.NewInt(table, "hide")
	o.PlatformName = field.NewString(table, "platform_name")
	o.PlatformUserId = field.NewString(table, "platform_user_id")
	o.Username = field.NewString(table, "username")
	o.PrimaryEmail = field.NewString(table, "primary_email")
	o.AccessToken = field.NewString(table, "access_token")
	o.RefreshToken = field.NewString(table, "refresh_token")
	o.ExpiredAt = field.NewInt64(table, "expired_at")
	o.AvatarURL = field.NewString(table, "avatar_url")
	o.HomepageLink = field.NewString(table, "homepage_link")
	o.Scopes = field.NewField(table, "scopes")
	o.Enable = field.NewInt(table, "enable")

	o.fillFieldMap()

	return o
}

func (o *oAuth2User) WithContext(ctx context.Context) IOAuth2UserDo {
	return o.oAuth2UserDo.WithContext(ctx)
}

func (o oAuth2User) TableName() string { return o.oAuth2UserDo.TableName() }

func (o oAuth2User) Alias() string { return o.oAuth2UserDo.Alias() }

func (o oAuth2User) Columns(cols ...field.Expr) gen.Columns { return o.oAuth2UserDo.Columns(cols...) }

func (o *oAuth2User) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := o.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (o *oAuth2User) fillFieldMap() {
	o.fieldMap = make(map[string]field.Expr, 15)
	o.fieldMap["id"] = o.ID
	o.fieldMap["created_at"] = o.CreatedAt
	o.fieldMap["updated_at"] = o.UpdatedAt
	o.fieldMap["hide"] = o.Hide
	o.fieldMap["platform_name"] = o.PlatformName
	o.fieldMap["platform_user_id"] = o.PlatformUserId
	o.fieldMap["username"] = o.Username
	o.fieldMap["primary_email"] = o.PrimaryEmail
	o.fieldMap["access_token"] = o.AccessToken
	o.fieldMap["refresh_token"] = o.RefreshToken
	o.fieldMap["expired_at"] = o.ExpiredAt
	o.fieldMap["avatar_url"] = o.AvatarURL
	o.fieldMap["homepage_link"] = o.HomepageLink
	o.fieldMap["scopes"] = o.Scopes
	o.fieldMap["enable"] = o.Enable
}

func (o oAuth2User) clone(db *gorm.DB) oAuth2User {
	o.oAuth2UserDo.ReplaceConnPool(db.Statement.ConnPool)
	return o
}

func (o oAuth2User) replaceDB(db *gorm.DB) oAuth2User {
	o.oAuth2UserDo.ReplaceDB(db)
	return o
}

type oAuth2UserDo struct{ gen.DO }

type IOAuth2UserDo interface {
	gen.SubQuery
	Debug() IOAuth2UserDo
	WithContext(ctx context.Context) IOAuth2UserDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IOAuth2UserDo
	WriteDB() IOAuth2UserDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IOAuth2UserDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IOAuth2UserDo
	Not(conds ...gen.Condition) IOAuth2UserDo
	Or(conds ...gen.Condition) IOAuth2UserDo
	Select(conds ...field.Expr) IOAuth2UserDo
	Where(conds ...gen.Condition) IOAuth2UserDo
	Order(conds ...field.Expr) IOAuth2UserDo
	Distinct(cols ...field.Expr) IOAuth2UserDo
	Omit(cols ...field.Expr) IOAuth2UserDo
	Join(table schema.Tabler, on ...field.Expr) IOAuth2UserDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IOAuth2UserDo
	RightJoin(table schema.Tabler, on ...field.Expr) IOAuth2UserDo
	Group(cols ...field.Expr) IOAuth2UserDo
	Having(conds ...gen.Condition) IOAuth2UserDo
	Limit(limit int) IOAuth2UserDo
	Offset(offset int) IOAuth2UserDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IOAuth2UserDo
	Unscoped() IOAuth2UserDo
	Create(values ...*model.OAuth2User) error
	CreateInBatches(values []*model.OAuth2User, batchSize int) error
	Save(values ...*model.OAuth2User) error
	First() (*model.OAuth2User, error)
	Take() (*model.OAuth2User, error)
	Last() (*model.OAuth2User, error)
	Find() ([]*model.OAuth2User, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.OAuth2User, err error)
	FindInBatches(result *[]*model.OAuth2User, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.OAuth2User) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IOAuth2UserDo
	Assign(attrs ...field.AssignExpr) IOAuth2UserDo
	Joins(fields ...field.RelationField) IOAuth2UserDo
	Preload(fields ...field.RelationField) IOAuth2UserDo
	FirstOrInit() (*model.OAuth2User, error)
	FirstOrCreate() (*model.OAuth2User, error)
	FindByPage(offset int, limit int) (result []*model.OAuth2User, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IOAuth2UserDo
	UnderlyingDB() *gorm.DB
	schema.Tabler

	SelectWithSorts(condList []model.WhereCond, sortList []model.SortColumn) (result []model.OAuth2User, err error)
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
func (o oAuth2UserDo) SelectWithSorts(condList []model.WhereCond, sortList []model.SortColumn) (result []model.OAuth2User, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("SELECT * FROM o_auth2_users ")
	if len(condList) > 0 {
		generateSQL.WriteString("WHERE ")
		for index, cond := range condList {
			if index < len(condList)-1 {
				params = append(params, cond.Val)
				generateSQL.WriteString(o.Quote(cond.Column) + " = ? AND ")
			} else {
				params = append(params, cond.Val)
				generateSQL.WriteString(o.Quote(cond.Column) + " = ? ")
			}
		}
	}
	if len(sortList) > 0 {
		generateSQL.WriteString("ORDER BY ")
		for index, sort := range sortList {
			if index < len(sortList)-1 {
				if sort.Desc {
					generateSQL.WriteString(o.Quote(sort.Column) + " DESC, ")
				} else {
					generateSQL.WriteString(o.Quote(sort.Column) + " ASC, ")
				}
			} else {
				if sort.Desc {
					generateSQL.WriteString(o.Quote(sort.Column) + " DESC ")
				} else {
					generateSQL.WriteString(o.Quote(sort.Column) + " ASC ")
				}
			}
		}
	}

	var executeSQL *gorm.DB
	executeSQL = o.UnderlyingDB().Raw(generateSQL.String(), params...).Find(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

func (o oAuth2UserDo) Debug() IOAuth2UserDo {
	return o.withDO(o.DO.Debug())
}

func (o oAuth2UserDo) WithContext(ctx context.Context) IOAuth2UserDo {
	return o.withDO(o.DO.WithContext(ctx))
}

func (o oAuth2UserDo) ReadDB() IOAuth2UserDo {
	return o.Clauses(dbresolver.Read)
}

func (o oAuth2UserDo) WriteDB() IOAuth2UserDo {
	return o.Clauses(dbresolver.Write)
}

func (o oAuth2UserDo) Session(config *gorm.Session) IOAuth2UserDo {
	return o.withDO(o.DO.Session(config))
}

func (o oAuth2UserDo) Clauses(conds ...clause.Expression) IOAuth2UserDo {
	return o.withDO(o.DO.Clauses(conds...))
}

func (o oAuth2UserDo) Returning(value interface{}, columns ...string) IOAuth2UserDo {
	return o.withDO(o.DO.Returning(value, columns...))
}

func (o oAuth2UserDo) Not(conds ...gen.Condition) IOAuth2UserDo {
	return o.withDO(o.DO.Not(conds...))
}

func (o oAuth2UserDo) Or(conds ...gen.Condition) IOAuth2UserDo {
	return o.withDO(o.DO.Or(conds...))
}

func (o oAuth2UserDo) Select(conds ...field.Expr) IOAuth2UserDo {
	return o.withDO(o.DO.Select(conds...))
}

func (o oAuth2UserDo) Where(conds ...gen.Condition) IOAuth2UserDo {
	return o.withDO(o.DO.Where(conds...))
}

func (o oAuth2UserDo) Order(conds ...field.Expr) IOAuth2UserDo {
	return o.withDO(o.DO.Order(conds...))
}

func (o oAuth2UserDo) Distinct(cols ...field.Expr) IOAuth2UserDo {
	return o.withDO(o.DO.Distinct(cols...))
}

func (o oAuth2UserDo) Omit(cols ...field.Expr) IOAuth2UserDo {
	return o.withDO(o.DO.Omit(cols...))
}

func (o oAuth2UserDo) Join(table schema.Tabler, on ...field.Expr) IOAuth2UserDo {
	return o.withDO(o.DO.Join(table, on...))
}

func (o oAuth2UserDo) LeftJoin(table schema.Tabler, on ...field.Expr) IOAuth2UserDo {
	return o.withDO(o.DO.LeftJoin(table, on...))
}

func (o oAuth2UserDo) RightJoin(table schema.Tabler, on ...field.Expr) IOAuth2UserDo {
	return o.withDO(o.DO.RightJoin(table, on...))
}

func (o oAuth2UserDo) Group(cols ...field.Expr) IOAuth2UserDo {
	return o.withDO(o.DO.Group(cols...))
}

func (o oAuth2UserDo) Having(conds ...gen.Condition) IOAuth2UserDo {
	return o.withDO(o.DO.Having(conds...))
}

func (o oAuth2UserDo) Limit(limit int) IOAuth2UserDo {
	return o.withDO(o.DO.Limit(limit))
}

func (o oAuth2UserDo) Offset(offset int) IOAuth2UserDo {
	return o.withDO(o.DO.Offset(offset))
}

func (o oAuth2UserDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IOAuth2UserDo {
	return o.withDO(o.DO.Scopes(funcs...))
}

func (o oAuth2UserDo) Unscoped() IOAuth2UserDo {
	return o.withDO(o.DO.Unscoped())
}

func (o oAuth2UserDo) Create(values ...*model.OAuth2User) error {
	if len(values) == 0 {
		return nil
	}
	return o.DO.Create(values)
}

func (o oAuth2UserDo) CreateInBatches(values []*model.OAuth2User, batchSize int) error {
	return o.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (o oAuth2UserDo) Save(values ...*model.OAuth2User) error {
	if len(values) == 0 {
		return nil
	}
	return o.DO.Save(values)
}

func (o oAuth2UserDo) First() (*model.OAuth2User, error) {
	if result, err := o.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.OAuth2User), nil
	}
}

func (o oAuth2UserDo) Take() (*model.OAuth2User, error) {
	if result, err := o.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.OAuth2User), nil
	}
}

func (o oAuth2UserDo) Last() (*model.OAuth2User, error) {
	if result, err := o.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.OAuth2User), nil
	}
}

func (o oAuth2UserDo) Find() ([]*model.OAuth2User, error) {
	result, err := o.DO.Find()
	return result.([]*model.OAuth2User), err
}

func (o oAuth2UserDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.OAuth2User, err error) {
	buf := make([]*model.OAuth2User, 0, batchSize)
	err = o.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (o oAuth2UserDo) FindInBatches(result *[]*model.OAuth2User, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return o.DO.FindInBatches(result, batchSize, fc)
}

func (o oAuth2UserDo) Attrs(attrs ...field.AssignExpr) IOAuth2UserDo {
	return o.withDO(o.DO.Attrs(attrs...))
}

func (o oAuth2UserDo) Assign(attrs ...field.AssignExpr) IOAuth2UserDo {
	return o.withDO(o.DO.Assign(attrs...))
}

func (o oAuth2UserDo) Joins(fields ...field.RelationField) IOAuth2UserDo {
	for _, _f := range fields {
		o = *o.withDO(o.DO.Joins(_f))
	}
	return &o
}

func (o oAuth2UserDo) Preload(fields ...field.RelationField) IOAuth2UserDo {
	for _, _f := range fields {
		o = *o.withDO(o.DO.Preload(_f))
	}
	return &o
}

func (o oAuth2UserDo) FirstOrInit() (*model.OAuth2User, error) {
	if result, err := o.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.OAuth2User), nil
	}
}

func (o oAuth2UserDo) FirstOrCreate() (*model.OAuth2User, error) {
	if result, err := o.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.OAuth2User), nil
	}
}

func (o oAuth2UserDo) FindByPage(offset int, limit int) (result []*model.OAuth2User, count int64, err error) {
	result, err = o.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = o.Offset(-1).Limit(-1).Count()
	return
}

func (o oAuth2UserDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = o.Count()
	if err != nil {
		return
	}

	err = o.Offset(offset).Limit(limit).Scan(result)
	return
}

func (o oAuth2UserDo) Scan(result interface{}) (err error) {
	return o.DO.Scan(result)
}

func (o oAuth2UserDo) Delete(models ...*model.OAuth2User) (result gen.ResultInfo, err error) {
	return o.DO.Delete(models)
}

func (o *oAuth2UserDo) withDO(do gen.Dao) *oAuth2UserDo {
	o.DO = *do.(*gen.DO)
	return o
}
