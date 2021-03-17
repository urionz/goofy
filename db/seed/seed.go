package seed

import "github.com/urionz/goofy/contracts"

var seederFiles []contracts.DBSeeder

func Register(files ...contracts.DBSeeder) {
	seederFiles = append(seederFiles, files...)
}
