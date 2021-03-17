package cache

import (
	"path"
	"runtime"

	"github.com/goava/di"
	"github.com/urionz/goofy/contracts"
)

func NewServiceProvider(app contracts.Application, conf contracts.Config) error {
	_, f, _, _ := runtime.Caller(0)
	if err := conf.LoadExists(path.Join(path.Dir(f), "cache.toml")); err != nil {
		return err
	}
	instance = NewManager(app, conf)
	if err := app.ProvideValue(instance, di.As(new(contracts.CacheFactory))); err != nil {
		return err
	}
	return nil
}
