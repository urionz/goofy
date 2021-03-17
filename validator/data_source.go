package validator

import (
	"fmt"
	"mime/multipart"
	"net/url"
	"reflect"

	"github.com/urionz/goofy/filter"
	"github.com/urionz/goutil/strutil"
)

type DataFace interface {
	Get(key string) (interface{}, bool)
	Set(field string, val interface{}) (interface{}, error)
}

type FormData struct {
	Form       url.Values
	Files      map[string]*multipart.FileHeader
	jsonBodies []byte
}

func newFormData() *FormData {
	return &FormData{
		Form:  make(map[string][]string),
		Files: make(map[string]*multipart.FileHeader),
	}
}

func (d *FormData) Add(key string, value string) {
	d.Form.Add(key, value)
}

func (d *FormData) Get(key string) (interface{}, bool) {
	if vs, ok := d.Form[key]; ok && len(vs) > 0 {
		return vs[0], true
	}
	if fh, ok := d.Files[key]; ok {
		return fh, true
	}
	return nil, false
}

func (d *FormData) AddValues(values url.Values) {
	for key, vals := range values {
		for _, val := range vals {
			d.Form.Add(key, val)
		}
	}
}

func (d *FormData) AddFiles(filesMap map[string][]*multipart.FileHeader) {
	for key, files := range filesMap {
		if len(files) != 0 {
			d.AddFile(key, files[0])
		}
	}
}

func (d *FormData) AddFile(key string, file *multipart.FileHeader) {
	d.Files[key] = file
}

func (d *FormData) Del(key string) {
	d.Form.Del(key)
}

// DelFile deletes the file associated with key (if any).
// If there is no file associated with key, it does nothing.
func (d *FormData) DelFile(key string) {
	delete(d.Files, key)
}

func (d *FormData) Set(field string, val interface{}) (newVal interface{}, err error) {
	newVal = val
	switch val.(type) {
	case string:
		d.Form.Set(field, val.(string))
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		newVal = strutil.MustString(val)
		d.Form.Set(field, newVal.(string))
	default:
		err = fmt.Errorf("set value failure for field: %s", field)
	}
	return
}

type MapDataSource struct {
	// Map the source map data
	Map map[string]interface{}
	// from reflect Map
	value reflect.Value
	// bodyJSON from the original JSON bytes/string.
	// available for FromJSONBytes(), FormJSON().
	bodyJSON []byte
	// map field reflect.Value caches
	// fields map[string]reflect.Value
}

func (d *MapDataSource) Set(field string, val interface{}) (interface{}, error) {
	d.Map[field] = val
	return val, nil
}

// Get value by key
func (d *MapDataSource) Get(field string) (interface{}, bool) {
	// if fv, ok := d.fields[field]; ok {
	// 	return fv, true
	// }

	return filter.GetByPath(field, d.Map)
}
