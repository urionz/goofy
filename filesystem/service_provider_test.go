package filesystem_test

import (
	"io/ioutil"
	"testing"

	"github.com/goava/di"
	"github.com/stretchr/testify/require"
	"github.com/urionz/goofy"
	"github.com/urionz/goofy/config"
	"github.com/urionz/goofy/contracts"
	"github.com/urionz/goofy/filesystem"
	"github.com/urionz/goutil/strutil"
)

func TestNewServiceProvider(t *testing.T) {
	require.NotPanics(t, func() {
		goofy.Default.AddServices(config.NewServiceProvider, filesystem.NewServiceProvider, func(fs *filesystem.Filesystem, app contracts.Application) {
			c1, err1 := fs.Get("./config.dev.toml")
			require.NoError(t, err1)
			c2, err2 := ioutil.ReadFile("./config.dev.toml")
			require.NoError(t, err2)
			require.Equal(t, c1, c2)

			_, err := fs.Get("./config.toml")
			require.Error(t, err)

			err = fs.Put("./config.put.toml", []byte("put test"))
			require.NoError(t, err)

			putContent, err := fs.Get("./config.put.toml")
			require.NoError(t, err)
			require.Equal(t, putContent, []byte("put test"))
			require.Equal(t, true, fs.Exists("./config.put.toml"))
			require.Equal(t, fs.Hash("./config.put.toml"), strutil.Md5("put test"))
			require.Equal(t, true, fs.Missing("./conf.toml"))

			var disk contracts.Filesystem
			require.NoError(t, app.Resolve(&disk, di.Name("filesystem.disk")))

			require.NoError(t, disk.Put("zhangsan/lisi.txt", []byte("zhangsan")))
			disk.Get("zhangsan/lisi.txt")
			// require.NoError(t, disk.Get(""))
		}).Run()
	})
}
