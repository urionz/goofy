package model

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jinzhu/inflection"
	"github.com/urionz/color"
	"github.com/urionz/goofy/command"
	"github.com/urionz/goofy/contracts"
	"github.com/urionz/goutil"
	"github.com/urionz/goutil/fsutil"
	"github.com/urionz/goutil/strutil"
)

var modelStub = `package models

import (
	"github.com/urionz/goofy/db/model"
)

type {{ .StructName }} struct {
	model.BaseModel
}

func (*{{ .StructName }}) TableName() string {
	return "{{ .TableName }}"
}

func (*{{ .StructName }}) Connection() string {
	return ""
}
`

var (
	name       string
	migration  bool
	repository bool
	service    bool
	mvcPath    string
	modName    string
)

func resolveConf(app contracts.Application) error {
	var conf contracts.Config
	if err := app.Resolve(&conf); err != nil {
		color.Errorln(err)
		return err
	}
	mvcPath = conf.String("app.mvc_path")
	modName = goutil.GetModName(app.Workspace())
	return nil
}

func Make(app contracts.Application) *command.Command {
	cmd := &command.Command{
		Name: "make-model",
		Desc: "创建数据模型",
		Config: func(c *command.Command) {
			c.BoolOpt(&migration, "migration", "m", false, "同时创建迁移文件")
			c.BoolOpt(&repository, "repo", "r", false, "同时创建repo文件")
			c.BoolOpt(&service, "service", "s", false, "同时创建service文件")
		},
		Func: func(c *command.Command, args []string) error {
			if err := resolveConf(app); err != nil {
				color.Errorln(err)
				return nil
			}
			if len(args) > 0 {
				name = args[0]
			}
			var prompt *survey.Input
			for {
				if name != "" || len(args) >= 1 {
					break
				}
				prompt = &survey.Input{
					Message: "请输入文件名称：",
				}
				survey.AskOne(prompt, &name)
			}
			tableName := strutil.ToSnake(inflection.Plural(name))

			if err := os.MkdirAll(path.Join(app.Workspace(), mvcPath, "models"), os.ModePerm); err != nil {
				color.Errorln(err)
				return nil
			}

			if !migration {
				migrateConfirm := &survey.Confirm{
					Message: "是否同时创建迁移文件？",
				}
				survey.AskOne(migrateConfirm, &migration)
			}

			if migration {
				if err := createMigration(app, tableName); err != nil {
					color.Errorln(err)
					return nil
				}
			}

			if err := createModel(name, tableName, app.Workspace()); err != nil {
				color.Errorln(err)
				return nil
			}

			if !repository {
				repoConfirm := &survey.Confirm{
					Message: "是否同时创建repo文件？",
				}
				survey.AskOne(repoConfirm, &repository)
			}

			if repository {
				if err := createRepository(app, name); err != nil {
					color.Errorln(err)
					return nil
				}
			}

			if !service {
				serviceConfirm := &survey.Confirm{
					Message: "是否同时创建service文件？",
				}
				survey.AskOne(serviceConfirm, &service)
			}

			if service {
				if err := createService(app, name); err != nil {
					color.Errorln(err)
					return nil
				}
			}
			return nil
		},
	}
	return cmd
}

func createModel(name, tableName, root string) error {
	var fHandle *os.File
	var err error
	var stubString string
	fileName := strings.ToLower(strutil.ToSnake(name))

	structName := strutil.ToCamel(name)

	filePath := path.Join(root, mvcPath, "models", fileName+".go")

	var cover bool

	if fsutil.FileExists(filePath) {
		confirm := &survey.Confirm{
			Message: "此文件已存在，是否覆盖？",
		}
		survey.AskOne(confirm, &cover)
	}

	if cover {
		fHandle, err = os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	} else {
		fHandle, err = os.OpenFile(filePath, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666)
	}

	if err != nil {
		return err
	}

	if stubString, err = modelPopulateStub(structName, tableName, modelStub); err != nil {
		return err
	}

	if _, err = fHandle.WriteString(stubString); err != nil {
		return err
	}

	return nil
}

func modelPopulateStub(structName, tableName, stub string) (string, error) {
	var templateBuffer bytes.Buffer
	tpl, err := template.New("model").Parse(stub)
	if err != nil {
		return templateBuffer.String(), err
	}

	if err := tpl.ExecuteTemplate(&templateBuffer, "model", map[string]interface{}{
		"StructName": structName,
		"TableName":  tableName,
	}); err != nil {
		return templateBuffer.String(), err
	}

	return templateBuffer.String(), nil
}

func createMigration(app contracts.Application, tableName string) error {
	app.Call("make-migration", fmt.Sprintf("create_%s_table", tableName))
	return nil
}

func createRepository(app contracts.Application, name string) error {
	app.Call("make-repo", name)
	return nil
}

func createService(app contracts.Application, name string) error {
	app.Call("make-service", name)
	return nil
}
