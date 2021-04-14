package model

import (
	"bytes"
	"os"
	"path"
	"text/template"

	"github.com/AlecAivazis/survey/v2"
	"github.com/urionz/color"
	"github.com/urionz/goofy/command"
	"github.com/urionz/goofy/contracts"
	"github.com/urionz/goutil/fsutil"
	"github.com/urionz/goutil/strutil"
)

var serviceStub = `package services

import (
	"{{ .ModName }}/{{ .MvcPath }}models"
	"github.com/urionz/goofy/pagination"
	repo "{{ .ModName }}/{{ .MvcPath }}repositories"
	
	"github.com/urionz/goofy/db"
	"github.com/urionz/goofy/web"
)

var {{ ToCamel .Name }}Service = new{{ ToCamel .Name }}Service()

type {{ ToLowerCamel .Name }}Service struct {
}

func new{{ ToCamel .Name }}Service() *{{ ToLowerCamel .Name }}Service {
	return &{{ ToLowerCamel .Name }}Service{}
}

func (s {{ ToLowerCamel .Name }}Service) Get(id uint64) *models.{{ ToCamel .Name }} {
	return repo.{{ ToCamel .Name }}Repo.Get(db.Model(&models.{{ ToCamel .Name }}{}), id)
}

func (s {{ ToLowerCamel .Name }}Service) Take(where ...interface{}) *models.{{ ToCamel .Name }} {
	return repo.{{ ToCamel .Name }}Repo.Take(db.Model(&models.{{ ToCamel .Name }}{}), where...)
}

func (s {{ ToLowerCamel .Name }}Service) Find(cnd *db.SqlCnd) []models.{{ ToCamel .Name }} {
	return repo.{{ ToCamel .Name }}Repo.Find(db.Model(&models.{{ ToCamel .Name }}{}), cnd)
}

func (s {{ ToLowerCamel .Name }}Service) FindOne(cnd *db.SqlCnd) *models.{{ ToCamel .Name }} {
	return repo.{{ ToCamel .Name }}Repo.FindOne(db.Model(&models.{{ ToCamel .Name }}{}), cnd)
}

func (s {{ ToLowerCamel .Name }}Service) FindPageByParams(params *web.QueryParams) (list []models.{{ ToCamel .Name }}, paging *pagination.Paging) {
	return repo.{{ ToCamel .Name }}Repo.FindPageByParams(db.Model(&models.{{ ToCamel .Name }}{}), params)
}

func (s {{ ToLowerCamel .Name }}Service) FindPageByCnd(cnd *db.SqlCnd) (list []models.{{ ToCamel .Name }}, paging *pagination.Paging) {
	return repo.{{ ToCamel .Name }}Repo.FindPageByCnd(db.Model(&models.{{ ToCamel .Name }}{}), cnd)
}

func (s {{ ToLowerCamel .Name }}Service) Update(t *models.{{ ToCamel .Name }}) error {
	err := repo.{{ ToCamel .Name }}Repo.Update(db.Model(&models.{{ ToCamel .Name }}{}), t)
	return err
}

func (s {{ ToLowerCamel .Name }}Service) Updates(id uint64, columns map[string]interface{}) error {
	err := repo.{{ ToCamel .Name }}Repo.Updates(db.Model(&models.{{ ToCamel .Name }}{}), id, columns)
	return err
}

func (s {{ ToLowerCamel .Name }}Service) UpdateColumn(id uint64, name string, value interface{}) error {
	err := repo.{{ ToCamel .Name }}Repo.UpdateColumn(db.Model(&models.{{ ToCamel .Name }}{}), id, name, value)
	return err
}

func (s {{ ToLowerCamel .Name }}Service) Delete(id uint64, soft ...bool) error {
	if len(soft) > 0 && soft[0] {
		return repo.{{ ToCamel .Name }}Repo.SoftDelete(db.Model(&models.{{ ToCamel .Name }}{}), id)
	}
	return repo.{{ ToCamel .Name }}Repo.Delete(db.Model(&models.{{ ToCamel .Name }}{}), id)
}
`

func MakeService(app contracts.Application) *command.Command {
	return &command.Command{
		Name: "make-service",
		Desc: "创建service",
		Config: func(c *command.Command) {
			c.BindArg(&command.Argument{
				Name: "name", Desc: "service文件名称",
			})
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

			if err := os.MkdirAll(path.Join(app.Workspace(), mvcPath, "services"), os.ModePerm); err != nil {
				color.Errorln(err)
				return nil
			}

			if err := createServiceTpl(name, app.Workspace()); err != nil {
				color.Errorln(err)
				return err
			}

			color.Infoln("创建成功")

			return nil
		},
	}
}

func createServiceTpl(name, root string) error {
	var fHandle *os.File
	var err error
	var stubString string
	filePath := path.Join(root, mvcPath, "services", strutil.ToSnake(name)+".go")

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

	if stubString, err = servicePopulateStub(); err != nil {
		return err
	}

	if _, err = fHandle.WriteString(stubString); err != nil {
		return err
	}

	return nil
}

func servicePopulateStub() (string, error) {
	var templateBuffer bytes.Buffer
	tpl, err := template.New("services").Funcs(template.FuncMap{
		"ToCamel":      strutil.ToCamel,
		"ToLowerCamel": strutil.ToLowerCamel,
	}).Parse(serviceStub)
	if err != nil {
		return templateBuffer.String(), err
	}

	if err := tpl.ExecuteTemplate(&templateBuffer, "services", map[string]interface{}{
		"ModName": modName,
		"MvcPath": mvcPath,
		"Name":    name,
	}); err != nil {
		return templateBuffer.String(), err
	}

	return templateBuffer.String(), nil
}
