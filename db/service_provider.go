package db

import (
	"github.com/goava/di"
	"github.com/urionz/goofy/contracts"
	"github.com/urionz/goofy/db/migrate"
	"github.com/urionz/goofy/db/model"
	"github.com/urionz/goofy/db/seed"
)

func NewServiceProvider(app contracts.Application, conf contracts.Config) error {
	app.AddCommanders(
		contracts.FuncCommander(migrate.Make), contracts.FuncCommander(migrate.Migrate),
		contracts.FuncCommander(migrate.Rollback), contracts.FuncCommander(migrate.Status),
		contracts.FuncCommander(migrate.Reset), contracts.FuncCommander(migrate.Refresh),
		contracts.FuncCommander(model.Make), contracts.FuncCommander(seed.Seed),
	)
	instance = NewManager(conf)
	return app.ProvideValue(instance, di.As(new(contracts.DBFactory)))
}
