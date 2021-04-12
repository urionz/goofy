package storage

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"path"
	"path/filepath"
	"time"

	"github.com/urionz/goofy/filesystem"
	"github.com/urionz/goutil/arrutil"
	"github.com/urionz/goutil/strutil"
)

type Storage struct {
	*filesystem.Manager
	name  string
	mimes []string
	disk  string
}

func NewStorage(manager *filesystem.Manager) *Storage {
	return &Storage{
		Manager: manager,
		name:    "file",
		mimes:   []string{"*"},
	}
}

func (s *Storage) FileName(name string) *Storage {
	s.name = name
	return s
}

func (s *Storage) AllowMime(mimes []string) *Storage {
	s.mimes = mimes
	return s
}

func (s *Storage) Disk(disk ...string) *Storage {
	if len(disk) >= 0 && disk[0] != "" {
		s.disk = disk[0]
	}
	return s
}

func (s *Storage) Upload(req *http.Request, savePath string) (string, error) {
	f, fh, err := req.FormFile(s.name)
	if err != nil {
		return "", err
	}
	uploadMineType, err := s.MimeType(f)
	if err != nil {
		return "", err
	}
	if !arrutil.StringsHas(s.mimes, uploadMineType) && !arrutil.StringsHas(s.mimes, "*") {
		return "", fmt.Errorf("the file mine type is illegal")
	}
	t := time.Now()
	date := t.Format("20060102")
	if ext := filepath.Ext(savePath); ext == "" {
		savePath = path.Join(date, savePath, strutil.Md5(fmt.Sprintf("%s%s", fh.Filename, date))+"."+filepath.Ext(fh.Filename))
	} else {
		savePath = path.Join(date, savePath)
	}
	return s.Manager.Disk(s.disk).WriteStream(savePath, f)
}

func (*Storage) MimeType(file multipart.File) (string, error) {
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return "", err
	}
	return http.DetectContentType(buffer), nil
}
