package contracts

import "gorm.io/gorm"

type DBFactory interface {
	Connection(...string) *gorm.DB
}

type DBConnection interface {
	Connection() string
}

type MigrateFile interface {
	MigrateTimestamp() int
	TableName() string
	Up(db *gorm.DB) error
	Down(db *gorm.DB) error
}

type DBSeeder interface {
	Handle(db *gorm.DB) error
}
