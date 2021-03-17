package validator

import "github.com/goava/di"

type FactoryContract interface {
	Make()
}

type Factory struct {
	di.Tags `name:"validator"`
}

var _ FactoryContract = (*Factory)(nil)

func NewFactory() *Factory {
	return &Factory{}
}

func (f *Factory) Make() {
}
