package codec

import (
	"context"
	"reflect"
)

// ValueDecoder :
type ValueDecoder func(context.Context, interface{}, reflect.Value) error

// ValueEncoder :
type ValueEncoder func(context.Context, reflect.Value) (interface{}, error)

// ValueCodec :
type ValueCodec interface {
	DecodeValue(interface{}, reflect.Value) error
	EncodeValue(reflect.Value) (interface{}, error)
}
