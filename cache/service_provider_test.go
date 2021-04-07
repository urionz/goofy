package cache_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/urionz/goofy"
	"github.com/urionz/goofy/contracts"
)

func TestNewServiceProvider(t *testing.T) {
	require.NotPanics(t, func() {
		goofy.Default.AddServices(func(conf contracts.Config, c contracts.CacheFactory) {
			store := c.Store()
			require.NoError(t, store.Set("testk", "testv", 0))
			require.Equal(t, "testv", store.Get("testk"))
			require.NoError(t, store.Forget("testk"))
			require.Equal(t, "123", store.Remember("test_remember", time.Second*1, func() interface{} {
				return "123"
			}))
			require.Equal(t, "123", store.Get("test_remember"))
			time.Sleep(time.Second * 2)
			require.Nil(t, store.Get("test_remember"))

			require.Equal(t, "321", store.RememberForever("test_remember_forever", func() interface{} {
				return "321"
			}))
			require.Equal(t, "321", store.Get("test_remember_forever"))
		}).Run()
	})
}
