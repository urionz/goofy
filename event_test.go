package goofy_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urionz/goofy"
)

func TestApplication_Emit(t *testing.T) {
	t.Run("test emit", func(t *testing.T) {
		app := goofy.New()
		err, _ := app.Emit("test", goofy.EventM{})
		require.NoError(t, err)
	})
}

func TestApplication_MustEmit(t *testing.T) {
	t.Run("test must emit", func(t *testing.T) {
		app := goofy.New()
		require.NotPanics(t, func() {
			app.MustEmit("test", goofy.EventM{})
		})
	})
}

func TestApplication_AddListener(t *testing.T) {
	t.Run("test add listener", func(t *testing.T) {
		app := goofy.New()
		testValue := "test"
		err := app.AddListeners(goofy.EventListeners{
			"test": goofy.Listeners(goofy.ListenerFunc(func(e goofy.Event) error {
				gotTestValue := e.Get("test")
				require.NotNil(t, gotTestValue)
				require.EqualValues(t, gotTestValue, testValue)
				return nil
			})),
		}).Error()
		require.NoError(t, err)

		app.Emit("test", goofy.EventM{
			"test": testValue,
		})
	})
}

func TestApplication_Dispatch(t *testing.T) {
	t.Run("test dispatch event", func(t *testing.T) {
		app := goofy.New()
		testValue := "test"
		err := app.AddListeners(goofy.EventListeners{
			"test": goofy.Listeners(goofy.ListenerFunc(func(e goofy.Event) error {
				gotTestValue := e.Get("test")
				require.NotNil(t, gotTestValue)
				require.EqualValues(t, gotTestValue, testValue)
				return nil
			})),
		}).Error()
		require.NoError(t, err)

		err = app.Dispatch("test", goofy.EventM{
			"test": testValue,
		}).Error()
		require.NoError(t, err)
	})
}
