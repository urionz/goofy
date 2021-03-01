package goofy_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urionz/cobra"
	"github.com/urionz/goofy"
)

type StructCommand struct {
}

func (*StructCommand) Handle(_ goofy.IApplication) *cobra.Command {
	return &cobra.Command{
		Use: "test",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
}

func TestApplication_AddCommander(t *testing.T) {
	t.Run("add commander", func(t *testing.T) {
		app := goofy.New()
		err := app.AddCommanders(goofy.FuncCommander(func(app goofy.IApplication) *cobra.Command {
			return &cobra.Command{
				Use: "test",
				Run: func(c *cobra.Command, args []string) {

				},
			}
		})).Error()
		require.NoError(t, err)

		err = app.AddCommanders(new(StructCommand)).Error()
		require.NoError(t, err)
	})
}
