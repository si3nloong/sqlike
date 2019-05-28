package core

import (
	// "encoding/base64"
	"bitbucket.org/SianLoong/sqlike/reflext"
	"bitbucket.org/SianLoong/sqlike/types"

	// "bitbucket.org/SianLoong/sqlike/core/codec"
	"reflect"
	// "strconv"
)

var typeOfKey = reflect.TypeOf(types.Key{})

// DefaultMapper :
var DefaultMapper = reflext.NewMapperFunc("sqlike", func(sf *reflext.StructField) bool {
	t := reflext.Deref(sf.Zero.Type())
	return t == typeOfKey
})
