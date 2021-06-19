package codec

import (
	"database/sql/driver"
	"reflect"
	"testing"

	"github.com/si3nloong/sqlike/x/reflext"
	"github.com/stretchr/testify/require"
)

type num struct {
}

func (num) Value() (driver.Value, error) {
	return int64(1000), nil
}

func TestRegistry(t *testing.T) {
	rg := NewRegistry()

	kind := reflect.String
	rg.RegisterKindCodec(kind, nil, nil)
	rg.RegisterKindEncoder(kind, nil)
	rg.RegisterKindDecoder(kind, nil)

	typeof := reflect.TypeOf([]byte{})
	rg.RegisterTypeCodec(typeof, nil, nil)
	rg.RegisterTypeEncoder(typeof, nil)
	rg.RegisterTypeDecoder(typeof, nil)

	tByte := reflect.TypeOf([]byte{})
	byteEncoder := func(_ reflext.StructFielder, v reflect.Value) (interface{}, error) {
		return v.Bytes(), nil
	}
	rg.RegisterTypeCodec(tByte, byteEncoder, nil)
	// encoder, err := rg.LookupEncoder(reflect.ValueOf([]byte{}))
	// require.NoError(t, err)

	// require.Same(t, ValueEncoder(byteEncoder), encoder)
}

func TestEncodeValue(t *testing.T) {
	{
		it, err := encodeValue(nil, reflect.ValueOf(nil))
		require.NoError(t, err)
		require.Nil(t, it)
	}

	{
		it, err := encodeValue(nil, reflect.ValueOf("hello world"))
		require.Error(t, err)
		require.Nil(t, it)
	}

	{
		var it interface{}
		it = num{}
		x := it.(driver.Valuer)
		it, err := encodeValue(nil, reflect.ValueOf(x))
		require.NoError(t, err)
		require.Equal(t, int64(1000), it)
	}
}

func TestNilEncoder(t *testing.T) {
	var (
		v   reflect.Value
		it  interface{}
		err error
	)

	{
		v = reflect.ValueOf(nil)
		it, err = NilEncoder(nil, v)
		require.NoError(t, err)
		require.Nil(t, it)
	}

	{
		str := "hello world"
		v = reflect.ValueOf(str)
		it, err = NilEncoder(nil, v)
		require.NoError(t, err)
		require.Nil(t, it)
	}
}
