package cache

import (
	"fmt"
	"sync"

	"github.com/goava/di"
	"github.com/urionz/goofy/contracts"
	"github.com/urionz/goofy/filesystem"
	"github.com/urionz/goofy/redis"
)

const (
	DvrFile  = "file"
	DvrRedis = "redis"
)

type Manager struct {
	di.Tags `name:"cache"`

	app    contracts.Application
	conf   contracts.Config
	stores sync.Map
}

var _ contracts.CacheFactory = new(Manager)

func NewManager(app contracts.Application, conf contracts.Config) *Manager {
	manager := &Manager{
		app:  app,
		conf: conf,
	}
	return manager
}

// Get a cache store instance by name, wrapped in a repository.
func (m *Manager) Store(name ...string) contracts.CacheRepository {
	var store contracts.CacheRepository
	var err error
	if len(name) == 0 {
		name = append(name, m.getDefaultDriver())
	}
	if store, ok := m.stores.Load(name[0]); ok {
		return store.(contracts.CacheRepository)
	}
	if store, err = m.resolve(name[0]); err != nil {
		return nil
	}
	m.stores.Store(name[0], store)
	return store
}

// Get a cache driver instance.
func (m *Manager) Driver(driver ...string) contracts.CacheRepository {
	return m.Store(driver...)
}

// Resolve the given store.
func (m *Manager) resolve(name string) (repo contracts.CacheRepository, err error) {
	conf := m.getConfig(name)
	if conf == nil {
		return nil, fmt.Errorf("cache store %s is not defined", name)
	}
	driver := conf.String("driver")
	switch driver {
	case DvrFile:
		repo = m.createFileDriver(conf)
		break
	case DvrRedis:
		repo = m.createRedisDriver(conf)
		break
	}
	return repo, nil
}

// Create an instance of the file cache driver.
func (m *Manager) createFileDriver(conf contracts.Config) *Repository {
	var files *filesystem.Filesystem
	if err := m.app.Resolve(&files); err != nil {
		return nil
	}
	return m.repository(NewFileStore(files, conf.String("path", "./")))
}

// Create an instance of the Redis cache driver.
func (m *Manager) createRedisDriver(conf contracts.Config) *Repository {
	var rdm *redis.Manager
	var err error
	if err = m.app.Resolve(&rdm); err != nil {
		return nil
	}
	connection := conf.String("connection", "default")
	return m.repository(NewRedisStore(rdm, m.getPrefix(conf), connection))
}

// Create a new cache repository with the given implementation.
func (m *Manager) repository(store contracts.Store) *Repository {
	return NewRepository(store)
}

func (m *Manager) getConfig(name string) contracts.Config {
	return m.conf.Object(fmt.Sprintf("cache.stores.%s", name))
}

func (m *Manager) getDefaultDriver() string {
	return m.conf.String("cache.default")
}

// Get the cache prefix.
func (m *Manager) getPrefix(conf contracts.Config) string {
	return conf.String("prefix", m.conf.String("cache.prefix"))
}
