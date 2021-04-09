package contracts

type (
	Commander interface {
		Handle(app Application) Command
	}
	FuncCommander func(app Application) Command
)

func (c FuncCommander) Handle(app Application) Command { return c(app) }
