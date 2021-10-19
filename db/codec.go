package db

import (
	"context"
	"reflect"
)

// Codecer :
type Codecer interface {
	RegisterTypeCodec(t reflect.Type, enc ValueEncoder, dec ValueDecoder)
	RegisterTypeEncoder(t reflect.Type, enc ValueEncoder)
	RegisterTypeDecoder(t reflect.Type, dec ValueDecoder)
	RegisterKindEncoder(k reflect.Kind, enc ValueEncoder)
	RegisterKindDecoder(k reflect.Kind, dec ValueDecoder)
	LookupEncoder(v reflect.Value) (ValueEncoder, error)
	LookupDecoder(t reflect.Type) (ValueDecoder, error)
}

// ValueDecoder :
type ValueDecoder func(interface{}, reflect.Value) error

// ValueEncoder :
type ValueEncoder func(context.Context, reflect.Value) (interface{}, error)

// ValueCodec :
type ValueCodec interface {
	DecodeValue(interface{}, reflect.Value) error
	EncodeValue(reflect.Value) (interface{}, error)
}
