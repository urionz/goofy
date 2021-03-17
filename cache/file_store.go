package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-module/carbon"
	"github.com/urionz/goofy/contracts"
	"github.com/urionz/goofy/filesystem"
	"github.com/urionz/goutil/strutil"
)

type FileStore struct {
	files     *filesystem.Filesystem
	directory string
	BaseStore
}

type DataPayload struct {
	Data interface{} `json:"data"`
	Time int64       `json:"time"`
}

var _ contracts.Store = new(FileStore)

// Create a new file cache store instance.
func NewFileStore(files *filesystem.Filesystem, dir string) *FileStore {
	return &FileStore{
		files:     files,
		directory: dir,
	}
}

// Retrieve an item from the cache by key.
func (f *FileStore) Get(key string) interface{} {
	return f.getPayload(key).Data
}

// Store an item in the cache for a given number of seconds.
func (f *FileStore) Put(key string, data interface{}, seconds time.Duration) error {
	var err error
	var payloadRaw []byte
	p := f.path(key)
	f.ensureCacheDirectoryExists(p)
	dataPayload := new(DataPayload)
	if seconds != 0 {
		dataPayload.Time = carbon.Now().AddDuration(seconds.String()).ToTimestamp()
	}
	dataPayload.Data = data
	if payloadRaw, err = json.Marshal(dataPayload); err != nil {
		return err
	}
	return f.files.Put(p, payloadRaw)
}

func (f *FileStore) ensureCacheDirectoryExists(p string) {
	if dir := filepath.Dir(p); !f.files.Exists(dir) {
		os.MkdirAll(dir, os.ModePerm)
	}
}

func (f *FileStore) Forget(key string) error {
	if file := f.path(key); f.files.Exists(file) {
		return os.RemoveAll(file)
	}
	return nil
}

func (f *FileStore) getPayload(key string) *DataPayload {
	var err error
	var raw []byte
	var unpackPayload DataPayload
	now := carbon.Now()
	p := f.path(key)
	if raw, err = f.files.Get(p); err != nil {
		return f.emptyPayload()
	}

	if err = json.Unmarshal(raw, &unpackPayload); err != nil {
		return f.emptyPayload()
	}

	if unpackPayload.Time != 0 {
		tsDiff := now.DiffInSeconds(carbon.CreateFromTimestamp(unpackPayload.Time))
		if tsDiff <= 0 {
			f.Forget(key)
			return f.emptyPayload()
		}
	}

	return &DataPayload{
		Data: unpackPayload.Data,
		Time: unpackPayload.Time - now.ToTimestamp(),
	}
}

func (f *FileStore) Forever(key string, value interface{}) error {
	return f.Put(key, value, 0)
}

func (f *FileStore) emptyPayload() *DataPayload {
	return new(DataPayload)
}

func (f *FileStore) path(key string) string {
	var hash string
	originHash := strutil.Sha1(key)
	if len(originHash) >= 4 {
		hash = fmt.Sprintf("%s/%s", originHash[0:2], originHash[2:4])
	}
	return f.directory + "/" + hash + "/" + originHash + ".data"
}
