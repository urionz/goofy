package container

import "github.com/goava/di"

func As(interfaces ...Interface) ProvideOption {
	return di.As(interfaces...)
}
