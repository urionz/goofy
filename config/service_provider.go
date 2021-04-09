package config

import (
	"fmt"
	"path"

	"github.com/urionz/config"
	"github.com/urionz/config/hcl"
	"github.com/urionz/config/ini"
	"github.com/urionz/config/json"
	"github.com/urionz/config/toml"
	"github.com/urionz/config/yaml"
	"github.com/urionz/goofy/container"
	"github.com/urionz/goofy/contracts"
	"github.com/urionz/ini/dotenv"
)

func NewServiceProvider(app contracts.Application) error {
	serve = &Configure{
		Config: config.New("goofy"),
	}

	serve.AddDriver(yaml.Driver)
	serve.AddDriver(json.Driver)
	serve.AddDriver(ini.Driver)
	serve.AddDriver(hcl.Driver)
	serve.AddDriver(toml.Driver)

	dotenv.LoadExists(app.Workspace(), ".env")

	envConfFile := dotenv.Get("APP_CONF", fmt.Sprintf("config.%s.toml", dotenv.Get("APP_ENV", "dev")))

	if err := serve.LoadExists(path.Join(app.Workspace(), "config.toml"), path.Join(app.Workspace(), envConfFile)); err != nil {
		return err
	}

	app.AddCommanders(contracts.FuncCommander(Command))

	return app.ProvideValue(serve, container.As(new(contracts.Config)), container.Tags{"name": "config"})
}
