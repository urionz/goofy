package utils

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/urionz/goofy/log"
)

// GetFileModTime returns unix timestamp of `os.File.ModTime` for the given path.
func GetFileModTime(path string) int64 {
	path = strings.Replace(path, "\\", "/", -1)
	f, err := os.Open(path)
	if err != nil {
		log.Errorf("Failed to open file on '%s': %s", path, err)
		return time.Now().Unix()
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		log.Errorf("Failed to get file stats: %s", err)
		return time.Now().Unix()
	}

	return fi.ModTime().Unix()
}

// SplitQuotedFields is like strings.Fields but ignores spaces
// inside areas surrounded by single quotes.
// To specify a single quote use backslash to escape it: '\''
func SplitQuotedFields(in string) []string {
	type stateEnum int
	const (
		inSpace stateEnum = iota
		inField
		inQuote
		inQuoteEscaped
	)
	state := inSpace
	r := []string{}
	var buf bytes.Buffer

	for _, ch := range in {
		switch state {
		case inSpace:
			if ch == '\'' {
				state = inQuote
			} else if !unicode.IsSpace(ch) {
				buf.WriteRune(ch)
				state = inField
			}

		case inField:
			if ch == '\'' {
				state = inQuote
			} else if unicode.IsSpace(ch) {
				r = append(r, buf.String())
				buf.Reset()
			} else {
				buf.WriteRune(ch)
			}

		case inQuote:
			if ch == '\'' {
				state = inField
			} else if ch == '\\' {
				state = inQuoteEscaped
			} else {
				buf.WriteRune(ch)
			}

		case inQuoteEscaped:
			buf.WriteRune(ch)
			state = inQuote
		}
	}

	if buf.Len() != 0 {
		r = append(r, buf.String())
	}

	return r
}

// loadPathsToWatch loads the paths that needs to be watched for changes
func LoadPathsToWatch(paths *[]string) error {
	directory, err := os.Getwd()
	if err != nil {
		return err
	}
	filepath.Walk(directory, func(path string, info os.FileInfo, _ error) error {
		if strings.HasSuffix(info.Name(), "docs") {
			return filepath.SkipDir
		}
		if strings.HasSuffix(info.Name(), "swagger") {
			return filepath.SkipDir
		}
		if strings.HasSuffix(info.Name(), "vendor") {
			return filepath.SkipDir
		}

		if filepath.Ext(info.Name()) == ".go" {
			*paths = append(*paths, path)
		}
		return nil
	})
	return nil
}

// GoCommand executes the passed command using Go tool
func GoCommand(command string, args ...string) error {
	allargs := []string{command}
	allargs = append(allargs, args...)
	goBuild := exec.Command("go", allargs...)
	goBuild.Stderr = os.Stderr
	return goBuild.Run()
}
