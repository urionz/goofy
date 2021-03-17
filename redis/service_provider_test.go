package redis_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/urionz/goofy"
	"github.com/urionz/goofy/config"
	"github.com/urionz/goofy/redis"
)

func TestNewServiceProvider(t *testing.T) {
	require.NotPanics(t, func() {
		goofy.Default.AddServices(config.NewServiceProvider, redis.NewServiceProvider, func(rdm *redis.Manager) {
			rds, err := rdm.Connection()
			require.NoError(t, err)
			resueRds, err2 := rdm.Connection()
			require.NoError(t, err2)
			require.Equal(t, rds, resueRds)
			require.Equal(t, "default", rds.GetName())
			require.NoError(t, rds.Set("set", "test set", 0))
			require.Equal(t, "test set", rds.Get("set"))
			require.NoError(t, rds.SetEX("set_ex", "test_ex", time.Second*2))
			require.Equal(t, "test_ex", rds.Get("set_ex"))
			time.Sleep(2 * time.Second)
			require.NotEqual(t, "test_ex", rds.Get("set_ex"))
			require.NoError(t, rds.SAdd("test_sadd", "sadd_test"))

			_, err3 := rdm.Connection("cache")
			require.Error(t, err3)
		}).Run()
	})
}
