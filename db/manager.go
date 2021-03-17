package db

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/goava/di"
	"github.com/urionz/goofy/contracts"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Manager struct {
	di.Tags `name:"db"`

	connections sync.Map
	conf        contracts.Config
}

var _ contracts.DBFactory = (*Manager)(nil)

func NewManager(conf contracts.Config) *Manager {
	return &Manager{
		conf: conf,
	}
}

func (m *Manager) Connection(names ...string) *gorm.DB {
	var conn *gorm.DB
	var err error
	if len(names) == 0 || names[0] == "" {
		names = append(names, m.getDefaultConnection())
	}
	if conn, ok := m.connections.Load(names[0]); ok {
		return conn.(*gorm.DB)
	}
	if conn, err = m.resolve(names[0]); err != nil {
		return nil
	}
	m.connections.Store(names[0], conn)
	return conn
}

func (m *Manager) resolve(name string) (conn *gorm.DB, err error) {
	var db *sql.DB
	conf := m.getConfig(name)
	writes := []io.Writer{
		os.Stdout,
	}
	if conn, err = gorm.Open(mysql.Open(
		fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			conf.String("user", "root"),
			conf.String("password", "root"),
			conf.String("host", "localhost"),
			conf.Int("port", 3306),
			conf.String("name", "test"),
			conf.String("charset", "utf8mb4"),
		),
	), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   conf.String("prefix", ""),
			SingularTable: conf.Bool("singular_table", false),
		},
		Logger: logger.New(
			log.New(io.MultiWriter(writes...), "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Duration(conf.Int("slow_threshold", 100)) * time.Millisecond,
				LogLevel:      logger.Info,
				Colorful:      true,
			},
		),
	}); err != nil {
		return nil, err
	}
	if db, err = conn.DB(); err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(conf.Int("max_open_conns", 10))
	db.SetMaxOpenConns(conf.Int("max_idle_conns", 10))

	return
}

func (m *Manager) getDefaultConnection() string {
	return m.conf.String("database.default")
}

func (m *Manager) getConfig(name string) contracts.Config {
	return m.conf.Object(fmt.Sprintf("database.conns.%s", name))
}
