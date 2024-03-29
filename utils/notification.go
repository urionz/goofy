package utils

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

const appName = "app"

func Notify(text, title string, enable bool) {
	if !enable {
		return
	}
	switch runtime.GOOS {
	case "darwin":
		osxNotify(text, title)
	case "linux":
		linuxNotify(text, title)
	case "windows":
		windowsNotify(text, title)
	}
}

func osxNotify(text, title string) {
	var cmd *exec.Cmd
	if existTerminalNotifier() {
		cmd = exec.Command("terminal-notifier", "-title", appName, "-message", text, "-subtitle", title)
	} else if MacOSVersionSupport() {
		notification := fmt.Sprintf("display notification \"%s\" with title \"%s\" subtitle \"%s\"", text, appName, title)
		cmd = exec.Command("osascript", "-e", notification)
	} else {
		cmd = exec.Command("growlnotify", "-n", appName, "-m", title)
	}
	cmd.Run()
}

func windowsNotify(text, title string) {
	exec.Command("growlnotify", "/i:", "", "/t:", title, text).Run()
}

func linuxNotify(text, title string) {
	exec.Command("notify-send", "-i", "", title, text).Run()
}

func existTerminalNotifier() bool {
	cmd := exec.Command("which", "terminal-notifier")
	err := cmd.Start()
	if err != nil {
		return false
	}
	err = cmd.Wait()
	return err != nil
}

func MacOSVersionSupport() bool {
	cmd := exec.Command("sw_vers", "-productVersion")
	check, _ := cmd.Output()
	version := strings.Split(string(check), ".")
	major, _ := strconv.Atoi(version[0])
	minor, _ := strconv.Atoi(version[1])
	if major < 10 || (major == 10 && minor < 9) {
		return false
	}
	return true
}
