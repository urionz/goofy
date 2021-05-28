package config

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/shima-park/agollo"
	"github.com/urionz/color"
	"github.com/urionz/config"
	"github.com/urionz/config/hcl"
	"github.com/urionz/config/ini"
	"github.com/urionz/config/json"
	"github.com/urionz/config/toml"
	"github.com/urionz/config/yaml"
	"github.com/urionz/goofy/container"
	"github.com/urionz/goofy/contracts"
	"github.com/urionz/goofy/log"
	"github.com/urionz/ini/dotenv"
)

const (
	DrvFile   = "file"
	DrvApollo = "apollo"
)

func NewServiceProvider(app contracts.Application) error {
	if err := dotenv.LoadExists(app.Workspace(), ".env"); err != nil {
		return err
	}

	confDriver := dotenv.Get("CONF_DRIVER", "file")
	envConf := dotenv.Get("APP_ENV", "dev")

	serve = &Configure{
		Config: config.New("goofy"),
	}

	serve.AddDriver(yaml.Driver)
	serve.AddDriver(json.Driver)
	serve.AddDriver(ini.Driver)
	serve.AddDriver(hcl.Driver)
	serve.AddDriver(toml.Driver)

	switch confDriver {
	case DrvFile:
		envConfFile := dotenv.Get("APP_CONF", fmt.Sprintf("config.%s.toml", envConf))

		if err := serve.LoadExists(path.Join(app.Workspace(), "config.toml"), path.Join(app.Workspace(), envConfFile)); err != nil {
			return err
		}
		break
	case DrvApollo:
		apolloServerURL := dotenv.Get("APOLLO_SERVER_URL", "")
		apolloAppId := dotenv.Get("APOLLO_APP_ID", "")
		apolloAccessKey := dotenv.Get("APOLLO_ACCESS_KEY")
		apollo, err := agollo.New(
			apolloServerURL, apolloAppId,
			agollo.Cluster(envConf),
			agollo.AutoFetchOnCacheMiss(),
			agollo.AccessKey(apolloAccessKey),
			agollo.WithLogger(agollo.NewLogger(agollo.LoggerWriter(os.Stdout))),
		)
		if err != nil {
			return err
		}

		var loadConf = func(kv agollo.Configurations) {
			for k, value := range kv {
				if strings.HasSuffix(k, "]") {
					fieldKeys := strings.Split(k, ".")
					target := fieldKeys[len(fieldKeys)-1]
					lastKeys := strings.Split(target, "[")
					arrKey := lastKeys[0]
					key := strings.Join(append(fieldKeys[0:len(fieldKeys)-1], arrKey), ".")
					existsValue := serve.Config.Strings(key)
					existsValue = append(existsValue, value.(string))
					serve.Set(key, existsValue, true)
					continue
				}
				serve.Set(k, value, true)
			}
		}

		loadConf(apollo.GetNameSpace("application"))

		errCh := apollo.Start()
		watchCh := apollo.Watch()

		go func() {
			for {
				select {
				case err := <-errCh:
					log.Panic(err)
				case resp := <-watchCh:
					color.Warnln("配置文件被更新")
					loadConf(resp.NewValue)
					if err := app.DynamicConf(serve); err != nil {
						log.Panic(err)
					}
				}
			}
		}()
		break
	}

	app.AddCommanders(contracts.FuncCommander(Command))

	return app.ProvideValue(serve, container.As(new(contracts.Config)), container.Tags{"name": "config"})
}
