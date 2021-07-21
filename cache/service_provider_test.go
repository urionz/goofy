package cache_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urionz/goofy"
	"github.com/urionz/goofy/contracts"
)

type TestData struct {
	Name string
	Age  uint
}

func TestNewServiceProvider(t *testing.T) {
	require.NotPanics(t, func() {
		goofy.Default.AddServices(func(conf contracts.Config, c contracts.CacheFactory) {
			// store := c.Store("redis")
			// var dst TestData
			// // err := store.Sear("teststruct", func() interface{} {
			// // 	return TestData{
			// // 		Name: "zhangsan",
			// // 		Age:  11,
			// // 	}
			// // }, &dst)
			// err := store.Set("teststructset", TestData{
			// 	Name: "lisi", Age: 20,
			// }, 0)
			// fmt.Println(err, dst)
			// err = store.Scan("teststructset", &dst)
			// fmt.Println(err, dst)
			// var dst string
			// err := store.Sear("testsear", func() interface{} {
			// 	return "testsear"
			// }, &dst)
			// fmt.Println(err, dst)
			// store.Set("testsear", "testsear", 0)
			// err := store.Scan("testsear", &dst)
			// fmt.Println(err, dst == "testsear", "-------")
			// require.NoError(t, store.Set("testk", "testv", 0))
			// var value string
			// require.NoError(t, store.Scan("testk", &value))
			// require.Equal(t, "testv", value)
			// require.NoError(t, store.Forget("testk"))
			//
			// require.NoError(t, store.Remember("test_remember", time.Second*1, func() interface{} {
			// 	return "123"
			// }, &value))
			// require.Equal(t, "123", value)
			// value = ""
			// require.NoError(t, store.Scan("test_remember", &value))
			// require.Equal(t, "123", value)
			// time.Sleep(time.Second * 2)
			// require.Error(t, store.Scan("test_remember", &value))
			//
			// value = ""
			// require.NoError(t, store.RememberForever("test_remember_forever", func() interface{} {
			// 	return "321"
			// }, &value))
			// require.Equal(t, "321", value)
			// value = ""
			// require.NoError(t, store.Scan("test_remember_forever", &value))
			// require.Equal(t, "321", value)
			//
			// var data TestData
			// require.NoError(t, store.RememberForever("test_remember_forever_struct", func() interface{} {
			// 	return &TestData{
			// 		Name: "test",
			// 	}
			// }, &data, true))
			// require.Equal(t, "test", data.Name)
		}).Run()
	})
}
