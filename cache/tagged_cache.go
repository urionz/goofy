package cache

import "github.com/urionz/goofy/contracts"

type TaggedCache struct {
	Repository
	tags *TagSet
}

type TaggableStore struct {
	BaseStore
}

func NewTaggedCache(store contracts.Store, tags *TagSet) *TaggedCache {
	taggedCache := &TaggedCache{
		tags: tags,
	}
	taggedCache.store = store
	return taggedCache
}

func (tag *TaggableStore) Tags(names ...string) (contracts.TaggableStore, error) {
	return NewTaggedCache(tag, NewTagSet(tag, names...)), nil
}
