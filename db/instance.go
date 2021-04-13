package db

import (
	"fmt"

	"github.com/urionz/goofy/contracts"
	"gorm.io/gorm/schema"
)

var instance *Manager

func Model(model contracts.DBConnection) *DB {
	return instance.Connection(model.Connection()).Model(model)
}

func Truncate(model schema.Tabler, db *DB) error {
	return db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", model.TableName())).Error
}

func Default() *DB {
	return instance.Connection()
}

func M() *Manager {
	return instance
}

func Tx(txFunc func(tx *DB) error, connections ...*DB) (err error) {
	conn := Default()
	if len(connections) > 0 && connections[0] != nil {
		conn = connections[0]
	}
	tx := conn.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit().Error
		}
	}()

	err = txFunc(tx)

	return err
}

func GetAllTable(connections ...*DB) []string {
	conn := Default()
	if len(connections) > 0 && connections[0] != nil {
		conn = connections[0]
	}
	var tables []string
	if err := conn.Raw("show tables").Scan(&tables).Error; err != nil {
		return tables
	}
	return tables
}

func DropAllTable(connections ...*DB) error {
	conn := Default()
	if len(connections) > 0 && connections[0] != nil {
		conn = connections[0]
	}
	tables := GetAllTable()
	for _, table := range tables {
		if err := conn.Migrator().DropTable(table); err != nil {
			return err
		}
	}
	return nil
}
