package web

import (
	"github.com/gookit/goutil/jsonutil"
	"github.com/urionz/goofy/errors"
	"github.com/urionz/goofy/pagination"
	"github.com/urionz/goutil"
)

type JsonResult struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Success bool        `json:"success"`
}

func Json(code int, message string, data interface{}, success bool) *JsonResult {
	return &JsonResult{
		Code:    code,
		Message: message,
		Data:    data,
		Success: success,
	}
}

func JsonDataCode(code int, data interface{}) *JsonResult {
	return &JsonResult{
		Code:    code,
		Data:    data,
		Success: true,
	}
}

func JsonData(data interface{}) *JsonResult {
	return &JsonResult{
		Code:    0,
		Data:    data,
		Success: true,
	}
}

func JsonPageData(results interface{}, page *pagination.Paging) *JsonResult {
	return JsonData(&pagination.PageResult{
		Results: results,
		Page:    page,
	})
}

func JsonPageMapData(results interface{}, page interface{}) *JsonResult {
	var paging pagination.Paging
	pageByte, err := jsonutil.Encode(page)
	if err != nil {
		return JsonError(errors.FromError(err))
	}
	if err = jsonutil.Decode(pageByte, &paging); err != nil {
		return JsonError(errors.FromError(err))
	}
	return JsonData(&pagination.PageResult{
		Results: results,
		Page:    &paging,
	})
}

func JsonCursorData(results interface{}, cursor string) *JsonResult {
	return JsonData(&pagination.CursorResult{
		Results: results,
		Cursor:  cursor,
	})
}

func JsonSuccess() *JsonResult {
	return &JsonResult{
		Code:    0,
		Data:    nil,
		Success: true,
	}
}

func JsonError(err *errors.CodeError) *JsonResult {
	return &JsonResult{
		Code:    err.Code,
		Message: err.Message,
		Data:    err.Data,
		Success: false,
	}
}

func JsonErrorMsg(message string) *JsonResult {
	return &JsonResult{
		Code:    0,
		Message: message,
		Data:    nil,
		Success: false,
	}
}

func JsonErrorCode(code int, message string) *JsonResult {
	return &JsonResult{
		Code:    code,
		Message: message,
		Data:    nil,
		Success: false,
	}
}

func JsonErrorData(code int, message string, data interface{}) *JsonResult {
	return &JsonResult{
		Code:    code,
		Message: message,
		Data:    data,
		Success: false,
	}
}

type RspBuilder struct {
	Data map[string]interface{}
}

func NewEmptyRspBuilder() *RspBuilder {
	return &RspBuilder{Data: make(map[string]interface{})}
}

func NewRspBuilder(obj interface{}) *RspBuilder {
	return NewRspBuilderExcludes(obj)
}

func NewRspBuilderExcludes(obj interface{}, excludes ...string) *RspBuilder {
	return &RspBuilder{Data: goutil.StructToMap(obj, excludes...)}
}

func (builder *RspBuilder) Put(key string, value interface{}) *RspBuilder {
	builder.Data[key] = value
	return builder
}

func (builder *RspBuilder) Build() map[string]interface{} {
	return builder.Data
}

func (builder *RspBuilder) JsonResult() *JsonResult {
	return JsonData(builder.Data)
}
