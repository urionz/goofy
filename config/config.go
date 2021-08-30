package config

import (
	"github.com/urionz/config"
	"github.com/urionz/goofy/contracts"
	"github.com/urionz/goutil"
	"github.com/urionz/ini/dotenv"
)

var _ contracts.Config = new(Configure)

type Configure struct {
	*config.Config
}

func (*Configure) Env(key string, defVal interface{}) interface{} {
	switch defVal.(type) {
	case bool:
		return dotenv.Bool(key, defVal.(bool))
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return dotenv.Int(key, defVal.(int))
	case string:
		break
	}
	return dotenv.Get(key, defVal.(string))
}

func (c *Configure) Object(key string, findByPath ...bool) contracts.Config {
	conf := config.New(key)
	val, ok := c.GetValue(key, findByPath...)
	if !ok {
		conf.SetData(make(map[string]interface{}))
	} else {
		conf.SetData(val.(map[string]interface{}))
	}
	return &Configure{
		Config: conf,
	}
}

func (c *Configure) Strings(key string, defVal ...string) goutil.Strings {
	val := c.Config.Strings(key)
	if len(val) == 0 && len(defVal) > 0 {
		return defVal
	}
	return val
}
