package dlv

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/go-delve/delve/pkg/terminal"
	"github.com/go-delve/delve/service"
	"github.com/go-delve/delve/service/debugger"
	"github.com/go-delve/delve/service/rpc2"
	"github.com/go-delve/delve/service/rpccommon"
	"github.com/gookit/gcli/v3"
	"github.com/urionz/goofy/contracts"
	"github.com/urionz/goofy/log"
	"github.com/urionz/goofy/utils"
)

var (
	packageName string
	verbose     bool
	port        int
	args        string
)

func Command(app contracts.Application) *gcli.Command {
	cmd := &gcli.Command{
		Name: "dlv",
		Desc: "使用delve工具开启调试模式",
		Config: func(c *gcli.Command) {
			c.StrOpt(&packageName, "package", "", "", "The package to debug (Must have a main package)")
			c.BoolOpt(&verbose, "verbose", "v", false, "Enable verbose mode")
			c.IntOpt(&port, "port", "p", 8181, "Port to listen to for clients")
			c.StrOpt(&args, "args", "", "", "Port to listen to for clients")
		},
		Func: func(cmd *gcli.Command, args []string) error {
			var conf contracts.Config
			if err := app.Resolve(&conf); err != nil {
				return err
			}
			runDlv()
			return nil
		},
	}

	return cmd
}

func runDlv() {
	var (
		addr       = fmt.Sprintf("127.0.0.1:%d", port)
		paths      = make([]string, 0)
		notifyChan = make(chan int)
	)

	if err := utils.LoadPathsToWatch(&paths); err != nil {
		log.Fatalf("Error while loading paths to watch: %v", err.Error())
	}
	go startWatcher(paths, notifyChan)
	startDelveDebugger(addr, notifyChan)
}

// buildDebug builds a debug binary in the current working directory
func buildDebug() (string, error) {
	args := []string{"-gcflags", "-N -l", "-o", "debug"}
	args = append(args, utils.SplitQuotedFields("-ldflags='-linkmode internal'")...)
	args = append(args, packageName)
	if err := utils.GoCommand("build", args...); err != nil {
		return "", err
	}

	fp, err := filepath.Abs("./debug")
	if err != nil {
		return "", err
	}
	return fp, nil
}

// startDelveDebugger starts the Delve debugger server
func startDelveDebugger(addr string, ch chan int) int {
	log.Info("Starting Delve Debugger...")

	fp, err := buildDebug()
	if err != nil {
		log.Fatalf("Error while building debug binary: %v", err)
	}
	defer os.Remove(fp)

	abs, err := filepath.Abs("./debug")
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Create and start the debugger server
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Could not start listener: %s", err)
	}
	defer listener.Close()

	processArgs := []string{abs}
	processArgs = append(processArgs, strings.Split(args, " ")...)

	server := rpccommon.NewServer(&service.Config{
		Listener:    listener,
		AcceptMulti: true,
		APIVersion:  2,
		ProcessArgs: processArgs,
		Debugger: debugger.Config{
			AttachPid:  0,
			WorkingDir: ".",
			Backend:    "default",
		},
	})
	if err := server.Run(); err != nil {
		log.Fatalf("Could not start debugger server: %v", err)
	}

	// Start the Delve client REPL
	client := rpc2.NewClient(addr)
	// Make sure the client is restarted when new changes are introduced
	go func() {
		for {
			if val := <-ch; val == 0 {
				if _, err := client.Restart(true); err != nil {
					utils.Notify("Error while restarting the client: "+err.Error(), "bee", true)
				} else {
					if verbose {
						utils.Notify("Delve Debugger Restarted", "bee", true)
					}
				}
			}
		}
	}()

	// Create the terminal and connect it to the client debugger
	term := terminal.New(client, nil)
	status, err := term.Run()
	if err != nil {
		log.Fatalf("Could not start Delve REPL: %v", err)
	}

	// Stop and kill the debugger server once user quits the REPL
	if err := server.Stop(); err != nil {
		log.Fatalf("Could not stop Delve server: %v", err)
	}
	return status
}

var eventsModTime = make(map[string]int64)

// startWatcher starts the fsnotify watcher on the passed paths
func startWatcher(paths []string, ch chan int) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Could not start the watcher: %v", err)
	}
	defer watcher.Close()

	// Feed the paths to the watcher
	for _, path := range paths {
		if err := watcher.Add(path); err != nil {
			log.Fatalf("Could not set a watch on path: %v", err)
		}
	}

	for {
		select {
		case evt := <-watcher.Events:
			build := true
			if filepath.Ext(evt.Name) != ".go" {
				continue
			}

			mt := utils.GetFileModTime(evt.Name)
			if t := eventsModTime[evt.Name]; mt == t {
				build = false
			}
			eventsModTime[evt.Name] = mt

			if build {
				go func() {
					if verbose {
						utils.Notify("Rebuilding application with the new changes", "bee", true)
					}

					// Wait 1s before re-build until there is no file change
					scheduleTime := time.Now().Add(1 * time.Second)
					time.Sleep(time.Until(scheduleTime))
					_, err := buildDebug()
					if err != nil {
						utils.Notify("Build Failed: "+err.Error(), "bee", true)
					} else {
						ch <- 0 // Notify listeners
					}
				}()
			}
		case err := <-watcher.Errors:
			if err != nil {
				ch <- -1
			}
		}
	}
}
