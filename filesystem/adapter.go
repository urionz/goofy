package filesystem

import (
	"os"

	"github.com/urionz/goofy/contracts"
)

type Adapter struct {
	driver contracts.Filesystem
}

func (a *Adapter) Exists(path string) bool {
	return a.driver.Exists(path)
}

func (a *Adapter) Get(path string) ([]byte, error) {
	return a.driver.Get(path)
}

func (a *Adapter) GetFile(path string) (*os.File, error) {
	return a.driver.GetFile(path)
}

func (a *Adapter) Put(path string, contents []byte) error {
	return a.driver.Put(path, contents)
}

func (a *Adapter) Delete(path ...string) error {
	return a.driver.Delete(path...)
}

func (a *Adapter) Copy(from, to string) error {
	return a.driver.Copy(from, to)
}

func (a *Adapter) Move(from, to string) error {
	return a.driver.Move(from, to)
}

func (a *Adapter) Size(path string) int64 {
	return a.driver.Size(path)
}

func (a *Adapter) Files(dir string, recursive bool, child ...bool) []contracts.FileInfo {
	return a.driver.Files(dir, recursive, child...)
}

func (a *Adapter) AllFiles(dir string) []contracts.FileInfo {
	return a.driver.AllFiles(dir)
}

var _ contracts.Filesystem = (*Adapter)(nil)

func NewAdapter(driver contracts.Filesystem) *Adapter {
	return &Adapter{
		driver: driver,
	}
}
