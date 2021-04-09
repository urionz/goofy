package log

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Field = zap.Field

func Binary(key string, val []byte) Field {
	return zap.Binary(key, val)
}

func Bool(key string, val bool) Field {
	return zap.Bool(key, val)
}

func Boolp(key string, val *bool) Field {
	return zap.Boolp(key, val)
}

func ByteString(key string, val []byte) Field {
	return zap.ByteString(key, val)
}

func Complex128(key string, val complex128) Field {
	return zap.Complex128(key, val)
}

func Complex128p(key string, val *complex128) Field {
	return zap.Complex128p(key, val)
}

func Complex64(key string, val complex64) Field {
	return zap.Complex64(key, val)
}

func Complex64p(key string, val *complex64) Field {
	return zap.Complex64p(key, val)
}

func Float64(key string, val float64) Field {
	return zap.Float64(key, val)
}

func Float64p(key string, val *float64) Field {
	return zap.Float64p(key, val)
}

func Float32(key string, val float32) Field {
	return zap.Float32(key, val)
}

func Float32p(key string, val *float32) Field {
	return zap.Float32p(key, val)
}

func Int(key string, val int) Field {
	return zap.Int(key, val)
}

func Intp(key string, val *int) Field {
	return zap.Intp(key, val)
}

func Int64(key string, val int64) Field {
	return zap.Int64(key, val)
}

func Int64p(key string, val *int64) Field {
	return zap.Int64p(key, val)
}

func Int32(key string, val int32) Field {
	return zap.Int32(key, val)
}

func Int32p(key string, val *int32) Field {
	return zap.Int32p(key, val)
}

func Int16(key string, val int16) Field {
	return zap.Int16(key, val)
}

func Int16p(key string, val *int16) Field {
	return zap.Int16p(key, val)
}

func Int8(key string, val int8) Field {
	return zap.Int8(key, val)
}

func Int8p(key string, val *int8) Field {
	return zap.Int8p(key, val)
}

func Stringp(key string, val *string) Field {
	return zap.Stringp(key, val)
}

func Uint(key string, val uint) Field {
	return zap.Uint(key, val)
}

func Uintp(key string, val *uint) Field {
	return zap.Uintp(key, val)
}

func Uint64(key string, val uint64) Field {
	return zap.Uint64(key, val)
}

func Uint64p(key string, val *uint64) Field {
	return zap.Uint64p(key, val)
}

func Uint32(key string, val uint32) Field {
	return zap.Uint32(key, val)
}

func Uint32p(key string, val *uint32) Field {
	return zap.Uint32p(key, val)
}

func Uint16(key string, val uint16) Field {
	return zap.Uint16(key, val)
}

func Uint16p(key string, val *uint16) Field {
	return zap.Uint16p(key, val)
}

func Uint8(key string, val uint8) Field {
	return zap.Uint8(key, val)
}

func Uint8p(key string, val *uint8) Field {
	return zap.Uint8p(key, val)
}

func Uintptr(key string, val uintptr) Field {
	return zap.Uintptr(key, val)
}

func Uintptrp(key string, val *uintptr) Field {
	return zap.Uintptrp(key, val)
}

func Reflect(key string, val interface{}) Field {
	return zap.Reflect(key, val)
}

func Namespace(key string) Field {
	return zap.Namespace(key)
}

func Stringer(key string, val fmt.Stringer) Field {
	return zap.Stringer(key, val)
}

func Time(key string, val time.Time) Field {
	return zap.Time(key, val)
}

func Timep(key string, val *time.Time) Field {
	return zap.Timep(key, val)
}

func Stack(key string) Field {
	return zap.Stack(key)
}

func StackSkip(key string, skip int) Field {
	return zap.StackSkip(key, skip)
}

func Duration(key string, val time.Duration) Field {
	return zap.Duration(key, val)
}

func Durationp(key string, val *time.Duration) Field {
	return zap.Durationp(key, val)
}

func Object(key string, val zapcore.ObjectMarshaler) Field {
	return zap.Object(key, val)
}

func Any(key string, value interface{}) Field {
	return zap.Any(key, value)
}
