// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"checkin/query/model"
)

func newUserCheckinMachineInfo(db *gorm.DB, opts ...gen.DOOption) userCheckinMachineInfo {
	_userCheckinMachineInfo := userCheckinMachineInfo{}

	_userCheckinMachineInfo.userCheckinMachineInfoDo.UseDB(db, opts...)
	_userCheckinMachineInfo.userCheckinMachineInfoDo.UseModel(&model.UserCheckinMachineInfo{})

	tableName := _userCheckinMachineInfo.userCheckinMachineInfoDo.TableName()
	_userCheckinMachineInfo.ALL = field.NewAsterisk(tableName)
	_userCheckinMachineInfo.ID = field.NewInt(tableName, "id")
	_userCheckinMachineInfo.Sn = field.NewString(tableName, "sn")
	_userCheckinMachineInfo.Enrollid = field.NewInt(tableName, "enrollid")
	_userCheckinMachineInfo.Name = field.NewString(tableName, "name")
	_userCheckinMachineInfo.Backupnum = field.NewInt(tableName, "backupnum")
	_userCheckinMachineInfo.Admin = field.NewInt(tableName, "admin")
	_userCheckinMachineInfo.Record = field.NewString(tableName, "record")
	_userCheckinMachineInfo.Status = field.NewInt(tableName, "status")
	_userCheckinMachineInfo.CreatedAt = field.NewTime(tableName, "created_at")
	_userCheckinMachineInfo.UpdatedAt = field.NewTime(tableName, "updated_at")
	_userCheckinMachineInfo.DeletedAt = field.NewField(tableName, "deleted_at")

	_userCheckinMachineInfo.fillFieldMap()

	return _userCheckinMachineInfo
}

type userCheckinMachineInfo struct {
	userCheckinMachineInfoDo userCheckinMachineInfoDo

	ALL       field.Asterisk
	ID        field.Int
	Sn        field.String
	Enrollid  field.Int
	Name      field.String
	Backupnum field.Int
	Admin     field.Int
	Record    field.String
	Status    field.Int
	CreatedAt field.Time
	UpdatedAt field.Time
	DeletedAt field.Field

	fieldMap map[string]field.Expr
}

func (u userCheckinMachineInfo) Table(newTableName string) *userCheckinMachineInfo {
	u.userCheckinMachineInfoDo.UseTable(newTableName)
	return u.updateTableName(newTableName)
}

func (u userCheckinMachineInfo) As(alias string) *userCheckinMachineInfo {
	u.userCheckinMachineInfoDo.DO = *(u.userCheckinMachineInfoDo.As(alias).(*gen.DO))
	return u.updateTableName(alias)
}

func (u *userCheckinMachineInfo) updateTableName(table string) *userCheckinMachineInfo {
	u.ALL = field.NewAsterisk(table)
	u.ID = field.NewInt(table, "id")
	u.Sn = field.NewString(table, "sn")
	u.Enrollid = field.NewInt(table, "enrollid")
	u.Name = field.NewString(table, "name")
	u.Backupnum = field.NewInt(table, "backupnum")
	u.Admin = field.NewInt(table, "admin")
	u.Record = field.NewString(table, "record")
	u.Status = field.NewInt(table, "status")
	u.CreatedAt = field.NewTime(table, "created_at")
	u.UpdatedAt = field.NewTime(table, "updated_at")
	u.DeletedAt = field.NewField(table, "deleted_at")

	u.fillFieldMap()

	return u
}

func (u *userCheckinMachineInfo) WithContext(ctx context.Context) IUserCheckinMachineInfoDo {
	return u.userCheckinMachineInfoDo.WithContext(ctx)
}

func (u userCheckinMachineInfo) TableName() string { return u.userCheckinMachineInfoDo.TableName() }

func (u userCheckinMachineInfo) Alias() string { return u.userCheckinMachineInfoDo.Alias() }

