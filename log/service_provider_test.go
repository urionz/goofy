package log_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urionz/goofy"
	"github.com/urionz/goofy/config"
	"github.com/urionz/goofy/log"
)

func TestNewServiceProvider(t *testing.T) {
	require.NotPanics(t, func() {
		goofy.Default.AddServices(config.NewServiceProvider, log.NewServiceProvider, func(logger *log.Logger) {
			logger.Sugar().Info("info")
			logger.Sugar().Error("err")
			logger.Sugar().Warn("warn")
			logger.Sugar().Debug("debug")
			require.Panics(t, func() {
				logger.Sugar().Panic("panic")
			})

			log.Info("info")
			log.Infof("%s", "infof")
			log.Error("err")
			log.Errorf("%s", "errf")
			log.Warn("warn")
			log.Warnf("%s", "warnf")
			log.Debug("debug")
			log.Debugf("%s", "debugf")

			require.Panics(t, func() {
				log.Panic("panic")
				log.Panicf("%s", "panicf")
			})
		}).Run()
	})
}
