package jsonb

import (
	"testing"
	"time"

	"github.com/si3nloong/sqlike/reflext"
	"github.com/stretchr/testify/require"
)

func TestDecodeByte(t *testing.T) {
	var (
		dec = Decoder{}
		r   *Reader
		x   []byte
		b   []byte
		err error
	)

	v := reflext.ValueOf(&x).Elem()

	r = NewReader([]byte(`""`))
	err = dec.DecodeByte(r, v)
	require.NoError(t, err)
	require.Equal(t, make([]byte, 0), x)

	r = NewReader([]byte(`null`))
	err = dec.DecodeByte(r, v)
	require.NoError(t, err)
	require.Equal(t, []byte(nil), x)

	b = []byte(`"VGhlIGlubGluZSB0YWJsZXMgYWJvdmUgYXJlIGlkZW50aWNhbCB0byB0aGUgZm9sbG93aW5nIHN0YW5kYXJkIHRhYmxlIGRlZmluaXRpb25zOg=="`)
	r = NewReader(b)
	err = dec.DecodeByte(r, v)
	require.NoError(t, err)
	require.Equal(t, []byte(`The inline tables above are identical to the following standard table definitions:`), x)
}

func TestDecodeTime(t *testing.T) {
	var (
		dec = Decoder{}
		r   *Reader
		dt  time.Time
		x   time.Time
		err error
	)

	v := reflext.ValueOf(&x).Elem()

	r = NewReader([]byte(`""`))
	err = dec.DecodeTime(r, v)
	require.NoError(t, err)
	require.Equal(t, time.Time{}, x)

	dt, _ = time.Parse("2006-01-02", "2018-01-02")
	r = NewReader([]byte(`"2018-01-02"`))
	err = dec.DecodeTime(r, v)
	require.NoError(t, err)
	require.Equal(t, dt, x)

	dt, _ = time.Parse("2006-01-02 15:04:05", "2018-01-02 13:58:26")
	r = NewReader([]byte(`"2018-01-02 13:58:26"`))
	err = dec.DecodeTime(r, v)
	require.NoError(t, err)
	require.Equal(t, dt, x)

	r = NewReader([]byte(`"2018-01-02 13:65:66"`))
	err = dec.DecodeTime(r, v)
	require.Error(t, err)
}

func TestDecodeMap(t *testing.T) {
	var (
		dec = Decoder{registry: buildRegistry()}
		r   *Reader
		x   map[string]interface{}
		err error
	)

	v := reflext.ValueOf(&x).Elem()
	r = NewReader([]byte(`null`))
	err = dec.DecodeMap(r, v)
	require.NoError(t, err)
	require.Equal(t, map[string]interface{}(nil), x)

	r = NewReader([]byte(`{}`))
	err = dec.DecodeMap(r, v)
	require.NoError(t, err)
	require.Equal(t, make(map[string]interface{}), x)

	r = NewReader([]byte(`
	{
		"a":"123", 
		"b":   108213312, 
		"c": true, 
		"d": "alSLKaj28173-021@#$%^&*\"",
		"e": 0.3127123
	}`))
	err = dec.DecodeMap(r, v)
	require.NoError(t, err)
	require.Equal(t, map[string]interface{}{
		"a": "123",
		"b": float64(108213312),
		"c": true,
		"d": `alSLKaj28173-021@#$%^&*"`,
		"e": float64(0.3127123),
	}, x)
}
