package migrate

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"os"
	"path"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/golang-module/carbon"
	"github.com/gookit/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/urionz/collection"
	"github.com/urionz/goofy/command"
	"github.com/urionz/goofy/contracts"
	"github.com/urionz/goutil/fsutil"
	"github.com/urionz/goutil/strutil"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var createStub = `package migrations

import (
	"github.com/urionz/goofy/db/migrate"
	"github.com/urionz/goofy/db/model"
	"gorm.io/gorm"
)

func init() {
	migrate.Register(&{{ .StructName }}{})
}

type {{ .StructName }} struct {
	model.BaseModel
}

func (table *{{ .StructName }}) TableName() string {
	return "{{ .TableName }}"
}

func (table *{{ .StructName }}) MigrateTimestamp() int {
	return {{ .Timestamp }}
}

func (table *{{ .StructName }}) Up(db *gorm.DB) error {
	if !db.Migrator().HasTable(table) {
		return db.Migrator().CreateTable(table)
	}
	return nil
}

func (table *{{ .StructName }}) Down(db *gorm.DB) error {
	return db.Migrator().DropTable(table)
}
`

var blankStub = `package migrations

import (
	"github.com/urionz/goofy/db/migrate"
	"github.com/urionz/goofy/db/model"
	"gorm.io/gorm"
)

func init() {
	migrate.Register(&{{.StructName}}{})
}

type {{ .StructName }} struct {
	model.BaseModel
}

func (table *{{ .StructName }}) MigrateTimestamp() int {
	return {{ .Timestamp }}
}

func (table *{{ .StructName }}) TableName() string {
	return "{{ .TableName }}"
}

func (table *{{ .StructName }}) Up(db *gorm.DB) error {
	return nil
}

func (table *{{ .StructName }}) Down(db *gorm.DB) error {
	return nil
}
`

