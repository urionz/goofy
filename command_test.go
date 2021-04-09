package goofy_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urionz/goofy"
	"github.com/urionz/goofy/command"
	"github.com/urionz/goofy/contracts"
)

type StructCommand struct {
}

func (*StructCommand) Handle(_ contracts.Application) *command.Command {
	return &command.Command{
		Name: "test",
		Func: func(cmd *command.Command, args []string) error {
			return nil
		},
	}
}

func TestApplication_AddCommander(t *testing.T) {
	t.Run("add commander", func(t *testing.T) {
		app := goofy.New()
		err := app.AddCommanders(contracts.FuncCommander(func(app contracts.Application) *command.Command {
			return &command.Command{
				Name: "test",
				Func: func(cmd *command.Command, args []string) error {
					return nil
				},
			}
		})).Error()
		require.NoError(t, err)

		err = app.AddCommanders(new(StructCommand)).Error()
		require.NoError(t, err)
	})
}
