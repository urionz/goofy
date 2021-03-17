package cache_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urionz/goofy"
	"github.com/urionz/goofy/cache"
	"github.com/urionz/goofy/config"
	"github.com/urionz/goofy/contracts"
	"github.com/urionz/goofy/filesystem"
	"github.com/urionz/goofy/redis"
)

func TestNewServiceProvider(t *testing.T) {
	require.NotPanics(t, func() {
		goofy.Default.AddServices(config.NewServiceProvider, filesystem.NewServiceProvider, redis.NewServiceProvider, cache.NewServiceProvider, func(conf contracts.Config, c contracts.CacheFactory) {
			store := c.Store()
			require.NoError(t, store.Set("testk", "testv", 0))
			require.Equal(t, "testv", store.Get("testk"))
			require.NoError(t, store.Forget("testk"))
		}).Run()
	})
}
