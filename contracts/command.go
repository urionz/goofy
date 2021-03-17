package contracts

import (
	"github.com/gookit/gcli/v3"
)

type (
	Commander interface {
		Handle(app Application) *gcli.Command
	}
	FuncCommander func(app Application) *gcli.Command
)

func (c FuncCommander) Handle(app Application) *gcli.Command { return c(app) }
