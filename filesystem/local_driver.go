package filesystem

import (
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/urionz/goofy/contracts"
)

const (
	FilePublicPermission  os.FileMode = 0644
	FilePrivatePermission os.FileMode = 0600
	DirPublicPermission   os.FileMode = 0755
	DirPrivatePermission  os.FileMode = 0700

	DefaultFileFlags = os.O_CREATE | os.O_WRONLY | os.O_APPEND
)

type LocalDriver struct {
	Adapter
	pathPrefix string
	conf       contracts.Config
}

var _ contracts.Filesystem = (*LocalDriver)(nil)

func NewLocalDriver(prefix string, conf contracts.Config) *LocalDriver {
	return &LocalDriver{
		pathPrefix: prefix,
		conf:       conf,
	}
}

func (l *LocalDriver) Unmarshal(p string, ptr interface{}, unmarshaler ...contracts.FileUnmarshaler) error {
	content, err := l.Get(p)
	if err != nil {
		return err
	}
	return unmarshaler[0](content, ptr)
}

func (l *LocalDriver) Url(path string) string {
	return fmt.Sprintf("%s/%s", strings.TrimRight(l.conf.String("url"), "/"), l.applyPathPrefix(path))
}

func (l *LocalDriver) WriteStream(path string, stream io.Reader) (string, error) {
	var handler *os.File
	var err error
	location := l.applyPathPrefix(path)
	if err := l.ensureDirectory(filepath.Dir(location)); err != nil {
		return "", err
	}
	handler, _ = os.OpenFile(location, DefaultFileFlags, os.ModePerm)
	defer handler.Close()
	if _, err = io.Copy(handler, stream); err != nil {
		return "", err
	}
	if err = handler.Sync(); err != nil {
		return "", err
	}
	if err = handler.Chmod(FilePublicPermission); err != nil {
		return "", err
	}
	return l.Url(path), nil
}

func (l *LocalDriver) Upload(path string, fh *multipart.FileHeader) (string, error) {
	file, err := fh.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()
	return l.WriteStream(path, file)
}

func (l *LocalDriver) Put(path string, contents []byte) (string, error) {
	var handler *os.File
	var err error
	location := l.applyPathPrefix(path)
	if err := l.ensureDirectory(filepath.Dir(location)); err != nil {
		return "", err
	}
	handler, _ = os.OpenFile(location, DefaultFileFlags, os.ModePerm)
	defer handler.Close()
	if _, err = handler.Write(contents); err != nil {
		return "", err
	}
	if err = handler.Sync(); err != nil {
		return "", err
	}
	if err = handler.Chmod(FilePublicPermission); err != nil {
		return "", err
	}
	return l.Url(path), nil
}

func (l *LocalDriver) Delete(path ...string) error {
	for _, p := range path {
		if err := os.RemoveAll(l.applyPathPrefix(p)); err != nil {
			return err
		}
	}
	return nil
}

func (l *LocalDriver) Copy(from, to string) (string, error) {
	var fromHandler, toHandler *os.File
	var err error
	if !l.Exists(from) {
		return "", fmt.Errorf("the %s is not exists", from)
	}
	if fromHandler, err = l.GetFile(from); err != nil {
		return "", err
	}
	defer fromHandler.Close()
	toLocation := l.applyPathPrefix(to)
	if err := l.ensureDirectory(filepath.Dir(toLocation)); err != nil {
		return "", err
	}
	if toHandler, err = os.Create(toLocation); err != nil {
		return "", err
	}
	defer toHandler.Close()

	if _, err = io.Copy(toHandler, fromHandler); err != nil {
		return "", err
	}

	return l.Url(to), nil
}

func (l *LocalDriver) Move(from, to string) (string, error) {
	u, err := l.Copy(from, to)
	if err != nil {
		return "", err
	}
	if err := os.Remove(l.applyPathPrefix(from)); err != nil {
		return "", err
	}
	return u, nil
}

func (l *LocalDriver) Size(path string) int64 {
	var handler *os.File
	var stat os.FileInfo
	var err error
	if handler, err = l.GetFile(path); err != nil {
		return 0
	}
	if stat, err = handler.Stat(); err != nil {
		return 0
	}

	return stat.Size()
}

func (l *LocalDriver) Files(dir string, recursive bool, child ...bool) []contracts.FileInfo {
	var results []contracts.FileInfo
	var filesInfo []os.FileInfo
	var err error
	var op string
	if len(child) == 0 || !child[0] {
		dir = l.applyPathPrefix(dir)
	}
	if op, err = filepath.Abs(dir); err != nil {
		return results
	}
	if filesInfo, err = ioutil.ReadDir(op); err != nil {
		return results
	}

	for _, info := range filesInfo {
		if info.IsDir() {
			if recursive {
				results = append(results, l.Files(path.Join(dir, info.Name()), true, true)...)
			}
		} else {
			if filePtr, err := os.Open(path.Join(dir, info.Name())); err == nil {
				results = append(results, contracts.FileInfo{
					Name: info.Name(),
					File: filePtr,
				})
			}
		}
	}
	return results
}

func (l *LocalDriver) AllFiles(dir string) []contracts.FileInfo {
	return l.Files(dir, true)
}

func (l *LocalDriver) Exists(path string) bool {
	if path == "" {
		return false
	}
	if fi, err := os.Stat(l.applyPathPrefix(path)); err == nil {
		return !fi.IsDir()
	}
	return false
}

func (l *LocalDriver) Get(path string) ([]byte, error) {
	var file *os.File
	var err error
	if file, err = l.GetFile(path); err != nil {
		return []byte{}, err
	}
	return ioutil.ReadAll(file)
}

func (l *LocalDriver) GetFile(path string) (*os.File, error) {
	return os.Open(l.applyPathPrefix(path))
}

func (l *LocalDriver) ensureDirectory(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(path, DirPublicPermission)
		}
		return err
	}
	return nil
}

func (l *LocalDriver) applyPathPrefix(p string) string {
	return path.Join(l.getPathPrefix(), strings.TrimLeft(p, "/"))
}

func (l *LocalDriver) getPathPrefix() string {
	return l.pathPrefix
}
