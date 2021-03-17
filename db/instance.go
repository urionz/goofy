package db

import (
	"github.com/urionz/goofy/contracts"
	"gorm.io/gorm"
)

var instance *Manager

func Model(model contracts.DBConnection) *gorm.DB {
	return instance.Connection(model.Connection())
}
