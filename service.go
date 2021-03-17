package goofy

import "github.com/urionz/goofy/contracts"

func (app *Application) AddServices(services ...interface{}) contracts.Application {
	app.services = append(app.services, services...)
	return app
}
