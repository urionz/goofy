package config_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urionz/goofy"
	"github.com/urionz/goofy/config"
	"github.com/urionz/goofy/contracts"
)

func TestNewConfigServiceProvider(t *testing.T) {
	require.NotPanics(t, func() {
		goofy.Default.AddServices(config.NewServiceProvider, func(conf contracts.Config) {
			require.NoError(t, config.LoadExists("./test.toml"))
			require.Equal(t, "test", config.String("string", "test"))
			require.Equal(t, 22, config.Int("int", 22))
			require.Equal(t, true, config.Bool("bool", true))
			require.Equal(t, "no_env", config.Env("app.env", "no_env"))
			require.Equal(t, 12, config.Env("int", 12))
			require.Equal(t, true, config.Env("bool", true))
			require.Equal(t, "redis", config.Object("cache.stores").String("redis.driver"))
			require.Equal(t, map[string]interface{}{}, config.Object("no.found").Data())
			require.Equal(t, []int{1, 2, 3}, config.Ints("app.ints"))
			require.Equal(t, []string{"1", "2", "3"}, config.Strings("app.strings"))
			require.Equal(t, "test", config.Get("app.test"))
			require.Equal(t, false, config.Exists("a.b.c.d"))
			require.Equal(t, uint(1), config.Uint("app.uint"))
			require.Equal(t, int64(1), config.Int64("app.uint"))
		}).Run()
	})
}