var (
	step      int
	driver    string
	tableName string
	create    string
	conn      *gorm.DB
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

func Migrate(app contracts.Application) *command.Command {
	cmd := &command.Command{
		Name: "migrate",
		Desc: "运行迁移",
		Config: func(c *command.Command) {
			c.StrOpt(&driver, "conn", "", "", "指定数据库连接")
			c.IntOpt(&step, "step", "s", 0, "指定迁移阶段")
		},
		Func: func(c *command.Command, args []string) error {

			if err := SwitchDBConnection(app); err != nil {
				color.Warnln(err)
				return nil
			}

			if err := RunMigrate(step); err != nil {
				color.Warnln(err)
				return nil
			}
			return nil
		},
	}

	return cmd
}

func Rollback(app contracts.Application) *command.Command {
	cmd := &command.Command{
		Name:     "migrate-rollback",
		Category: "migrate",
		Desc:     "迁移回滚",
		Config: func(c *command.Command) {
			c.StrOpt(&driver, "conn", "", "", "指定数据库连接")
			c.IntOpt(&step, "step", "s", 0, "指定迁移阶段")
		},
		Func: func(c *command.Command, args []string) error {
			if err := SwitchDBConnection(app); err != nil {
				color.Errorln(err)
				return nil
			}

			if err := RunRollback(step); err != nil {
				color.Errorln(err)
				return nil
			}
			return nil
		},
	}

	return cmd
}

func Refresh(app contracts.Application) *command.Command {
	cmd := &command.Command{
		Name: "migrate-refresh",
		Desc: "刷新迁移",
		Config: func(c *command.Command) {
			c.IntOpt(&step, "step", "s", 0, "指定迁移阶段")
		},
		Func: func(c *command.Command, args []string) error {
			var err error
			if err = SwitchDBConnection(app); err != nil {
				color.Errorln(err)
				return nil
			}

			if step > 0 {
				err = RunRollback(step)
			} else {
				err = RunReset()
			}

			if err != nil {
				color.Errorln(err)
				return nil
			}

			if err = RunMigrate(step); err != nil {
				color.Errorln(err)
				return nil
			}
			return nil
		},
	}

	return cmd
}

func Fresh(app contracts.Application) *command.Command {
	cmd := &command.Command{
		Name: "migrate-fresh",
		Desc: "migrate fresh",
		Config: func(c *command.Command) {
			c.IntOpt(&step, "step", "s", 0, "指定迁移阶段")
		},
		Func: func(c *command.Command, args []string) error {
			var err error
			var tables []string

			if err = SwitchDBConnection(app); err != nil {
				color.Errorln(err)
				return nil
			}

			if err := conn.Raw("show tables").Scan(&tables).Error; err != nil {
				color.Errorln(err)
				return nil
			}

			for _, t := range tables {
				if err := conn.Migrator().DropTable(t); err != nil {
					color.Errorln(err)
					return nil
				}
			}

			color.Infoln("Dropped all tables successfully.")

			if err := RunMigrate(step); err != nil {
				color.Errorln(err)
				return nil
			}
			return nil
		},
	}

	return cmd
}

func Status(app contracts.Application) *command.Command {
	cmd := &command.Command{
		Name: "migrate-status",
		Desc: "查看迁移状态",
		Func: func(c *command.Command, args []string) error {
			var batches map[string]int
			var ran []*Model
			var err error

			if err = SwitchDBConnection(app); err != nil {
				color.Errorln(err)
				return nil
			}

			repository := NewDBMigration(conn)
			ran, err = repository.GetRan()
			if err != nil {
				color.Errorln(err)
				return nil
			}
			ranNameCollection := collection.NewObjPointCollection(ran).Pluck("Migration")

			batches, err = repository.GetMigrationBatches()
			if err != nil {
				color.Errorln(err)
				return nil
			}

			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.AppendHeader(table.Row{"Ran?", "Migration", "Batch"})
			for _, migrateFile := range GetMigrationFiles() {
				migrationName := GetMigrationName(migrateFile)
				if ranNameCollection.Contains(migrationName) {
					t.AppendRow(table.Row{
						color.String("<green>Yes</>"),
						migrationName,
						strconv.Itoa(batches[migrationName]),
					})
				} else {
					t.AppendRow(table.Row{
						color.String("<red>No</>"),
						migrationName,
						"",
					})
				}
			}
			t.Render()
			return nil
		},
	}

	return cmd
}

func Reset(app contracts.Application) *command.Command {
	cmd := &command.Command{
		Name: "migrate-reset",
		Desc: "重置迁移",
		Func: func(c *command.Command, args []string) error {
			if err := SwitchDBConnection(app); err != nil {
				color.Errorln(err)
				return nil
			}
			if err := RunReset(); err != nil {
				color.Errorln(err)
				return nil
			}
			return nil
		},
	}

	return cmd
}

type MakeCommand struct {
	table  string
	create string
}

func Make(app contracts.Application) *command.Command {
	cmd := &command.Command{
		Name:     "make-migration",
		Category: "make",
		Desc:     "创建迁移文件",
		Config: func(c *command.Command) {
			c.StrOpt(&tableName, "table", "t", "", "The table to migrate")
			c.StrOpt(&create, "create", "c", "", "The table to be created")
			c.BindArg(&command.Argument{
				Name: "name", Desc: "迁移文件名称",
			})
		},
		Func: func(c *command.Command, args []string) error {
			var prompt *survey.Input
			var name string
			var isCreate bool
			if len(args) > 0 {
				name = args[0]
			}
			for {
				if name != "" || len(args) >= 1 {
					break
				}
				prompt = &survey.Input{
					Message: "请输入文件名称：",
				}
				survey.AskOne(prompt, &name)
			}

			if err := os.MkdirAll(path.Join(app.Database(), "migrations"), os.ModePerm); err != nil {
				color.Errorln(err)
				return nil
			}

			name = strutil.ToSnake(name)

			if tableName == "" && create != "" {
				tableName = create
				isCreate = true
			}

			if tableName == "" {
				tableName, isCreate = NewTableGuesser().Guess(name)
			}

			generatePath := path.Join(app.Database(), "migrations")

			if err := WriteMigration(name, tableName, generatePath, isCreate); err != nil {
				color.Errorln(err)
				return nil
			}

			color.Infoln("执行完毕")

			return nil
		},
	}

	return cmd
}

func getStub(isCreate bool) string {
	if isCreate {
		return createStub
	}

	return blankStub
}

func WriteMigration(name, table, generatePath string, isCreate bool) error {
	stub := getStub(isCreate)

	filePath := path.Join(generatePath, strings.ToLower(name)+".go")

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

	stubString, err := populateStub(stub, table)

	if err != nil {
		return err
	}

	if f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666); err == nil {
		f.WriteString(stubString)
	} else {
		return err
	}

	return nil
}

func populateStub(stub, table string) (string, error) {
	var templateBuffer bytes.Buffer
	tpl, err := template.New("migration").Parse(stub)
	if err != nil {
		return templateBuffer.String(), err
	}

	if err := tpl.ExecuteTemplate(&templateBuffer, "migration", map[string]interface{}{
		"StructName": strings.ToUpper(strutil.RandomChars(6)),
		"TableName":  table,
		"Timestamp":  carbon.Now().ToTimestamp(),
	}); err != nil {
		return templateBuffer.String(), err
	}

	return templateBuffer.String(), nil
}

