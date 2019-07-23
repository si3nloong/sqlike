package jsonb

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncodeMap(t *testing.T) {
	var (
		v   reflect.Value
		enc Encoder
		w   *Writer
		err error
	)

	// uninitialized map
	{
		w = new(Writer)
		var uninitmap map[string]interface{}
		v = reflect.ValueOf(uninitmap)
		err = enc.EncodeMap(w, v)
		require.NoError(t, err)
		require.Equal(t, []byte(`null`), w.Bytes())
	}

	// initiliazed map
	{
		w = new(Writer)
		initmap := make(map[string]interface{})
		v = reflect.ValueOf(initmap)
		err = enc.EncodeMap(w, v)
		require.NoError(t, err)
		require.Equal(t, []byte(`{}`), w.Bytes())
	}

	{
		w = new(Writer)
		m := make(map[string]interface{})
		m["a"] = "v1"
		m["b"] = 123
		m["c"] = float64(10.888)
		m["d"] = true

		v = reflect.ValueOf(m)
		err = enc.EncodeMap(w, v)
		require.NoError(t, err)
		// require.Equal(t, []byte(`{}`), w.Bytes())
	}
}
