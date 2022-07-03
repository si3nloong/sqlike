package codec

import (
	"context"
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

	tByte := reflect.TypeOf([]byte{})
	byteEncoder := func(_ context.Context, v reflect.Value) (any, error) {
		return v.Bytes(), nil
	}
	rg.RegisterTypeCodec(tByte, byteEncoder, nil)
	// encoder, err := rg.LookupEncoder(reflect.ValueOf([]byte{}))
	// require.NoError(t, err)

	// require.Same(t, ValueEncoder(byteEncoder), encoder)
}

func TestEncodeValue(t *testing.T) {
	{
		it, err := encodeDriverValue(context.TODO(), reflect.ValueOf(nil))
		require.NoError(t, err)
		require.Nil(t, it)
	}

	{
		it, err := encodeDriverValue(context.TODO(), reflect.ValueOf("hello world"))
		require.Error(t, err)
		require.Nil(t, it)
	}

	{
		var it any
		it = num{}
		x := it.(driver.Valuer)
		it, err := encodeDriverValue(context.TODO(), reflect.ValueOf(x))
		require.NoError(t, err)
		require.Equal(t, int64(1000), it)
	}
}

func TestNilEncoder(t *testing.T) {
	var (
		v   reflect.Value
		it  any
		err error
	)

	{
		v = reflect.ValueOf(nil)
		it, err = NilEncoder(context.TODO(), v)
		require.NoError(t, err)
		require.Nil(t, it)
	}

	{
		str := "hello world"
		v = reflect.ValueOf(str)
		it, err = NilEncoder(context.TODO(), v)
		require.NoError(t, err)
		require.Nil(t, it)
	}
}
