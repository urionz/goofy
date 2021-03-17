package goofy_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urionz/goofy"
)

type TestImpl struct {
}

func (*TestImpl) Test() string {
	return "test"
}

func TestApplication_AddServices(t *testing.T) {
	t.Run("add service not err", func(t *testing.T) {
		app := goofy.New()
		require.NotPanics(t, func() {
			app.AddServices(func(a goofy.Application) {
				a.Provide(func() *TestImpl {
					return new(TestImpl)
				})
			}).AddServices(func(test *TestImpl) {
				require.Equal(t, "test", test.Test())
			}).Run()
		})
	})

	t.Run("add service with err should panic", func(t *testing.T) {
		app := goofy.New()
		require.Panics(t, func() {
			app.AddServices(func(a goofy.Application) error {
				return fmt.Errorf("test add service with err")
			}).Run()
		})
	})

	t.Run("when add service use not exists dep service should panic", func(t *testing.T) {
		app := goofy.New()
		require.Panics(t, func() {
			app.AddServices(func(test *TestImpl) {
				require.Equal(t, "test", test.Test())
			}).Run()
		})
	})
}
