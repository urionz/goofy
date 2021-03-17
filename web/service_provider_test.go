package web_test

import (
	"testing"

	"github.com/urionz/goofy/config"

	"github.com/stretchr/testify/require"
	"github.com/urionz/goofy"
	"github.com/urionz/goofy/web"
)

func TestNewServiceProvider(t *testing.T) {
	require.NotPanics(t, func() {
		goofy.Default.AddServices(config.NewServiceProvider, web.NewServiceProvider)
	})
}
