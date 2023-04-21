package contracts

import (
	"io"
	"mime/multipart"
	"os"
	"reflect"
)

type FileInfo struct {
	Name string
	File *os.File
}

type PluginConstructor interface{}

type FilesystemFactory interface {
	Disk(name ...string) Filesystem
}

type PluginValue struct {
	ArgLen int
	Value  reflect.Value
}

type FileUnmarshaler func(data []byte, ptr interface{}) error

type Filesystem interface {
	Cloud
	// determine if a file exists
	Exists(path string) bool
	// get the os.File ptr of a file
	GetFile(path string) (*os.File, error)
	// get the file contents
	Get(path string) ([]byte, error)
	// Unmarshal .
	Unmarshal(path string, ptr interface{}, unmarshaler ...FileUnmarshaler) error
	// write the contents of a file
	Put(path string, contents []byte) (string, error)
	// Write a new file using a stream.
	WriteStream(path string, stream io.Reader) (string, error)
	// Upload a file
	Upload(path string, fh *multipart.FileHeader) (string, error)
	// delete the file at a given path
	Delete(path ...string) error
	// copy a file to a new location
	Copy(from, to string) (string, error)
	// move a file to a new location
	Move(from, to string) (string, error)
	// get the file size of a given file
	Size(path string) int64
	// get an array of all files name in a directory
	Files(dir string, recursive bool, child ...bool) []FileInfo
	// get an array of all files name in a directory (recursive)
	AllFiles(dir string) []FileInfo
	// add a plugin
	AddPlugin(constructor PluginConstructor, name string)
	// call a plugin method
	CallMethod(name string, args ...interface{}) []interface{}
	// get all plugins
	GetPlugins() map[string]PluginValue
}

type Cloud interface {
	Url(path string) string
	GetPreSignedUrl(path string) string
}

type CanOverwriteFiles interface {
}

type Adapter interface {
	Write(path string, content []byte) (string, error)
	WriteStream(path string, stream io.Reader) (string, error)
	Update(path string, content []byte) error
	UpdateSteam(path string, stream io.Reader) error
	Rename(path, newPath string) (string, error)
	Copy(from, to string) (string, error)
	Delete(path ...string) error
	DeleteDir(dir string) error
	CreateDir(dir string) error
	SetVisibility(path, visibility string) error
}
