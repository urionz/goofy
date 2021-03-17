package contracts

type Config interface {
	Get(key string, findByPath ...bool) interface{}
	Set(key string, val interface{}, setByPath ...bool) error
	String(key string, defVal ...string) string
	Strings(key string) (arr []string)
	Int(key string, defVal ...int) int
	Ints(key string) (arr []int)
	Int64(key string, defVal ...int64) (value int64)
	Uint(key string, defVal ...uint) (value uint)
	Bool(key string, defVal ...bool) bool
	Env(key string, defVal interface{}) interface{}
	Exists(key string, findByPath ...bool) bool
	LoadExists(...string) error
	Object(key string, findByPath ...bool) Config
	Data() map[string]interface{}
}
