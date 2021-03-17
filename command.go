package goofy

import "github.com/urionz/goofy/contracts"

func (app *Application) AddCommanders(commander ...contracts.Commander) contracts.Application {
	for _, command := range commander {
		app.App.AddCommand(command.Handle(app))
	}
	return app
}
