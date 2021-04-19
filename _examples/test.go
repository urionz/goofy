package main

import (
	"fmt"

	"github.com/zouyx/agollo/v4"
	"github.com/zouyx/agollo/v4/component/log"
	"github.com/zouyx/agollo/v4/env/config"
)

func main() {
	c := &config.AppConfig{
		AppID:          "ifanni-server",
		Cluster:        "dev",
		IP:             "http://apollo-wh.shengtian.com:8070/",
		NamespaceName:  "application",
		IsBackupConfig: true,
		Secret:         "aed71bf9485c40f1b39be7cc1dcaca55",
	}
	agollo.SetLogger(&log.DefaultLogger{})
	apollo, err := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return c, nil
	})
	if err != nil {
		panic(err)
	}

	cache := apollo.GetConfigCache(c.NamespaceName)
	fmt.Println(cache.Get("config"))
}
