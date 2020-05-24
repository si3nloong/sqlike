package codec

import (
	"database/sql/driver"
	"reflect"
	"testing"

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
}

func TestEncodeValue(t *testing.T) {
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
