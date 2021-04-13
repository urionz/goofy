package validation

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/urionz/goofy/errors"
	"github.com/urionz/goofy/govalidator"
	"github.com/urionz/goofy/web/context"
	"github.com/urionz/goutil/arrutil"
)

type Validation struct {
}

func NewValidation() *Validation {
	validate := &Validation{}
	validate.Register()
	return validate
}

func getMimeType(file multipart.File) (string, error) {
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return "", err
	}
	return http.DetectContentType(buffer), nil
}

func (v *Validation) Register() {
	govalidator.InterfaceParamTagMap["ext"] = func(in interface{}, all reflect.Value, params ...string) bool {
		var filename string
		if len(params) == 0 {
			return false
		}
		shouldExt := strings.Split(params[0], "|")
		switch reflect.TypeOf(in).Kind() {
		case reflect.TypeOf(&multipart.FileHeader{}).Kind():
			fh := in.(*multipart.FileHeader)
			filename = fh.Filename
			break
		case reflect.String:
			filename = in.(string)
			break
		}

		for _, ext := range shouldExt {
			if strings.HasSuffix(filename, ext) {
				return true
			}
		}

		return false
	}
	govalidator.InterfaceParamTagRegexMap["ext"] = regexp.MustCompile("^ext\\((.*)\\)")

	govalidator.InterfaceParamTagMap["mime"] = func(in interface{}, all reflect.Value, params ...string) bool {
		var mime string
		if len(params) == 0 {
			return false
		}
		shouldMimes := strings.Split(params[0], "|")
		switch reflect.TypeOf(in).Kind() {
		case reflect.TypeOf(&multipart.FileHeader{}).Kind():
			fh := in.(*multipart.FileHeader)
			file, err := fh.Open()
			if err != nil {
				return false
			}
			if mType, err := getMimeType(file); err != nil {
				return false
			} else {
				mime = mType
			}
			break
		}

		if arrutil.StringsHas(shouldMimes, mime) {
			return true
		}

		return false
	}
	govalidator.InterfaceParamTagRegexMap["mime"] = regexp.MustCompile("^mime\\((.*)\\)")
}

func (v *Validation) Validate(ctx *context.Context, reqPtr interface{}) *errors.CodeError {
	_ = ctx.ReadParams(reqPtr)
	_ = ctx.ReadQuery(reqPtr)
	_ = ctx.ReadForm(reqPtr)
	_ = ctx.ReadJSON(reqPtr)
	val := reflect.ValueOf(reqPtr)

	if val.Kind() != reflect.Ptr {
		return errors.NewErrorMsg(fmt.Sprintf("function only accepts structs ptr; got %s", reqPtr))
	}
	typ := reflect.TypeOf(reqPtr).Elem()
	val = val.Elem()
	numField := val.NumField()

	for i := 0; i < numField; i++ {
		if field := val.Field(i); field.Type().String() == "*multipart.FileHeader" {
			if name, exists := typ.Field(i).Tag.Lookup("form"); exists {
				if ctx.FileExists(name) {
					_, fh, _ := ctx.FormFile(name)
					field.Set(reflect.ValueOf(fh))
				}
			}
		}
	}

	_, err := govalidator.ValidateStruct(reqPtr)
	if e, ok := err.(govalidator.Errors); ok {
		if len(e) > 0 {
			return errors.FromError(e[0])
		}
	}
	return errors.FromError(err)
}
