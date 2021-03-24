package db

import (
	"fmt"

	"github.com/urionz/goofy/contracts"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var instance *Manager

func Model(model contracts.DBConnection) *gorm.DB {
	return instance.Connection(model.Connection()).Model(&model)
}

func Truncate(model schema.Tabler, db *gorm.DB) error {
	return db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", model.TableName())).Error
}

func Default() *gorm.DB {
	return instance.Connection()
}
