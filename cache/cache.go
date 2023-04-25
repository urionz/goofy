package cache

import (
	"time"

	"github.com/urionz/goofy/contracts"
)

type BaseCache struct {
}

func (*BaseCache) Scan(_ string, ptr interface{}, _ ...interface{}) error {
	return nil
}
func (*BaseCache) Set(_ string, _ interface{}, _ time.Duration) error {
	return nil
}
func (*BaseCache) Delete(_ string) error {
	return nil
}
func (*BaseCache) Clear() error {
	return nil
}
func (*BaseCache) GetMultiple(_ []string, _ interface{}) map[string]interface{} {
	return map[string]interface{}{}
}
func (*BaseCache) SetMultiple(_ map[string]interface{}, _ ...time.Duration) error {
	return nil
}
func (*BaseCache) DelMultiple(_ []string) error {
	return nil
}
func (*BaseCache) Has(_ string) bool {
	return false
}

type BaseStore struct {
}

func (*BaseStore) Get(_ string) (interface{}, error) {
	return nil, nil
}
func (*BaseStore) Many(_ []string) []interface{} {
	return nil
}
func (*BaseStore) Put(_ string, _ interface{}, _ time.Duration) error {
	return nil
}
func (*BaseStore) PutPure(_ string, _ interface{}, _ time.Duration) error {
	return nil
}
func (*BaseStore) PutMany(_ map[string]interface{}, _ int) error {
	return nil
}
func (*BaseStore) Increment(_ string, _ ...int64) error {
	return nil
}
func (*BaseStore) Decrement(_ string, _ ...int64) error {
	return nil
}
func (*BaseStore) Forever(_ string, _ interface{}) error {
	return nil
}
func (*BaseStore) ForeverPure(_ string, _ interface{}) error {
	return nil
}
func (*BaseStore) Forget(_ string) error {
	return nil
}
func (*BaseStore) Flush() error {
	return nil
}
func (*BaseStore) Has(_ string) bool {
	return false
}
func (*BaseStore) GetPrefix() string {
	return ""
}
func (*BaseStore) ItemKey(key string) string {
	return key
}

type BaseRepository struct {
	BaseCache
}

func (*BaseRepository) Pull(_ string, _ ...interface{}) interface{} {
	return nil
}
func (*BaseRepository) Put(_ string, _ interface{}, _ time.Duration) error {
	return nil
}
func (*BaseRepository) Add(_ string, _ interface{}, _ ...time.Duration) error {
	return nil
}
func (*BaseRepository) Increment(_ string, _ ...int64) error {
	return nil
}
func (*BaseRepository) Decrement(_ string, _ ...int64) error {
	return nil
}
func (*BaseRepository) Forever(_ string, _ interface{}) error {
	return nil
}
func (*BaseRepository) Remember(_ string, _ time.Duration, _ contracts.CacheClosure, _ interface{}) error {
	return nil
}
func (*BaseRepository) Sear(_ string, _ contracts.CacheClosure, _ interface{}) error {
	return nil
}
func (*BaseRepository) RememberForever(_ string, _ contracts.CacheClosure, _ interface{}) error {
	return nil
}
func (*BaseRepository) Forget(_ string) error {
	return nil
}
func (*BaseRepository) GetStore() contracts.Cache {
	return nil
}

type BaseFactory struct {
}

func (*BaseFactory) Store(_ ...string) (contracts.CacheRepository, error) {
	return nil, nil
}
