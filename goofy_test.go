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

func TName(name string) {

}

func TestApplication_Run(t *testing.T) {
	t.Run("run application", func(t *testing.T) {
		// app := goofy.New(goofy.SetWorkspace("./"))
		//
		// var storage string
		// var database string
		// require.NoError(t, app.Resolve(&storage, di.Tags{"name": "path.storage"}))
		// require.Equal(t, "storage", storage)
		// require.NoError(t, app.Resolve(&database, di.Tags{"name": "path.database"}))
		// require.Equal(t, "database", database)
		//
		// err := app.Run().Error()
		// require.NoError(t, err)
	})
}
