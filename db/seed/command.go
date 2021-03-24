package seed

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/gookit/gcli/v3"
	"github.com/urionz/color"
	"github.com/urionz/goofy/contracts"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	driver string
	conn   *gorm.DB
)

func SwitchDBConnection(app contracts.Application) error {
	var conf contracts.Config
	var manager contracts.DBFactory
	if err := app.Resolve(&conf); err != nil {
		color.Errorln(err)
		return err
	}
	if err := app.Resolve(&manager); err != nil {
		color.Errorln(err)
		return err
	}
	if driver == "" {
		sw := false
		prompt := &survey.Confirm{
			Message: "当前设置数据库连接为空，将使用默认连接，是否切换连接？",
		}
		survey.AskOne(prompt, &sw)
		if sw {
			var options []string
			conns := conf.Object("database.conns").Data()
			for conn, _ := range conns {
				options = append(options, conn)
			}
			prompt := &survey.Select{
				Message: "请选择一个数据库连接:",
				Options: options,
			}
			survey.AskOne(prompt, conn)
		}
	}

	conn = manager.Connection(driver)

	conn.Config.Logger = logger.Discard

	return nil
}

func Seed(app contracts.Application) *gcli.Command {
	command := &gcli.Command{
		Name: "db-seed",
		Desc: "生成数据",
		Func: func(cmd *gcli.Command, args []string) error {
			if err := SwitchDBConnection(app); err != nil {
				color.Warnln(err)
				return err
			}
			for _, seeder := range seederFiles {
				if err := seeder.Handle(conn); err != nil {
					color.Errorln(err)
					return err
				}
			}
			return nil
		},
	}
	return command
}
