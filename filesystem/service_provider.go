package filesystem

import (
	"github.com/goava/di"
	"github.com/urionz/goofy/contracts"
)

func init() {
	contracts.AddConfTpl(`[filesystems]
default = "local"
    [filesystems.disks.local]
    driver = "local"
    root = "storage/local"
    url = "http://localhost"
`)
}

func NewServiceProvider(app contracts.Application, conf contracts.Config) error {
	var err error
	var provide = func(value di.Value, options ...di.ProvideOption) {
		if err != nil {
			return
		}
		err = app.ProvideValue(value, options...)
	}
	provide(NewFilesystem())
	provide(NewManager(conf))

	if err = app.Provide(func(filesystem *Manager, conf contracts.Config) contracts.Filesystem {
		return filesystem.Disk(conf.String("filesystems.default"))
	}, di.Tags{"name": "filesystem.disk"}); err != nil {
		return err
	}

	return err
}
