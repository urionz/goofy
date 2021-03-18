package repl

import (
	"os"

	"github.com/gookit/gcli/v3"
	"github.com/urionz/goofy/contracts"
)

var (
	outWriter  = os.Stdout
	errWriter  = os.Stderr
	autoImport bool
	extFiles   string
	pkg        string
)

func Command(_ contracts.Application) *gcli.Command {
	return &gcli.Command{
		Name: "repl",
		Desc: "交互式命令工具",
		Config: func(c *gcli.Command) {
			c.BoolOpt(&autoImport, "autoimport", "", true, "自动import")
			c.StrOpt(&extFiles, "context", "", "", "import packages, functions, variables and constants from external golang source files")
			c.StrOpt(&pkg, "pkg", "", "", "the package where the session will be run inside")
		},
		Func: func(c *gcli.Command, args []string) error {
			return New(
				AutoImport(autoImport),
				ExtFiles(extFiles),
				PackageName(pkg),
				OutWriter(outWriter),
				ErrWriter(errWriter),
			).Run()
		},
	}
}
