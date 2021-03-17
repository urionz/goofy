package validator

import "github.com/goava/di"

type Inject struct {
	di.Tags
}

type injectable interface {
	isInjectable()
}

func (Inject) isInjectable() {
}
