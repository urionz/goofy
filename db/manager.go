package db

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/urionz/goofy/container"
	"github.com/urionz/goofy/contracts"
	"github.com/urionz/goofy/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Manager struct {
	container.Tags `name:"db"`

	connections sync.Map
	conf        contracts.Config
}

var _ contracts.DBFactory = (*Manager)(nil)

func NewManager(conf contracts.Config) *Manager {
	return &Manager{
		conf: conf,
	}
}

func (m *Manager) DynamicConf(_ contracts.Application, conf contracts.Config) error {
	m.connections.Range(func(key, value interface{}) bool {
		c := conf.Object(fmt.Sprintf("database.conns.%s", key))
		if db, ok := value.(*gorm.DB); ok {
			db.Logger = gormLogger.New(
				log.Log(),
				gormLogger.Config{
					SlowThreshold: time.Duration(c.Int("slow_threshold", 100)) * time.Millisecond,
					LogLevel:      m.parseLogLevel(c.String("log_level", conf.String("database.log_level", conf.String("logger.level", "debug")))),
					Colorful:      c.Bool("log_color", conf.Bool("database.log_color", conf.Bool("logger.color", true))),
				},
			)
		}
		return true
	})
	return nil
}

func (m *Manager) Connection(names ...string) *gorm.DB {
	var conn *gorm.DB
	var err error

	driver := m.getDefaultConnection()
	if len(names) > 0 && names[0] != "" {
		driver = names[0]
	}

	if conn, ok := m.connections.Load(driver); ok {
		return conn.(*gorm.DB)
	}
	if conn, err = m.resolve(driver); err != nil {
		return nil
	}
	m.connections.Store(driver, conn)
	return conn
}

func (m *Manager) resolve(name string) (conn *gorm.DB, err error) {
	var db *sql.DB
	conf := m.getConfig(name)
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
		DisableForeignKeyConstraintWhenMigrating: !conf.Bool("auto_migrate_constraint", false),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   conf.String("prefix", ""),
			SingularTable: conf.Bool("singular_table", false),
		},
		Logger: gormLogger.New(
			log.Log(),
			gormLogger.Config{
				SlowThreshold: time.Duration(conf.Int("slow_threshold", 100)) * time.Millisecond,
				LogLevel:      m.parseLogLevel(conf.String("log_level", m.conf.String("database.log_level", m.conf.String("logger.level", "debug")))),
				Colorful:      conf.Bool("log_color", m.conf.Bool("database.log_color", m.conf.Bool("logger.color", true))),
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

func (m *Manager) parseLogLevel(level string) gormLogger.LogLevel {
	switch level {
	case contracts.DebugLevel:
		return gormLogger.Info
	case contracts.ErrorLevel:
		return gormLogger.Error
	case contracts.InfoLevel:
		return gormLogger.Info
	case contracts.WarnLevel:
		return gormLogger.Warn
	case contracts.PanicLevel:
		return gormLogger.Error
	case contracts.FatalLevel:
		return gormLogger.Error
	}
	return gormLogger.Info
}
