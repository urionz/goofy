package redis

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/goava/di"
	"github.com/urionz/goofy/contracts"
)

type Manager struct {
	di.Tags `name:"redis"`

	app         contracts.Application
	driver      string
	conf        contracts.Config
	connections sync.Map
}

var _ contracts.RedisFactory = (*Manager)(nil)

func NewRedisManager(app contracts.Application, conf contracts.Config) *Manager {
	manager := &Manager{
		app:  app,
		conf: conf,
	}
	return manager
}

func (m *Manager) Connection(names ...string) (contracts.RedisConnection, error) {
	var err error
	var conn *Connection

	driver := m.getDefaultConnection()
	if len(names) > 0 && names[0] != "" {
		driver = names[0]
	}

	if conn, ok := m.connections.Load(driver); ok {
		return conn.(*Connection), nil
	}
	if conn, err = m.configure(m.getConfig(driver), driver); err != nil {
		return nil, err
	}
	m.connections.Store(driver, conn)

	return conn, nil
}

func (m *Manager) getDefaultConnection() string {
	return m.conf.String("redis.default")
}

func (m *Manager) configure(conf contracts.Config, name string) (*Connection, error) {
	conn := NewConnection(redis.NewClient(
		&redis.Options{
			Addr:     conf.String("address", "localhost:6379"),
			Password: conf.String("password", ""),
			DB:       conf.Int("db", 0),
		},
	)).SetName(name)
	return conn, conn.client.Ping(context.Background()).Err()
}

func (m *Manager) getConfig(name string) contracts.Config {
	return m.conf.Object(fmt.Sprintf("redis.conns.%s", name))
}
