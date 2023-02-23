package db

import (
	"reflect"
)

type SqlValuer interface {
	SqlValue(SqlDriver, map[string]string) (string, []any, error)
}

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

// ValueEncoder :
type ValueEncoder func(SqlDriver, reflect.Value, map[string]string) (string, []any, error)

// ValueDecoder :
type ValueDecoder func(any, reflect.Value) error

// ValueCodec :
type ValueCodec interface {
	DecodeValue(SqlDriver, reflect.Value, map[string]string) (string, []any, error)
	EncodeValue(reflect.Value) (any, error)
}
