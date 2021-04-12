package filesystem_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/goava/di"
	"github.com/stretchr/testify/require"
	"github.com/urionz/goofy"
	"github.com/urionz/goofy/contracts"
	"github.com/urionz/goofy/filesystem"
	"github.com/urionz/goutil/strutil"
)

func TestNewServiceProvider(t *testing.T) {
	require.NotPanics(t, func() {
		goofy.Default.AddServices(func(fs *filesystem.Filesystem, manager *filesystem.Manager, app contracts.Application) {
			var err error
			c1, err1 := fs.Get("./config.dev.toml")
			require.NoError(t, err1)
			c2, err2 := ioutil.ReadFile("./config.dev.toml")
			require.NoError(t, err2)
			require.Equal(t, c1, c2)

			_, err = fs.Get("./config.toml")
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

			_, err = disk.Put("zhangsan/lisi.txt", []byte("zhangsan"))
			require.NoError(t, err)
			_, err = disk.Get("zhangsan/lisi.txt")
			require.NoError(t, err)

			f, err := os.Open("./test.png")
			fmt.Println(err)
			fmt.Println(disk.WriteStream("storage/n.png", f))
		}).Run()
	})
}
