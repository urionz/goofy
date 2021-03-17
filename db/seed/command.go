package seed

import (
	"github.com/gookit/gcli/v3"
	"github.com/urionz/color"
	"github.com/urionz/goofy/contracts"
	"gorm.io/gorm"
)

func Seed(app contracts.Application) *gcli.Command {
	var db *gorm.DB
	command := &gcli.Command{
		Name: "db-seed",
		Desc: "生成数据",
		Func: func(cmd *gcli.Command, args []string) error {
			if err := app.Resolve(&db); err != nil {
				return err
			}
			for _, seeder := range seederFiles {
				if err := seeder.Handle(db); err != nil {
					color.Errorln(err)
					return err
				}
			}
			return nil
		},
	}
	return command
}
