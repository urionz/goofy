package cache

import (
	"strings"

	"github.com/urionz/goofy/contracts"
	"github.com/urionz/goutil/strutil"
)

type TagSet struct {
	store contracts.Store
	names []string
}

func NewTagSet(store contracts.Store, names ...string) *TagSet {
	return &TagSet{
		store: store,
		names: names,
	}
}

func (tag *TagSet) Reset() {

}

func (tag *TagSet) ResetTag(name string) string {
	id := strings.Replace(strutil.NewUniqId(strutil.UniqIdParams{
		Prefix:      "",
		MoreEntropy: true,
	}), ".", "", -1)
	tag.store.Forever(tag.TagKey(name), id)
	return id
}

func (tag *TagSet) GetNamespace() string {
	return strings.Join(tag.tagIds(), "|")
}

func (tag *TagSet) tagIds() []string {
	var ids []string
	for _, name := range tag.names {
		ids = append(ids, tag.TagId(name))
	}
	return ids
}

func (tag *TagSet) TagId(name string) string {
	storeGet := tag.store.Get(tag.TagKey(name))
	if storeGet == nil {
		return tag.ResetTag(name)
	}
	return storeGet.(string)
}

func (tag *TagSet) TagKey(name string) string {
	return "tag:" + name + ":key"
}
