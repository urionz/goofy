package cache_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/urionz/goofy"
	"github.com/urionz/goofy/contracts"
)

type TestData struct {
	Name string
}

func TestNewServiceProvider(t *testing.T) {
	require.NotPanics(t, func() {
		goofy.Default.AddServices(func(conf contracts.Config, c contracts.CacheFactory) {
			store := c.Store()
			require.NoError(t, store.Set("testk", "testv", 0))
			var value string
			require.NoError(t, store.Scan("testk", &value))
			require.Equal(t, "testv", value)
			require.NoError(t, store.Forget("testk"))

			require.NoError(t, store.Remember("test_remember", time.Second*1, func() interface{} {
				return "123"
			}, &value))
			require.Equal(t, "123", value)
			value = ""
			require.NoError(t, store.Scan("test_remember", &value))
			require.Equal(t, "123", value)
			time.Sleep(time.Second * 2)
			require.Error(t, store.Scan("test_remember", &value))

			value = ""
			require.NoError(t, store.RememberForever("test_remember_forever", func() interface{} {
				return "321"
			}, &value))
			require.Equal(t, "321", value)
			value = ""
			require.NoError(t, store.Scan("test_remember_forever", &value))
			require.Equal(t, "321", value)

			var data TestData
			require.NoError(t, store.RememberForever("test_remember_forever_struct", func() interface{} {
				fmt.Println("sotre")
				return &TestData{
					Name: "test",
				}
			}, &data, true))
			require.Equal(t, "test", data.Name)
		}).Run()
	})
}