func (u userCheckinMachineInfo) Columns(cols ...field.Expr) gen.Columns {
	return u.userCheckinMachineInfoDo.Columns(cols...)
}

func (u *userCheckinMachineInfo) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := u.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (u *userCheckinMachineInfo) fillFieldMap() {
	u.fieldMap = make(map[string]field.Expr, 11)
	u.fieldMap["id"] = u.ID
	u.fieldMap["sn"] = u.Sn
	u.fieldMap["enrollid"] = u.Enrollid
	u.fieldMap["name"] = u.Name
	u.fieldMap["backupnum"] = u.Backupnum
	u.fieldMap["admin"] = u.Admin
	u.fieldMap["record"] = u.Record
	u.fieldMap["status"] = u.Status
	u.fieldMap["created_at"] = u.CreatedAt
	u.fieldMap["updated_at"] = u.UpdatedAt
	u.fieldMap["deleted_at"] = u.DeletedAt
}

func (u userCheckinMachineInfo) clone(db *gorm.DB) userCheckinMachineInfo {
	u.userCheckinMachineInfoDo.ReplaceConnPool(db.Statement.ConnPool)
	return u
}

func (u userCheckinMachineInfo) replaceDB(db *gorm.DB) userCheckinMachineInfo {
	u.userCheckinMachineInfoDo.ReplaceDB(db)
	return u
}

type userCheckinMachineInfoDo struct{ gen.DO }

