package migrations

import (
	"github.com/urionz/goofy/db/migrate"
	"github.com/urionz/goofy/db/model"
	"gorm.io/gorm"
)

func init() {
	migrate.Register(&FQQDQD{})
}

type FQQDQD struct {
	model.BaseModel
}

func (table *FQQDQD) MigrateTimestamp() int {
	return 1615450379
}

func (table *FQQDQD) TableName() string {
	return "tests"
}

func (table *FQQDQD) Connection() string {
	return ""
}

func (table *FQQDQD) Up(db *gorm.DB) error {
	return db.Migrator().CreateTable(table)
}

func (table *FQQDQD) Down(db *gorm.DB) error {
	return db.Migrator().DropTable(table)
}