func RunPending(migrations []contracts.MigrateFile, step int) error {
	repository := &Model{
		DB: conn,
	}

	if len(migrations) == 0 {
		err := errors.New("nothing to migrate")
		return err
	}

	batch, err := repository.GetNextBatchNumber()

	if err != nil {
		return err
	}

	for _, migrationFile := range migrations {
		if err := RunUp(migrationFile, batch); err != nil {
			return err
		}
		if step > 0 {
			batch++
		}
	}

	return nil
}

func RunUp(file contracts.MigrateFile, batch int) error {
	repository := &Model{
		DB: conn,
	}
	name := GetMigrationName(file)

	color.Infoln("Migrating:", name)

	if err := file.Up(conn); err != nil {
		return err
	}
	if err := repository.Log(GetMigrationName(file), batch); err != nil {
		return err
	}

	color.Infoln("Migrated:", name)

	return nil
}

func RunRollback(step int) error {
	repository := NewDBMigration(conn)
	dbMigrations, err := GetMigrationsForRollback(step, repository)

	if err != nil {
		return err
	}

	return RollbackMigrations(dbMigrations)
}

func RunReset() error {
	var migrations []*Model
	var err error
	migrations, err = NewDBMigration(conn).GetRan()
	if err != nil {
		return err
	}

	if len(migrations) == 0 || len(GetMigrationFiles()) == 0 {
		err = errors.New("nothing to rollback")
		return err
	}

	return RollbackMigrations(migrations)
}

func RunMigrate(step int) error {
	var ran []*Model
	var err error

	if ran, err = NewDBMigration(conn).GetRan(); err != nil {
		return err
	}

	migrateFiles := GetMigrationFiles()

	pendingMigrateFiles := GetPendingMigrations(migrateFiles, ran)

	SortFileMigrations(pendingMigrateFiles)

	if err := RunPending(pendingMigrateFiles, step); err != nil {
		return err
	}
	return nil
}

func GetMigrationFiles() []contracts.MigrateFile {
	return migrationFiles
}

func GetMigrationName(migrateFile contracts.MigrateFile) string {
	migrationNames := strings.Split(reflect.TypeOf(migrateFile).String(), ".")
	return fmt.Sprintf("%d_%s", migrateFile.MigrateTimestamp(), strutil.ToSnake(migrationNames[len(migrationNames)-1]))
}

func SortFileMigrations(files []contracts.MigrateFile) {
	sort.Slice(files, func(i, j int) bool {
		return files[i].MigrateTimestamp() < files[j].MigrateTimestamp()
	})
}

func GetPendingMigrations(files []contracts.MigrateFile, ran []*Model) []contracts.MigrateFile {
	var pendingMigrations []contracts.MigrateFile
	ranNameCollection := collection.NewObjPointCollection(ran).Pluck("Migration")
	for _, migrateFile := range files {
		if !ranNameCollection.Contains(GetMigrationName(migrateFile)) {
			pendingMigrations = append(pendingMigrations, migrateFile)
		}
	}
	return pendingMigrations
}

func GetMigrationsForRollback(step int, repository *Model) ([]*Model, error) {
	var dbMigrates []*Model
	var err error
	if step > 0 {
		dbMigrates, err = repository.GetMigrations(step)
	} else {
		dbMigrates, err = repository.GetLast()
	}
	return dbMigrates, err
}

func RollbackMigrations(migrations []*Model) error {
	files := GetMigrationFiles()

	existsFileMigrates := func(dbMigrate *Model) (contracts.MigrateFile, bool) {
		for _, migrateFile := range files {
			if dbMigrate.Migration == GetMigrationName(migrateFile) {
				return migrateFile, true
			}
		}
		return nil, false
	}

	for _, migration := range migrations {
		file, exists := existsFileMigrates(migration)

		if !exists {
			color.Warnln("Migration not found:", migration.Migration)
			continue
		}

		if err := RunDown(file, migration); err != nil {
			return err
		}
	}

	return nil
}

func RunDown(file contracts.MigrateFile, migration *Model) (err error) {
	repository := &Model{
		DB: conn,
	}

	name := GetMigrationName(file)

	color.Infoln("Rolling back:", name)

	if err := file.Down(conn); err != nil {
		return err
	}

	if err := repository.Delete(migration); err != nil {
		return err
	}

	color.Infoln("Rolled back:", name)

	return nil
}
