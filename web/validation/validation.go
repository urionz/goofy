package validation

import (
	"sort"

	"github.com/urionz/goofy/errors"
	"github.com/urionz/goofy/validator"
	"github.com/urionz/goofy/web/context"
)

type IValidation interface {
	Rules(ctx *context.Context) validator.MapData
	Messages(ctx *context.Context) validator.MapData
}

type BaseValidator struct {
}

func (*BaseValidator) Rules(_ *context.Context) validator.MapData {
	return nil
}

func (*BaseValidator) Messages(_ *context.Context) validator.MapData {
	return nil
}

type Validation struct {
}

func validate(bindingFunc func(rdp IValidation) error, ctx *context.Context, reqDataPtr IValidation, maxMemoryLimit ...int64) *errors.CodeError {
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

func (v *Validation) Validate(ctx *context.Context, reqDataPtr IValidation, maxMemoryLimit ...int64) *errors.CodeError {
	return validate(func(rdp IValidation) error {
		_ = ctx.ReadParams(reqDataPtr)
		_ = ctx.ReadQuery(reqDataPtr)
		_ = ctx.ReadForm(reqDataPtr)
		_ = ctx.ReadJSON(reqDataPtr)
		return nil
	}, ctx, reqDataPtr, maxMemoryLimit...)
}
