package validator

import (
	"github.com/goava/di"
	"github.com/urionz/goofy/contracts"
)

func NewServiceProvider(app contracts.Application) error {
	return app.ProvideValue(NewFactory(), di.As(new(FactoryContract)))
}
