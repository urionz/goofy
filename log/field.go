package log

import (
	"fmt"
	"time"

	"github.com/urionz/goofy/contracts"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Binary(key string, val []byte) contracts.LogField {
	return zap.Binary(key, val)
}

func Bool(key string, val bool) contracts.LogField {
	return zap.Bool(key, val)
}

func Boolp(key string, val *bool) contracts.LogField {
	return zap.Boolp(key, val)
}

func ByteString(key string, val []byte) contracts.LogField {
	return zap.ByteString(key, val)
}

func Complex128(key string, val complex128) contracts.LogField {
	return zap.Complex128(key, val)
}

func Complex128p(key string, val *complex128) contracts.LogField {
	return zap.Complex128p(key, val)
}

func Complex64(key string, val complex64) contracts.LogField {
	return zap.Complex64(key, val)
}

func Complex64p(key string, val *complex64) contracts.LogField {
	return zap.Complex64p(key, val)
}

func Float64(key string, val float64) contracts.LogField {
	return zap.Float64(key, val)
}

func Float64p(key string, val *float64) contracts.LogField {
	return zap.Float64p(key, val)
}

func Float32(key string, val float32) contracts.LogField {
	return zap.Float32(key, val)
}

func Float32p(key string, val *float32) contracts.LogField {
	return zap.Float32p(key, val)
}

func Int(key string, val int) contracts.LogField {
	return zap.Int(key, val)
}

func Intp(key string, val *int) contracts.LogField {
	return zap.Intp(key, val)
}

func Int64(key string, val int64) contracts.LogField {
	return zap.Int64(key, val)
}

func Int64p(key string, val *int64) contracts.LogField {
	return zap.Int64p(key, val)
}

func Int32(key string, val int32) contracts.LogField {
	return zap.Int32(key, val)
}

func Int32p(key string, val *int32) contracts.LogField {
	return zap.Int32p(key, val)
}

func Int16(key string, val int16) contracts.LogField {
	return zap.Int16(key, val)
}

func Int16p(key string, val *int16) contracts.LogField {
	return zap.Int16p(key, val)
}

func Int8(key string, val int8) contracts.LogField {
	return zap.Int8(key, val)
}

func Int8p(key string, val *int8) contracts.LogField {
	return zap.Int8p(key, val)
}

func Stringp(key string, val *string) contracts.LogField {
	return zap.Stringp(key, val)
}

func Uint(key string, val uint) contracts.LogField {
	return zap.Uint(key, val)
}

func Uintp(key string, val *uint) contracts.LogField {
	return zap.Uintp(key, val)
}

func Uint64(key string, val uint64) contracts.LogField {
	return zap.Uint64(key, val)
}

func Uint64p(key string, val *uint64) contracts.LogField {
	return zap.Uint64p(key, val)
}

func Uint32(key string, val uint32) contracts.LogField {
	return zap.Uint32(key, val)
}

func Uint32p(key string, val *uint32) contracts.LogField {
	return zap.Uint32p(key, val)
}

func Uint16(key string, val uint16) contracts.LogField {
	return zap.Uint16(key, val)
}

func Uint16p(key string, val *uint16) contracts.LogField {
	return zap.Uint16p(key, val)
}

func Uint8(key string, val uint8) contracts.LogField {
	return zap.Uint8(key, val)
}

func Uint8p(key string, val *uint8) contracts.LogField {
	return zap.Uint8p(key, val)
}

func Uintptr(key string, val uintptr) contracts.LogField {
	return zap.Uintptr(key, val)
}

func Uintptrp(key string, val *uintptr) contracts.LogField {
	return zap.Uintptrp(key, val)
}

func Reflect(key string, val interface{}) contracts.LogField {
	return zap.Reflect(key, val)
}

func Namespace(key string) contracts.LogField {
	return zap.Namespace(key)
}

func Stringer(key string, val fmt.Stringer) contracts.LogField {
	return zap.Stringer(key, val)
}

func Time(key string, val time.Time) contracts.LogField {
	return zap.Time(key, val)
}

func Timep(key string, val *time.Time) contracts.LogField {
	return zap.Timep(key, val)
}

func Stack(key string) contracts.LogField {
	return zap.Stack(key)
}

func StackSkip(key string, skip int) contracts.LogField {
	return zap.StackSkip(key, skip)
}

func Duration(key string, val time.Duration) contracts.LogField {
	return zap.Duration(key, val)
}

func Durationp(key string, val *time.Duration) contracts.LogField {
	return zap.Durationp(key, val)
}

func Object(key string, val zapcore.ObjectMarshaler) contracts.LogField {
	return zap.Object(key, val)
}

func Any(key string, value interface{}) contracts.LogField {
	return zap.Any(key, value)
}
