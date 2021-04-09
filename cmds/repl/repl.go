package repl

import (
	"fmt"
	"os"
	"path"
	"runtime/debug"

	"github.com/c-bata/go-prompt"
	"github.com/urionz/goofy/cmds/repl/interpreter"
	"github.com/urionz/goofy/command"
	"github.com/urionz/goofy/contracts"
	"github.com/urionz/goofy/log"
)

var (
	currentInterpreter *interpreter.Interpreter
	suggestions        = []prompt.Suggest{
		{
			Text: "json.Marshal",
		},
	}
)

func Command(_ contracts.Application) *command.Command {
	return &command.Command{
		Name: "repl",
		Desc: "交互式命令工具",
		Func: func(c *command.Command, args []string) error {
			wd, err := os.Getwd()
			if err != nil {
				log.Panic(err)
			}
			currentInterpreter, err = interpreter.NewSession(wd)
			if err != nil {
				log.Panic(err)
			}
			_, err = currentInterpreter.Eval(":e 1")
			if err != nil {
				log.Panic(err)
			}
			p := prompt.New(handler, completer, prompt.OptionPrefix(">>> "))
			p.Run()
			return os.RemoveAll(path.Join(wd, ".repl"))
		},
	}
}

func completer(d prompt.Document) []prompt.Suggest {
	w := d.GetWordBeforeCursor()
	if w == "" {
		return []prompt.Suggest{}
	}
	return prompt.FilterHasPrefix(suggestions, w, true)
}

func handler(input string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Panic: %v\n%s", err, debug.Stack())
		}
	}()
	out, err := currentInterpreter.Eval(input)
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Print(out)
}