type IUserCheckinMachineInfoDo interface {
	gen.SubQuery
	Debug() IUserCheckinMachineInfoDo
	WithContext(ctx context.Context) IUserCheckinMachineInfoDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IUserCheckinMachineInfoDo
	WriteDB() IUserCheckinMachineInfoDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IUserCheckinMachineInfoDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IUserCheckinMachineInfoDo
	Not(conds ...gen.Condition) IUserCheckinMachineInfoDo
	Or(conds ...gen.Condition) IUserCheckinMachineInfoDo
	Select(conds ...field.Expr) IUserCheckinMachineInfoDo
	Where(conds ...gen.Condition) IUserCheckinMachineInfoDo
	Order(conds ...field.Expr) IUserCheckinMachineInfoDo
	Distinct(cols ...field.Expr) IUserCheckinMachineInfoDo
	Omit(cols ...field.Expr) IUserCheckinMachineInfoDo
	Join(table schema.Tabler, on ...field.Expr) IUserCheckinMachineInfoDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IUserCheckinMachineInfoDo
	RightJoin(table schema.Tabler, on ...field.Expr) IUserCheckinMachineInfoDo
	Group(cols ...field.Expr) IUserCheckinMachineInfoDo
	Having(conds ...gen.Condition) IUserCheckinMachineInfoDo
	Limit(limit int) IUserCheckinMachineInfoDo
	Offset(offset int) IUserCheckinMachineInfoDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IUserCheckinMachineInfoDo
	Unscoped() IUserCheckinMachineInfoDo
	Create(values ...*model.UserCheckinMachineInfo) error
	CreateInBatches(values []*model.UserCheckinMachineInfo, batchSize int) error
	Save(values ...*model.UserCheckinMachineInfo) error
	First() (*model.UserCheckinMachineInfo, error)
	Take() (*model.UserCheckinMachineInfo, error)
	Last() (*model.UserCheckinMachineInfo, error)
	Find() ([]*model.UserCheckinMachineInfo, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.UserCheckinMachineInfo, err error)
	FindInBatches(result *[]*model.UserCheckinMachineInfo, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.UserCheckinMachineInfo) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IUserCheckinMachineInfoDo
	Assign(attrs ...field.AssignExpr) IUserCheckinMachineInfoDo
	Joins(fields ...field.RelationField) IUserCheckinMachineInfoDo
	Preload(fields ...field.RelationField) IUserCheckinMachineInfoDo
	FirstOrInit() (*model.UserCheckinMachineInfo, error)
	FirstOrCreate() (*model.UserCheckinMachineInfo, error)
	FindByPage(offset int, limit int) (result []*model.UserCheckinMachineInfo, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IUserCheckinMachineInfoDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (u userCheckinMachineInfoDo) Debug() IUserCheckinMachineInfoDo {
	return u.withDO(u.DO.Debug())
}

func (u userCheckinMachineInfoDo) WithContext(ctx context.Context) IUserCheckinMachineInfoDo {
	return u.withDO(u.DO.WithContext(ctx))
}

func (u userCheckinMachineInfoDo) ReadDB() IUserCheckinMachineInfoDo {
	return u.Clauses(dbresolver.Read)
}

func (u userCheckinMachineInfoDo) WriteDB() IUserCheckinMachineInfoDo {
	return u.Clauses(dbresolver.Write)
}

func (u userCheckinMachineInfoDo) Session(config *gorm.Session) IUserCheckinMachineInfoDo {
	return u.withDO(u.DO.Session(config))
}

func (u userCheckinMachineInfoDo) Clauses(conds ...clause.Expression) IUserCheckinMachineInfoDo {
	return u.withDO(u.DO.Clauses(conds...))
}

func (u userCheckinMachineInfoDo) Returning(value interface{}, columns ...string) IUserCheckinMachineInfoDo {
	return u.withDO(u.DO.Returning(value, columns...))
}

func (u userCheckinMachineInfoDo) Not(conds ...gen.Condition) IUserCheckinMachineInfoDo {
	return u.withDO(u.DO.Not(conds...))
}

func (u userCheckinMachineInfoDo) Or(conds ...gen.Condition) IUserCheckinMachineInfoDo {
	return u.withDO(u.DO.Or(conds...))
}

func (u userCheckinMachineInfoDo) Select(conds ...field.Expr) IUserCheckinMachineInfoDo {
	return u.withDO(u.DO.Select(conds...))
}

func (u userCheckinMachineInfoDo) Where(conds ...gen.Condition) IUserCheckinMachineInfoDo {
	return u.withDO(u.DO.Where(conds...))
}

func (u userCheckinMachineInfoDo) Order(conds ...field.Expr) IUserCheckinMachineInfoDo {
	return u.withDO(u.DO.Order(conds...))
}

func (u userCheckinMachineInfoDo) Distinct(cols ...field.Expr) IUserCheckinMachineInfoDo {
	return u.withDO(u.DO.Distinct(cols...))
}

func (u userCheckinMachineInfoDo) Omit(cols ...field.Expr) IUserCheckinMachineInfoDo {
	return u.withDO(u.DO.Omit(cols...))
}

func (u userCheckinMachineInfoDo) Join(table schema.Tabler, on ...field.Expr) IUserCheckinMachineInfoDo {
	return u.withDO(u.DO.Join(table, on...))
}

func (u userCheckinMachineInfoDo) LeftJoin(table schema.Tabler, on ...field.Expr) IUserCheckinMachineInfoDo {
	return u.withDO(u.DO.LeftJoin(table, on...))
}

func (u userCheckinMachineInfoDo) RightJoin(table schema.Tabler, on ...field.Expr) IUserCheckinMachineInfoDo {
	return u.withDO(u.DO.RightJoin(table, on...))
}

func (u userCheckinMachineInfoDo) Group(cols ...field.Expr) IUserCheckinMachineInfoDo {
	return u.withDO(u.DO.Group(cols...))
}

func (u userCheckinMachineInfoDo) Having(conds ...gen.Condition) IUserCheckinMachineInfoDo {
	return u.withDO(u.DO.Having(conds...))
}

func (u userCheckinMachineInfoDo) Limit(limit int) IUserCheckinMachineInfoDo {
	return u.withDO(u.DO.Limit(limit))
}

func (u userCheckinMachineInfoDo) Offset(offset int) IUserCheckinMachineInfoDo {
	return u.withDO(u.DO.Offset(offset))
}

func (u userCheckinMachineInfoDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IUserCheckinMachineInfoDo {
	return u.withDO(u.DO.Scopes(funcs...))
}

func (u userCheckinMachineInfoDo) Unscoped() IUserCheckinMachineInfoDo {
	return u.withDO(u.DO.Unscoped())
}

func (u userCheckinMachineInfoDo) Create(values ...*model.UserCheckinMachineInfo) error {
	if len(values) == 0 {
		return nil
	}
	return u.DO.Create(values)
}

func (u userCheckinMachineInfoDo) CreateInBatches(values []*model.UserCheckinMachineInfo, batchSize int) error {
	return u.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (u userCheckinMachineInfoDo) Save(values ...*model.UserCheckinMachineInfo) error {
	if len(values) == 0 {
		return nil
	}
	return u.DO.Save(values)
}

func (u userCheckinMachineInfoDo) First() (*model.UserCheckinMachineInfo, error) {
	if result, err := u.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserCheckinMachineInfo), nil
	}
}

func (u userCheckinMachineInfoDo) Take() (*model.UserCheckinMachineInfo, error) {
	if result, err := u.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserCheckinMachineInfo), nil
	}
}

