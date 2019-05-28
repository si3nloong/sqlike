package core

import (
	// "encoding/base64"
	"github.com/si3nloong/sqlike/reflext"
	"github.com/si3nloong/sqlike/types"

	// "github.com/si3nloong/sqlike/core/codec"
	"reflect"
	// "strconv"
)

var typeOfKey = reflect.TypeOf(types.Key{})

// DefaultMapper :
var DefaultMapper = reflext.NewMapperFunc("sqlike", func(sf *reflext.StructField) bool {
	t := reflext.Deref(sf.Zero.Type())
	return t == typeOfKey
})
