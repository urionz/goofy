package filesystem

import (
	"io"
	"mime/multipart"
	"os"
	"reflect"

	"github.com/urionz/goofy/contracts"
)

const (
	VisibilityPublic  = "public"
	VisibilityPrivate = "private"
)

type Adapter struct {
	driver  contracts.Filesystem
	plugins map[string]contracts.PluginValue
}

var _ contracts.Filesystem = (*Adapter)(nil)

func NewAdapter(driver contracts.Filesystem) *Adapter {
	return &Adapter{
		driver:  driver,
		plugins: map[string]contracts.PluginValue{},
	}
}

func (a *Adapter) GetPlugins() map[string]contracts.PluginValue {
	return a.plugins
}

func (a *Adapter) AddPlugin(constructor contracts.PluginConstructor, name string) {
	if _, ok := a.plugins[name]; !ok {
		a.plugins[name] = contracts.PluginValue{
			ArgLen: reflect.TypeOf(constructor).NumIn(),
			Value:  reflect.ValueOf(constructor),
		}
	}
}

func (a *Adapter) CallMethod(name string, args ...interface{}) []interface{} {
	method, exists := a.plugins[name]
	if !exists {
		return []interface{}{}
	}
	input := make([]reflect.Value, method.ArgLen)
	for index, arg := range args {
		input[index] = reflect.ValueOf(arg)
	}
	result := method.Value.Call(input)
	var output []interface{}
	for _, o := range result {
		output = append(output, o.Interface())
	}
	return output
}

func (a *Adapter) Url(path string) string {
	return a.driver.Url(path)
}

func (a *Adapter) GetPreSignedUrl(path string) string {
	return a.driver.GetPreSignedUrl(path)
}

func (a *Adapter) WriteStream(path string, stream io.Reader) (string, error) {
	return a.driver.WriteStream(path, stream)
}

func (a *Adapter) Upload(path string, fh *multipart.FileHeader) (string, error) {
	return a.driver.Upload(path, fh)
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

func (a *Adapter) Put(path string, contents []byte) (string, error) {
	return a.driver.Put(path, contents)
}

func (a *Adapter) Delete(path ...string) error {
	return a.driver.Delete(path...)
}

func (a *Adapter) Copy(from, to string) (string, error) {
	return a.driver.Copy(from, to)
}

func (a *Adapter) Move(from, to string) (string, error) {
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
