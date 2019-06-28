package reflext

import (
	"reflect"
)

// ValueOf : this is the replacement for reflect.ValueOf()
func ValueOf(i interface{}) reflect.Value {
	if x, isOk := i.(reflect.Value); isOk {
		return x
	}
	return reflect.ValueOf(i)
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
		k == reflect.Map) && v.IsNil()
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
	case reflect.Func, reflect.Map, reflect.Slice:
		return v.IsNil()
	case reflect.Array:
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
