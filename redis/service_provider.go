package redis

import (
	"github.com/urionz/goofy/contracts"
)

func NewServiceProvider(app contracts.Application, conf contracts.Config) error {
	instance = NewRedisManager(app, conf)
	return app.ProvideValue(instance)
}
