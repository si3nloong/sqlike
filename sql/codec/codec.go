package codec

import (
	"reflect"

	"github.com/si3nloong/sqlike/reflext"
)

// ValueDecoder :
type ValueDecoder func(interface{}, reflect.Value) error

// ValueEncoder :
type ValueEncoder func(reflext.StructFielder, reflect.Value) (interface{}, error)

// ValueCodec :
type ValueCodec interface {
	DecodeValue(interface{}, reflect.Value) error
	EncodeValue(reflect.Value) (interface{}, error)
}
