package contracts

import "github.com/urionz/goofy/command"

type (
	Commander interface {
		Handle(app Application) *command.Command
	}
	FuncCommander func(app Application) *command.Command
)

func (c FuncCommander) Handle(app Application) *command.Command { return c(app) }
