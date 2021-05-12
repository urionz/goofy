package cache

import (
	"fmt"
	"reflect"
	"time"

	"github.com/golang-module/carbon"
	"github.com/urionz/goofy/contracts"
	"github.com/urionz/goutil/jsonutil"
	"github.com/urionz/goutil/refutil"
)

type Repository struct {
	store contracts.Store
	BaseRepository
}

var _ contracts.Cache = new(Repository)

func NewRepository(store contracts.Store) *Repository {
	return &Repository{
		store: store,
	}
}

func (repo *Repository) Scan(key string, ptr interface{}, defVal ...interface{}) error {
	value := repo.store.Get(key)
	if len(defVal) > 0 && (value == nil || refutil.IsBlank(value)) {
		if closure, ok := defVal[0].(contracts.CacheClosure); ok {
			value = closure()
		} else {
			value = defVal[0]
		}
	}
	if value == nil {
		return fmt.Errorf("the key %s is not found", key)
	}
	b, err := jsonutil.Encode(value)
	if err != nil {
		return err
	}
	if err := jsonutil.Decode(b, ptr); err != nil {
		return err
	}
	return nil
}

func (repo *Repository) Set(key string, value interface{}, ttl time.Duration) error {
	return repo.Put(key, value, ttl)
}

func (repo *Repository) Put(key string, value interface{}, ttl time.Duration) error {
	if ttl == 0 {
		return repo.Forever(key, value)
	}

	seconds := repo.getSeconds(ttl)
	if seconds <= 0 {
		return repo.Forget(key)
	}
	return repo.store.Put(repo.store.ItemKey(key), value, seconds)
}

func (repo *Repository) Tags(names ...string) (contracts.TaggableStore, error) {
	typeof := reflect.TypeOf(repo.store)
	if _, exists := typeof.MethodByName("Tags"); !exists {
		return nil, fmt.Errorf("this cache store does not support tagging")
	}
	inputs := make([]reflect.Value, len(names))
	for index, name := range names {
		inputs[index] = reflect.ValueOf(name)
	}
	results := reflect.ValueOf(repo.store).MethodByName("Tags").Call(inputs)
	if !results[1].IsNil() {
		return nil, results[1].Interface().(error)
	}
	return results[0].Interface().(contracts.TaggableStore), nil
}

func (repo *Repository) getSeconds(ttl time.Duration) time.Duration {
	duration := carbon.Now().AddDuration(ttl.String())
	diffSeconds := carbon.Now().DiffInSeconds(duration)
	if diffSeconds > 0 {
		return time.Duration(diffSeconds) * time.Second
	}
	return 0
}

func (repo *Repository) Forever(key string, value interface{}) error {
	return repo.store.Forever(repo.store.ItemKey(key), value)
}

func (repo *Repository) Forget(key string) error {
	return repo.store.Forget(repo.store.ItemKey(key))
}

func (repo *Repository) Remember(key string, ttl time.Duration, callback contracts.CacheClosure, ptr interface{}, force ...bool) error {
	if len(force) <= 0 || !force[0] {
		if err := repo.Scan(key, ptr); err == nil {
			return nil
		}
	}
	value := callback()
	b, err := jsonutil.Encode(value)
	if err != nil {
		return err
	}
	if err := jsonutil.Decode(b, ptr); err != nil {
		return err
	}
	if err := repo.Put(key, ptr, ttl); err != nil {
		return err
	}
	return nil
}

func (repo *Repository) Sear(key string, callback contracts.CacheClosure, ptr interface{}, force ...bool) error {
	return repo.RememberForever(key, callback, ptr, force...)
}

func (repo *Repository) RememberForever(key string, callback contracts.CacheClosure, ptr interface{}, force ...bool) error {
	if len(force) <= 0 || !force[0] {
		if err := repo.Scan(key, ptr); err == nil {
			return nil
		}
	}

	value := callback()
	b, err := jsonutil.Encode(value)
	if err != nil {
		return err
	}
	if err := jsonutil.Decode(b, ptr); err != nil {
		return err
	}
	if err := repo.Forever(key, string(b)); err != nil {
		return err
	}

	return nil
}
