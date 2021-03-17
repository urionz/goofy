package cache

import (
	"strings"
	"time"

	"github.com/urionz/goofy/contracts"
	"github.com/urionz/goutil/strutil"
)

const (
	ReferenceKeyForever  = "forever_ref"
	ReferenceKeyStandard = "standard_ref"
)

type RedisTaggedCache struct {
	store *RedisStore
	TaggedCache
}

var _ contracts.TaggableStore = new(RedisTaggedCache)

func NewRedisTaggedCache(store *RedisStore, tags *TagSet) *RedisTaggedCache {
	r := new(RedisTaggedCache)
	r.TaggedCache.store = store
	r.store = store
	r.tags = tags
	return r
}

func (r *RedisTaggedCache) Get(key string, defVal ...interface{}) interface{} {
	return r.TaggedCache.Get(r.ItemKey(key), defVal...)
}

func (r *RedisTaggedCache) Set(key string, value interface{}, ttl time.Duration) error {
	return r.Put(key, value, ttl)
}

func (r *RedisTaggedCache) Put(key string, value interface{}, ttl time.Duration) error {
	if ttl == 0 {
		return r.Forever(key, value)
	}
	r.pushStandardKeys(r.tags.GetNamespace(), key)
	return r.TaggedCache.Put(r.ItemKey(key), value, ttl)
}

func (r *RedisTaggedCache) Forever(key string, value interface{}) error {
	r.pushForeverKeys(r.tags.GetNamespace(), key)
	return r.TaggedCache.Forever(r.ItemKey(key), value)
}

// Store standard key references into store.
func (r *RedisTaggedCache) pushStandardKeys(namespace, key string) error {
	return r.pushKeys(namespace, key, ReferenceKeyStandard)
}

// Store forever key references into store.
func (r *RedisTaggedCache) pushForeverKeys(namespace, key string) error {
	return r.pushKeys(namespace, key, ReferenceKeyForever)
}

// Store a reference to the cache key against the reference key.
func (r *RedisTaggedCache) pushKeys(namespace, key, reference string) error {
	fullKey := r.store.GetPrefix() + strutil.Sha1(namespace) + ":" + key
	segments := strings.Split(namespace, "|")
	for _, segment := range segments {
		if err := r.store.Connection().SAdd(r.referenceKey(segment, reference), fullKey); err != nil {
			return err
		}
	}
	return nil
}

// Get the reference key for the segment.
func (r *RedisTaggedCache) referenceKey(segment, suffix string) string {
	return r.store.GetPrefix() + segment + ":" + suffix
}

func (r *RedisTaggedCache) ItemKey(key string) string {
	return r.taggedItemKey(key)
}

func (r *RedisTaggedCache) taggedItemKey(key string) string {
	return strutil.Sha1(r.tags.GetNamespace()) + ":" + key
}
