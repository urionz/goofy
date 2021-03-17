package contracts

import "os"

type FileInfo struct {
	Name string
	File *os.File
}

type FilesystemFactory interface {
	Disk(name ...string) Filesystem
}

type Filesystem interface {
	// determine if a file exists
	Exists(path string) bool
	// get the os.File ptr of a file
	GetFile(path string) (*os.File, error)
	// get the file contents
	Get(path string) ([]byte, error)
	// write the contents of a file
	Put(path string, contents []byte) error
	// delete the file at a given path
	Delete(path ...string) error
	// copy a file to a new location
	Copy(from, to string) error
	// move a file to a new location
	Move(from, to string) error
	// get the file size of a given file
	Size(path string) int64
	// get an array of all files name in a directory
	Files(dir string, recursive bool, child ...bool) []FileInfo
	// get an array of all files name in a directory (recursive)
	AllFiles(dir string) []FileInfo
}
