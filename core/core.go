package core

import (
	"reflect"

	"github.com/si3nloong/sqlike/reflext"
	"github.com/si3nloong/sqlike/types"
)

var typeOfKey = reflect.TypeOf(types.Key{})

// DefaultMapper :
var DefaultMapper = reflext.NewMapperFunc(
	"sqlike",
	func(sf *reflext.StructField) bool {
		t := reflext.Deref(sf.Zero.Type())
		return t == typeOfKey
	}, nil)
