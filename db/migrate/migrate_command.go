package migrate

import (
	"bytes"
	"errors"
	"html/template"
	"os"
	"path"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/golang-module/carbon"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/progress"
	"github.com/gookit/gcli/v3/show"
	"github.com/urionz/collection"
	"github.com/urionz/color"
	"github.com/urionz/goofy/contracts"
	"github.com/urionz/goofy/log"
	"github.com/urionz/goutil/fsutil"
	"github.com/urionz/goutil/strutil"
	"gorm.io/gorm"
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
	step   int
	conn   string
	table  string
	create string
)

func SwitchDBConnection(app contracts.Application) error {
	var conf contracts.Config
	if err := app.Resolve(&conf); err != nil {
		log.Fatal(err)
		return err
	}
	if conn == "" {
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
	return nil
}

func Migrate(app contracts.Application) *gcli.Command {
	command := &gcli.Command{
		Name: "migrate",
		Desc: "运行迁移",
		Config: func(c *gcli.Command) {
			c.StrOpt(&conn, "conn", "", "", "指定数据库连接")
			c.IntOpt(&step, "step", "s", 0, "指定迁移阶段")
		},
		Func: func(c *gcli.Command, args []string) error {
			var manager contracts.DBFactory
			if err := app.Resolve(&manager); err != nil {
				return err
			}

			if err := SwitchDBConnection(app); err != nil {
				log.Fatal(err)
			}

			changeDir(path.Join(app.Database(), "migrations"))
			if err := RunMigrate(step, manager.Connection(conn)); err != nil {
				log.Fatal(err)
			}
			return nil
		},
	}

	return command
}

func Rollback(app contracts.Application) *gcli.Command {
	command := &gcli.Command{
		Name:     "migrate-rollback",
		Category: "migrate",
		Desc:     "迁移回滚",
		Config: func(c *gcli.Command) {
			c.IntOpt(&step, "step", "s", 0, "指定迁移阶段")
		},
		Func: func(c *gcli.Command, args []string) error {
			var manager contracts.DBFactory
			if err := app.Resolve(&manager); err != nil {
				return err
			}

			if err := SwitchDBConnection(app); err != nil {
				log.Fatal(err)
			}

			if err := RunRollback(step, manager.Connection(conn)); err != nil {
				log.Error(err)
			}
			return nil
		},
	}

	return command
}

func Refresh(app contracts.Application) *gcli.Command {
	command := &gcli.Command{
		Name: "migrate-refresh",
		Desc: "刷新迁移",
		Config: func(c *gcli.Command) {
			c.IntOpt(&step, "step", "s", 0, "指定迁移阶段")
		},
		Func: func(c *gcli.Command, args []string) error {
			var manager contracts.DBFactory
			var err error
			if err = app.Resolve(&manager); err != nil {
				return err
			}

			if err = SwitchDBConnection(app); err != nil {
				log.Fatal(err)
			}

			if step > 0 {
				err = RunRollback(step, manager.Connection(conn))
			} else {
				err = RunReset(manager.Connection(conn))
			}

			if err != nil {
				log.Error(err)
				return err
			}

			if err = RunMigrate(step, manager.Connection()); err != nil {
				log.Error(err)
				return err
			}
			return nil
		},
	}

	return command
}

func Fresh(app contracts.Application) *gcli.Command {
	command := &gcli.Command{
		Name: "migrate-fresh",
		Desc: "migrate fresh",
		Config: func(c *gcli.Command) {
			c.IntOpt(&step, "step", "s", 0, "指定迁移阶段")
		},
		Func: func(c *gcli.Command, args []string) error {
			var manager contracts.DBFactory
			if err := app.Resolve(&manager); err != nil {
				return err
			}
			var tables []string
			db := manager.Connection()
			if err := db.Raw("show tables").Scan(&tables).Error; err != nil {
				return err
			}

			for _, table := range tables {
				if err := db.Migrator().DropTable(table); err != nil {
					return err
				}
			}

			log.Info("Dropped all tables successfully.")

			if err := RunMigrate(step, manager.Connection()); err != nil {
				log.Error(err)
				return err
			}
			return nil
		},
	}

	return command
}

func Status(app contracts.Application) *gcli.Command {
	var db *gorm.DB
	command := &gcli.Command{
		Name: "migrate-status",
		Desc: "查看迁移状态",
		Func: func(c *gcli.Command, args []string) error {
			if err := app.Resolve(&db); err != nil {
				return err
			}
			var batches map[string]int
			var ran []*Model
			var err error
			repository := NewDBMigration(db)
			ran, err = repository.GetRan()
			if err != nil {
				log.Error(err)
				return nil
			}
			ranNameCollection := collection.NewObjPointCollection(ran).Pluck("Migration")

			batches, err = repository.GetMigrationBatches()
			if err != nil {
				log.Error(err)
				return err
			}
			table := show.NewTable("migrate status")
			table.Cols = []string{"Ran?", "Migration", "Batch"}
			for _, migrateFile := range GetMigrationFiles() {
				migrationName := GetMigrationName(migrateFile)
				if ranNameCollection.Contains(migrationName) {
					table.Cols = []string{
						color.String("<green>Yes</>"),
						migrationName,
						strconv.Itoa(batches[migrationName]),
					}
				} else {
					table.Cols = []string{
						color.String("<red>No</>"),
						migrationName,
						"",
					}
				}
			}
			table.SetOutput(os.Stdout)
			table.Println()
			return nil
		},
	}

	return command
}

func Reset(app contracts.Application) *gcli.Command {
	command := &gcli.Command{
		Name: "migrate-reset",
		Desc: "重置迁移",
		Func: func(c *gcli.Command, args []string) error {
			var manager contracts.DBFactory
			if err := app.Resolve(&manager); err != nil {
				return err
			}
			if err := SwitchDBConnection(app); err != nil {
				log.Fatal(err)
			}
			if err := RunReset(manager.Connection(conn)); err != nil {
				log.Error(err)
				return err
			}
			return nil
		},
	}

	return command
}

type MakeCommand struct {
	table  string
	create string
}

func Make(app contracts.Application) *gcli.Command {
	command := &gcli.Command{
		Name:     "make-migration",
		Category: "make",
		Desc:     "创建迁移文件",
		Config: func(c *gcli.Command) {
			c.StrOpt(&table, "table", "t", "", "The table to migrate")
			c.StrOpt(&create, "create", "c", "", "The table to be created")
		},
		Func: func(c *gcli.Command, args []string) error {
			var prompt *survey.Input
			var name string
			var isCreate bool
			for {
				if name != "" || len(args) >= 1 {
					break
				}
				prompt = &survey.Input{
					Message: "请输入文件名称：",
				}
				survey.AskOne(prompt, &name)
			}
			if name == "" {
				name = args[0]
			}

			if err := os.MkdirAll(path.Join(app.Database(), "migrations"), os.ModePerm); err != nil {
				log.Error(err)
			}

			name = strutil.ToSnake(name)

			if table == "" && create != "" {
				table = create
				isCreate = true
			}

			if table == "" {
				table, isCreate = NewTableGuesser().Guess(name)
			}

			generatePath := path.Join(app.Database(), "migrations")

			if err := WriteMigration(name, table, generatePath, isCreate); err != nil {
				log.Error(err)
				return err
			}

			log.Info("执行完毕")

			return nil
		},
	}

	return command
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

func RunPending(migrations []contracts.MigrateFile, step int, db *gorm.DB) error {
	repository := &Model{
		DB: db,
	}

	if len(migrations) == 0 {
		err := errors.New("nothing to migrate")
		return err
	}

	batch, err := repository.GetNextBatchNumber()

	if err != nil {
		return err
	}

	p := progress.Bar(len(migrations))
	p.Start()

	for _, migrationFile := range migrations {

		if err := RunUp(migrationFile, batch, db); err != nil {
			return err
		}

		if step > 0 {
			batch++
		}
		p.Advance()
	}

	p.Finish()

	return nil
}

func RunUp(file contracts.MigrateFile, batch int, db *gorm.DB) error {
	repository := &Model{
		DB: db,
	}
	name := GetMigrationName(file)

	log.Info("Migrating:", name)

	if err := file.Up(db); err != nil {
		return err
	}
	if err := repository.Log(name, batch); err != nil {
		return err
	}

	log.Info("Migrated:", name)

	return nil
}

func RunRollback(step int, db *gorm.DB) error {
	repository := NewDBMigration(db)
	dbMigrations, err := GetMigrationsForRollback(step, repository)

	if err != nil {
		return err
	}

	return RollbackMigrations(dbMigrations, db)
}

func RunReset(db *gorm.DB) error {
	var migrations []*Model
	var err error
	migrations, err = NewDBMigration(db).GetRan()
	if err != nil {
		return err
	}

	if len(migrations) == 0 || len(GetMigrationFiles()) == 0 {
		err = errors.New("nothing to rollback")
		return err
	}

	return RollbackMigrations(migrations, db)
}

func RunMigrate(step int, db *gorm.DB) error {
	var ran []*Model
	var err error

	if ran, err = NewDBMigration(db).GetRan(); err != nil {
		return err
	}

	migrateFiles := GetMigrationFiles()

	pendingMigrateFiles := GetPendingMigrations(migrateFiles, ran)

	SortFileMigrations(pendingMigrateFiles)

	if err := RunPending(pendingMigrateFiles, step, db); err != nil {
		return err
	}
	return nil
}

func GetMigrationFiles() []contracts.MigrateFile {
	return migrationFiles
}

func GetMigrationName(migrateFile contracts.MigrateFile) string {
	migrationNames := strings.Split(reflect.TypeOf(migrateFile).String(), ".")
	return strutil.ToSnake(migrationNames[len(migrationNames)-1])
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

func RollbackMigrations(migrations []*Model, db *gorm.DB) error {
	files := GetMigrationFiles()

	existsFileMigrates := func(dbMigrate *Model) (contracts.MigrateFile, bool) {
		for _, migrateFile := range files {
			migrationNames := strings.Split(reflect.TypeOf(migrateFile).String(), ".")
			migrationName := strutil.ToSnake(migrationNames[len(migrationNames)-1])
			if dbMigrate.Migration == migrationName {
				return migrateFile, true
			}
		}
		return nil, false
	}

	for _, migration := range migrations {
		file, exists := existsFileMigrates(migration)

		if !exists {
			log.Warn("Migration not found:", migration.Migration)
			continue
		}

		if err := RunDown(file, migration, db); err != nil {
			return err
		}
	}

	return nil
}

func RunDown(file contracts.MigrateFile, migration *Model, db *gorm.DB) (err error) {
	repository := &Model{
		DB: db,
	}

	name := GetMigrationName(file)

	log.Info("Rolling back:", name)

	if err := file.Down(db); err != nil {
		return err
	}

	if err := repository.Delete(migration); err != nil {
		return err
	}

	log.Info("Rolled back:", name)

	return nil
}

func changeDir(dir string) {
	if err := os.Chdir(dir); err != nil {
		log.Fatal(err)
	}
}
