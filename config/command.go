package config

import (
	"fmt"
	"os"
	"path"

	"github.com/AlecAivazis/survey/v2"
	"github.com/gookit/goutil/fsutil"
	"github.com/urionz/goofy/command"
	"github.com/urionz/goofy/contracts"
	"github.com/urionz/goofy/log"
)

var (
	mode = "dev"
)

func Command(app contracts.Application) *command.Command {
	cmd := &command.Command{
		Name: "make-conf",
		Desc: "生产环境配置文件",
		Config: func(c *command.Command) {
			c.StrOpt(&mode, "mode", "", "dev", "生成环境")
		},
		Func: func(cmd *command.Command, args []string) error {
			filePath := path.Join(app.Workspace(), fmt.Sprintf("config.%s.toml", mode))
			if fsutil.FileExists(filePath) {
				cover := false
				prompt := &survey.Confirm{
					Message: "该迁移文件已存在，是否覆盖？",
				}
				survey.AskOne(prompt, &cover)
				if !cover {
					return nil
				}
			}

			if f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666); err == nil {
				f.WriteString(contracts.GetConfTpl())
				f.Close()
			} else {
				return err
			}

			envPath := path.Join(app.Workspace(), ".env")

			sw := false
			prompt := &survey.Confirm{
				Message: "是否切换当前运行环境？",
			}
			survey.AskOne(prompt, &sw)
			if !sw {
				return nil
			}

			if f, err := os.OpenFile(envPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666); err == nil {
				f.Close()
			} else {
				log.Error(err)
				return err
			}

			return nil
		},
	}
	return cmd
}
