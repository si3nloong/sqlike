package expr

import (
	"reflect"

	"github.com/si3nloong/sqlike/x/primitive"
)

// String :
func String(val string) (o primitive.TypeSafe) {
	o.Type = reflect.String
	o.Value = val
	return
}

// Bool :
func Bool(val bool) (o primitive.TypeSafe) {
	o.Type = reflect.Bool
	o.Value = val
	return
}

// Int64 :
func Int64(val int64) (o primitive.TypeSafe) {
	o.Type = reflect.Int64
	o.Value = val
	return
}

// Int32 :
func Int32(val int32) (o primitive.TypeSafe) {
	o.Type = reflect.Int32
	o.Value = val
	return
}

// Int16 :
func Int16(val int16) (o primitive.TypeSafe) {
	o.Type = reflect.Int16
	o.Value = val
	return
}

// Int8 :
func Int8(val int8) (o primitive.TypeSafe) {
	o.Type = reflect.Int8
	o.Value = val
	return
}

// Int :
func Int(val int) (o primitive.TypeSafe) {
	o.Type = reflect.Int
	o.Value = val
	return
}

// Uint64 :
func Uint64(val uint64) (o primitive.TypeSafe) {
	o.Type = reflect.Uint64
	o.Value = val
	return
}

// Uint32 :
func Uint32(val uint32) (o primitive.TypeSafe) {
	o.Type = reflect.Uint32
	o.Value = val
	return
}

// Uint16 :
func Uint16(val uint16) (o primitive.TypeSafe) {
	o.Type = reflect.Uint16
	o.Value = val
	return
}

// Uint8 :
func Uint8(val uint8) (o primitive.TypeSafe) {
	o.Type = reflect.Uint8
	o.Value = val
	return
}

// Uint :
func Uint(val uint) (o primitive.TypeSafe) {
	o.Type = reflect.Uint
	o.Value = val
	return
}

// Float32 :
func Float32(val float32) (o primitive.TypeSafe) {
	o.Type = reflect.Float32
	o.Value = val
	return
}

// Float64 :
func Float64(val float64) (o primitive.TypeSafe) {
	o.Type = reflect.Float64
	o.Value = val
	return
}
