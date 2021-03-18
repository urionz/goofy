package interpreter

import (
	"errors"
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Interpreter struct {
	shellCmdOutput   string
	imports          ImportDatas
	types            map[string]string
	funcs            map[string]string
	vars             Vars
	tmpCodes         []int
	code             []string
	tmpIfusedAsValue string
	sessionDir       string
	Writer           io.Writer
	continueMode     bool
	indents          int
}

const helpText = `
List of REPL commands:
:help => shows help
:doc => shows go documentation of package/function
:e => evaluates expression
:pop => pop latest code from interpreter
:dump => dumps current interpreter
:file => prints go file generated from interpreter
:vars => shows only vars of current state
:types => shows only types of current state
:funcs => shows only funcs of current state
:imports => shows only imports of current state
`
const moduleTemplate = `module shell
go 1.13
%s
`

func wrapInPrint(code string) string {
	return fmt.Sprintf(`fmt.Printf("<%%T> %%+v\n", %s, %s)`, code, code)
}
func (s *Interpreter) importsForSource() string {
	return s.imports.String()
}

func (s *Interpreter) addImport(im []ImportData) {
	s.imports = append(s.imports, im...)
}

func (s *Interpreter) appendToLastCode(code string) {
	if len(s.code) == 0 {
		s.code = append(s.code, code)
		return
	}
	s.code[len(s.code)-1] += "\n" + code
	return
}

func (s *Interpreter) addCode(t Type, code string) (string, error) {
	if s.continueMode {
		s.appendToLastCode(code)
		indents, shouldContinue := ShouldContinue(s.code[len(s.code)-1])
		s.indents = indents
		if !shouldContinue {
			s.continueMode = false
			code = s.code[len(s.code)-1]
			s.code = s.code[:len(s.code)-1]
			return s.Eval(code)
		}
		return "", nil
	}
	indents, shouldContinue := ShouldContinue(code)
	s.indents = indents
	if shouldContinue {
		s.continueMode = true
		s.code = append(s.code, code)
		return "", nil
	}
	switch t {
	case Shell:
		return s.handleShellCommands(code)
	case Import:
		s.addImport(ExtractImportData(code))
		return "", nil
	case Print:
		s.code = append(s.code, code)
		s.tmpCodes = append(s.tmpCodes, len(s.code)-1)
		return "", nil
	case TypeDecl:
		s.addType(ExtractTypeName(code), code)
		return "", nil
	case FuncDecl:
		s.addFunc(ExtractFuncName(code), code)
		return "", nil
	case VarDecl:
		s.addVar(NewVar(code))
		return "", nil
	case Empty:
		return "", nil
	case Expr:
		s.tmpIfusedAsValue = code
		return s.addCode(Print, wrapInPrint(code))
	default:
		s.code = append(s.code, code)
		return "", nil
	}
}

func (s *Interpreter) Eval(code string) (string, error) {
	if code == "exit" {
		fmt.Println("Bye ...")
		os.Exit(0)
	}
	s.removeTmpCodes()
	typ, err := Parse(code)
	if err != nil {
		return "", err
	}
	_, err = s.addCode(typ, code)
	if err != nil {
		return "", err
	}
	if s.continueMode {
		return strings.Repeat("...", s.indents), nil
	}
	if typ != Shell {
		if err := checkIfHasParsingError(s.String()); err != nil {
			s.removeLastCode()
			return "", errors.New(err.Error() + "\n")
		}
	}
	return s.eval(), nil
}

// used as value
func createTmpDir(workingDirectory string) (string, error) {
	sessionDir := workingDirectory + "/.repl/sessions/" + fmt.Sprint(time.Now().Nanosecond())
	err := os.MkdirAll(sessionDir, 500)
	if err != nil {
		return sessionDir, err
	}
	return sessionDir, nil
}

func (s *Interpreter) removeTmpCodes() {
	for _, t := range s.tmpCodes {
		s.code[t] = ""
	}
	s.tmpCodes = s.tmpCodes[:0]
	for idx, c := range s.code {
		if c == "" {
			s.code = append(s.code[:idx], s.code[idx+1:]...)
		}
	}
}

func NewSession(workingDirectory string) (*Interpreter, error) {
	sessionDir, err := createTmpDir(workingDirectory)
	if err != nil {
		return nil, err
	}
	err = os.Chdir(sessionDir)
	if err != nil {
		panic(err)
	}
	session := &Interpreter{
		shellCmdOutput: "",
		imports:        ImportDatas{},
		types:          map[string]string{},
		funcs:          map[string]string{},
		vars:           map[string]Var{},
		tmpCodes:       []int{},
		code:           []string{},
		sessionDir:     sessionDir,
		Writer:         nil,
		continueMode:   false,
		indents:        0,
	}
	currentModule := getModuleNameOfCurrentProject(workingDirectory)
	if err = session.createModule(workingDirectory, currentModule); err != nil {
		return nil, err
	}
	err = session.writeToFile()
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (s *Interpreter) removeTmpDir() error {
	return os.RemoveAll(s.sessionDir)
}

func (s *Interpreter) writeToFile() error {
	return ioutil.WriteFile(s.sessionDir+"/main.go", []byte(s.String()), 500)
}

func (s *Interpreter) removeLastCode() {
	if len(s.code) == 0 {
		s.code = []string{}
		return
	}
	idx := len(s.code) - 1
	for tmpIdx, t := range s.tmpCodes {
		if t == idx {
			s.tmpCodes = append(s.tmpCodes[:tmpIdx], s.tmpCodes[tmpIdx+1:]...)
		}
	}
	s.code = s.code[:len(s.code)-1]
}

func checkIfErrIsNotDecl(err string) bool {
	return strings.Contains(err, "not used") && !strings.Contains(err, "evaluated")
}

func checkIfHasParsingError(code string) error {
	fs := token.NewFileSet()
	_, err := parser.ParseFile(fs, "", code, parser.AllErrors)
	if err != nil {
		return err
	}
	return nil
}
func checkIfErrIsUsedAsValue(err string) bool {
	return strings.Contains(err, "used as value")
}

func (s *Interpreter) eval() string {
	if s.shellCmdOutput != "" {
		output := s.shellCmdOutput
		s.shellCmdOutput = ""
		return output + "\n"
	}
	if len(s.code) == 0 {
		return ""
	}
	if s.continueMode {
		return strings.Repeat("...", s.indents)
	}
	if err := checkIfHasParsingError(s.String()); err != nil {
		s.removeLastCode()
		return err.Error() + "\n"
	}
	err := s.writeToFile()
	if err != nil {
		return err.Error()
	}
	err = os.Chdir(s.sessionDir)
	if err != nil {
		panic(err)
	}
	cmdImport := exec.Command("goimports", "-w", "main.go")
	out, err := cmdImport.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("%s %s\n", string(out), err.Error())
	}
	cmdRun := exec.Command("go", "run", "main.go")
	out, err = cmdRun.CombinedOutput()
	if err != nil {
		if checkIfErrIsNotDecl(string(out)) {
			return fmt.Sprintf("%s %s\n", string(out), err.Error())
		} else if checkIfErrIsUsedAsValue(string(out)) {
			s.removeLastCode()
			if _, err = s.addCode(Unknown, s.tmpIfusedAsValue); err != nil {
				return err.Error()
			}
			return s.eval()
		} else {
			s.removeLastCode()
			return fmt.Sprintf("%s %s\n", string(out), err.Error())
		}
	}
	return fmt.Sprintf("%s", out)
}

func (s *Interpreter) String() string {
	code := "package main\n%s\n%s\n%s\n%s\nfunc main() {\n%s\n}"
	return fmt.Sprintf(code, s.importsForSource(), s.typesForSource(), s.funcsForSource(), "var(\n"+s.vars.String()+"\n)", strings.Join(s.code, "\n"))
}
