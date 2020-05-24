package codec

import (
	"database/sql"
	"encoding/base64"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeByte(t *testing.T) {
	var (
		dd  = DefaultDecoders{}
		err error
	)

	{
		var b []byte
		v := reflect.ValueOf(&b)
		name := "john doe"
		b64 := base64.StdEncoding.EncodeToString([]byte(name))
		err = dd.DecodeByte(b64, v.Elem())
		require.NoError(t, err)
		require.Equal(t, name, string(b))
	}

	{
		var b []byte
		v := reflect.ValueOf(&b)
		num := "88"
		b64 := base64.StdEncoding.EncodeToString([]byte(num))
		err = dd.DecodeByte([]byte(b64), v.Elem())
		require.NoError(t, err)
		require.Equal(t, num, string(b))
	}

	{
		var b []byte
		v := reflect.ValueOf(&b)
		err = dd.DecodeByte(nil, v.Elem())
		require.NoError(t, err)
		require.NotNil(t, b)
		require.True(t, len(b) == 0)
	}
}

func TestDecodeRawBytes(t *testing.T) {
	var (
		dd = DefaultDecoders{}
	)

	{
		var b sql.RawBytes
		str := "JOHN Cena"
		v := reflect.ValueOf(&b)
		err := dd.DecodeRawBytes(str, v.Elem())
		require.NoError(t, err)
		require.Equal(t, sql.RawBytes(str), b)
	}

	{
		var b sql.RawBytes
		i64 := int64(1231298738213812)
		v := reflect.ValueOf(&b)
		err := dd.DecodeRawBytes(i64, v.Elem())
		require.NoError(t, err)
		require.Equal(t, sql.RawBytes("1231298738213812"), b)
	}

	{
		var b sql.RawBytes
		flag := true
		v := reflect.ValueOf(&b)
		err := dd.DecodeRawBytes(flag, v.Elem())
		require.NoError(t, err)
		require.Equal(t, sql.RawBytes("true"), b)

		flag = false
		err = dd.DecodeRawBytes(flag, v.Elem())
		require.NoError(t, err)
		require.Equal(t, sql.RawBytes("false"), b)
	}
}
