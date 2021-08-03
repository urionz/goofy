package contracts

import (
	"github.com/urionz/goofy/pagination"
	"gorm.io/gorm"
)

type DBFactory interface {
	Connection(...string) *gorm.DB
}

type DBConnection interface {
	Connection() string
}

type MigrateFile interface {
	MigrateTimestamp() int
	Filename() string
	TableName() string
	Up(db *gorm.DB) error
	Down(db *gorm.DB) error
}

type DBSeeder interface {
	Handle(db *gorm.DB) error
}

type SqlCondition interface {
	AppendOrders(orders ...pagination.OrderByCol)
	Cols(selectCols ...string) SqlCondition
	WhereDate(column, operator, date string) SqlCondition
	Eq(column string, args ...interface{}) SqlCondition
	IsNull(column string) SqlCondition
	NotEq(column string, args ...interface{}) SqlCondition
	Gt(column string, args ...interface{}) SqlCondition
	Gte(column string, args ...interface{}) SqlCondition
	Lt(column string, args ...interface{}) SqlCondition
	Lte(column string, args ...interface{}) SqlCondition
	Like(column string, str string) SqlCondition
	Starting(column string, str string) SqlCondition
	Ending(column string, str string) SqlCondition
	In(column string, params interface{}) SqlCondition
	Where(query string, args ...interface{}) SqlCondition
	Asc(column string) SqlCondition
	Desc(column string) SqlCondition
	Limit(limit int) SqlCondition
	Page(page, limit int) SqlCondition
	Build(db *gorm.DB) *gorm.DB
	Find(db *gorm.DB, out interface{}) error
	FindOne(db *gorm.DB, out interface{}) error
	Count(db *gorm.DB, model interface{}) (int64, error)
}
