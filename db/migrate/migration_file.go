package migrate

import "github.com/urionz/goofy/contracts"

var migrationFiles []contracts.MigrateFile

func Register(migrateFile ...contracts.MigrateFile) {
	migrationFiles = append(migrationFiles, migrateFile...)
}
