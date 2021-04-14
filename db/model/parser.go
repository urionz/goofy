package model

import (
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/urionz/goutil/arrutil"
)

var (
	blankReg            = regexp.MustCompile("\\s?")
	multiLineCommentReg = regexp.MustCompile("/\\*+\\s?.*?\\s?\\*+/")
	singleCommentReg    = regexp.MustCompile("//[^\\r\\n]*")
	structReg           = regexp.MustCompile("type\\s+\\w+\\s+struct\\s+{\\s+[\\w.]+BaseModel\\s+([\\w\\s`\\[\\]{:\"_.;*|,]+)\\s+}")
	fieldReg            = regexp.MustCompile("\\s?(\\w+)\\s+([\\w*\\[\\].]+)\\s+.*")
)

func ParseGoFileStruct(filepath string) map[string]string {
	result := make(map[string]string)
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		return result
	}
	content := string(b)
	content = multiLineCommentReg.ReplaceAllString(content, "")
	content = singleCommentReg.ReplaceAllString(content, "")

	structMatch := structReg.FindStringSubmatch(content)
	fieldMatch := strings.TrimSpace(fieldReg.ReplaceAllString(structMatch[1], "$1"))
	typeMatch := strings.TrimSpace(fieldReg.ReplaceAllString(structMatch[1], "$2"))

	fields := strings.Split(fieldMatch, "\n")
	types := strings.Split(typeMatch, "\n")

	for index, field := range fields {
		typ := blankReg.ReplaceAllString(types[index], "")

		if strings.Contains(typ, "*") {
			typ = "ptr"
		} else if strings.Contains(typ, "[") {
			typ = "array"
		} else if strings.Contains(typ, "JSON") {
			typ = "ptr"
		} else if arrutil.StringsHas([]string{
			"int", "int8", "int16", "int32", "int64",
			"uint", "uint8", "uint16", "uint32", "uint64",
			"float32", "float64"}, typ) {
			typ = "number"
		}
		result[blankReg.ReplaceAllString(field, "")] = typ
	}
	return result
}
