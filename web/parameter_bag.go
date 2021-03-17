package web

import (
	"sync"
)

type Parameters map[string][]interface{}

type ParameterBag struct {
	parameters Parameters
	lock       sync.RWMutex
}

func NewParameterBag(parameters Parameters) *ParameterBag {
	return &ParameterBag{
		parameters: parameters,
	}
}

func NewEmptyParameterBag() *ParameterBag {
	return &ParameterBag{
		parameters: map[string][]interface{}{},
	}
}

func (params *ParameterBag) All() Parameters {
	return params.parameters
}

func (params *ParameterBag) Keys() []string {
	var keys []string
	for k, _ := range params.parameters {
		keys = append(keys, k)
	}
	return keys
}

func (params *ParameterBag) Replace(parameters Parameters) {
	params.lock.Lock()
	defer params.lock.Unlock()
	params.parameters = parameters
}

func (params *ParameterBag) Add(key string, value interface{}) {
	params.lock.Lock()
	defer params.lock.Unlock()
	params.parameters[key] = append(params.parameters[key], value)
}

func (params *ParameterBag) Get(key string, def ...interface{}) interface{} {
	params.lock.RLock()
	defer params.lock.RUnlock()
	if len(def) == 0 {
		def = append(def, nil)
	}
	if params.parameters == nil {
		return def[0]
	}
	vs := params.parameters[key]
	if len(vs) == 0 {
		return def[0]
	}
	return vs[0]
}

func (params *ParameterBag) Set(key string, value interface{}) {
	params.lock.Lock()
	defer params.lock.Unlock()
	params.parameters[key] = []interface{}{value}
}

func (params *ParameterBag) Remove(key string) {
	params.lock.Lock()
	defer params.lock.Unlock()
	delete(params.parameters, key)
}

func (params *ParameterBag) Has(key string) bool {
	params.lock.RLock()
	defer params.lock.RUnlock()
	_, ok := params.parameters[key]
	return ok
}
