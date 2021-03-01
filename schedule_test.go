package goofy_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urionz/goofy"
)

func TestApplication_AddSchedules(t *testing.T) {
	t.Run("add schedule func type", func(t *testing.T) {
		app := goofy.New()
		err := app.AddSchedules(goofy.ScheduleJob{
			"0/1 * * * *": goofy.Jobs(func() {}, func() {}),
		}).Error()
		require.NoError(t, err)
	})
	t.Run("add schedule job func type", func(t *testing.T) {
		app := goofy.New()
		var Job goofy.FuncJob = func() {

		}
		err := app.AddSchedules(goofy.ScheduleJob{
			"0/1 * * * *": goofy.Jobs(Job),
		}).Error()
		require.NoError(t, err)
	})
	t.Run("add schedule job other type", func(t *testing.T) {
		app := goofy.New()
		type other struct{}
		err := app.AddSchedules(goofy.ScheduleJob{
			"0/1 * * * *": goofy.Jobs(new(other)),
		}).Error()
		require.Error(t, err)
	})
}
