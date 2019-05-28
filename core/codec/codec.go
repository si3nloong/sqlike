package codec

import (
	"reflect"
)

// ValueWriter :
type ValueWriter interface {
	WriteRune(rune) (int, error)
	Write([]byte) (int, error)
	WriteByte(byte) error
	WriteString(string) (int, error)
	Len() int
}

// ValueReader :
type ValueReader interface {
	ReadByte() (byte, error)
	Reset()
	Len() int
	Bytes() []byte
	String() string
}

// ValueDecoder :
type ValueDecoder func(ValueReader, reflect.Value) error

// ValueEncoder :
type ValueEncoder func(ValueWriter, reflect.Value) error

// ValueCodec :
type ValueCodec interface {
	DecodeValue(ValueReader, reflect.Value) error
	EncodeValue(ValueWriter, reflect.Value) error
}
