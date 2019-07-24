package codec

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncodeMap(t *testing.T) {
	var (
		enc = DefaultEncoders{}
		it  interface{}
		err error
	)

	// Nil map
	{
		var nilmap map[string]string
		it, err = enc.EncodeMap(nil, reflect.ValueOf(nilmap))
		require.NoError(t, err)
		require.Equal(t, `null`, it)
	}

	// Initialized map
	{
		initmap := make(map[string]int)
		it, err = enc.EncodeMap(nil, reflect.ValueOf(initmap))
		require.NoError(t, err)
		require.Equal(t, []byte(`{}`), it)
	}

	// Initialized map
	{
		initmap := make(map[string]int)
		it, err = enc.EncodeMap(nil, reflect.ValueOf(initmap))
		require.NoError(t, err)
		require.Equal(t, []byte(`{}`), it)
	}

	// Unsupported data type of map's key
	{
		intmap := make(map[int]interface{})
		it, err = enc.EncodeMap(nil, reflect.ValueOf(intmap))
		require.Error(t, err)
	}

}
