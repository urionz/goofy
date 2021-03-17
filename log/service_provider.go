package log

import (
	"github.com/urionz/goofy/contracts"
)

func NewServiceProvider(app contracts.Application, conf contracts.Config) error {
	return app.ProvideValue(NewLogger(app, conf))
}