func (u userCheckinMachineInfoDo) Last() (*model.UserCheckinMachineInfo, error) {
	if result, err := u.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserCheckinMachineInfo), nil
	}
}

func (u userCheckinMachineInfoDo) Find() ([]*model.UserCheckinMachineInfo, error) {
	result, err := u.DO.Find()
	return result.([]*model.UserCheckinMachineInfo), err
}

func (u userCheckinMachineInfoDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.UserCheckinMachineInfo, err error) {
	buf := make([]*model.UserCheckinMachineInfo, 0, batchSize)
	err = u.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (u userCheckinMachineInfoDo) FindInBatches(result *[]*model.UserCheckinMachineInfo, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return u.DO.FindInBatches(result, batchSize, fc)
}

func (u userCheckinMachineInfoDo) Attrs(attrs ...field.AssignExpr) IUserCheckinMachineInfoDo {
	return u.withDO(u.DO.Attrs(attrs...))
}

func (u userCheckinMachineInfoDo) Assign(attrs ...field.AssignExpr) IUserCheckinMachineInfoDo {
	return u.withDO(u.DO.Assign(attrs...))
}

func (u userCheckinMachineInfoDo) Joins(fields ...field.RelationField) IUserCheckinMachineInfoDo {
	for _, _f := range fields {
		u = *u.withDO(u.DO.Joins(_f))
	}
	return &u
}

func (u userCheckinMachineInfoDo) Preload(fields ...field.RelationField) IUserCheckinMachineInfoDo {
	for _, _f := range fields {
		u = *u.withDO(u.DO.Preload(_f))
	}
	return &u
}

func (u userCheckinMachineInfoDo) FirstOrInit() (*model.UserCheckinMachineInfo, error) {
	if result, err := u.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserCheckinMachineInfo), nil
	}
}

func (u userCheckinMachineInfoDo) FirstOrCreate() (*model.UserCheckinMachineInfo, error) {
	if result, err := u.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserCheckinMachineInfo), nil
	}
}

func (u userCheckinMachineInfoDo) FindByPage(offset int, limit int) (result []*model.UserCheckinMachineInfo, count int64, err error) {
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

func (u userCheckinMachineInfoDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = u.Count()
	if err != nil {
		return
	}

	err = u.Offset(offset).Limit(limit).Scan(result)
	return
}

func (u userCheckinMachineInfoDo) Scan(result interface{}) (err error) {
	return u.DO.Scan(result)
}

func (u userCheckinMachineInfoDo) Delete(models ...*model.UserCheckinMachineInfo) (result gen.ResultInfo, err error) {
	return u.DO.Delete(models)
}

func (u *userCheckinMachineInfoDo) withDO(do gen.Dao) *userCheckinMachineInfoDo {
	u.DO = *do.(*gen.DO)
	return u
}
