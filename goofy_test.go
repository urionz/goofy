package goofy_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urionz/goofy"
)

func TestOption(t *testing.T) {
	t.Run("set workspace option", func(t *testing.T) {
		workspace := "./"
		newApp := goofy.New(goofy.SetWorkspace(workspace))
		require.NotNil(t, newApp)
		require.Equal(t, workspace, newApp.Workspace())
	})
}

func TestApplication_Run(t *testing.T) {
	t.Run("run application", func(t *testing.T) {
		app := goofy.New()
		err := app.Run().Error()
		require.NoError(t, err)
	})
}
