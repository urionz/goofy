package web

import (
	"bytes"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/urionz/goofy/log"
	"github.com/urionz/goofy/utils"
)

var (
	state        sync.Mutex
	cmd          *exec.Cmd
	filesModTime = make(map[string]int64)
	scheduleTime time.Time
)

func NewWatcher(paths []string, files []string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Failed to create watcher: %s", err)
	}
	go func() {
		for {
			select {
			case e := <-watcher.Events:
				isBuild := true
				mt := utils.GetFileModTime(e.Name)
				if t := filesModTime[e.Name]; mt == t {
					isBuild = false
				}
				filesModTime[e.Name] = mt
				if isBuild {
					log.Infof("Event fired: %s", e)
					go func() {
						scheduleTime = time.Now().Add(1 * time.Second)
						time.Sleep(time.Until(scheduleTime))
						AutoBuild(files)
					}()
				}
			case err := <-watcher.Errors:
				log.Warnf("Watcher error: %s", err.Error())
			}
		}
	}()

	for _, path := range paths {
		if err = watcher.Add(path); err != nil {
			log.Fatalf("Failed to watch directory: %s", err)
		}
	}
}

func AutoBuild(files []string) {
	state.Lock()
	defer state.Unlock()

	var (
		err    error
		stderr bytes.Buffer
	)

	cmdName := "go"

	appName := "goofy"
	if runtime.GOOS == "windows" {
		appName += ".exe"
	}
	args := []string{"build"}
	args = append(args, "-o", appName)
	args = append(args, files...)

	bcmd := exec.Command(cmdName, args...)
	bcmd.Env = append(os.Environ(), "GOGC=off")
	bcmd.Stderr = &stderr
	if err = bcmd.Run(); err != nil {
		log.Errorf("Failed to build the application: %s", stderr.String())
		return
	}
	Restart(appName)
}

func Kill() {
	defer func() {
		// 尝试recover进程
		if e := recover(); e != nil {
			log.Infof("Kill recover: %s", e)
		}
	}()
	if cmd != nil && cmd.Process != nil {
		// 向当前进程发送关闭信号
		if runtime.GOOS == "windows" {
			cmd.Process.Signal(os.Kill)
		} else {
			cmd.Process.Signal(os.Interrupt)
		}

		// 等待进程结束完毕
		ch := make(chan struct{}, 1)
		go func() {
			cmd.Wait()
			ch <- struct{}{}
		}()

		select {
		// 正常接收关闭后 recover 重启进程
		case <-ch:
			return
		// 如果等待进程执行超过10秒则强制kill掉进程进入recover
		case <-time.After(10 * time.Second):
			log.Info("Timout. Force kill cmd process")
			if err := cmd.Process.Kill(); err != nil {
				log.Errorf("Error while killing cmd process: %s", err)
			}
			return
		}
	}
}

func Restart(appname string) {
	Kill()
	Start(appname)
}

func Start(appname string) {
	log.Infof("Restarting '%s'...", appname)
	if !strings.Contains(appname, "./") {
		appname = "./" + appname
	}
	cmd = exec.Command(appname, "web")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()
}
