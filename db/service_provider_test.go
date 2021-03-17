package db_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urionz/goofy"
	"github.com/urionz/goofy/config"
	"github.com/urionz/goofy/contracts"
	"github.com/urionz/goofy/db"
)

func TestNewServiceProvider(t *testing.T) {
	require.NotPanics(t, func() {
		goofy.Default.AddServices(config.NewServiceProvider, db.NewServiceProvider, func(manager contracts.DBFactory) {
			manager.Connection().Migrator().CreateTable()
		}).Run()
	})
}
