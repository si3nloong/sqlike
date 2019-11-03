package reflext

import (
	"reflect"
)

// Init :
func Init(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Ptr && v.IsNil() {
		v.Set(reflect.New(v.Type().Elem()))
	}
	if v.Kind() == reflect.Map && v.IsNil() {
		v.Set(reflect.MakeMap(v.Type()))
	}
	if v.Kind() == reflect.Slice && v.IsNil() {
		v.Set(reflect.MakeSlice(v.Type(), 0, 0))
	}
	return v
}

// IndirectInit :
func IndirectInit(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}
	return v
}

// ValueOf : this is the replacement for reflect.ValueOf()
func ValueOf(i interface{}) reflect.Value {
	if x, ok := i.(reflect.Value); ok {
		return x
	}
	return reflect.ValueOf(i)
}

// TypeOf : this is the replacement for reflect.TypeOf()
func TypeOf(i interface{}) reflect.Type {
	if x, ok := i.(reflect.Type); ok {
		return x
	}
	return reflect.TypeOf(i)
}

// Deref : this is the replacement for reflect.Elem()
func Deref(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

// Indirect : this is the replacement for reflect.Indirect()
func Indirect(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			break
		}
		v = v.Elem()
	}
	return v
}

// IsNull :
func IsNull(v reflect.Value) bool {
	k := v.Kind()
	return (k == reflect.Ptr ||
		k == reflect.Slice ||
		k == reflect.Map ||
		k == reflect.Interface ||
		k == reflect.Func) && v.IsNil()
}

// Zero :
func Zero(t reflect.Type) (v reflect.Value) {
	v = reflect.New(t)
	vi := v.Elem()
	for vi.Kind() == reflect.Ptr && vi.IsNil() {
		vi.Set(reflect.New(vi.Type().Elem()))
		vi = vi.Elem()
	}
	return v.Elem()
}

// IsNullable :
func IsNullable(t reflect.Type) bool {
	k := t.Kind()
	return k == reflect.Ptr ||
		k == reflect.Slice ||
		k == reflect.Map ||
		k == reflect.Func ||
		k == reflect.Interface
}

// IsKind :
func IsKind(t reflect.Type, k reflect.Kind) bool {
	return t.Kind() == k
}

// IsZero :
func IsZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Func, reflect.Map:
		return v.IsNil()
	case reflect.Slice:
		return v.IsNil() || v.Len() == 0
	case reflect.Array:
		if v.Len() == 0 {
			return true
		}
		z := true
		for i := 0; i < v.Len(); i++ {
			z = z && IsZero(v.Index(i))
		}
		return z
	case reflect.Struct:
		z := true
		for i := 0; i < v.NumField(); i++ {
			z = z && IsZero(v.Field(i))
		}
		return z
	}
	// Compare other types directly:
	z := reflect.Zero(v.Type())
	return v.Interface() == z.Interface()
}

// Set :
func Set(src, v reflect.Value) {
	IndirectInit(src).Set(v)
}
