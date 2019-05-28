package codec

import (
	"reflect"
)

// ValueDecoder :
type ValueDecoder func(interface{}, reflect.Value) error

// ValueEncoder :
type ValueEncoder func(reflect.Value) (interface{}, error)

// ValueCodec :
type ValueCodec interface {
	DecodeValue(interface{}, reflect.Value) error
	EncodeValue(reflect.Value) (interface{}, error)
}
