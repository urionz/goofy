package filesystem

import (
	"testing"
)

// func init() {
// 	_ = os.RemoveAll("./storage")
// }
//
// func TestLocalDriver_Put(t *testing.T) {
// 	var err error
// 	local := NewLocalDriver("storage")
// 	_, err = local.Put("put/test.txt", []byte("test"))
// 	require.NoError(t, err)
// 	_ = os.MkdirAll("storage/only_read", 0400)
//
// 	forbidRootLocal := NewLocalDriver("read_only_dir")
// 	_ = os.MkdirAll("read_only_dir", 0400)
// 	_, err = forbidRootLocal.Put("test/test.txt", []byte("forbid"))
// 	require.Error(t, err)
//
// 	_, err = local.Put("read_only.txt", []byte("read_only"))
// 	require.NoError(t, err)
// 	os.Chmod("storage/read_only.txt", 0400)
// 	_, err = local.Put("read_only.txt", []byte("append"))
// 	require.Error(t, err)
// }
//
// func TestLocalDriver_Delete(t *testing.T) {
// 	var err error
// 	local := NewLocalDriver("storage")
// 	_, err = local.Put("delete/test.txt", []byte("test"))
// 	require.NoError(t, err)
// 	require.NoError(t, local.Delete("delete/test.txt"))
// }
//
// func TestLocalDriver_Copy(t *testing.T) {
// 	var err error
// 	local := NewLocalDriver("storage")
// 	_, err = local.Put("copy/test.txt", []byte("copy"))
// 	require.NoError(t, err)
// 	_, err = local.Copy("copy/test.txt", "copy/test2.txt")
// 	require.NoError(t, err)
//
// 	require.Equal(t, true, local.Exists("copy/test.txt"))
// 	require.Equal(t, true, local.Exists("copy/test2.txt"))
// 	require.Equal(t, false, local.Exists(""))
// 	_, err = local.Copy("a/b/c/d/d/d/d/move.txt", "a/b/c/d/copy.txt")
// 	require.Error(t, err)
// }
//
// func TestLocalDriver_Move(t *testing.T) {
//
// }
//
// func TestLocalDriver_Size(t *testing.T) {
//
// }
//
// func TestLocalDriver_Exists(t *testing.T) {
//
// }
//
// func TestLocalDriver_Get(t *testing.T) {
//
// }

func TestLocalDriver(t *testing.T) {
	// os.RemoveAll("./storage")
	// readOnlyDir := "read_only_dir"
	// writeOnlyDir := "write_only_dir"
	//
	// _ = os.Mkdir("storage", 0777)
	// _ = os.Mkdir(path.Join("storage", readOnlyDir), 0400)
	// _ = os.Mkdir(path.Join("storage", writeOnlyDir), 0200)
	// local := NewLocalDriver("storage")
	//
	// t.Run("get not exists file", func(t *testing.T) {
	// 	_, getErr := local.Get("test.txt")
	// 	require.Error(t, getErr)
	// })
	//
	// t.Run("get exists file", func(t *testing.T) {
	// 	require.NoError(t, local.Put("a/b/c/test.txt", []byte("test")))
	// 	_, getErr2 := local.Get("a/b/c/test.txt")
	// 	require.NoError(t, getErr2)
	// 	require.NoError(t, local.Put("a/b/c/test.txt", []byte("test")))
	// })
	//
	// t.Run("append test", func(t *testing.T) {
	// 	getHandler, _ := local.Get("a/b/c/test.txt")
	// 	content, _ := ioutil.ReadAll(getHandler)
	// 	require.NotEqual(t, []byte("test"), content)
	// 	require.Equal(t, []byte("testtest"), content)
	// })
	//
	// t.Run("delete", func(t *testing.T) {
	// 	require.NoError(t, local.Delete("a/b/c/test.txt"))
	// 	local.pathPrefix = ""
	// 	require.Error(t, local.Delete("."))
	// 	local.pathPrefix = "storage"
	// })

	// require.NoError(t, local.Put("a/b/c/test.txt", []byte("test")))
	// getHandler2, _ := local.Get("a/b/c/test.txt")
	// content2, _ := ioutil.ReadAll(getHandler2)
	// require.Equal(t, []byte("test"), content2)
	//
	// require.NotEqual(t, 0, local.Size("a/b/c/test.txt"))
	//
	// require.NoError(t, local.Move("a/b/c/test.txt", "a/b/c/d/move.txt"))
	// require.Equal(t, false, local.Exists("a/b/c/test.txt"))
	// require.Equal(t, true, local.Exists("a/b/c/d/move.txt"))
	// _ = os.Mkdir("storage/forbidw", 0444)
	// require.Error(t, local.Move("a/b/c/d/move.txt", "forbidw/move2.txt"))
	//
	// require.NoError(t, local.Copy("a/b/c/d/move.txt", "a/b/c/d/copy.txt"))
	// require.Equal(t, true, local.Exists("a/b/c/d/copy.txt"))
	// require.Equal(t, true, local.Exists("a/b/c/d/move.txt"))
	// require.Equal(t, false, local.Exists(""))
	// require.Error(t, local.Copy("a/b/c/d/d/d/d/move.txt", "a/b/c/d/copy.txt"))
	//
	// _ = os.Mkdir("forbidw", 0444)
	// bandLocal := NewLocalDriver("forbidw")
	// require.Error(t, bandLocal.Put("a/b/write.txt", []byte("write")))
}
