package mysql

import (
	"database/sql"
	"encoding/base64"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDecodeByte(a *testing.T) {
	var (
		dd  = DefaultDecoders{}
		err error
	)

	a.Run("Decode Byte with empty value", func(b *testing.T) {
		var raw []byte
		v := reflect.ValueOf(&raw)
		err = dd.DecodeByte(nil, v.Elem())
		require.NoError(b, err)
		require.NotNil(b, b)
		require.True(b, len(raw) == 0)
	})

	a.Run("Decode Byte with string value", func(b *testing.T) {
		var raw []byte
		v := reflect.ValueOf(&raw)
		name := "john doe"
		b64 := base64.StdEncoding.EncodeToString([]byte(name))
		err = dd.DecodeByte(b64, v.Elem())
		require.NoError(b, err)
		require.Equal(b, name, string(raw))
	})

	a.Run("Decode Byte with number value", func(b *testing.T) {
		var raw []byte
		v := reflect.ValueOf(&raw)
		num := "88"
		b64 := base64.StdEncoding.EncodeToString([]byte(num))
		err = dd.DecodeByte([]byte(b64), v.Elem())
		require.NoError(b, err)
		require.Equal(b, num, string(raw))
	})

}

func TestDecodeRawBytes(a *testing.T) {
	var (
		dd = DefaultDecoders{}
	)

	a.Run("Decode RawBytes to String", func(b *testing.T) {
		var raw sql.RawBytes
		str := "JOHN Cena"
		v := reflect.ValueOf(&raw)
		err := dd.DecodeRawBytes(str, v.Elem())
		require.NoError(b, err)
		require.Equal(b, sql.RawBytes(str), raw)
	})

	a.Run("Decode RawBytes to Int64", func(b *testing.T) {
		var raw sql.RawBytes
		v := reflect.ValueOf(&raw)
		i64 := int64(-1231298738213812)
		err := dd.DecodeRawBytes(i64, v.Elem())
		require.NoError(b, err)
		require.Equal(b, sql.RawBytes("-1231298738213812"), raw)
	})

	a.Run("Decode RawBytes to Uint64", func(b *testing.T) {
		var raw sql.RawBytes
		v := reflect.ValueOf(&raw)
		ui64 := uint64(1231298738213812)
		err := dd.DecodeRawBytes(ui64, v.Elem())
		require.NoError(b, err)
		require.Equal(b, sql.RawBytes("1231298738213812"), raw)
	})

	a.Run("Decode RawBytes to Boolean", func(b *testing.T) {
		var raw sql.RawBytes
		flag := true
		v := reflect.ValueOf(&raw)
		err := dd.DecodeRawBytes(flag, v.Elem())
		require.NoError(b, err)
		require.Equal(b, sql.RawBytes("true"), raw)

		flag = false
		err = dd.DecodeRawBytes(flag, v.Elem())
		require.NoError(b, err)
		require.Equal(b, sql.RawBytes("false"), raw)
	})

	a.Run("Decode RawBytes to Time", func(b *testing.T) {
		var raw sql.RawBytes
		v := reflect.ValueOf(&raw)
		tm := time.Time{}
		err := dd.DecodeRawBytes(tm, v.Elem())
		require.NoError(b, err)
		require.Equal(b, sql.RawBytes("0001-01-01T00:00:00Z"), raw)
	})
}

func TestDecodeTime(a *testing.T) {
	var (
		tm  time.Time
		err error
	)

	a.Run("Time with YYYY-MM-DD", func(b *testing.T) {
		tm, err = decodeTime("2021-10-17")
		require.NoError(b, err)
		require.Equal(b, "2021-10-17T00:00:00Z", tm.Format(time.RFC3339Nano))
	})

	a.Run("Time with 1 digit milliseconds", func(b *testing.T) {
		tm, err = decodeTime("2021-10-17 07:15:04.3")
		require.NoError(b, err)
		require.Equal(b, "2021-10-17T07:15:04.3Z", tm.Format(time.RFC3339Nano))

		tm, err = decodeTime("2021-10-17 07:15:04.30")
		require.NoError(b, err)
		require.Equal(b, "2021-10-17T07:15:04.3Z", tm.Format(time.RFC3339Nano))
	})

	a.Run("Time with 2 digit milliseconds", func(b *testing.T) {
		tm, err := decodeTime("2021-10-17 07:15:04.36")
		require.NoError(b, err)
		require.Equal(b, "2021-10-17T07:15:04.36Z", tm.Format(time.RFC3339Nano))
	})

	a.Run("Time with 3 digit milliseconds", func(b *testing.T) {
		tm, err := decodeTime("2021-10-17 07:15:04.366")
		require.NoError(b, err)
		require.Equal(b, "2021-10-17T07:15:04.366Z", tm.Format(time.RFC3339Nano))
	})

	a.Run("Time with 4 digit milliseconds", func(b *testing.T) {
		tm, err := decodeTime("2021-10-17 07:15:04.3661")
		require.NoError(b, err)
		require.Equal(b, "2021-10-17T07:15:04.3661Z", tm.Format(time.RFC3339Nano))
	})

	a.Run("Time with 5 digit milliseconds", func(b *testing.T) {
		tm, err := decodeTime("2021-10-17 07:15:04.36617")
		require.NoError(b, err)
		require.Equal(b, "2021-10-17T07:15:04.36617Z", tm.Format(time.RFC3339Nano))
	})

	a.Run("Time with 6 digit milliseconds", func(b *testing.T) {
		tm, err = decodeTime("2021-10-17 07:15:04.366170")
		require.NoError(b, err)
		require.Equal(b, "2021-10-17T07:15:04.36617Z", tm.Format(time.RFC3339Nano))

		tm, err = decodeTime("2021-10-17 07:15:04.366176")
		require.NoError(b, err)
		require.Equal(b, "2021-10-17T07:15:04.366176Z", tm.Format(time.RFC3339Nano))
	})
}
