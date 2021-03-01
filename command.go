package goofy

import (
	"github.com/urionz/cobra"
)

type (
	Commander interface {
		Handle(app IApplication) *cobra.Command
	}
	FuncCommander func(app IApplication) *cobra.Command
)

func (c FuncCommander) Handle(app IApplication) *cobra.Command { return c(app) }

func (app *Application) AddCommanders(commander ...Commander) IApplication {
	for _, command := range commander {
		app.Command.AddCommand(command.Handle(app))
	}
	return app
}
