package web

import (
	"sort"

	"github.com/urionz/goofy/errors"
	"github.com/urionz/goofy/validator"
)

type IValidation interface {
	Rules(ctx *Context) validator.MapData
	Messages(ctx *Context) validator.MapData
}

type Validation struct {
}

type ValidationOption struct {
}

func validate(bindingFunc func(rdp IValidation) error, ctx *Context, reqDataPtr IValidation, maxMemoryLimit ...int64) *errors.CodeError {
	ctx.RecordRequestBody(true)
	_ = bindingFunc(reqDataPtr)
	ctx.RecordRequestBody(false)
	opts := validator.Options{
		Request:         ctx.Request(),
		Rules:           reqDataPtr.Rules(ctx),
		Messages:        reqDataPtr.Messages(ctx),
		RequiredDefault: true,
	}
	v := validator.New(opts)
	if err := v.Validate(); len(err) > 0 {
		sortedErrorFields := make([]string, 0)
		for k := range err {
			sortedErrorFields = append(sortedErrorFields, k)
		}
		sort.Strings(sortedErrorFields)
		return errors.NewErrorMsg(err.Get(sortedErrorFields[len(sortedErrorFields)-1]))
	}
	return nil
}

func (v *Validation) Validate(ctx *Context, reqDataPtr IValidation, maxMemoryLimit ...int64) *errors.CodeError {
	return validate(func(rdp IValidation) error {
		_ = ctx.ReadParams(reqDataPtr)
		_ = ctx.ReadQuery(reqDataPtr)
		_ = ctx.ReadForm(reqDataPtr)
		_ = ctx.ReadJSON(reqDataPtr)
		return nil
	}, ctx, reqDataPtr, maxMemoryLimit...)
}
