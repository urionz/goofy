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

var repoStub = `package repo

import (
	"github.com/urionz/goofy/db"
	"{{ PathJoin .ModName .MvcPath "models" }}"
	"github.com/urionz/goofy/pagination"
	"github.com/urionz/goofy/web"

	"gorm.io/gorm"
)

var {{ ToCamel .Name }}Repo = new{{ ToCamel .Name }}Repo()

type {{ ToLowerCamel .Name }}Repo struct {
}

func new{{ ToCamel .Name }}Repo() *{{ ToLowerCamel .Name }}Repo {
	return &{{ ToLowerCamel .Name }}Repo{}
}

func (r *{{ ToLowerCamel .Name }}Repo) Get(db *gorm.DB, id uint) *models.{{ ToCamel .Name }} {
	ret := &models.{{ ToCamel .Name }}{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *{{ ToLowerCamel .Name }}Repo) Take(db *gorm.DB, where ...interface{}) *models.{{ ToCamel .Name }} {
	ret := &models.{{ ToCamel .Name }}{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *{{ ToLowerCamel .Name }}Repo) Find(db *gorm.DB, cnd *db.SqlCnd) (list []models.{{ ToCamel .Name }}) {
	cnd.Find(db, &list)
	return
}

func (r *{{ ToLowerCamel .Name }}Repo) FindOne(db *gorm.DB, cnd *db.SqlCnd) *models.{{ ToCamel .Name }} {
	ret := &models.{{ ToCamel .Name }}{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *{{ ToLowerCamel .Name }}Repo) FindPageByParams(db *gorm.DB, params *web.QueryParams) (list []models.{{ ToCamel .Name }}, paging *pagination.Paging) {
	return r.FindPageByCnd(db, &params.SqlCnd)
}

func (r *{{ ToLowerCamel .Name }}Repo) FindPageByCnd(db *gorm.DB, cnd *db.SqlCnd) (list []models.{{ ToCamel .Name }}, paging *pagination.Paging) {
	cnd.Find(db, &list)
	count, _ := cnd.Count(db, &models.{{ ToCamel .Name }}{})

	paging = &pagination.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *{{ ToLowerCamel .Name }}Repo) Create(db *gorm.DB, t *models.{{ ToCamel .Name }}) (err error) {
	err = db.Create(t).Error
	return
}

func (r *{{ ToLowerCamel .Name }}Repo) Update(db *gorm.DB, t *models.{{ ToCamel .Name }}) (err error) {
	err = db.Save(t).Error
	return
}

func (r *{{ ToLowerCamel .Name }}Repo) Updates(db *gorm.DB, id uint, columns map[string]interface{}) (err error) {
	err = db.Model(&models.{{ ToCamel .Name }}{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *{{ ToLowerCamel .Name }}Repo) UpdateColumn(db *gorm.DB, id uint, name string, value interface{}) (err error) {
	err = db.Model(&models.{{ ToCamel .Name }}{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *{{ ToLowerCamel .Name }}Repo) SoftDelete(db *gorm.DB, id uint) (err error) {
	err = db.Delete(&models.{{ ToCamel .Name }}{}, "id = ?", id).Error
	return
}

func (r *{{ ToLowerCamel .Name }}Repo) Delete(db *gorm.DB, id uint) (err error) {
	err = db.Unscoped().Delete(&models.{{ ToCamel .Name }}{}, "id = ?", id).Error
	return
}
`

func MakeRepo(app contracts.Application) *command.Command {
	return &command.Command{
		Name: "make-repo",
		Desc: "创建repo",
		Config: func(c *command.Command) {
			c.BindArg(&command.Argument{
				Name: "name", Desc: "repo文件名称",
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

			if err := os.MkdirAll(path.Join(app.Workspace(), mvcPath, "repositories"), os.ModePerm); err != nil {
				color.Errorln(err)
				return nil
			}

			if err := createRepo(name, app.Workspace()); err != nil {
				color.Errorln(err)
				return err
			}

			color.Infoln("创建成功")

			return nil
		},
	}
}

func createRepo(name, root string) error {
	var fHandle *os.File
	var err error
	var stubString string

	filePath := path.Join(root, mvcPath, "repositories", strutil.ToSnake(name)+".go")

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

	if stubString, err = repoPopulateStub(); err != nil {
		return err
	}

	if _, err = fHandle.WriteString(stubString); err != nil {
		return err
	}

	return nil
}

func repoPopulateStub() (string, error) {
	var templateBuffer bytes.Buffer
	tpl, err := template.New("Repo").Funcs(template.FuncMap{
		"ToCamel":      strutil.ToCamel,
		"ToLowerCamel": strutil.ToLowerCamel,
		"ToSnake":      strutil.ToSnake,
		"PathJoin":     path.Join,
	}).Parse(repoStub)
	if err != nil {
		return templateBuffer.String(), err
	}

	if err := tpl.ExecuteTemplate(&templateBuffer, "Repo", map[string]interface{}{
		"ModName": modName,
		"MvcPath": mvcPath,
		"Name":    name,
	}); err != nil {
		return templateBuffer.String(), err
	}

	return templateBuffer.String(), nil
}
