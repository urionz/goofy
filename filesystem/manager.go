package filesystem

import (
	"fmt"
	"sync"

	"github.com/goava/di"
	"github.com/urionz/goofy/contracts"
)

const (
	LocalDrv = "local"
	FtpDrv   = "ftp"
	SFtpDrv  = "sftp"
)

type Manager struct {
	di.Tags `name:"filesystem"`

	conf           contracts.Config
	disks          sync.Map
	customCreators sync.Map
}

type customCreator func(conf contracts.Config) interface{}

var _ contracts.FilesystemFactory = (*Manager)(nil)

func NewManager(conf contracts.Config) *Manager {
	return &Manager{
		conf: conf,
	}
}

func (m *Manager) Disk(name ...string) contracts.Filesystem {
	if len(name) == 0 {
		name = append(name, m.getDefaultDriver())
	}
	return m.get(name[0])
}

func (m *Manager) get(name string) contracts.Filesystem {
	if disk, ok := m.disks.Load(name); ok {
		return disk.(contracts.Filesystem)
	}
	disk, err := m.resolve(name)
	if err != nil {
		panic(err)
	}
	m.disks.Store(name, disk)
	return disk
}

func (m *Manager) resolve(name string) (contracts.Filesystem, error) {
	var driver contracts.Filesystem
	var err error
	conf := m.getConfig(name)
	if _, exists := m.customCreators.Load(conf.String("driver")); exists {
		return m.callCustomCreator(m.getConfig(name))
	}

	switch conf.String("driver") {
	case LocalDrv:
		driver = m.createLocalDriver(conf)
		break
	}
	return driver, err
}

func (m *Manager) createLocalDriver(conf contracts.Config) contracts.Filesystem {
	return NewLocalDriver(conf.String("root"))
}

func (m *Manager) callCustomCreator(conf contracts.Config) (contracts.Filesystem, error) {
	customCreatorDriver, _ := m.customCreators.Load(conf.String("driver"))
	creator, ok := customCreatorDriver.(customCreator)
	if !ok {
		return nil, fmt.Errorf("the creator %+v is not support", customCreatorDriver)
	}
	driver := creator(conf)
	if drive, ok := driver.(contracts.Filesystem); ok {
		return m.adapt(drive), nil
	}
	return nil, fmt.Errorf("the creator %+v is not support", customCreatorDriver)
}

func (m *Manager) adapt(filesystem contracts.Filesystem) contracts.Filesystem {
	return NewAdapter(filesystem)
}

func (m *Manager) getConfig(name string) contracts.Config {
	return m.conf.Object(fmt.Sprintf("filesystems.disks.%s", name))
}

func (m *Manager) getDefaultDriver() string {
	return m.conf.String("filesystems.default")
}
